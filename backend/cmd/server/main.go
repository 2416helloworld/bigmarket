package main

import (
	"log/slog"
	"strconv"

	"bigmarket/internal/api"
	"bigmarket/internal/app"
	"bigmarket/internal/config"
	"bigmarket/internal/infra"
	"bigmarket/internal/store"
)

func main() {
	cfg := config.Load()
	clients, err := infra.NewClients(cfg)
	if err != nil {
		slog.Error("init mysql/redis failed", "err", err)
		panic(err)
	}
	st := store.New(clients.DB, clients.Redis)
	mq := infra.NewMQ(cfg.RabbitURL)
	application := app.New(st, mq, cfg)
	application.Start()

	router := api.NewRouter(application, cfg)
	slog.Info("http server listening", "port", cfg.Port)
	if err = router.Run(":" + strconv.Itoa(cfg.Port)); err != nil {
		slog.Error("http server exit", "err", err)
	}
}
