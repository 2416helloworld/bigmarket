package api

import (
	"strconv"
	"time"

	"bigmarket/internal/model"
	"bigmarket/internal/pkg"

	"github.com/gin-gonic/gin"
)

type resourceConf struct {
	table   string
	orderBy string
	fields  []string
}

var adminResources = map[string]resourceConf{
	"activities":       {"raffle_activity", "activity_id desc", []string{"activity_id", "activity_name", "activity_desc", "begin_date_time", "end_date_time", "strategy_id", "state"}},
	"counts":           {"raffle_activity_count", "activity_count_id asc", []string{"activity_count_id", "total_count", "day_count", "month_count"}},
	"skus":             {"raffle_activity_sku", "sku asc", []string{"sku", "activity_id", "activity_count_id", "stock_count", "stock_count_surplus", "product_amount"}},
	"strategies":       {"strategy", "strategy_id asc", []string{"strategy_id", "strategy_desc", "rule_models"}},
	"strategy-awards":  {"strategy_award", "strategy_id asc, sort asc", []string{"strategy_id", "award_id", "award_title", "award_subtitle", "award_count", "award_count_surplus", "award_rate", "rule_models", "sort"}},
	"strategy-rules":   {"strategy_rule", "strategy_id asc", []string{"strategy_id", "award_id", "rule_type", "rule_model", "rule_value", "rule_desc"}},
	"awards":           {"award", "award_id asc", []string{"award_id", "award_key", "award_config", "award_desc"}},
	"rule-trees":       {"rule_tree", "tree_id asc", []string{"tree_id", "tree_name", "tree_desc", "tree_node_rule_key"}},
	"rule-tree-nodes":  {"rule_tree_node", "tree_id asc, id asc", []string{"tree_id", "rule_key", "rule_desc", "rule_value"}},
	"rule-tree-lines":  {"rule_tree_node_line", "tree_id asc, id asc", []string{"tree_id", "rule_node_from", "rule_node_to", "rule_limit_type", "rule_limit_value"}},
	"rebates":          {"daily_behavior_rebate", "id asc", []string{"behavior_type", "rebate_desc", "rebate_type", "rebate_config", "state"}},
}

func normalizeFields(m map[string]interface{}, conf resourceConf) map[string]interface{} {
	out := map[string]interface{}{}
	for _, field := range conf.fields {
		v, ok := m[field]
		if !ok {
			continue
		}
		switch field {
		case "begin_date_time", "end_date_time":
			if s, ok := v.(string); ok && s != "" {
				out[field] = normalizeTime(s)
			}
		default:
			out[field] = v
		}
	}
	return out
}

func normalizeTime(s string) string {
	layouts := []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05", time.RFC3339, "2006-01-02"}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t.Format("2006-01-02 15:04:05")
		}
	}
	return s
}

func (a *AppHandlers) adminList(c *gin.Context) {
	resource := c.Param("resource")
	conf, ok := adminResources[resource]
	if !ok {
		render(c, nil, pkg.NewBiz("0002"))
		return
	}
	var rows []map[string]interface{}
	q := a.App.Store.DB.Table(conf.table).Order(conf.orderBy)
	if strategyID := c.Query("strategyId"); strategyID != "" {
		q = q.Where("strategy_id = ?", strategyID)
	}
	if activityID := c.Query("activityId"); activityID != "" && resource == "skus" {
		q = q.Where("activity_id = ?", activityID)
	}
	if treeID := c.Query("treeId"); treeID != "" && (resource == "rule-tree-nodes" || resource == "rule-tree-lines") {
		q = q.Where("tree_id = ?", treeID)
	}
	if err := q.Limit(500).Find(&rows).Error; err != nil {
		render(c, nil, err)
		return
	}
	render(c, rows, nil)
}

func (a *AppHandlers) adminCreate(c *gin.Context) {
	resource := c.Param("resource")
	conf, ok := adminResources[resource]
	if !ok {
		render(c, nil, pkg.NewBiz("0002"))
		return
	}
	m := params(c)
	values := normalizeFields(m, conf)
	if len(values) == 0 {
		render(c, nil, pkg.NewBiz("0002"))
		return
	}
	if err := a.App.Store.DB.Table(conf.table).Create(values).Error; err != nil {
		render(c, nil, pkg.NewBizInfo("0001", err.Error()))
		return
	}
	a.App.Store.FlushCaches()
	render(c, true, nil)
}

func (a *AppHandlers) adminUpdate(c *gin.Context) {
	resource := c.Param("resource")
	conf, ok := adminResources[resource]
	if !ok {
		render(c, nil, pkg.NewBiz("0002"))
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		render(c, nil, pkg.NewBiz("0002"))
		return
	}
	m := params(c)
	values := normalizeFields(m, conf)
	delete(values, "id")
	if len(values) == 0 {
		render(c, true, nil)
		return
	}
	if err = a.App.Store.DB.Table(conf.table).Where("id = ?", id).Updates(values).Error; err != nil {
		render(c, nil, pkg.NewBizInfo("0001", err.Error()))
		return
	}
	a.App.Store.FlushCaches()
	render(c, true, nil)
}

func (a *AppHandlers) adminDelete(c *gin.Context) {
	resource := c.Param("resource")
	conf, ok := adminResources[resource]
	if !ok {
		render(c, nil, pkg.NewBiz("0002"))
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		render(c, nil, pkg.NewBiz("0002"))
		return
	}
	if err = a.App.Store.DB.Exec("delete from "+conf.table+" where id = ?", id).Error; err != nil {
		render(c, nil, pkg.NewBizInfo("0001", err.Error()))
		return
	}
	a.App.Store.FlushCaches()
	render(c, true, nil)
}

func (a *AppHandlers) dashboard(c *gin.Context) {
	db := a.App.Store.DB
	var activityCount, strategyCount, activityOrderCount, awardCount, creditOrderCount, taskCount int64
	db.Model(&model.RaffleActivity{}).Count(&activityCount)
	db.Model(&model.Strategy{}).Count(&strategyCount)
	db.Model(&model.RaffleActivityOrder{}).Count(&activityOrderCount)
	db.Model(&model.UserAwardRecord{}).Count(&awardCount)
	db.Model(&model.UserCreditOrder{}).Count(&creditOrderCount)
	db.Model(&model.Task{}).Count(&taskCount)
	var drawToday int64
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	db.Model(&model.UserRaffleOrder{}).Where("order_time >= ?", start).Count(&drawToday)
	var awardToday int64
	db.Model(&model.UserAwardRecord{}).Where("create_time >= ?", start).Count(&awardToday)

	var recentAwards []model.UserAwardRecord
	db.Order("id desc").Limit(10).Find(&recentAwards)
	var recentOrders []model.RaffleActivityOrder
	db.Order("id desc").Limit(10).Find(&recentOrders)

	render(c, gin.H{
		"activityCount": activityCount, "strategyCount": strategyCount,
		"activityOrderCount": activityOrderCount, "awardCount": awardCount,
		"creditOrderCount": creditOrderCount, "taskCount": taskCount,
		"drawToday": drawToday, "awardToday": awardToday,
		"recentAwards": recentAwards, "recentOrders": recentOrders,
	}, nil)
}

func (a *AppHandlers) adminOrders(c *gin.Context) {
	db := a.App.Store.DB
	orderType := c.Query("type")
	userID := c.Query("userId")
	limit := 100
	var rows interface{}
	var err error
	switch orderType {
	case "raffle":
		var list []model.UserRaffleOrder
		q := db.Order("id desc").Limit(limit)
		if userID != "" {
			q = q.Where("user_id = ?", userID)
		}
		err = q.Find(&list).Error
		rows = list
	case "award":
		var list []model.UserAwardRecord
		q := db.Order("id desc").Limit(limit)
		if userID != "" {
			q = q.Where("user_id = ?", userID)
		}
		err = q.Find(&list).Error
		rows = list
	case "credit":
		var list []model.UserCreditOrder
		q := db.Order("id desc").Limit(limit)
		if userID != "" {
			q = q.Where("user_id = ?", userID)
		}
		err = q.Find(&list).Error
		rows = list
	case "rebate":
		var list []model.UserBehaviorRebateOrder
		q := db.Order("id desc").Limit(limit)
		if userID != "" {
			q = q.Where("user_id = ?", userID)
		}
		err = q.Find(&list).Error
		rows = list
	default:
		var list []model.RaffleActivityOrder
		q := db.Order("id desc").Limit(limit)
		if userID != "" {
			q = q.Where("user_id = ?", userID)
		}
		err = q.Find(&list).Error
		rows = list
	}
	render(c, rows, err)
}

func (a *AppHandlers) adminTasks(c *gin.Context) {
	list, err := a.App.Store.TaskList(100)
	render(c, list, err)
}

func (a *AppHandlers) adminUsers(c *gin.Context) {
	var activityAccounts []model.RaffleActivityAccount
	a.App.Store.DB.Order("id desc").Limit(100).Find(&activityAccounts)
	var creditAccounts []model.UserCreditAccount
	a.App.Store.DB.Order("id desc").Limit(100).Find(&creditAccounts)
	render(c, gin.H{"activityAccounts": activityAccounts, "creditAccounts": creditAccounts}, nil)
}

func (a *AppHandlers) adminCreditAdjust(c *gin.Context) {
	m := params(c)
	userID := pStr(m, "userId")
	amount := pFloat(m, "amount")
	tradeType := pStr(m, "tradeType")
	if tradeType == "" {
		tradeType = pkg.TradeForward
	}
	outBizNo := "admin_" + pkg.GenOrderID()
	err := a.App.CreditAdjust(userID, "admin", tradeType, amount, outBizNo)
	render(c, true, err)
}

func (a *AppHandlers) adminArmory(c *gin.Context) {
	m := params(c)
	err := a.App.AssembleActivity(pInt64(m, "activityId"))
	render(c, true, err)
}

func (a *AppHandlers) adminRuleTreeDetail(c *gin.Context) {
	treeID := c.Param("treeId")
	var tree model.RuleTree
	if err := a.App.Store.DB.Where("tree_id = ?", treeID).First(&tree).Error; err != nil {
		render(c, nil, err)
		return
	}
	var nodes []model.RuleTreeNode
	var lines []model.RuleTreeNodeLine
	a.App.Store.DB.Where("tree_id = ?", treeID).Find(&nodes)
	a.App.Store.DB.Where("tree_id = ?", treeID).Find(&lines)
	render(c, gin.H{"tree": tree, "nodes": nodes, "lines": lines}, nil)
}
