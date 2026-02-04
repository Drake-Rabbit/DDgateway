package dto

// 修改密码
type ChangePasswordInput struct {
	UserID      string `json:"user_id"`
	OldPassword string `json:"oldPwd"`
	NewPassword string `json:"newPwd"`
}
