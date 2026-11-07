package app

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"bigmarket/internal/infra"
	"bigmarket/internal/model"
	"bigmarket/internal/pkg"

	"gorm.io/gorm"
)

// SignRebate 签到返利：同一用户同一天每种返利仅发一次
func (a *App) SignRebate(userID string) (bool, error) {
	rebates, err := a.Store.BehaviorRebateList("sign")
	if err != nil || len(rebates) == 0 {
		return false, pkg.NewBiz("0002")
	}
	for _, rebate := range rebates {
		outBusinessNo := fmt.Sprintf("%s_%s_%s_%s", userID, rebate.BehaviorType, rebate.RebateType, time.Now().Format("20060102"))
		var existed model.UserBehaviorRebateOrder
		if err := a.Store.DB.Where("user_id = ? and out_business_no = ?", userID, outBusinessNo).First(&existed).Error; err == nil {
			continue
		}
		amount := 0.0
		if rebate.RebateType == pkg.RebateTypeIntegral {
			amount, _ = strconv.ParseFloat(rebate.RebateConfig, 64)
		}
		order := model.UserBehaviorRebateOrder{
			UserID: userID, OrderID: pkg.GenOrderID(), BehaviorType: rebate.BehaviorType,
			RebateDesc: rebate.RebateDesc, RebateType: rebate.RebateType, RebateAmount: amount,
			BizID: outBusinessNo, OutBusinessNo: outBusinessNo,
		}
		var task model.Task
		err = a.Store.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&order).Error; err != nil {
				return err
			}
			event := infra.Event{
				ID:        pkg.GenMessageID(),
				Timestamp: time.Now().UnixMilli(),
				Data: SendRebateMessage{
					UserID: userID, RebateType: rebate.RebateType,
					RebateConfig: rebate.RebateConfig, BizID: outBusinessNo,
				},
			}
			message, _ := jsonMarshal(event)
			task = model.Task{UserID: userID, Topic: pkg.TopicSendRebate, MessageID: event.ID, Message: message, State: pkg.StateCreate}
			return tx.Create(&task).Error
		})
		if err != nil {
			return false, pkg.NewBiz("0002")
		}
		if err = a.MQ.PublishRaw(task.Topic, task.Message); err != nil {
			a.Store.UpdateTaskState(task.ID, pkg.StateFail)
		} else {
			a.Store.UpdateTaskState(task.ID, pkg.StateCompleted)
		}
	}
	return true, nil
}

func (a *App) IsSignRebate(userID string) bool {
	start := time.Now()
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	var count int64
	a.Store.DB.Model(&model.UserBehaviorRebateOrder{}).
		Where("user_id = ? and behavior_type = ? and create_time >= ?", userID, "sign", start).Count(&count)
	return count > 0
}

func (a *App) handleSendRebate(body []byte) error {
	message, err := decodeEvent[SendRebateMessage](body)
	if err != nil {
		slog.Error("decode rebate message failed", "err", err)
		return nil
	}
	switch message.RebateType {
	case pkg.RebateTypeSku:
		sku, err := strconv.ParseInt(message.RebateConfig, 10, 64)
		if err != nil {
			return nil
		}
		_, err = a.CreateQuotaOrder(message.UserID, sku, "rebate_no_pay", message.BizID)
		if err != nil {
			slog.Error("rebate sku order failed", "err", err)
			return err
		}
	case pkg.RebateTypeIntegral:
		amount, err := strconv.ParseFloat(message.RebateConfig, 64)
		if err != nil {
			return nil
		}
		return a.CreditAdjust(message.UserID, pkg.TradeNameRebate, pkg.TradeForward, amount, message.BizID)
	}
	return nil
}
