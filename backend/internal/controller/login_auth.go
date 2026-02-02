package controller

import (
	"fmt"
	"gateway-service/internal/config"
	"gateway-service/internal/service"
	"gateway-service/pkg/jwt"
	"gateway-service/pkg/response"
	"net/http"
	"strconv"
	"strings"
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
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
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
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	token, err := jwt.GenerateToken(strconv.Itoa(int(user.ID)), user.Username)
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
			//response.Unauthorized(ctx, err.Error())
			response.Unauthorized(ctx, "登录失败,请检查用户名和密码")
		}
		return
	}

	now := time.Now()
	user.LastLogin = &now
	_ = a.userService.UpdateLastLogin(user)

	token, err := jwt.GenerateToken(strconv.Itoa(int(user.ID)), user.Username)
	if err != nil {
		response.InternalError(ctx, "Failed to generate token")
		return
	}

	fmt.Println("auth token:", token)
	response.Success(ctx, AuthResponse{Token: token})
}

func (a *AuthController) TokenUserinfo(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		ctx.Abort()
		return
	}

	// Extract token from "Bearer <token>" format
	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if tokenString == "" || tokenString == authHeader {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
		ctx.Abort()
		return
	}

	claims, err := jwt.ValidateToken(tokenString)
	if err != nil {
		response.Unauthorized(ctx, err.Error())
		return
	}
	fmt.Println("tokenUserinfo:", claims)

	response.Success(ctx, gin.H{
		"username": claims.Username,
		"user_id":  claims.UserID,
	})
}

// 注销
func (a *AuthController) Logout(ctx *gin.Context) {
	//TODO   加入黑名单功能

	response.Success(ctx, "注销登陆成功")
}
