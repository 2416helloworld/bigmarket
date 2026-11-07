package infra

import (
	"context"
	"log/slog"
	"time"

	"bigmarket/internal/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Clients struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewClients(cfg config.Config) (*Clients, error) {
	db, err := openDB(cfg.DBDSN)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr, Password: cfg.RedisPassword})
	for i := 0; i < 30; i++ {
		if err = rdb.Ping(context.Background()).Err(); err == nil {
			break
		}
		slog.Warn("waiting for redis", "err", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, err
	}
	return &Clients{DB: db, Redis: rdb}, nil
}

func openDB(dsn string) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)
	for i := 0; i < 30; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Warn),
			DisableForeignKeyConstraintWhenMigrating: true,
			NamingStrategy:                           schema.NamingStrategy{SingularTable: true},
		})
		if err == nil {
			return db, nil
		}
		slog.Warn("waiting for mysql", "err", err)
		time.Sleep(2 * time.Second)
	}
	return nil, err
}
