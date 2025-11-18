package constants

import (
	topError "github.com/stardustagi/TopLib/libs/errors"
)

var (
	ErrInternalServer = topError.New("Internal server error", 500)
	ErrInvalidParams  = topError.New("无效的请求参数", 501)
	// 1000-1500 系统错误
	ErrUserRegFailed    = topError.New("User registration failed", 1001)
	ErrUserNotFound     = topError.New("User not found", 1002)
	ErrUserActiveFailed = topError.New("User activation failed", 1003)
	ErrLoginFailed      = topError.New("User login failed", 1004)
	ErrNotDataSet       = topError.New("Not data set", 1005)
)
