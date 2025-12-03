package node_llm

import (
	"github.com/stardustagi/TopLib/libs/server"
	"github.com/stardustagi/TopModelsNode/backend"
)

func (n *NodeHttpService) initialization() {
	n.app.AddGroup("node", server.Request(), backend.NodeAccess())
	n.app.AddGroup("node/public", server.Request())
	n.app.AddPostHandler("node/public", server.NewHandler(
		"nodeLogin",
		[]string{"nodeLogin", "public"},
		n.NodeLogin))
	n.app.AddPostHandler("node", server.NewHandler(
		"keepLive",
		[]string{"llm", "node"},
		n.KeepLive))
	n.app.AddPostHandler("node", server.NewHandler(
		"ListNodeInfos",
		[]string{"llm", "node"},
		n.ListNodeInfos))
	n.app.AddPostHandler("node", server.NewHandler(
		"nodeBillingUsage",
		[]string{"llm", "node"},
		n.NodeBillingUsage))
	n.app.AddPostHandler("node", server.NewHandler(
		"NodeUnregister",
		[]string{"llm", "node"},
		n.NodeUnregister))
}
