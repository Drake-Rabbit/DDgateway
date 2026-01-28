package controller

import (
	"fmt"
	"gateway-service/internal/config"
	"gateway-service/internal/service"
	"gateway-service/pkg/jwt"
	"gateway-service/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	userService *service.UserService
	cfg         *config.Config
}

// NewAuthController 创建认证控制器
func NewAuthController(cfg *config.Config) *AuthController {
	return &AuthController{
		userService: &service.UserService{},
		cfg:         cfg,
	}
}

type RegisterRequest struct {
	Username   string `json:"username" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	TenantCode string `json:"tenant_code" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

// Register 用户注册
func (a *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	user, err := a.userService.Register(service.RegisterRequest{
		Username:   req.Username,
		Email:      req.Email,
		Password:   req.Password,
		TenantCode: req.TenantCode,
	})
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	token, err := jwt.GenerateToken(strconv.Itoa(int(user.ID)), user.Username,
		a.cfg.JWT.Secret, a.cfg.JWT.ExpireHours)
	if err != nil {
		response.InternalError(ctx, "Failed to generate token")
		return
	}

	response.Success(ctx, AuthResponse{Token: token})
}

// Login 用户登录
func (a *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	user, err := a.userService.Login(service.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		if err.Error() == "user account is inactive" {
			response.Forbidden(ctx, err.Error())
		} else {
			response.Unauthorized(ctx, err.Error())
		}
		return
	}

	now := time.Now()
	user.LastLogin = &now
	_ = a.userService.UpdateLastLogin(user)

	token, err := jwt.GenerateToken(strconv.Itoa(int(user.ID)), user.Username, a.cfg.JWT.Secret, a.cfg.JWT.ExpireHours)
	if err != nil {
		response.InternalError(ctx, "Failed to generate token")
		return
	}

	fmt.Println("login token:", token)
	response.Success(ctx, AuthResponse{Token: token})
}
