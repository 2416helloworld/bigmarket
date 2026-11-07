package store

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Store struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func New(db *gorm.DB, rdb *redis.Client) *Store { return &Store{DB: db, RDB: rdb} }

func (s *Store) setJSON(key string, v interface{}, ttl time.Duration) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.RDB.Set(context.Background(), key, b, ttl).Err()
}

func (s *Store) getJSON(key string, dest interface{}) bool {
	b, err := s.RDB.Get(context.Background(), key).Bytes()
	if err != nil {
		return false
	}
	return json.Unmarshal(b, dest) == nil
}

func (s *Store) Del(keys ...string) {
	if len(keys) > 0 {
		s.RDB.Del(context.Background(), keys...)
	}
}

// PopDueItems 从 ZSET 延迟队列取出到期任务
func (s *Store) PopDueItems(key string, limit int64) []string {
	ctx := context.Background()
	max := strconv.FormatInt(time.Now().Unix(), 10)
	items, err := s.RDB.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:   "-inf",
		Max:   max,
		Count: limit,
	}).Result()
	if len(items) == 0 || err != nil {
		return nil
	}
	for _, item := range items {
		s.RDB.ZRem(ctx, key, item)
	}
	return items
}

func (s *Store) PushDelayed(key, member string, delay time.Duration) {
	s.RDB.ZAdd(context.Background(), key, redis.Z{
		Score:  float64(time.Now().Add(delay).Unix()),
		Member: member,
	})
}
