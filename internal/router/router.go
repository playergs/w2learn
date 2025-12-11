package router

import (
	"log"
	"w2learn/internal/config"
	"w2learn/internal/controller"
	"w2learn/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, healthCtrl controller.HealthController) *gin.Engine {
	if cfg == nil {
		log.Fatal("config is nil")
		return nil
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// 配置 Gin 中间件
	r.Use(
		middleware.CORS(),
		middleware.Logger(),
	)

	// 配置 /health 的路由
	healthGroup := r.Group("/health")

	healthGroup.GET("/", healthCtrl.HealthCheck)
	healthGroup.GET("/:flag", healthCtrl.HealthCheckWithFlag)

	return r
}
