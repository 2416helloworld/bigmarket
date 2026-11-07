package app

import (
	"log/slog"
	"strconv"
	"strings"
	"time"

	"bigmarket/internal/model"
	"bigmarket/internal/pkg"
)

type SkuStockMessage struct {
	Sku        int64 `json:"sku"`
	ActivityID int64 `json:"activityId"`
}

// runStockUpdateJob 延迟队列 -> 最终一致更新数据库库存
func (a *App) runStockUpdateJob() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		for _, item := range a.Store.PopDueItems(pkg.ActivitySkuUpdateQueueKey, 100) {
			parts := strings.Split(item, ":")
			if len(parts) < 2 {
				continue
			}
			sku, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				continue
			}
			a.Store.UpdateSkuStock(sku)
		}
		for _, item := range a.Store.PopDueItems(pkg.StrategyAwardUpdateQueue, 200) {
			parts := strings.Split(item, ":")
			if len(parts) < 2 {
				continue
			}
			strategyID, err1 := strconv.ParseInt(parts[0], 10, 64)
			awardID, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil {
				continue
			}
			a.Store.UpdateStrategyAwardStock(strategyID, awardID)
		}
	}
}

// runTaskRetryJob 扫描未发送成功的任务，重新投递 MQ
func (a *App) runTaskRetryJob() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		tasks, err := a.Store.RetryTasks(20)
		if err != nil {
			continue
		}
		for _, task := range tasks {
			if err = a.MQ.PublishRaw(task.Topic, task.Message); err != nil {
				a.Store.UpdateTaskState(task.ID, pkg.StateFail)
				slog.Error("retry task publish failed", "taskId", task.ID, "err", err)
				continue
			}
			a.Store.UpdateTaskState(task.ID, pkg.StateCompleted)
			slog.Info("retry task published", "taskId", task.ID, "topic", task.Topic)
		}
	}
}

// handleSkuStockZero 商品库存耗尽时同步数据库库存
func (a *App) handleSkuStockZero(body []byte) error {
	message, err := decodeEvent[SkuStockMessage](body)
	if err != nil {
		return nil
	}
	a.Store.DB.Model(&model.RaffleActivitySku{}).Where("sku = ?", message.Sku).Update("stock_count_surplus", 0)
	return nil
}
