package app

import (
	"log/slog"
	"time"

	"bigmarket/internal/infra"
	"bigmarket/internal/model"
	"bigmarket/internal/pkg"

	"gorm.io/gorm"
)

type CreditAdjustMessage struct {
	UserID        string `json:"userId"`
	OutBusinessNo string `json:"outBusinessNo"`
}

type SendAwardMessage struct {
	UserID        string `json:"userId"`
	AwardID       int    `json:"awardId"`
	AwardConfig   string `json:"awardConfig"`
	OutBusinessNo string `json:"outBusinessNo"`
}

type SendRebateMessage struct {
	UserID       string `json:"userId"`
	RebateType   string `json:"rebateType"`
	RebateConfig string `json:"rebateConfig"`
	BizID        string `json:"bizId"`
}

// CreditAdjust 积分调整（正向增加 / 反向扣减），幂等
func (a *App) CreditAdjust(userID, tradeName, tradeType string, amount float64, outBusinessNo string) error {
	if userID == "" || outBusinessNo == "" || amount == 0 {
		return pkg.NewBiz("0002")
	}
	var existed model.UserCreditOrder
	if err := a.Store.DB.Where("user_id = ? and out_business_no = ?", userID, outBusinessNo).First(&existed).Error; err == nil {
		return nil
	}
	if tradeType == pkg.TradeReverse {
		account := a.Store.CreditAccount(userID)
		if account == nil || account.AvailableAmount < amount {
			return pkg.NewBiz("0006")
		}
	}
	order := model.UserCreditOrder{
		UserID: userID, OrderID: pkg.GenOrderID(), TradeName: tradeName, TradeType: tradeType,
		TradeAmount: amount, OutBusinessNo: outBusinessNo,
	}
	if tradeType == pkg.TradeReverse {
		order.TradeAmount = -amount
	}
	var task model.Task
	err := a.Store.DB.Transaction(func(tx *gorm.DB) error {
		var account model.UserCreditAccount
		err := tx.Where("user_id = ?", userID).First(&account).Error
		if err == gorm.ErrRecordNotFound {
			if tradeType == pkg.TradeReverse {
				return pkg.NewBiz("0006")
			}
			account = model.UserCreditAccount{UserID: userID, TotalAmount: amount, AvailableAmount: amount, AccountStatus: "open"}
			if err = tx.Create(&account).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			var res *gorm.DB
			if tradeType == pkg.TradeReverse {
				res = tx.Exec(`update user_credit_account set total_amount = total_amount - ?, available_amount = available_amount - ?
					where user_id = ? and available_amount >= ?`, amount, amount, userID, amount)
			} else {
				res = tx.Exec(`update user_credit_account set total_amount = total_amount + ?, available_amount = available_amount + ?
					where user_id = ?`, amount, amount, userID)
			}
			if res.RowsAffected != 1 {
				return pkg.NewBiz("0006")
			}
		}
		if err = tx.Create(&order).Error; err != nil {
			return pkg.NewBiz("0001")
		}
		event := infra.Event{
			ID:        pkg.GenMessageID(),
			Timestamp: time.Now().UnixMilli(),
			Data:      CreditAdjustMessage{UserID: userID, OutBusinessNo: outBusinessNo},
		}
		message, _ := jsonMarshal(event)
		task = model.Task{UserID: userID, Topic: pkg.TopicCreditAdjust, MessageID: event.ID, Message: message, State: pkg.StateCreate}
		return tx.Create(&task).Error
	})
	if err != nil {
		return err
	}
	if err = a.MQ.PublishRaw(task.Topic, task.Message); err != nil {
		a.Store.UpdateTaskState(task.ID, pkg.StateFail)
	} else {
		a.Store.UpdateTaskState(task.ID, pkg.StateCompleted)
	}
	return nil
}

// QueryCreditAccount 查询用户积分账户
func (a *App) QueryCreditAccount(userID string) map[string]interface{} {
	available, total := 0.0, 0.0
	if account := a.Store.CreditAccount(userID); account != nil {
		available, total = account.AvailableAmount, account.TotalAmount
	}
	return map[string]interface{}{"userId": userID, "availableAmount": available, "totalAmount": total}
}

func (a *App) handleCreditAdjustSuccess(body []byte) error {
	message, err := decodeEvent[CreditAdjustMessage](body)
	if err != nil {
		slog.Error("decode credit adjust message failed", "err", err)
		return nil
	}
	slog.Info("credit adjust success, deliver order", "userId", message.UserID, "outBusinessNo", message.OutBusinessNo)
	err = a.UpdateOrder(message.UserID, message.OutBusinessNo)
	if err != nil {
		if biz, ok := err.(*pkg.BizError); ok && biz.Code == "0002" {
			// 非积分兑换订单（如签到积分），无需激活活动额度
			slog.Info("no activity order for credit adjust", "outBusinessNo", message.OutBusinessNo)
			return nil
		}
		return err
	}
	return nil
}
