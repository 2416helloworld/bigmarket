package store

import (
	"context"
	"fmt"
	"time"

	"bigmarket/internal/model"
	"bigmarket/internal/pkg"
)

type SkuProduct struct {
	Sku               int64   `json:"sku"`
	ActivityID        int64   `json:"activityId"`
	ActivityName      string  `json:"activityName"`
	ActivityCountID   int64   `json:"activityCountId"`
	StockCount        int     `json:"stockCount"`
	StockCountSurplus int     `json:"stockCountSurplus"`
	ProductAmount     float64 `json:"productAmount"`
	TotalCount        int     `json:"totalCount"`
	DayCount          int     `json:"dayCount"`
	MonthCount        int     `json:"monthCount"`
}

func (s *Store) ActivityByID(activityID int64) (*model.RaffleActivity, error) {
	cacheKey := fmt.Sprintf("%s%d", pkg.ActivityKey, activityID)
	var activity model.RaffleActivity
	if s.getJSON(cacheKey, &activity) {
		return &activity, nil
	}
	if err := s.DB.Where("activity_id = ?", activityID).First(&activity).Error; err != nil {
		return nil, err
	}
	_ = s.setJSON(cacheKey, activity, 10*time.Minute)
	return &activity, nil
}

func (s *Store) SkuBySku(sku int64) (*model.RaffleActivitySku, error) {
	cacheKey := fmt.Sprintf("%s%d", pkg.SkuKey, sku)
	var entity model.RaffleActivitySku
	if s.getJSON(cacheKey, &entity) {
		return &entity, nil
	}
	if err := s.DB.Where("sku = ?", sku).First(&entity).Error; err != nil {
		return nil, err
	}
	_ = s.setJSON(cacheKey, entity, 10*time.Minute)
	return &entity, nil
}

func (s *Store) SkuListByActivityID(activityID int64) ([]model.RaffleActivitySku, error) {
	var list []model.RaffleActivitySku
	err := s.DB.Where("activity_id = ?", activityID).Order("sku asc").Find(&list).Error
	return list, err
}

func (s *Store) SkuProductsByActivityID(activityID int64) ([]SkuProduct, error) {
	var list []SkuProduct
	err := s.DB.Raw(`select s.sku, s.activity_id, a.activity_name, s.activity_count_id, s.stock_count, s.stock_count_surplus,
			s.product_amount, c.total_count, c.day_count, c.month_count
		from raffle_activity_sku s
		join raffle_activity a on a.activity_id = s.activity_id
		left join raffle_activity_count c on c.activity_count_id = s.activity_count_id
		where s.activity_id = ? order by s.sku asc`, activityID).Scan(&list).Error
	return list, err
}

func (s *Store) CountByID(activityCountID int64) (*model.RaffleActivityCount, error) {
	var count model.RaffleActivityCount
	if err := s.DB.Where("activity_count_id = ?", activityCountID).First(&count).Error; err != nil {
		return nil, err
	}
	return &count, nil
}

func (s *Store) Account(userID string, activityID int64) *model.RaffleActivityAccount {
	var account model.RaffleActivityAccount
	if err := s.DB.Where("user_id = ? and activity_id = ?", userID, activityID).First(&account).Error; err != nil {
		return nil
	}
	return &account
}

func (s *Store) AccountDay(userID string, activityID int64) *model.RaffleActivityAccountDay {
	var account model.RaffleActivityAccountDay
	if err := s.DB.Where("user_id = ? and activity_id = ? and day = ?", userID, activityID, pkg.Today()).First(&account).Error; err != nil {
		return nil
	}
	return &account
}

func (s *Store) AccountMonth(userID string, activityID int64) *model.RaffleActivityAccountMonth {
	var account model.RaffleActivityAccountMonth
	if err := s.DB.Where("user_id = ? and activity_id = ? and month = ?", userID, activityID, pkg.ThisMonth()).First(&account).Error; err != nil {
		return nil
	}
	return &account
}

func (s *Store) NoUsedRaffleOrder(userID string, activityID int64) *model.UserRaffleOrder {
	var order model.UserRaffleOrder
	if err := s.DB.Where("user_id = ? and activity_id = ? and order_state = ?", userID, activityID, pkg.StateCreate).
		Order("id desc").First(&order).Error; err != nil {
		return nil
	}
	return &order
}

func (s *Store) UnpaidOrder(userID string, sku int64) *model.RaffleActivityOrder {
	var order model.RaffleActivityOrder
	if err := s.DB.Where("user_id = ? and sku = ? and state = ? and order_time >= ?", userID, sku, pkg.StateWaitPay, time.Now().AddDate(0, -1, 0)).
		Order("id desc").First(&order).Error; err != nil {
		return nil
	}
	return &order
}

func (s *Store) OrderByOutBusinessNo(userID, outBusinessNo string) *model.RaffleActivityOrder {
	var order model.RaffleActivityOrder
	if err := s.DB.Where("user_id = ? and out_business_no = ?", userID, outBusinessNo).First(&order).Error; err != nil {
		return nil
	}
	return &order
}

func (s *Store) CacheSkuStock(sku int64, surplus int) bool {
	key := fmt.Sprintf("%s%d", pkg.SkuStockKey, sku)
	return s.RDB.SetNX(context.Background(), key, surplus, 0).Val()
}

func (s *Store) SubtractSkuStock(sku int64) int {
	key := fmt.Sprintf("%s%d", pkg.SkuStockKey, sku)
	v, err := s.SubtractStock(key)
	if err != nil {
		return -1
	}
	if v < 0 {
		return v
	}
	return v
}

func (s *Store) QueueSkuStockUpdate(sku int64, activityID int64) {
	member := fmt.Sprintf("%d:%d:%d", sku, activityID, time.Now().UnixNano())
	s.PushDelayed(pkg.ActivitySkuUpdateQueueKey, member, 3*time.Second)
}

func (s *Store) UpdateSkuStock(sku int64) {
	s.DB.Exec("update raffle_activity_sku set stock_count_surplus = stock_count_surplus - 1 where sku = ? and stock_count_surplus > 0", sku)
}

func (s *Store) ClearActivityCache(activityID int64) {
	s.Del(fmt.Sprintf("%s%d", pkg.ActivityKey, activityID))
}

func (s *Store) ClearSkuCache(sku int64) {
	s.Del(fmt.Sprintf("%s%d", pkg.SkuKey, sku), fmt.Sprintf("%s%d", pkg.SkuStockKey, sku))
}
