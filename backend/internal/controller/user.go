package controller

import (
	"gateway-service/internal/service"
	"gateway-service/pkg/response"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService *service.UserService
}

// NewUserController 创建用户控制器
func NewUserController() *UserController {
	return &UserController{
		userService: &service.UserService{},
	}
}

// ListUsers 获取用户列表
func (c *UserController) ListUsers(ctx *gin.Context) {
	tenantID := ctx.GetUint("tenant_id")

	users, err := c.userService.ListUsers(int(tenantID))
	if err != nil {
		response.InternalError(ctx, "Failed to fetch users")
		return
	}

	response.Success(ctx, users)
}

// GetUser 获取用户详情
func (c *UserController) GetUser(ctx *gin.Context) {
	tenantID := ctx.GetUint("tenant_id")
	id := ctx.Param("id")

	user, err := c.userService.GetUser(tenantID, id)
	if err != nil {
		response.NotFound(ctx, "User not found")
		return
	}

	response.Success(ctx, user)
}

// UpdateUser 更新用户
func (c *UserController) UpdateUser(ctx *gin.Context) {
	tenantID := ctx.GetUint("tenant_id")
	id := ctx.Param("id")

	var req struct {
		Email  string `json:"email"`
		Role   string `json:"role"`
		Status string `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	user, err := c.userService.UpdateUser(tenantID, id, req.Email, req.Role, req.Status)
	if err != nil {
		response.NotFound(ctx, "User not found")
		return
	}

	response.Success(ctx, user)
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx *gin.Context) {
	tenantID := ctx.GetUint("tenant_id")
	id := ctx.Param("id")

	err := c.userService.DeleteUser(tenantID, id)
	if err != nil {
		response.NotFound(ctx, "User not found")
		return
	}

	response.Success(ctx, gin.H{"message": "User deleted successfully"})
}
