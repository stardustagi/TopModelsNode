package requests

import "time"

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

type StatusReportReq struct {
	TraceId          string    `json:"trace_id"`           //跟踪id
	NodeAddr         string    `json:"node_addr"`          //llm代理地址
	Model            string    `json:"model"`              //模型名字
	ModelID          int       `json:"model_id"`           // 模型id（计费使用）
	ActualModel      string    `json:"actual_model"`       // 实际使用的模型
	Provider         string    `json:"provider"`           //虚拟provider
	ActualProvider   string    `json:"actual_provider"`    // 实际服务商
	ActualProviderId string    `json:"actual_provider_id"` // 实际服务商id
	UserId           int64     `json:"user_id"`            //用户id
	CallerKey        string    `json:"caller_key"`         //客户端key
	Stream           bool      `json:"stream"`             //是否流式访问
	ReportType       string    `json:"report_type"`        //text/image/video
	TokensPerSec     int       `json:"tokens_per_sec"`     //每秒输出token
	Latency          float64   `json:"latency"`            //请求延迟
	Step             string    `json:"step"`               //调用环节
	StatusCode       string    `json:"status_code"`        // 状态码（非空为失败）
	StatusMessage    string    `json:"status_message"`     //状态消息 （状态码非空的时候有值）
	CreatedAt        time.Time `json:"created_at"`         //请求时间
}
