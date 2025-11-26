package nodeUsers

import (
	"github.com/stardustagi/TopLib/libs/server"
	"github.com/stardustagi/TopModelsNode/backend"
)

func (nus *NodeUsersHttpService) initialization() {
	nus.app.AddGroup("users/public", server.Request())
	nus.app.AddGroup("users", server.Request(), server.Cors(), backend.NodeUserAccess())

	// 注册公开接口：用户注册和登录
	nus.app.AddPostHandler("users/public", server.NewHandler(
		"register",
		[]string{"Users", "用户注册"},
		nus.NodeUserRegister,
	))

	nus.app.AddPostHandler("users/public", server.NewHandler(
		"login",
		[]string{"Users", "用户登录"},
		nus.LoginFromByEmail,
	))

	// 注册需要鉴权的接口
	nus.app.AddPostHandler("users", server.NewHandler(
		"logout",
		[]string{"Users", "用户登出"},
		nus.LogoutUser,
	))

	nus.app.AddGetHandler("users", server.NewHandler(
		"info",
		[]string{"Users", "获取用户信息"},
		nus.GetUserInfo,
	))

	nus.app.AddGetHandler("users/public", server.NewHandler(
		"nodeUserActive",
		[]string{"user", "active"},
		nus.NodeUserActive,
	))

	nus.app.AddPostHandler("users", server.NewHandler(
		"update",
		[]string{"Users", "更新用户信息"},
		nus.UpdateUserInfo,
	))

	nus.app.AddPostHandler("users", server.NewHandler(
		"list",
		[]string{"Users", "用户列表"},
		nus.ListUsers,
	))
	nus.app.AddPostHandler("users", server.NewHandler(
		"upsetModelsInfos",
		[]string{"node", "llm", "upsetModelsInfos"},
		nus.UpsetModelsInfos))
	nus.app.AddPostHandler("users", server.NewHandler(
		"listModelsInfos",
		[]string{"node", "llm", "upsetModelsInfos"},
		nus.ListModelsInfos))
	nus.app.AddPostHandler("users", server.NewHandler(
		"upsetModelsProvider",
		[]string{"node", "llm", "upsetModelsProvider"},
		nus.UpsetModelsProvider))
	nus.app.AddPostHandler("users", server.NewHandler(
		"listModelsProviderInfos",
		[]string{"node", "llm", "listModelsProviderInfos"},
		nus.ListModelsProviderInfos))
	nus.app.AddPostHandler("users", server.NewHandler(
		"upsetNodeInfos",
		[]string{"node", "llm", "upsetNodeInfos"},
		nus.UpsetNodeInfos))
	nus.app.AddPostHandler("users", server.NewHandler(
		"listNodeInfos",
		[]string{"node", "llm", "upsetNodeInfos"},
		nus.ListNodeInfos))
	nus.app.AddPostHandler("users", server.NewHandler(
		"MapModelsProviderInfoToNode",
		[]string{"node", "llm", "MapModelsProviderInfoToNode"},
		nus.MapModelsProviderInfoToNode))
	nus.app.AddPostHandler("users", server.NewHandler(
		"ListNodeModelsProviderInfos",
		[]string{"node", "llm", "ListNodeModelsProviderInfos"},
		nus.ListNodeModelsProviderInfos))
}
