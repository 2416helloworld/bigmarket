package app

import (
	"context"
	"strconv"
	"strings"
	"time"

	"bigmarket/internal/model"
	"bigmarket/internal/pkg"
)

type RaffleAward struct {
	AwardID     int    `json:"awardId"`
	AwardTitle  string `json:"awardTitle"`
	AwardConfig string `json:"awardConfig"`
	RuleModel   string `json:"ruleModel"`
}

// AssembleStrategy 策略装配：奖品库存、概率表、权重概率表写入 Redis
func (a *App) AssembleStrategy(strategyID int64) ([]model.StrategyAward, error) {
	strategy, err := a.Store.StrategyByID(strategyID)
	if err != nil {
		return nil, pkg.NewBiz("0004")
	}
	awards, err := a.Store.StrategyAwardList(strategyID)
	if err != nil || len(awards) == 0 {
		return nil, pkg.NewBiz("0004")
	}
	for _, award := range awards {
		a.Store.CacheAwardStock(strategyID, award.AwardID, award.AwardCountSurplus)
	}
	rateAwards := make([]pkg.RateAward, 0, len(awards))
	for _, award := range awards {
		rateAwards = append(rateAwards, pkg.NewRateAward(award.AwardID, award.AwardRate))
	}
	table, rangeValue := pkg.BuildRateTable(rateAwards)
	key := strconv.FormatInt(strategyID, 10)
	a.Store.CacheRateTable(key, table, rangeValue)

	if strings.Contains(strategy.RuleModels, pkg.RuleWeight) {
		ruleValue := a.Store.RuleValue(strategyID, pkg.RuleWeight)
		for _, weight := range pkg.ParseRuleWeight(ruleValue) {
			var filtered []pkg.RateAward
			for _, award := range awards {
				for _, awardID := range weight.AwardIDs {
					if award.AwardID == awardID {
						filtered = append(filtered, pkg.NewRateAward(award.AwardID, award.AwardRate))
					}
				}
			}
			if len(filtered) == 0 {
				continue
			}
			wTable, wRange := pkg.BuildRateTable(filtered)
			a.Store.CacheRateTable(key+"_"+strconv.Itoa(weight.Weight), wTable, wRange)
		}
	}
	return awards, nil
}

type chainResult struct {
	LogicModel  string
	AwardID     int
	AwardConfig string
}

type treeResult struct {
	Level       string
	HasAward    bool
	AwardID     int
	AwardConfig string
}

// PerformRaffle 责任链 + 规则树抽奖
func (a *App) PerformRaffle(userID string, strategyID int64) (*RaffleAward, error) {
	strategy, err := a.Store.StrategyByID(strategyID)
	if err != nil {
		return nil, pkg.NewBiz("0004")
	}
	awardID := 0
	ruleModel := pkg.DefaultRuleLogicModel
	ruleConfig := ""
	if strings.TrimSpace(strategy.RuleModels) != "" {
		chain, err := a.performChain(userID, strategyID, strategy.RuleModels)
		if err != nil {
			return nil, err
		}
		awardID = chain.AwardID
		ruleModel = chain.LogicModel
		ruleConfig = chain.AwardConfig
		// 规则树准入
		finalAwardID, finalConfig, err := a.performTree(userID, strategyID, awardID)
		if err != nil {
			return nil, err
		}
		if finalAwardID > 0 {
			awardID = finalAwardID
			ruleConfig = finalConfig
			ruleModel = pkg.RuleLuckAward
		}
	} else {
		awardID, err = a.Store.RandomAwardByTable(strconv.FormatInt(strategyID, 10))
		if err != nil {
			return nil, err
		}
	}
	title := a.Store.AwardTitle(strategyID, awardID)
	return &RaffleAward{AwardID: awardID, AwardTitle: title, AwardConfig: ruleConfig, RuleModel: ruleModel}, nil
}

func (a *App) performChain(userID string, strategyID int64, ruleModels string) (*chainResult, error) {
	models := strings.Split(ruleModels, ",")
	models = append(models, pkg.DefaultRuleLogicModel)
	for _, ruleModel := range models {
		switch strings.TrimSpace(ruleModel) {
		case pkg.RuleBlacklist:
			if awardID, config, hit := a.checkBlacklist(userID, strategyID); hit {
				return &chainResult{LogicModel: pkg.RuleBlacklist, AwardID: awardID, AwardConfig: config}, nil
			}
		case pkg.RuleWeight:
			if result, hit := a.checkWeight(userID, strategyID); hit {
				return result, nil
			}
		case pkg.DefaultRuleLogicModel:
			awardID, err := a.Store.RandomAwardByTable(strconv.FormatInt(strategyID, 10))
			if err != nil {
				return nil, err
			}
			return &chainResult{LogicModel: pkg.DefaultRuleLogicModel, AwardID: awardID}, nil
		}
	}
	return nil, pkg.NewBiz("0004")
}

func (a *App) checkBlacklist(userID string, strategyID int64) (int, string, bool) {
	ctx := context.Background()
	key := pkg.BlacklistKey + userID
	if a.Store.RDB.Exists(ctx, key).Val() == 1 {
		ruleValue := a.Store.RuleValue(strategyID, pkg.RuleBlacklist)
		awardID, _ := pkg.ParseRuleBlacklist(ruleValue)
		return awardID, "1,100", true
	}
	ruleValue := a.Store.RuleValue(strategyID, pkg.RuleBlacklist)
	if ruleValue == "" {
		return 0, "", false
	}
	awardID, users := pkg.ParseRuleBlacklist(ruleValue)
	for _, u := range users {
		if strings.TrimSpace(u) == userID {
			a.Store.RDB.Set(ctx, key, "1", 24*time.Hour)
			return awardID, "1,100", true
		}
	}
	return 0, "", false
}

func (a *App) checkWeight(userID string, strategyID int64) (*chainResult, bool) {
	ruleValue := a.Store.RuleValue(strategyID, pkg.RuleWeight)
	if ruleValue == "" {
		return nil, false
	}
	score := a.Store.AccountTotalUseCount(userID, strategyID)
	for _, weight := range pkg.ParseRuleWeight(ruleValue) {
		if weight.Weight != score {
			continue
		}
		key := strconv.FormatInt(strategyID, 10) + "_" + strconv.Itoa(weight.Weight)
		awardID, err := a.Store.RandomAwardByTable(key)
		if err != nil {
			return nil, false
		}
		return &chainResult{LogicModel: pkg.RuleWeight, AwardID: awardID}, true
	}
	return nil, false
}

func (a *App) performTree(userID string, strategyID int64, awardID int) (int, string, error) {
	treeID := a.Store.AwardRuleModels(strategyID, awardID)
	if treeID == "" {
		return 0, "", nil
	}
	tree, err := a.Store.RuleTreeVO(treeID)
	if err != nil {
		return 0, "", nil
	}
	nodeKey := tree.TreeRoot
	lockCount := a.lockCount(treeID)
	for i := 0; i < 10; i++ {
		node, ok := tree.Nodes[nodeKey]
		if !ok {
			return 0, "", nil
		}
		var result treeResult
		switch node.RuleKey {
		case pkg.RuleLock:
			required := 0
			if v, ok := lockCount[treeID]; ok {
				required = v
			} else if v, err := strconv.Atoi(strings.TrimSpace(node.RuleValue)); err == nil {
				required = v
			}
			today := a.Store.TodayRaffleCount(userID, strategyID)
			if today < required {
				result = treeResult{Level: pkg.ChainResultTakeOver}
			} else {
				result = treeResult{Level: pkg.ChainResultAllow}
			}
		case pkg.RuleStock:
			if a.Store.SubtractAwardStock(strategyID, awardID) {
				result = treeResult{Level: pkg.ChainResultTakeOver, HasAward: true, AwardID: awardID}
			} else {
				result = treeResult{Level: pkg.ChainResultAllow}
			}
		case pkg.RuleLuckAward:
			luckAwardID, config := pkg.ParseLuckAward(node.RuleValue)
			if a.Store.SubtractAwardStock(strategyID, luckAwardID) {
				result = treeResult{Level: pkg.ChainResultTakeOver, HasAward: true, AwardID: luckAwardID, AwardConfig: config}
			} else {
				return 0, "", pkg.NewBiz("0014")
			}
		default:
			return 0, "", nil
		}
		if result.HasAward {
			return result.AwardID, result.AwardConfig, nil
		}
		nextKey := ""
		for _, line := range node.Lines {
			if line.LimitType == "EQUAL" && line.LimitValue == result.Level {
				nextKey = line.To
				break
			}
		}
		if nextKey == "" {
			break
		}
		nodeKey = nextKey
	}
	return 0, "", nil
}

func (a *App) lockCount(treeID string) map[string]int {
	return a.Store.RuleLockCount([]string{treeID})
}

// QueryAwardList 查询活动奖品列表（含缓存库存）
func (a *App) QueryAwardList(activityID int64) ([]model.StrategyAward, error) {
	activity, err := a.Store.ActivityByID(activityID)
	if err != nil {
		return nil, pkg.NewBiz("0002")
	}
	return a.Store.StrategyAwardList(activity.StrategyID)
}

// QueryAwardListByStrategy 查询策略奖品列表
func (a *App) QueryAwardListByStrategy(strategyID int64) ([]model.StrategyAward, error) {
	return a.Store.StrategyAwardList(strategyID)
}

// QueryWeightRule 查询策略权重规则
func (a *App) QueryWeightRule(activityID int64) (map[string]interface{}, error) {
	activity, err := a.Store.ActivityByID(activityID)
	if err != nil {
		return nil, pkg.NewBiz("0002")
	}
	ruleValue := a.Store.RuleValue(activity.StrategyID, pkg.RuleWeight)
	result := map[string]interface{}{"strategyId": activity.StrategyID, "ruleValue": ruleValue}
	var weights []map[string]interface{}
	for _, weight := range pkg.ParseRuleWeight(ruleValue) {
		var ids []map[string]interface{}
		for _, awardID := range weight.AwardIDs {
			ids = append(ids, map[string]interface{}{"awardId": awardID, "awardTitle": a.Store.AwardTitle(activity.StrategyID, awardID)})
		}
		weights = append(weights, map[string]interface{}{"weight": weight.Weight, "awardIds": ids})
	}
	result["weights"] = weights
	return result, nil
}

// RateLimiter 抽奖限流：同一用户每秒最多 N 次，超出则拉黑 24 小时
func (a *App) RateLimiter(userID string) error {
	ctx := context.Background()
	if a.Store.RDB.Exists(ctx, pkg.BlacklistKey+userID).Val() == 1 {
		return pkg.NewBiz("0005")
	}
	key := pkg.RateLimiterKey + userID
	count, err := a.Store.RDB.Incr(ctx, key).Result()
	if err != nil {
		return nil
	}
	if count == 1 {
		a.Store.RDB.Expire(ctx, key, time.Second)
	}
	if int(count) > a.Cfg.DrawRateLimit {
		a.Store.RDB.Set(ctx, pkg.BlacklistKey+userID, "1", 24*time.Hour)
		return pkg.NewBiz("0005")
	}
	return nil
}
