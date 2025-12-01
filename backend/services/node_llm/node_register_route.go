package node_llm

import (
	"github.com/stardustagi/TopLib/libs/server"
	"github.com/stardustagi/TopModelsNode/backend"
)

func (n *NodeHttpService) initialization() {
	n.app.AddGroup("node/llm", server.Request(), backend.NodeAccess())
	n.app.AddGroup("node/public", server.Request())
	n.app.AddPostHandler("node/public", server.NewHandler(
		"nodeLogin",
		[]string{"nodeLogin", "public"},
		n.NodeLogin))
	n.app.AddPostHandler("node/llm", server.NewHandler(
		"keepLive",
		[]string{"llm", "node"},
		n.KeepLive))
	n.app.AddPostHandler("node/llm", server.NewHandler(
		"ListNodeInfos",
		[]string{"llm", "node"},
		n.ListNodeInfos))
	n.app.AddPostHandler("node/llm", server.NewHandler(
		"NodeBillingUsage",
		[]string{"llm", "node"},
		n.NodeBillingUsage))
	n.app.AddPostHandler("node/llm", server.NewHandler(
		"NodeUnregister",
		[]string{"llm", "node"},
		n.NodeUnregister))
}
