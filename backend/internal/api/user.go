package api

import (
	"encoding/json"
	"errors"
	"strconv"

	"bigmarket/internal/pkg"

	"github.com/gin-gonic/gin"
)

func params(c *gin.Context) map[string]interface{} {
	m := map[string]interface{}{}
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			m[k] = v[0]
		}
	}
	if c.Request.Method == "POST" || c.Request.Method == "PUT" {
		raw, err := c.GetRawData()
		if err == nil && len(raw) > 0 {
			var body map[string]interface{}
			if json.Unmarshal(raw, &body) == nil {
				for k, v := range body {
					m[k] = v
				}
			}
		}
	}
	return m
}

func pStr(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case string:
			return val
		case float64:
			return strconv.FormatInt(int64(val), 10)
		}
	}
	return ""
}

func pInt64(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return int64(val)
		case string:
			n, _ := strconv.ParseInt(val, 10, 64)
			return n
		}
	}
	return 0
}

func pFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return val
		case string:
			n, _ := strconv.ParseFloat(val, 64)
			return n
		}
	}
	return 0
}

func render(c *gin.Context, data interface{}, err error) {
	if err != nil {
		var biz *pkg.BizError
		if errors.As(err, &biz) {
			c.JSON(200, pkg.FailWith(biz.Code, biz.Info))
			return
		}
		c.JSON(200, pkg.FailWith("0002", err.Error()))
		return
	}
	c.JSON(200, pkg.OK(data))
}

func (a *AppHandlers) activityArmory(c *gin.Context) {
	m := params(c)
	activityID := pInt64(m, "activityId")
	if activityID == 0 {
		activityID = 100301
	}
	if err := a.App.AssembleActivity(activityID); err != nil {
		render(c, nil, err)
		return
	}
	render(c, true, nil)
}

func (a *AppHandlers) draw(c *gin.Context) {
	m := params(c)
	userID := pStr(m, "userId")
	activityID := pInt64(m, "activityId")
	result, err := a.App.Draw(userID, activityID)
	render(c, result, err)
}

func (a *AppHandlers) activityList(c *gin.Context) {
	list, err := a.App.ActivityList()
	render(c, list, err)
}

func (a *AppHandlers) calendarSignRebate(c *gin.Context) {
	m := params(c)
	ok, err := a.App.SignRebate(pStr(m, "userId"))
	render(c, ok, err)
}

func (a *AppHandlers) isCalendarSignRebate(c *gin.Context) {
	userID := c.Query("userId")
	render(c, a.App.IsSignRebate(userID), nil)
}

func (a *AppHandlers) queryUserAccount(c *gin.Context) {
	m := params(c)
	account, err := a.App.QueryUserAccount(pStr(m, "userId"), pInt64(m, "activityId"))
	render(c, account, err)
}

func (a *AppHandlers) querySkuProducts(c *gin.Context) {
	m := params(c)
	list, err := a.App.QuerySkuProducts(pInt64(m, "activityId"))
	render(c, list, err)
}

func (a *AppHandlers) queryCreditAccount(c *gin.Context) {
	m := params(c)
	render(c, a.App.QueryCreditAccount(pStr(m, "userId")), nil)
}

func (a *AppHandlers) creditPayExchangeSku(c *gin.Context) {
	m := params(c)
	order, err := a.App.CreateQuotaOrder(pStr(m, "userId"), pInt64(m, "sku"), "credit_pay", "")
	render(c, order, err)
}

func (a *AppHandlers) queryAwardRecords(c *gin.Context) {
	list, err := a.App.Store.UserAwardRecords(c.Query("userId"), 50)
	render(c, list, err)
}

func (a *AppHandlers) queryRaffleOrders(c *gin.Context) {
	list, err := a.App.Store.UserRaffleOrders(c.Query("userId"), 50)
	render(c, list, err)
}

func (a *AppHandlers) queryCreditOrders(c *gin.Context) {
	list, err := a.App.Store.CreditOrdersByUser(c.Query("userId"), 50)
	render(c, list, err)
}

func (a *AppHandlers) queryRebateOrders(c *gin.Context) {
	list, err := a.App.Store.BehaviorRebateOrders(c.Query("userId"), 50)
	render(c, list, err)
}

func (a *AppHandlers) strategyArmory(c *gin.Context) {
	m := params(c)
	_, err := a.App.AssembleStrategy(pInt64(m, "strategyId"))
	render(c, true, err)
}

func (a *AppHandlers) queryRaffleAwardList(c *gin.Context) {
	m := params(c)
	list, err := a.App.QueryAwardList(pInt64(m, "activityId"))
	render(c, list, err)
}

func (a *AppHandlers) queryRuleWeight(c *gin.Context) {
	m := params(c)
	data, err := a.App.QueryWeightRule(pInt64(m, "activityId"))
	render(c, data, err)
}

func (a *AppHandlers) randomRaffle(c *gin.Context) {
	m := params(c)
	userID := pStr(m, "userId")
	if userID == "" {
		userID = "mock-user"
	}
	award, err := a.App.PerformRaffle(userID, pInt64(m, "strategyId"))
	render(c, award, err)
}
