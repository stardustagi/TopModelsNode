package node_bill

import (
	"github.com/stardustagi/TopLib/libs/server"
	"github.com/stardustagi/TopModelsNode/backend"
)

func (n *NodeHttpBillService) initialization() {
	n.app.AddGroup("node/bill", server.Request(), backend.NodeAccess())

	// 模型折扣相关路由
	n.app.AddPostHandler("node/bill", server.NewHandler(
		"createModelsDiscount",
		[]string{"bill", "discount"},
		n.CreateModelsDiscount))

	n.app.AddPostHandler("node/bill", server.NewHandler(
		"getModelsDiscount",
		[]string{"bill", "discount"},
		n.GetModelsDiscount))

	n.app.AddPostHandler("node/bill", server.NewHandler(
		"listModelsDiscount",
		[]string{"bill", "discount"},
		n.ListModelsDiscount))

	n.app.AddPostHandler("node/bill", server.NewHandler(
		"updateModelsDiscount",
		[]string{"bill", "discount"},
		n.UpdateModelsDiscount))

	n.app.AddPostHandler("node/bill", server.NewHandler(
		"deleteModelsDiscount",
		[]string{"bill", "discount"},
		n.DeleteModelsDiscount))

	n.app.AddPostHandler("node/bill", server.NewHandler(
		"getModelDiscountByModelId",
		[]string{"bill", "discount"},
		n.GetModelDiscountByModelId))
}
