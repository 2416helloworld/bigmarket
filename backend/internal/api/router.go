package api

import (
	"net/http"

	"bigmarket/internal/app"
	"bigmarket/internal/config"
	"bigmarket/internal/pkg"

	"github.com/gin-gonic/gin"
)

type AppHandlers struct {
	App *app.App
	Cfg config.Config
}

func NewRouter(application *app.App, cfg config.Config) *gin.Engine {
	h := &AppHandlers{App: application, Cfg: cfg}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, pkg.OK("ok")) })

	v1 := r.Group("/api/v1/raffle")
	{
		v1.POST("/activity/armory", h.activityArmory)
		v1.GET("/activity/armory", h.activityArmory)
		v1.POST("/activity/draw", h.draw)
		v1.GET("/activity/list", h.activityList)
		v1.POST("/activity/calendar_sign_rebate", h.calendarSignRebate)
		v1.GET("/activity/is_calendar_sign_rebate", h.isCalendarSignRebate)
		v1.POST("/activity/query_user_activity_account", h.queryUserAccount)
		v1.GET("/activity/query_user_activity_account", h.queryUserAccount)
		v1.GET("/activity/query_sku_product_list_by_activity_id", h.querySkuProducts)
		v1.POST("/activity/query_user_credit_account", h.queryCreditAccount)
		v1.GET("/activity/query_user_credit_account", h.queryCreditAccount)
		v1.POST("/activity/credit_pay_exchange_sku", h.creditPayExchangeSku)
		v1.GET("/activity/query_user_award_records", h.queryAwardRecords)
		v1.GET("/activity/query_user_raffle_orders", h.queryRaffleOrders)
		v1.GET("/activity/query_user_credit_orders", h.queryCreditOrders)
		v1.GET("/activity/query_user_rebate_orders", h.queryRebateOrders)

		v1.GET("/strategy/strategy_armory", h.strategyArmory)
		v1.POST("/strategy/strategy_armory", h.strategyArmory)
		v1.GET("/strategy/query_raffle_award_list", h.queryRaffleAwardList)
		v1.POST("/strategy/query_raffle_strategy_rule_weight", h.queryRuleWeight)
		v1.POST("/strategy/random_raffle", h.randomRaffle)
	}

	admin := r.Group("/api/v1/admin")
	admin.POST("/login", h.adminLogin)
	authed := admin.Group("", jwtMiddleware(cfg.JWTSecret))
	{
		authed.GET("/me", h.adminMe)
		authed.GET("/dashboard", h.dashboard)
		authed.GET("/users", h.adminUsers)
		authed.POST("/users/credit", h.adminCreditAdjust)
		authed.POST("/armory", h.adminArmory)
		authed.GET("/orders", h.adminOrders)
		authed.GET("/tasks", h.adminTasks)
		authed.GET("/rule-trees/:treeId/detail", h.adminRuleTreeDetail)
		authed.POST("/crud/:resource", h.adminCreate)
		authed.PUT("/crud/:resource/:id", h.adminUpdate)
		authed.DELETE("/crud/:resource/:id", h.adminDelete)
		authed.GET("/crud/:resource", h.adminList)
	}
	return r
}
