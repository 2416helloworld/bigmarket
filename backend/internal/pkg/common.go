package pkg

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	Code string      `json:"code"`
	Info string      `json:"info"`
	Data interface{} `json:"data,omitempty"`
}

func OK(data interface{}) Response { return Response{Code: "0000", Info: "成功", Data: data} }

func Fail(code string) Response {
	info, ok := CodeMap[code]
	if !ok {
		info = "系统异常"
	}
	return Response{Code: code, Info: info}
}

func FailWith(code, info string) Response { return Response{Code: code, Info: info} }

var CodeMap = map[string]string{
	"0000": "成功",
	"0001": "唯一索引冲突",
	"0002": "非法的参数",
	"0003": "需要登录",
	"0004": "未配置抽奖策略",
	"0005": "抽奖过于频繁",
	"0006": "账户额度不足",
	"0007": "活动未开启",
	"0008": "不在活动时间范围",
	"0009": "没有可用的抽奖次数",
	"0010": "扣减库存失败",
	"0011": "用户无可用抽奖单",
	"0012": "任务重试",
	"0013": "重复下单",
	"0014": "奖品库存不足",
	"0015": "积分账户不可用",
}

type BizError struct {
	Code string
	Info string
}

func NewBiz(code string) *BizError {
	info, ok := CodeMap[code]
	if !ok {
		info = code
	}
	return &BizError{Code: code, Info: info}
}

func NewBizInfo(code, info string) *BizError { return &BizError{Code: code, Info: info} }

func (e *BizError) Error() string { return fmt.Sprintf("%s: %s", e.Code, e.Info) }

const (
	StrategyAwardStockKey     = "big_market_strategy_award_stock_key_"
	AwardCountKey             = "big_market_strategy_award_count_key_"
	RateTableKey              = "big_market_strategy_rate_table_key_"
	RateRangeKey              = "big_market_strategy_rate_range_key_"
	RuleTreeVOKey             = "big_market_rule_tree_vo_key_"
	ActivityKey               = "big_market_raffle_activity_key_"
	SkuKey                    = "big_market_raffle_activity_sku_key_"
	SkuStockKey               = "big_market_raffle_activity_sku_stock_key_"
	SkuStockLockKey           = "big_market_raffle_activity_sku_stock_lock_key_"
	BlacklistKey              = "big_market_raffle_strategy_user_blacklist_"
	RateLimiterKey            = "big_market_raffle_strategy_user_rate_limiter_"
	ActivitySkuUpdateQueueKey = "big_market_activity_sku_stock_update_queue"
	StrategyAwardUpdateQueue  = "big_market_strategy_award_stock_update_queue"
	ActivitySkuUpdateLockKey  = "big_market_activity_sku_stock_update_lock"
)

const (
	TopicSendAward          = "big_market_send_award"
	TopicSendRebate         = "big_market_send_rebate"
	TopicCreditAdjust       = "big_market_credit_adjust_success"
	TopicActivitySkuStock   = "big_market_activity_sku_stock_zero"
	ExchangeName            = "big_market_exchange"
	DefaultRuleLogicModel   = "rule_default"
	RuleBlacklist           = "rule_blacklist"
	RuleWeight              = "rule_weight"
	RuleLock                = "rule_lock"
	RuleStock               = "rule_stock"
	RuleLuckAward           = "rule_luck_award"
	StateCreate             = "create"
	StateCompleted          = "completed"
	StateFail               = "fail"
	StateOpen               = "open"
	StateClose              = "close"
	StateWaitPay            = "wait_pay"
	StateUsed               = "used"
	TradeForward            = "forward"
	TradeReverse            = "reverse"
	RebateTypeSku           = "sku"
	RebateTypeIntegral      = "integral"
	TradeNameRebate         = "rebate"
	TradeNameAward          = "award"
	ChainResultAllow        = "ALLOW"
	ChainResultTakeOver     = "TAKE_OVER"
	ChainResultDeny         = "DENY"
)

type WeightRule struct {
	Weight   int
	AwardIDs []int
}

// ParseRuleWeight 解析 "10:102,103 70:106,107"
func ParseRuleWeight(value string) []WeightRule {
	var result []WeightRule
	for _, part := range strings.Fields(value) {
		idx := strings.Index(part, ":")
		if idx <= 0 {
			continue
		}
		weight, err := strconv.Atoi(part[:idx])
		if err != nil {
			continue
		}
		var ids []int
		for _, s := range strings.Split(part[idx+1:], ",") {
			if id, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
				ids = append(ids, id)
			}
		}
		result = append(result, WeightRule{Weight: weight, AwardIDs: ids})
	}
	return result
}

// ParseRuleBlacklist 解析 "101:user001,user002"
func ParseRuleBlacklist(value string) (int, []string) {
	idx := strings.Index(value, ":")
	if idx <= 0 {
		return 0, nil
	}
	awardID, _ := strconv.Atoi(value[:idx])
	users := strings.Split(value[idx+1:], ",")
	return awardID, users
}

// ParseLuckAward 解析 "101:1,100"
func ParseLuckAward(value string) (int, string) {
	idx := strings.Index(value, ":")
	if idx <= 0 {
		return 0, ""
	}
	awardID, _ := strconv.Atoi(value[:idx])
	return awardID, value[idx+1:]
}

// BuildRateTable 按概率构建抽奖概率表，使用最小概率精度放大后填充
func BuildRateTable(awards []RateAward) ([]int, int) {
	if len(awards) == 0 {
		return nil, 0
	}
	minRate := awards[0].Rate
	for _, a := range awards {
		if a.Rate < minRate {
			minRate = a.Rate
		}
	}
	rangeValue := 1
	if minRate > 0 {
		s := strconv.FormatFloat(minRate, 'f', -1, 64)
		if i := strings.Index(s, "."); i >= 0 {
			digits := len(s) - i - 1
			rangeValue = int(math.Pow10(digits))
			if rangeValue < 1 {
				rangeValue = 1
			}
		} else {
			rangeValue = 1
		}
	}
	var table []int
	for _, a := range awards {
		count := int(math.Round(float64(rangeValue) * a.Rate))
		if a.Rate > 0 && count < 1 {
			count = 1
		}
		for i := 0; i < count; i++ {
			table = append(table, a.AwardID)
		}
	}
	if len(table) == 0 {
		for _, a := range awards {
			table = append(table, a.AwardID)
		}
	}
	// 打乱概率表
	for i := len(table) - 1; i > 0; i-- {
		j, _ := randInt(i + 1)
		table[i], table[j] = table[j], table[i]
	}
	return table, len(table)
}

type RateAward struct {
	AwardID int
	Rate    float64
}

func NewRateAward(awardID int, rate float64) RateAward { return RateAward{AwardID: awardID, Rate: rate} }

func randInt(n int) (int, error) {
	if n <= 0 {
		return 0, nil
	}
	v, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return int(time.Now().UnixNano()) % n, nil
	}
	return int(v.Int64()), nil
}

func RandomInt(max int) int {
	v, _ := randInt(max)
	return v
}

func GenOrderID() string {
	return time.Now().Format("20060102150405") + fmt.Sprintf("%06d", RandomInt(1000000))
}

func GenMessageID() string {
	return fmt.Sprintf("%d%06d", time.Now().Unix(), RandomInt(1000000))
}

func Today() string { return time.Now().Format("2006-01-02") }

func ThisMonth() string { return time.Now().Format("2006-01") }
