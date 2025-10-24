// 更新路由注册，支持 API 和静态文件服务
package routes

import (
	"net/http"

	"github.com/axiaoxin-com/investool/api"
	"github.com/gin-gonic/gin"
)

// Routes 注册 API URL 路由
func Routes(app *gin.Engine) {
	// 创建 API 控制器
	fundController := api.NewFundController()

	// API 路由组
	apiGroup := app.Group("/api")
	{
		// 基金相关 API
		apiGroup.GET("/fund", fundController.GetFundIndex)
		apiGroup.GET("/fund/filter", fundController.GetFundFilter)
		apiGroup.POST("/fund/check", fundController.CheckFund)
		apiGroup.GET("/fund/managers", fundController.GetFundManagers)
		apiGroup.GET("/fund/similarity", fundController.GetFundSimilarity)
		apiGroup.POST("/fund/query_by_stock", fundController.QueryByStock)

		// 健康检查
		apiGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "API服务正常运行",
			})
		})
	}
}
