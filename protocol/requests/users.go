package requests

// RegisterUserRequest 用户注册请求
type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,min=11,max=15"`
	Nickname string `json:"nickname"`
}

// LoginUserRequest 用户登录请求
type LoginUserRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UpdateUserInfoRequest 更新用户信息请求
type UpdateUserInfoRequest struct {
	UserID   int64  `json:"user_id" validate:"required"`
	Nickname string `json:"nickname" validate:"omitempty,max=50"`
	Email    string `json:"email" validate:"omitempty,email"`
	Phone    string `json:"phone" validate:"omitempty,min=11,max=15"`
	Avatar   string `json:"avatar" validate:"omitempty,url"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	UserID      int64  `json:"user_id" validate:"required"`
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=100"`
}

// RefreshTokenRequest 刷新Token请求
type RefreshTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

// GetUserInfoRequest 获取用户信息请求
type GetUserInfoRequest struct {
	UserID int64 `json:"user_id" validate:"required"`
}

// ListUsersRequest 用户列表请求
type ListUsersRequest struct {
	Page     PageReq `json:"page"`
	Username string  `json:"username" validate:"omitempty"`
	Email    string  `json:"email" validate:"omitempty"`
	Status   int     `json:"status" validate:"omitempty,oneof=0 1 2"` // 0:全部 1:正常 2:禁用
}

// LogoutUserRequest 用户登出请求
type LogoutUserRequest struct {
	UserID int64 `json:"user_id" validate:"required"`
}

type GetNodeUserAkSkReq struct {
	UserId int64 `json:"user_id" validate:"required"`
}

type DeleteNodeAkSkRequest struct {
	UserId int64  `json:"user_id" validate:"required"`
	NodeId string `json:"node_id" validate:"required"`
	Ak     string `json:"ak" validate:"required"`
	Sk     string `json:"sk" validate:"required"`
}

type UserBalanceReq struct {
	UserID     int64  `json:"user_id" validate:"required"`
	WalletType string `json:"wallet_type" validate:"required"`
}
