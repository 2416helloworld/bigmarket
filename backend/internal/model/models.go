package model

import "time"

type Award struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	AwardID     int       `json:"awardId"`
	AwardKey    string    `json:"awardKey"`
	AwardConfig string    `json:"awardConfig"`
	AwardDesc   string    `json:"awardDesc"`
	CreateTime  time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime  time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (Award) TableName() string { return "award" }

type DailyBehaviorRebate struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	BehaviorType string    `json:"behaviorType"`
	RebateDesc   string    `json:"rebateDesc"`
	RebateType   string    `json:"rebateType"`
	RebateConfig string    `json:"rebateConfig"`
	State        string    `json:"state"`
	CreateTime   time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime   time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (DailyBehaviorRebate) TableName() string { return "daily_behavior_rebate" }

type RaffleActivity struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	ActivityID    int64     `json:"activityId"`
	ActivityName  string    `json:"activityName"`
	ActivityDesc  string    `json:"activityDesc"`
	BeginDateTime time.Time `gorm:"column:begin_date_time" json:"beginDateTime"`
	EndDateTime   time.Time `gorm:"column:end_date_time" json:"endDateTime"`
	StrategyID    int64     `json:"strategyId"`
	State         string    `json:"state"`
	CreateTime    time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime    time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (RaffleActivity) TableName() string { return "raffle_activity" }

type RaffleActivityCount struct {
	ID              int64     `gorm:"primaryKey" json:"id"`
	ActivityCountID int64     `json:"activityCountId"`
	TotalCount      int       `json:"totalCount"`
	DayCount        int       `json:"dayCount"`
	MonthCount      int       `json:"monthCount"`
	CreateTime      time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime      time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (RaffleActivityCount) TableName() string { return "raffle_activity_count" }

type RaffleActivitySku struct {
	ID                int64     `gorm:"primaryKey" json:"id"`
	Sku               int64     `json:"sku"`
	ActivityID        int64     `json:"activityId"`
	ActivityCountID   int64     `json:"activityCountId"`
	StockCount        int       `json:"stockCount"`
	StockCountSurplus int       `json:"stockCountSurplus"`
	ProductAmount     float64   `json:"productAmount"`
	CreateTime        time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime        time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (RaffleActivitySku) TableName() string { return "raffle_activity_sku" }

type Strategy struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	StrategyID   int64     `json:"strategyId"`
	StrategyDesc string    `json:"strategyDesc"`
	RuleModels   string    `json:"ruleModels"`
	CreateTime   time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime   time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (Strategy) TableName() string { return "strategy" }

type StrategyAward struct {
	ID                int64     `gorm:"primaryKey" json:"id"`
	StrategyID        int64     `json:"strategyId"`
	AwardID           int       `json:"awardId"`
	AwardTitle        string    `json:"awardTitle"`
	AwardSubtitle     string    `json:"awardSubtitle"`
	AwardCount        int       `json:"awardCount"`
	AwardCountSurplus int       `json:"awardCountSurplus"`
	AwardRate         float64   `json:"awardRate"`
	RuleModels        string    `json:"ruleModels"`
	Sort              int       `json:"sort"`
	CreateTime        time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime        time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (StrategyAward) TableName() string { return "strategy_award" }

type StrategyRule struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	StrategyID int64     `json:"strategyId"`
	AwardID    int       `json:"awardId"`
	RuleType   int       `json:"ruleType"`
	RuleModel  string    `json:"ruleModel"`
	RuleValue  string    `json:"ruleValue"`
	RuleDesc   string    `json:"ruleDesc"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (StrategyRule) TableName() string { return "strategy_rule" }

type RuleTree struct {
	ID               int64     `gorm:"primaryKey" json:"id"`
	TreeID           string    `json:"treeId"`
	TreeName         string    `json:"treeName"`
	TreeDesc         string    `json:"treeDesc"`
	TreeNodeRuleKey  string    `gorm:"column:tree_node_rule_key" json:"treeNodeRuleKey"`
	CreateTime       time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime       time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (RuleTree) TableName() string { return "rule_tree" }

type RuleTreeNode struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	TreeID     string    `json:"treeId"`
	RuleKey    string    `json:"ruleKey"`
	RuleDesc   string    `json:"ruleDesc"`
	RuleValue  string    `json:"ruleValue"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (RuleTreeNode) TableName() string { return "rule_tree_node" }

type RuleTreeNodeLine struct {
	ID             int64     `gorm:"primaryKey" json:"id"`
	TreeID         string    `json:"treeId"`
	RuleNodeFrom   string    `json:"ruleNodeFrom"`
	RuleNodeTo     string    `json:"ruleNodeTo"`
	RuleLimitType  string    `json:"ruleLimitType"`
	RuleLimitValue string    `json:"ruleLimitValue"`
	CreateTime     time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime     time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (RuleTreeNodeLine) TableName() string { return "rule_tree_node_line" }

type RaffleActivityAccount struct {
	ID                int64     `gorm:"primaryKey" json:"id"`
	UserID            string    `json:"userId"`
	ActivityID        int64     `json:"activityId"`
	TotalCount        int       `json:"totalCount"`
	TotalCountSurplus int       `json:"totalCountSurplus"`
	DayCount          int       `json:"dayCount"`
	DayCountSurplus   int       `json:"dayCountSurplus"`
	MonthCount        int       `json:"monthCount"`
	MonthCountSurplus int       `json:"monthCountSurplus"`
	CreateTime        time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime        time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (RaffleActivityAccount) TableName() string { return "raffle_activity_account" }

type RaffleActivityAccountDay struct {
	ID              int64     `gorm:"primaryKey" json:"id"`
	UserID          string    `json:"userId"`
	ActivityID      int64     `json:"activityId"`
	Day             string    `json:"day"`
	DayCount        int       `json:"dayCount"`
	DayCountSurplus int       `json:"dayCountSurplus"`
	CreateTime      time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime      time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (RaffleActivityAccountDay) TableName() string { return "raffle_activity_account_day" }

type RaffleActivityAccountMonth struct {
	ID                int64     `gorm:"primaryKey" json:"id"`
	UserID            string    `json:"userId"`
	ActivityID        int64     `json:"activityId"`
	Month             string    `json:"month"`
	MonthCount        int       `json:"monthCount"`
	MonthCountSurplus int       `json:"monthCountSurplus"`
	CreateTime        time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime        time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (RaffleActivityAccountMonth) TableName() string { return "raffle_activity_account_month" }

type RaffleActivityOrder struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	UserID        string    `json:"userId"`
	Sku           int64     `json:"sku"`
	ActivityID    int64     `json:"activityId"`
	ActivityName  string    `json:"activityName"`
	StrategyID    int64     `json:"strategyId"`
	OrderID       string    `json:"orderId"`
	OrderTime     time.Time `gorm:"column:order_time" json:"orderTime"`
	TotalCount    int       `json:"totalCount"`
	DayCount      int       `json:"dayCount"`
	MonthCount    int       `json:"monthCount"`
	PayAmount     float64   `json:"payAmount"`
	State         string    `json:"state"`
	OutBusinessNo string    `json:"outBusinessNo"`
	CreateTime    time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime    time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (RaffleActivityOrder) TableName() string { return "raffle_activity_order" }

type UserRaffleOrder struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	UserID     string    `json:"userId"`
	ActivityID int64     `json:"activityId"`
	StrategyID int64     `json:"strategyId"`
	OrderID    string    `json:"orderId"`
	OrderTime  time.Time `gorm:"column:order_time" json:"orderTime"`
	OrderState string    `json:"orderState"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (UserRaffleOrder) TableName() string { return "user_raffle_order" }

type UserAwardRecord struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	UserID      string    `json:"userId"`
	ActivityID  int64     `json:"activityId"`
	StrategyID  int64     `json:"strategyId"`
	OrderID     string    `json:"orderId"`
	AwardID     int       `json:"awardId"`
	AwardTitle  string    `json:"awardTitle"`
	AwardConfig string    `json:"awardConfig"`
	AwardState  string    `json:"awardState"`
	CreateTime  time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime  time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (UserAwardRecord) TableName() string { return "user_award_record" }

type UserBehaviorRebateOrder struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	UserID        string    `json:"userId"`
	OrderID       string    `json:"orderId"`
	BehaviorType  string    `json:"behaviorType"`
	RebateDesc    string    `json:"rebateDesc"`
	RebateType    string    `json:"rebateType"`
	RebateAmount  float64   `json:"rebateAmount"`
	BizID         string    `json:"bizId"`
	OutBusinessNo string    `json:"outBusinessNo"`
	CreateTime    time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime    time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (UserBehaviorRebateOrder) TableName() string { return "user_behavior_rebate_order" }

type UserCreditAccount struct {
	ID              int64     `gorm:"primaryKey" json:"id"`
	UserID          string    `json:"userId"`
	TotalAmount     float64   `json:"totalAmount"`
	AvailableAmount float64   `json:"availableAmount"`
	AccountStatus   string    `json:"accountStatus"`
	CreateTime      time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime      time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (UserCreditAccount) TableName() string { return "user_credit_account" }

type UserCreditOrder struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	UserID        string    `json:"userId"`
	OrderID       string    `json:"orderId"`
	TradeName     string    `json:"tradeName"`
	TradeType     string    `json:"tradeType"`
	TradeAmount   float64   `json:"tradeAmount"`
	OutBusinessNo string    `json:"outBusinessNo"`
	CreateTime    time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime    time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (UserCreditOrder) TableName() string { return "user_credit_order" }

type Task struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	UserID     string    `json:"userId"`
	Topic      string    `json:"topic"`
	MessageID  string    `json:"messageId"`
	Message    string    `json:"message"`
	State      string    `json:"state"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (Task) TableName() string { return "task" }

type AdminUser struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"-"`
	Nickname   string    `json:"nickname"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (AdminUser) TableName() string { return "admin_user" }
