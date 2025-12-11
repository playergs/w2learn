package controller

import (
	"strconv"
	"w2learn/internal/model"
	"w2learn/internal/service"
	"w2learn/pkg/response"

	"github.com/gin-gonic/gin"
)

var _ HealthController = (*healthController)(nil)

type HealthController interface {
	HealthCheck(ctx *gin.Context)
	HealthCheckWithFlag(ctx *gin.Context)
}

type healthController struct {
	healthService service.HealthService
}

func NewHealthController(healthService service.HealthService) HealthController {
	return &healthController{
		healthService: healthService,
	}
}

func (ctrl *healthController) HealthCheckWithFlag(c *gin.Context) {
	flag := c.Param("flag")

	flagInt, err := strconv.Atoi(flag)

	if err != nil {
		response.Error(c, "convert flag fail")
	}

	health, err := ctrl.healthService.GetHealth(c, flagInt)

	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, health)
}

func (ctrl *healthController) HealthCheck(c *gin.Context) {
	response.Success(c, model.GetDefaultHealthModel())
}
