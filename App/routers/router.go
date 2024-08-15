package routers

import (
	"Demo/App/internal"

	"github.com/gin-gonic/gin"
)

// 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/data/realTime", internal.ShowRealTimeData)
	r.GET("/data/history", internal.ShowHistoryData)
	r.POST("/option/settings", internal.ParmSet)
	r.POST("/option/remoteControl", internal.RemoteControl)

	return r
}
