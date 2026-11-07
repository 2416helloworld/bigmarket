package app

import (
	"context"
	"encoding/json"
	"log/slog"

	"bigmarket/internal/config"
	"bigmarket/internal/infra"
	"bigmarket/internal/pkg"
	"bigmarket/internal/store"
)

type App struct {
	Store *store.Store
	MQ    *infra.MQ
	Cfg   config.Config
}

func New(st *store.Store, mq *infra.MQ, cfg config.Config) *App {
	return &App{Store: st, MQ: mq, Cfg: cfg}
}

func jsonMarshal(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}

func ctxBackground() context.Context { return context.Background() }

func decodeEvent[T any](body []byte) (T, error) {
	var zero T
	var event infra.Event
	if err := json.Unmarshal(body, &event); err != nil {
		return zero, err
	}
	raw, err := json.Marshal(event.Data)
	if err != nil {
		return zero, err
	}
	var data T
	if err = json.Unmarshal(raw, &data); err != nil {
		return zero, err
	}
	return data, nil
}

// Start 启动 MQ 消费者与后台任务
func (a *App) Start() {
	go a.MQ.Subscribe(pkg.TopicSendAward, a.handleSendAward)
	go a.MQ.Subscribe(pkg.TopicSendRebate, a.handleSendRebate)
	go a.MQ.Subscribe(pkg.TopicCreditAdjust, a.handleCreditAdjustSuccess)
	go a.MQ.Subscribe(pkg.TopicActivitySkuStock, a.handleSkuStockZero)
	go a.runStockUpdateJob()
	go a.runTaskRetryJob()
	slog.Info("big-market-go started")
}
