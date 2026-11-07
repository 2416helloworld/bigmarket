package store

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"bigmarket/internal/model"
	"bigmarket/internal/pkg"

	"github.com/redis/go-redis/v9"
)

type RuleTreeVO struct {
	TreeID   string                       `json:"treeId"`
	TreeName string                       `json:"treeName"`
	TreeRoot string                       `json:"treeRoot"`
	Nodes    map[string]*RuleTreeNodeVO   `json:"nodes"`
}

type RuleTreeNodeVO struct {
	RuleKey   string                   `json:"ruleKey"`
	RuleValue string                   `json:"ruleValue"`
	Lines     []RuleTreeNodeLineVO     `json:"lines"`
}

type RuleTreeNodeLineVO struct {
	From       string `json:"from"`
	To         string `json:"to"`
	LimitType  string `json:"limitType"`
	LimitValue string `json:"limitValue"`
}

func (s *Store) StrategyAwardList(strategyID int64) ([]model.StrategyAward, error) {
	var list []model.StrategyAward
	err := s.DB.Where("strategy_id = ?", strategyID).Order("sort asc").Find(&list).Error
	return list, err
}

func (s *Store) StrategyByID(strategyID int64) (*model.Strategy, error) {
	var st model.Strategy
	err := s.DB.Where("strategy_id = ?", strategyID).First(&st).Error
	if err != nil {
		return nil, err
	}
	return &st, nil
}

func (s *Store) StrategyList() ([]model.Strategy, error) {
	var list []model.Strategy
	err := s.DB.Order("strategy_id asc").Find(&list).Error
	return list, err
}

func (s *Store) RulesByStrategy(strategyID int64) ([]model.StrategyRule, error) {
	var list []model.StrategyRule
	err := s.DB.Where("strategy_id = ?", strategyID).Order("id asc").Find(&list).Error
	return list, err
}

func (s *Store) RuleValue(strategyID int64, ruleModel string) string {
	var rule model.StrategyRule
	err := s.DB.Where("strategy_id = ? and rule_model = ?", strategyID, ruleModel).First(&rule).Error
	if err != nil {
		return ""
	}
	return rule.RuleValue
}

func (s *Store) TodayRaffleCount(userID string, strategyID int64) int {
	var dayCount, daySurplus int
	err := s.DB.Raw(`select d.day_count, d.day_count_surplus
		from raffle_activity_account_day d join raffle_activity ra on ra.activity_id = d.activity_id
		where d.user_id = ? and ra.strategy_id = ? and d.day = ? limit 1`, userID, strategyID, pkg.Today()).
		Row().Scan(&dayCount, &daySurplus)
	if err != nil {
		return 0
	}
	return dayCount - daySurplus
}

func (s *Store) AccountTotalUseCount(userID string, strategyID int64) int {
	var useCount int
	s.DB.Raw(`select a.total_count - a.total_count_surplus as use_count
		from raffle_activity_account a join raffle_activity ra on ra.activity_id = a.activity_id
		where a.user_id = ? and ra.strategy_id = ? limit 1`, userID, strategyID).Scan(&useCount)
	return useCount
}

func (s *Store) RuleTreeVO(treeID string) (*RuleTreeVO, error) {
	if treeID == "" {
		return nil, fmt.Errorf("rule tree id is empty")
	}
	cacheKey := pkg.RuleTreeVOKey + treeID
	var vo RuleTreeVO
	if s.getJSON(cacheKey, &vo) {
		return &vo, nil
	}
	var tree model.RuleTree
	if err := s.DB.Where("tree_id = ?", treeID).First(&tree).Error; err != nil {
		return nil, err
	}
	var nodes []model.RuleTreeNode
	if err := s.DB.Where("tree_id = ?", treeID).Find(&nodes).Error; err != nil {
		return nil, err
	}
	var lines []model.RuleTreeNodeLine
	if err := s.DB.Where("tree_id = ?", treeID).Find(&lines).Error; err != nil {
		return nil, err
	}
	result := &RuleTreeVO{TreeID: tree.TreeID, TreeName: tree.TreeName, TreeRoot: tree.TreeNodeRuleKey, Nodes: map[string]*RuleTreeNodeVO{}}
	for _, n := range nodes {
		result.Nodes[n.RuleKey] = &RuleTreeNodeVO{RuleKey: n.RuleKey, RuleValue: n.RuleValue}
	}
	for _, l := range lines {
		if node, ok := result.Nodes[l.RuleNodeFrom]; ok {
			node.Lines = append(node.Lines, RuleTreeNodeLineVO{From: l.RuleNodeFrom, To: l.RuleNodeTo, LimitType: l.RuleLimitType, LimitValue: l.RuleLimitValue})
		}
	}
	_ = s.setJSON(cacheKey, result, 10*time.Minute)
	return result, nil
}

func (s *Store) RuleLockCount(treeIDs []string) map[string]int {
	result := map[string]int{}
	if len(treeIDs) == 0 {
		return result
	}
	var nodes []model.RuleTreeNode
	s.DB.Where("tree_id in ? and rule_key = ?", treeIDs, pkg.RuleLock).Find(&nodes)
	for _, n := range nodes {
		if v, err := strconv.Atoi(strings.TrimSpace(n.RuleValue)); err == nil {
			result[n.TreeID] = v
		}
	}
	return result
}

func (s *Store) AwardTitle(strategyID int64, awardID int) string {
	var award model.StrategyAward
	s.DB.Where("strategy_id = ? and award_id = ?", strategyID, awardID).First(&award)
	return award.AwardTitle
}

func (s *Store) AwardRuleModels(strategyID int64, awardID int) string {
	var award model.StrategyAward
	if err := s.DB.Where("strategy_id = ? and award_id = ?", strategyID, awardID).First(&award).Error; err != nil {
		return ""
	}
	return award.RuleModels
}

func (s *Store) ClearStrategyCache(strategyID int64) {
	s.Del(pkg.RateTableKey+strconv.FormatInt(strategyID, 10), pkg.RateRangeKey+strconv.FormatInt(strategyID, 10))
	for _, weight := range []string{"10", "60", "70", "200", "1000"} {
		s.Del(pkg.RateTableKey + strconv.FormatInt(strategyID, 10) + "_" + weight)
	}
}

func (s *Store) ClearRuleTreeCache(treeID string) {
	s.Del(pkg.RuleTreeVOKey + treeID)
}

// CacheAwardStock 缓存奖品库存，返回是否首次缓存
func (s *Store) CacheAwardStock(strategyID int64, awardID, surplus int) bool {
	key := fmt.Sprintf("%s%d_%d", pkg.AwardCountKey, strategyID, awardID)
	return s.RDB.SetNX(context.Background(), key, surplus, 0).Val()
}

func (s *Store) CacheRateTable(key string, table []int, rangeValue int) {
	b, _ := json.Marshal(table)
	ctx := context.Background()
	s.RDB.Set(ctx, pkg.RateTableKey+key, b, 0)
	s.RDB.Set(ctx, pkg.RateRangeKey+key, rangeValue, 0)
}

func (s *Store) RandomAwardByTable(key string) (int, error) {
	ctx := context.Background()
	rangeValue, err := s.RDB.Get(ctx, pkg.RateRangeKey+key).Int()
	if err != nil || rangeValue <= 0 {
		return 0, pkg.NewBiz("0004")
	}
	idx := pkg.RandomInt(rangeValue)
	b, err := s.RDB.Get(ctx, pkg.RateTableKey+key).Bytes()
	if err != nil {
		return 0, pkg.NewBiz("0004")
	}
	var table []int
	if err = json.Unmarshal(b, &table); err != nil || len(table) == 0 {
		return 0, pkg.NewBiz("0004")
	}
	if idx >= len(table) {
		idx = len(table) - 1
	}
	return table[idx], nil
}

var stockLua = redis.NewScript(`
local stock = redis.call('get', KEYS[1])
if not stock then return -1 end
stock = tonumber(stock)
if stock <= 0 then return 0 end
redis.call('decr', KEYS[1])
return 1
`)

func (s *Store) SubtractStock(key string) (int, error) {
	v, err := stockLua.Run(context.Background(), s.RDB, []string{key}).Int()
	return v, err
}

func (s *Store) QueryAwardCount(key string) (int, bool) {
	v, err := s.RDB.Get(context.Background(), pkg.AwardCountKey+key).Int()
	return v, err == nil
}

// SubtractAwardStock 扣减奖品缓存库存并写入延迟更新队列
func (s *Store) SubtractAwardStock(strategyID int64, awardID int) bool {
	key := fmt.Sprintf("%s%d_%d", pkg.AwardCountKey, strategyID, awardID)
	v, err := s.SubtractStock(key)
	if err != nil || v != 1 {
		if v < 0 {
			return false
		}
		return false
	}
	member := fmt.Sprintf("%d:%d:%d", strategyID, awardID, time.Now().UnixNano())
	s.PushDelayed(pkg.StrategyAwardUpdateQueue, member, 3*time.Second)
	return true
}

func (s *Store) UpdateStrategyAwardStock(strategyID int64, awardID int) {
	s.DB.Exec("update strategy_award set award_count_surplus = award_count_surplus - 1 where strategy_id = ? and award_id = ? and award_count_surplus > 0", strategyID, awardID)
}
