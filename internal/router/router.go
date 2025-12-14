package router

import (
	"log"
	"w2learn/internal/config"
	"w2learn/internal/controller"
	"w2learn/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	cfg *config.Config,
	healthCtrl controller.HealthController,
	userCtrl controller.UserController,
	habitCtrl controller.HabitController,
) *gin.Engine {
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

	// 配置 /user 路由
	userGroup := r.Group("/user")
	userGroup.GET("", userCtrl.ListUsers)
	userGroup.POST("", userCtrl.CreateUser)
	userGroup.GET("/i/:id", userCtrl.GetUser)
	userGroup.GET("/u/:username", userCtrl.GetUserByUsername)
	userGroup.PUT("/:id", userCtrl.UpdateUser)
	userGroup.DELETE("/:id", userCtrl.DeleteUser)

	// 配置 /habit 路由
	habitGroup := r.Group("/habit")
	habitGroup.GET("", habitCtrl.ListHabits)
	habitGroup.POST("", habitCtrl.CreateHabit)
	habitGroup.GET("/:id", habitCtrl.GetHabit)
	habitGroup.PUT("/:id", habitCtrl.UpdateHabit)
	habitGroup.DELETE("", habitCtrl.DeleteHabit)

	return r
}
