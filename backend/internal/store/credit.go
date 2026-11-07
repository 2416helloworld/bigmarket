package store

import (
	"bigmarket/internal/model"
)

func (s *Store) CreditAccount(userID string) *model.UserCreditAccount {
	var account model.UserCreditAccount
	if err := s.DB.Where("user_id = ?", userID).First(&account).Error; err != nil {
		return nil
	}
	return &account
}

func (s *Store) CreditOrdersByUser(userID string, limit int) ([]model.UserCreditOrder, error) {
	var list []model.UserCreditOrder
	err := s.DB.Where("user_id = ?", userID).Order("id desc").Limit(limit).Find(&list).Error
	return list, err
}

func (s *Store) BehaviorRebateList(behaviorType string) ([]model.DailyBehaviorRebate, error) {
	var list []model.DailyBehaviorRebate
	err := s.DB.Where("behavior_type = ? and state = ?", behaviorType, "open").Find(&list).Error
	return list, err
}

func (s *Store) BehaviorRebateOrders(userID string, limit int) ([]model.UserBehaviorRebateOrder, error) {
	var list []model.UserBehaviorRebateOrder
	q := s.DB.Order("id desc").Limit(limit)
	if userID != "" {
		q = q.Where("user_id = ?", userID)
	}
	err := q.Find(&list).Error
	return list, err
}

func (s *Store) RetryTasks(limit int) ([]model.Task, error) {
	var list []model.Task
	err := s.DB.Where("state = ?", "create").Order("id asc").Limit(limit).Find(&list).Error
	return list, err
}

func (s *Store) UpdateTaskState(id int64, state string) {
	s.DB.Model(&model.Task{}).Where("id = ?", id).Update("state", state)
}

func (s *Store) TaskList(limit int) ([]model.Task, error) {
	var list []model.Task
	err := s.DB.Order("id desc").Limit(limit).Find(&list).Error
	return list, err
}
