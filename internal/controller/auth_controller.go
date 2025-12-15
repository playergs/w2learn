package controller

import (
	"w2learn/internal/dto"
	"w2learn/internal/service"
	"w2learn/pkg/response"

	"github.com/gin-gonic/gin"
)

var _ AuthController = (*authController)(nil)

type AuthController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
}

type authController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

func (ctrl *authController) Register(c *gin.Context) {
	var req dto.RegisterRequest

	err := c.ShouldBindJSON(&req)

	if err != nil {
		response.Error(c, "Parameter binding failed: "+err.Error())
		return
	}

	err = ctrl.authService.Register(c, &req)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, "Register successfully")
}

func (ctrl *authController) Login(c *gin.Context) {
	// TODO: 当前登录接口仅支持用户名密码登录，并且会产生大量可登录的token，后续需要在产生新的Token时失效掉旧的token
	var req dto.LoginRequest

	err := c.ShouldBindJSON(&req)

	if err != nil {
		response.Error(c, "Parameter binding failed: "+err.Error())
		return
	}

	token, err := ctrl.authService.Login(c, &req)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, token)
}

func (ctrl *authController) Logout(c *gin.Context) {
	id := c.GetString("id")

	if id == "" {
		response.Error(c, "Invalid id")
	}

	err := ctrl.authService.Logout(c, id)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, "Logout successfully")
}
