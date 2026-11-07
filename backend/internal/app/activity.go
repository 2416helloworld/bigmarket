package app

import (
	"fmt"
	"log/slog"
	"time"

	"bigmarket/internal/infra"
	"bigmarket/internal/model"
	"bigmarket/internal/pkg"

	"gorm.io/gorm"
)

type DrawResult struct {
	OrderID     string `json:"orderId"`
	AwardID     int    `json:"awardId"`
	AwardTitle  string `json:"awardTitle"`
	AwardConfig string `json:"awardConfig"`
	AwardState  string `json:"awardState"`
}

type ActivityAccountVO struct {
	UserID            string `json:"userId"`
	ActivityID        int64  `json:"activityId"`
	ActivityName      string `json:"activityName"`
	TotalCount        int    `json:"totalCount"`
	TotalCountSurplus int    `json:"totalCountSurplus"`
	DayCount          int    `json:"dayCount"`
	DayCountSurplus   int    `json:"dayCountSurplus"`
	MonthCount        int    `json:"monthCount"`
	MonthCountSurplus int    `json:"monthCountSurplus"`
	Day               string `json:"day"`
	Month             string `json:"month"`
}

// AssembleActivity 活动装配：策略概率表 + SKU 缓存库存
func (a *App) AssembleActivity(activityID int64) error {
	activity, err := a.Store.ActivityByID(activityID)
	if err != nil {
		return pkg.NewBiz("0002")
	}
	if _, err = a.AssembleStrategy(activity.StrategyID); err != nil {
		return err
	}
	skus, err := a.Store.SkuListByActivityID(activityID)
	if err != nil {
		return err
	}
	for _, sku := range skus {
		a.Store.CacheSkuStock(sku.Sku, sku.StockCountSurplus)
	}
	return nil
}

// Draw 抽奖主流程：限流 -> 下单扣额度 -> 责任链/规则树 -> 中奖落库 -> 异步发奖
func (a *App) Draw(userID string, activityID int64) (*DrawResult, error) {
	if userID == "" {
		return nil, pkg.NewBiz("0002")
	}
	if err := a.RateLimiter(userID); err != nil {
		return nil, err
	}
	activity, err := a.Store.ActivityByID(activityID)
	if err != nil {
		return nil, pkg.NewBiz("0002")
	}
	order, err := a.Partake(userID, activity)
	if err != nil {
		return nil, err
	}
	raffleAward, err := a.PerformRaffle(userID, activity.StrategyID)
	if err != nil {
		_ = a.rollbackPartake(order)
		return nil, err
	}
	config := raffleAward.AwardConfig
	if config == "" {
		if award, err := a.Store.AwardByID(raffleAward.AwardID); err == nil {
			config = award.AwardConfig
		}
	}
	record := model.UserAwardRecord{
		UserID: userID, ActivityID: activityID, StrategyID: activity.StrategyID,
		OrderID: order.OrderID, AwardID: raffleAward.AwardID, AwardTitle: raffleAward.AwardTitle,
		AwardConfig: config, AwardState: pkg.StateCreate,
	}
	if err := a.saveAwardRecord(record); err != nil {
		_ = a.rollbackPartake(order)
		return nil, err
	}
	return &DrawResult{
		OrderID: order.OrderID, AwardID: raffleAward.AwardID, AwardTitle: raffleAward.AwardTitle,
		AwardConfig: config, AwardState: pkg.StateCreate,
	}, nil
}

// Partake 参与抽奖：校验活动、创建抽奖单并扣减额度
func (a *App) Partake(userID string, activity *model.RaffleActivity) (*model.UserRaffleOrder, error) {
	now := time.Now()
	if activity.State != pkg.StateOpen {
		return nil, pkg.NewBiz("0007")
	}
	if now.Before(activity.BeginDateTime) || now.After(activity.EndDateTime) {
		return nil, pkg.NewBiz("0008")
	}
	if order := a.Store.NoUsedRaffleOrder(userID, activity.ActivityID); order != nil {
		return order, nil
	}
	account := a.Store.Account(userID, activity.ActivityID)
	if account == nil || account.TotalCountSurplus <= 0 {
		return nil, pkg.NewBiz("0009")
	}
	order := &model.UserRaffleOrder{
		UserID: userID, ActivityID: activity.ActivityID, StrategyID: activity.StrategyID,
		OrderID: pkg.GenOrderID(), OrderTime: now, OrderState: pkg.StateCreate,
	}
	err := a.Store.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Exec(`update raffle_activity_account set total_count_surplus = total_count_surplus - 1
			where user_id = ? and activity_id = ? and total_count_surplus > 0`, userID, activity.ActivityID)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return pkg.NewBiz("0009")
		}
		// 月账户
		var month model.RaffleActivityAccountMonth
		err := tx.Where("user_id = ? and activity_id = ? and month = ?", userID, activity.ActivityID, pkg.ThisMonth()).First(&month).Error
		if err == gorm.ErrRecordNotFound {
			month = model.RaffleActivityAccountMonth{
				UserID: userID, ActivityID: activity.ActivityID, Month: pkg.ThisMonth(),
				MonthCount: account.MonthCountSurplus, MonthCountSurplus: account.MonthCountSurplus - 1,
			}
			if month.MonthCountSurplus < 0 {
				month.MonthCountSurplus = 0
			}
			if err = tx.Create(&month).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			res = tx.Exec(`update raffle_activity_account_month set month_count_surplus = month_count_surplus - 1
				where user_id = ? and activity_id = ? and month = ? and month_count_surplus > 0`, userID, activity.ActivityID, pkg.ThisMonth())
			if res.RowsAffected != 1 {
				return pkg.NewBiz("0009")
			}
		}
		// 日账户
		var day model.RaffleActivityAccountDay
		err = tx.Where("user_id = ? and activity_id = ? and day = ?", userID, activity.ActivityID, pkg.Today()).First(&day).Error
		if err == gorm.ErrRecordNotFound {
			day = model.RaffleActivityAccountDay{
				UserID: userID, ActivityID: activity.ActivityID, Day: pkg.Today(),
				DayCount: account.DayCountSurplus, DayCountSurplus: account.DayCountSurplus - 1,
			}
			if day.DayCountSurplus < 0 {
				day.DayCountSurplus = 0
			}
			if err = tx.Create(&day).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			res = tx.Exec(`update raffle_activity_account_day set day_count_surplus = day_count_surplus - 1
				where user_id = ? and activity_id = ? and day = ? and day_count_surplus > 0`, userID, activity.ActivityID, pkg.Today())
			if res.RowsAffected != 1 {
				return pkg.NewBiz("0009")
			}
		}
		return tx.Create(order).Error
	})
	if err != nil {
		if _, ok := err.(*pkg.BizError); ok {
			return nil, err
		}
		return nil, pkg.NewBiz("0002")
	}
	return order, nil
}

// rollbackPartake 抽奖失败时返还额度
func (a *App) rollbackPartake(order *model.UserRaffleOrder) error {
	if order == nil {
		return nil
	}
	return a.Store.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? and order_id = ?", order.UserID, order.OrderID).Delete(&model.UserRaffleOrder{}).Error; err != nil {
			return err
		}
		tx.Exec(`update raffle_activity_account set total_count_surplus = total_count_surplus + 1 where user_id = ? and activity_id = ?`, order.UserID, order.ActivityID)
		tx.Exec(`update raffle_activity_account_month set month_count_surplus = month_count_surplus + 1 where user_id = ? and activity_id = ? and month = ?`, order.UserID, order.ActivityID, pkg.ThisMonth())
		tx.Exec(`update raffle_activity_account_day set day_count_surplus = day_count_surplus + 1 where user_id = ? and activity_id = ? and day = ?`, order.UserID, order.ActivityID, pkg.Today())
		return nil
	})
}

func (a *App) saveAwardRecord(record model.UserAwardRecord) error {
	var task model.Task
	err := a.Store.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return pkg.NewBiz("0011")
		}
		event := infra.Event{
			ID:        pkg.GenMessageID(),
			Timestamp: time.Now().UnixMilli(),
			Data: map[string]interface{}{
				"userId": record.UserID, "awardId": record.AwardID,
				"awardConfig": record.AwardConfig, "outBusinessNo": record.OrderID,
			},
		}
		message, _ := jsonMarshal(event)
		task = model.Task{UserID: record.UserID, Topic: pkg.TopicSendAward, MessageID: event.ID, Message: message, State: pkg.StateCreate}
		if err := tx.Create(&task).Error; err != nil {
			return err
		}
		res := tx.Model(&model.UserRaffleOrder{}).
			Where("user_id = ? and order_id = ? and order_state = ?", record.UserID, record.OrderID, pkg.StateCreate).
			Update("order_state", pkg.StateUsed)
		if res.RowsAffected != 1 {
			return pkg.NewBiz("0011")
		}
		return nil
	})
	if err != nil {
		return err
	}
	// 事务外发送 MQ，失败由任务表补偿
	if err := a.MQ.PublishRaw(task.Topic, task.Message); err != nil {
		slog.Error("send award mq failed", "err", err)
		a.Store.UpdateTaskState(task.ID, pkg.StateFail)
	} else {
		a.Store.UpdateTaskState(task.ID, pkg.StateCompleted)
	}
	return nil
}

// QueryUserAccount 查询用户活动账户
func (a *App) QueryUserAccount(userID string, activityID int64) (*ActivityAccountVO, error) {
	activity, err := a.Store.ActivityByID(activityID)
	if err != nil {
		return nil, pkg.NewBiz("0002")
	}
	vo := &ActivityAccountVO{UserID: userID, ActivityID: activityID, ActivityName: activity.ActivityName, Day: pkg.Today(), Month: pkg.ThisMonth()}
	if account := a.Store.Account(userID, activityID); account != nil {
		vo.TotalCount = account.TotalCount
		vo.TotalCountSurplus = account.TotalCountSurplus
		vo.DayCount = account.DayCount
		vo.DayCountSurplus = account.DayCountSurplus
		vo.MonthCount = account.MonthCount
		vo.MonthCountSurplus = account.MonthCountSurplus
	}
	if day := a.Store.AccountDay(userID, activityID); day != nil {
		vo.DayCount = day.DayCount
		vo.DayCountSurplus = day.DayCountSurplus
	}
	if month := a.Store.AccountMonth(userID, activityID); month != nil {
		vo.MonthCount = month.MonthCount
		vo.MonthCountSurplus = month.MonthCountSurplus
	}
	return vo, nil
}

// CreateQuotaOrder 创建额度充值/兑换订单
func (a *App) CreateQuotaOrder(userID string, sku int64, tradeType string, outBusinessNo string) (*model.RaffleActivityOrder, error) {
	if outBusinessNo != "" {
		if order := a.Store.OrderByOutBusinessNo(userID, outBusinessNo); order != nil {
			return order, nil
		}
	}
	skuEntity, err := a.Store.SkuBySku(sku)
	if err != nil {
		return nil, pkg.NewBiz("0002")
	}
	activity, err := a.Store.ActivityByID(skuEntity.ActivityID)
	if err != nil {
		return nil, pkg.NewBiz("0002")
	}
	if activity.State != pkg.StateOpen || time.Now().After(activity.EndDateTime) {
		return nil, pkg.NewBiz("0007")
	}
	count, err := a.Store.CountByID(skuEntity.ActivityCountID)
	if err != nil {
		return nil, pkg.NewBiz("0002")
	}
	if tradeType == "credit_pay" {
		if order := a.Store.UnpaidOrder(userID, sku); order != nil {
			return order, nil
		}
		credit := a.Store.CreditAccount(userID)
		if credit == nil || credit.AvailableAmount < skuEntity.ProductAmount {
			return nil, pkg.NewBiz("0006")
		}
	}
	// 扣减 Redis 库存
	if a.Store.SubtractSkuStock(sku) != 1 {
		return nil, pkg.NewBiz("0010")
	}
	a.Store.QueueSkuStockUpdate(sku, skuEntity.ActivityID)
	order := &model.RaffleActivityOrder{
		UserID: userID, Sku: sku, ActivityID: skuEntity.ActivityID, ActivityName: activity.ActivityName,
		StrategyID: activity.StrategyID, OrderID: pkg.GenOrderID(), OrderTime: time.Now(),
		TotalCount: count.TotalCount, DayCount: count.DayCount, MonthCount: count.MonthCount,
		PayAmount: skuEntity.ProductAmount, State: pkg.StateWaitPay, OutBusinessNo: outBusinessNo,
	}
	if order.OutBusinessNo == "" {
		order.OutBusinessNo = pkg.GenOrderID()
	}
	if tradeType == "rebate_no_pay" {
		order.State = pkg.StateCompleted
		order.PayAmount = 0
	}
	err = a.Store.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return pkg.NewBiz("0013")
		}
		if order.State == pkg.StateCompleted {
			return a.addQuota(tx, userID, order.ActivityID, order.TotalCount, order.DayCount, order.MonthCount)
		}
		return nil
	})
	if err != nil {
		// 失败归还库存
		a.Store.RDB.Incr(ctxBackground(), fmt.Sprintf("%s%d", pkg.SkuStockKey, sku))
		return nil, err
	}
	if tradeType == "credit_pay" {
		if err = a.CreditAdjust(userID, "exchange", pkg.TradeReverse, skuEntity.ProductAmount, order.OutBusinessNo); err != nil {
			return nil, err
		}
	}
	return order, nil
}

// UpdateOrder 积分支付回调：订单完成并发货（充值额度）
func (a *App) UpdateOrder(userID, outBusinessNo string) error {
	order := a.Store.OrderByOutBusinessNo(userID, outBusinessNo)
	if order == nil {
		return pkg.NewBiz("0002")
	}
	if order.State == pkg.StateCompleted {
		return nil
	}
	return a.Store.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&model.RaffleActivityOrder{}).Where("user_id = ? and out_business_no = ? and state = ?", userID, outBusinessNo, pkg.StateWaitPay).
			Update("state", pkg.StateCompleted)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return nil
		}
		return a.addQuota(tx, userID, order.ActivityID, order.TotalCount, order.DayCount, order.MonthCount)
	})
}

// addQuota 增加用户活动额度
func (a *App) addQuota(tx *gorm.DB, userID string, activityID int64, total, day, month int) error {
	if err := tx.Exec(`insert into raffle_activity_account
		(user_id, activity_id, total_count, total_count_surplus, day_count, day_count_surplus, month_count, month_count_surplus)
		values (?, ?, ?, ?, ?, ?, ?, ?)
		on duplicate key update total_count = total_count + ?, total_count_surplus = total_count_surplus + ?,
			day_count = day_count + ?, day_count_surplus = day_count_surplus + ?,
			month_count = month_count + ?, month_count_surplus = month_count_surplus + ?`,
		userID, activityID, total, total, day, day, month, month, total, total, day, day, month, month).Error; err != nil {
		return err
	}
	if err := tx.Exec(`insert into raffle_activity_account_day
		(user_id, activity_id, day, day_count, day_count_surplus) values (?, ?, ?, ?, ?)
		on duplicate key update day_count = day_count + ?, day_count_surplus = day_count_surplus + ?`,
		userID, activityID, pkg.Today(), day, day, day, day).Error; err != nil {
		return err
	}
	return tx.Exec(`insert into raffle_activity_account_month
		(user_id, activity_id, month, month_count, month_count_surplus) values (?, ?, ?, ?, ?)
		on duplicate key update month_count = month_count + ?, month_count_surplus = month_count_surplus + ?`,
		userID, activityID, pkg.ThisMonth(), month, month, month, month).Error
}

// ActivityList 开放中的活动列表
func (a *App) ActivityList() ([]model.RaffleActivity, error) {
	var list []model.RaffleActivity
	now := time.Now()
	err := a.Store.DB.Where("state = ? and begin_date_time <= ? and end_date_time >= ?", pkg.StateOpen, now, now).
		Order("activity_id asc").Find(&list).Error
	return list, err
}

// QuerySkuProducts 活动商品列表（附带 Redis 剩余库存）
func (a *App) QuerySkuProducts(activityID int64) ([]map[string]interface{}, error) {
	products, err := a.Store.SkuProductsByActivityID(activityID)
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	for _, p := range products {
		stock := p.StockCountSurplus
		if v, err := a.Store.RDB.Get(ctxBackground(), fmt.Sprintf("%s%d", pkg.SkuStockKey, p.Sku)).Int(); err == nil {
			stock = v
		}
		result = append(result, map[string]interface{}{
			"sku": p.Sku, "activityId": p.ActivityID, "activityName": p.ActivityName,
			"stockCount": p.StockCount, "stockCountSurplus": stock, "productAmount": p.ProductAmount,
			"totalCount": p.TotalCount, "dayCount": p.DayCount, "monthCount": p.MonthCount,
		})
	}
	return result, nil
}
