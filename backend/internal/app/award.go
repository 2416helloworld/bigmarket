package app

import (
	"log/slog"
	"strconv"
	"strings"

	"bigmarket/internal/model"
	"bigmarket/internal/pkg"
)

func (a *App) handleSendAward(body []byte) error {
	message, err := decodeEvent[SendAwardMessage](body)
	if err != nil {
		slog.Error("decode send award message failed", "err", err)
		return nil
	}
	award, err := a.Store.AwardByID(message.AwardID)
	if err != nil {
		slog.Error("award not found", "awardId", message.AwardID)
		return nil
	}
	config := message.AwardConfig
	if config == "" {
		config = award.AwardConfig
	}
	switch award.AwardKey {
	case "user_credit_random":
		amount := randomAmount(config)
		if amount > 0 {
			if err = a.CreditAdjust(message.UserID, pkg.TradeNameAward, pkg.TradeForward, amount, message.OutBusinessNo); err != nil {
				slog.Error("award credit adjust failed", "userId", message.UserID, "err", err)
				return err
			}
		}
	case "user_credit_blacklist":
		amount := randomAmount(config)
		if amount > 0 {
			if err = a.CreditAdjust(message.UserID, pkg.TradeNameAward, pkg.TradeForward, amount, message.OutBusinessNo); err != nil {
				slog.Error("blacklist credit adjust failed", "userId", message.UserID, "err", err)
				return err
			}
		}
	default:
		// 外部权益类奖品（OpenAI 次数/模型、实物等），此处记录为已发放，便于扩展对接
		slog.Info("external award distributed", "awardKey", award.AwardKey, "config", config, "userId", message.UserID)
	}
	res := a.Store.DB.Model(&model.UserAwardRecord{}).
		Where("user_id = ? and order_id = ? and award_state = ?", message.UserID, message.OutBusinessNo, pkg.StateCreate).
		Update("award_state", pkg.StateCompleted)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func randomAmount(config string) float64 {
	if config == "" {
		return 0
	}
	parts := strings.Split(config, ",")
	if len(parts) == 1 {
		v, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err != nil {
			return 0
		}
		return v
	}
	minV, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	maxV, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err1 != nil || err2 != nil {
		return 0
	}
	if maxV <= minV {
		return minV
	}
	return minV + float64(pkg.RandomInt(int(maxV-minV)+1))
}
