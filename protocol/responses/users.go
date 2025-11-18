package responses

import "time"

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	UserID     int64     `json:"user_id"`
	Username   string    `json:"username"`
	Nickname   string    `json:"nickname"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Avatar     string    `json:"avatar"`
	Status     int       `json:"status"` // 1:正常 2:禁用
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

// RegisterUserResponse 用户注册响应
type RegisterUserResponse struct {
	UserInfo UserInfoResponse `json:"user_info"`
	Token    string           `json:"token"`
	ExpireAt time.Time        `json:"expire_at"`
}

// LoginUserResponse 用户登录响应
type LoginUserResponse struct {
	UserInfo UserInfoResponse `json:"user_info"`
	Token    string           `json:"token"`
	ExpireAt time.Time        `json:"expire_at"`
}

// RefreshTokenResponse 刷新Token响应
type RefreshTokenResponse struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
}

// UpdateUserInfoResponse 更新用户信息响应
type UpdateUserInfoResponse struct {
	UserInfo UserInfoResponse `json:"user_info"`
}

// GetUserInfoResponse 获取用户信息响应
type GetUserInfoResponse struct {
	UserInfo UserInfoResponse `json:"user_info"`
}

// ListUsersResponse 用户列表响应
type ListUsersResponse struct {
	Total int64              `json:"total"`
	List  []UserInfoResponse `json:"list"`
}

// LogoutUserResponse 用户登出响应
type LogoutUserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ChangePasswordResponse 修改密码响应
type ChangePasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type NodeUserAkSkInfos struct {
	NodeId       string `json:"node_id"`
	AccessKey    string `json:"access_key"`
	SecretKey    string `json:"secret_key"`
	LastupdateAt int64  `json:"lastupdate_at"`
}

type GetNodeUserAkSkResp struct {
	Keys []NodeUserAkSkInfos `json:"keys"`
}
