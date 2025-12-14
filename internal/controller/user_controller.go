package controller

import (
	"strconv"
	"w2learn/internal/dto"
	"w2learn/internal/service"
	"w2learn/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	CreateUser(c *gin.Context)
	GetUser(c *gin.Context)
	GetUserByUsername(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	ListUsers(c *gin.Context)
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

func (ctrl *userController) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "Parameter binding failed: "+err.Error())
		return
	}

	user, err := ctrl.userService.CreateUser(c.Request.Context(), &req)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, user)
}

func (ctrl *userController) GetUser(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	user, err := ctrl.userService.GetUserByID(c.Request.Context(), id)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, user)
}

func (ctrl *userController) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	if username == "" {
		response.Error(c, "Username can't be empty")
		return
	}

	user, err := ctrl.userService.GetUserByUsername(c.Request.Context(), username)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, user)
}

func (ctrl *userController) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, "Invalid user ID")
		return
	}

	var req dto.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "Parameter binding failed: "+err.Error())
		return
	}

	user, err := ctrl.userService.UpdateUser(c.Request.Context(), id, &req)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, user)
}

func (ctrl *userController) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		response.Error(c, "Invalid user ID")
		return
	}

	err = ctrl.userService.DeleteUser(c.Request.Context(), id)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *userController) ListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")

	page, err := strconv.Atoi(pageStr)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	pageSizeStr := c.DefaultQuery("page_size", "10")

	pageSize, err := strconv.Atoi(pageSizeStr)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	users, err := ctrl.userService.ListUsers(c, page, pageSize)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, users)
}
