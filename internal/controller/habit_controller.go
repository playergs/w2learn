package controller

import (
	"strconv"
	"w2learn/internal/dto"
	"w2learn/internal/service"
	"w2learn/pkg/response"

	"github.com/gin-gonic/gin"
)

var _ HabitController = (*habitController)(nil)

type HabitController interface {
	CreateHabit(c *gin.Context)
	GetHabit(c *gin.Context)
	UpdateHabit(c *gin.Context)
	DeleteHabit(c *gin.Context)
	ListHabits(c *gin.Context)
}

type habitController struct {
	habitService service.HabitService
}

func NewHabitsController(habitsService service.HabitService) HabitController {
	return &habitController{
		habitService: habitsService,
	}
}

func (ctrl *habitController) CreateHabit(c *gin.Context) {
	var req dto.CreateHabitRequest

	err := c.ShouldBindJSON(&req)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	habit, err := ctrl.habitService.CreateHabit(c.Request.Context(), &req)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, habit)
}

func (ctrl *habitController) GetHabit(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	h, err := ctrl.habitService.GetHabitByID(c.Request.Context(), id)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, h)
}

func (ctrl *habitController) UpdateHabit(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	var req dto.UpdateHabitRequest

	err = c.ShouldBindJSON(&req)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	h, err := ctrl.habitService.UpdateHabit(c.Request.Context(), id, &req)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, h)
}

func (ctrl *habitController) DeleteHabit(c *gin.Context) {
	var req dto.DeleteHabitRequest

	err := c.ShouldBindJSON(&req)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	err = ctrl.habitService.DeleteHabit(c.Request.Context(), &req)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *habitController) ListHabits(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")

	page, err := strconv.Atoi(pageStr)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	pageSizeStr := c.DefaultQuery("pageSize", "10")

	pageSize, err := strconv.Atoi(pageSizeStr)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	habits, err := ctrl.habitService.ListHabits(c.Request.Context(), page, pageSize)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, habits)
}
