package store

import (
	"bigmarket/internal/model"
)

func (s *Store) AwardByID(awardID int) (*model.Award, error) {
	var award model.Award
	if err := s.DB.Where("award_id = ?", awardID).First(&award).Error; err != nil {
		return nil, err
	}
	return &award, nil
}

func (s *Store) AwardList() ([]model.Award, error) {
	var list []model.Award
	err := s.DB.Order("award_id asc").Find(&list).Error
	return list, err
}

func (s *Store) UserAwardRecords(userID string, limit int) ([]model.UserAwardRecord, error) {
	var list []model.UserAwardRecord
	q := s.DB.Order("id desc").Limit(limit)
	if userID != "" {
		q = q.Where("user_id = ?", userID)
	}
	err := q.Find(&list).Error
	return list, err
}

func (s *Store) UserRaffleOrders(userID string, limit int) ([]model.UserRaffleOrder, error) {
	var list []model.UserRaffleOrder
	q := s.DB.Order("id desc").Limit(limit)
	if userID != "" {
		q = q.Where("user_id = ?", userID)
	}
	err := q.Find(&list).Error
	return list, err
}
