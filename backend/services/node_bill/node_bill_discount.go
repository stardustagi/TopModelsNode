package node_bill

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stardustagi/TopLib/protocol"
	"github.com/stardustagi/TopModelsNode/constants"
	"github.com/stardustagi/TopModelsNode/models"
	"github.com/stardustagi/TopModelsNode/protocol/requests"
	"github.com/stardustagi/TopModelsNode/protocol/responses"
	"go.uber.org/zap"
)

// CreateModelsDiscount 创建模型折扣
// @Summary 创建模型折扣
// @Description 创建模型折扣
// @Tags bill
// @Accept json
// @Produce json
// @Param request body requests.CreateModelsDiscountReq true "请求参数"
// @Success 200 {object} responses.CreateModelsDiscountResp
// @Router /node/bill/createModelsDiscount [post]
func (n *NodeHttpBillService) CreateModelsDiscount(ctx echo.Context,
	req requests.CreateModelsDiscountReq, resp responses.CreateModelsDiscountResp) error {
	n.logger.Info("CreateModelsDiscount called",
		zap.Int64("modelId", req.ModelId),
		zap.Int64("modelProviderId", req.ModelProviderId))

	if req.ModelId <= 0 {
		n.logger.Warn("创建模型折扣时，模型ID无效")
		return protocol.Response(ctx,
			constants.ErrInvalidParams.AppendErrors(fmt.Errorf("模型ID不能为空")),
			nil)
	}

	session := n.dao.NewSession()
	defer session.Close()

	now := time.Now().Unix()
	discount := &models.ModelsDiscount{
		ModelId:         req.ModelId,
		ModelProviderId: req.ModelProviderId,
		Type:            req.Type,
		Value:           req.Value,
		CreatedAt:       now,
		LastUpdate:      now,
	}

	_, err := session.Native().Insert(discount)
	if err != nil {
		n.logger.Error("创建模型折扣失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	resp.Id = discount.Id
	resp.ModelId = discount.ModelId
	resp.ModelProviderId = discount.ModelProviderId
	resp.Type = discount.Type
	resp.Value = discount.Value
	resp.CreatedAt = discount.CreatedAt

	return protocol.Response(ctx, nil, resp)
}

// GetModelsDiscount 获取模型折扣详情
// @Summary 获取模型折扣详情
// @Description 获取模型折扣详情
// @Tags bill
// @Accept json
// @Produce json
// @Param request body requests.GetModelsDiscountReq true "请求参数"
// @Success 200 {object} responses.GetModelsDiscountResp
// @Router /node/bill/getModelsDiscount [post]
func (n *NodeHttpBillService) GetModelsDiscount(ctx echo.Context,
	req requests.GetModelsDiscountReq, resp responses.GetModelsDiscountResp) error {
	n.logger.Info("GetModelsDiscount called", zap.Int64("id", req.Id))

	if req.Id <= 0 {
		n.logger.Warn("获取模型折扣时，ID无效")
		return protocol.Response(ctx,
			constants.ErrInvalidParams.AppendErrors(fmt.Errorf("折扣ID不能为空")),
			nil)
	}

	session := n.dao.NewSession()
	defer session.Close()

	discount := &models.ModelsDiscount{Id: req.Id}
	ok, err := session.FindOne(discount)
	if err != nil {
		n.logger.Error("查询模型折扣失败", zap.Error(err), zap.Int64("id", req.Id))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	if !ok {
		n.logger.Warn("模型折扣不存在", zap.Int64("id", req.Id))
		return protocol.Response(ctx, constants.ErrNotDataSet, nil)
	}

	resp.Id = discount.Id
	resp.ModelId = discount.ModelId
	resp.ModelProviderId = discount.ModelProviderId
	resp.Type = discount.Type
	resp.Value = discount.Value
	resp.CreatedAt = discount.CreatedAt
	resp.LastUpdate = discount.LastUpdate

	return protocol.Response(ctx, nil, resp)
}

// ListModelsDiscount 获取模型折扣列表
// @Summary 获取模型折扣列表
// @Description 获取模型折扣列表
// @Tags bill
// @Accept json
// @Produce json
// @Param request body requests.ListModelsDiscountReq true "请求参数"
// @Success 200 {object} responses.ListModelsDiscountResp
// @Router /node/bill/listModelsDiscount [post]
func (n *NodeHttpBillService) ListModelsDiscount(ctx echo.Context,
	req requests.ListModelsDiscountReq, resp responses.ListModelsDiscountResp) error {
	n.logger.Info("ListModelsDiscount called",
		zap.Int64("modelId", req.ModelId),
		zap.Int64("modelProviderId", req.ModelProviderId))

	session := n.dao.NewSession()
	defer session.Close()

	// 默认分页
	if req.PageInfo.Limit <= 0 {
		req.PageInfo.Limit = 20
	}
	if req.PageInfo.Sort == "" {
		req.PageInfo.Sort = "id desc"
	}

	query := session.Native().NewSession()
	defer query.Close()

	// 条件过滤
	if req.ModelId > 0 {
		query = query.Where("model_id = ?", req.ModelId)
	}
	if req.ModelProviderId > 0 {
		query = query.Where("model_provider_id = ?", req.ModelProviderId)
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	var discounts []models.ModelsDiscount
	err := query.
		OrderBy(req.PageInfo.Sort).
		Limit(req.PageInfo.Limit, req.PageInfo.Skip).
		Find(&discounts)
	if err != nil {
		n.logger.Error("查询模型折扣列表失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	// 获取总数
	total, err := session.Native().Count(&models.ModelsDiscount{})
	if err != nil {
		n.logger.Error("统计模型折扣数量失败", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	resp.Discounts = discounts
	resp.Total = int(total)

	return protocol.Response(ctx, nil, resp)
}

// UpdateModelsDiscount 更新模型折扣
// @Summary 更新模型折扣
// @Description 更新模型折扣
// @Tags bill
// @Accept json
// @Produce json
// @Param request body requests.UpdateModelsDiscountReq true "请求参数"
// @Success 200 {object} responses.UpdateModelsDiscountResp
// @Router /node/bill/updateModelsDiscount [post]
func (n *NodeHttpBillService) UpdateModelsDiscount(ctx echo.Context,
	req requests.UpdateModelsDiscountReq, resp responses.UpdateModelsDiscountResp) error {
	n.logger.Info("UpdateModelsDiscount called", zap.Int64("id", req.Id))

	if req.Id <= 0 {
		n.logger.Warn("更新模型折扣时，ID无效")
		return protocol.Response(ctx,
			constants.ErrInvalidParams.AppendErrors(fmt.Errorf("折扣ID不能为空")),
			nil)
	}

	session := n.dao.NewSession()
	defer session.Close()

	// 先查询是否存在
	discount := &models.ModelsDiscount{Id: req.Id}
	ok, err := session.FindOne(discount)
	if err != nil {
		n.logger.Error("查询模型折扣失败", zap.Error(err), zap.Int64("id", req.Id))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	if !ok {
		n.logger.Warn("模型折扣不存在", zap.Int64("id", req.Id))
		return protocol.Response(ctx, constants.ErrNotDataSet, nil)
	}

	// 更新字段
	if req.ModelId > 0 {
		discount.ModelId = req.ModelId
	}
	if req.ModelProviderId > 0 {
		discount.ModelProviderId = req.ModelProviderId
	}
	if req.Type != "" {
		discount.Type = req.Type
	}
	if req.Value > 0 {
		discount.Value = req.Value
	}
	discount.LastUpdate = time.Now().Unix()

	_, err = session.Native().ID(req.Id).Update(discount)
	if err != nil {
		n.logger.Error("更新模型折扣失败", zap.Error(err), zap.Int64("id", req.Id))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	resp.Id = discount.Id
	resp.ModelId = discount.ModelId
	resp.ModelProviderId = discount.ModelProviderId
	resp.Type = discount.Type
	resp.Value = discount.Value
	resp.LastUpdate = discount.LastUpdate

	return protocol.Response(ctx, nil, resp)
}

// DeleteModelsDiscount 删除模型折扣
// @Summary 删除模型折扣
// @Description 删除模型折扣
// @Tags bill
// @Accept json
// @Produce json
// @Param request body requests.DeleteModelsDiscountReq true "请求参数"
// @Success 200 {object} responses.DefaultResponse
// @Router /node/bill/deleteModelsDiscount [post]
func (n *NodeHttpBillService) DeleteModelsDiscount(ctx echo.Context,
	req requests.DeleteModelsDiscountReq, resp responses.DefaultResponse) error {
	n.logger.Info("DeleteModelsDiscount called", zap.Int64("id", req.Id))

	if req.Id <= 0 {
		n.logger.Warn("删除模型折扣时，ID无效")
		return protocol.Response(ctx,
			constants.ErrInvalidParams.AppendErrors(fmt.Errorf("折扣ID不能为空")),
			nil)
	}

	session := n.dao.NewSession()
	defer session.Close()

	// 先查询是否存在
	discount := &models.ModelsDiscount{Id: req.Id}
	ok, err := session.FindOne(discount)
	if err != nil {
		n.logger.Error("查询模型折扣失败", zap.Error(err), zap.Int64("id", req.Id))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	if !ok {
		n.logger.Warn("模型折扣不存在", zap.Int64("id", req.Id))
		return protocol.Response(ctx, constants.ErrNotDataSet, nil)
	}

	_, err = session.Native().ID(req.Id).Delete(&models.ModelsDiscount{})
	if err != nil {
		n.logger.Error("删除模型折扣失败", zap.Error(err), zap.Int64("id", req.Id))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	return protocol.Response(ctx, nil, "删除成功")
}

// GetModelDiscountByModelId 根据模型ID获取折扣
// @Summary 根据模型ID获取折扣
// @Description 根据模型ID获取折扣
// @Tags bill
// @Accept json
// @Produce json
// @Param request body requests.GetModelDiscountByModelIdReq true "请求参数"
// @Success 200 {object} responses.GetModelsDiscountResp
// @Router /node/bill/getModelDiscountByModelId [post]
func (n *NodeHttpBillService) GetModelDiscountByModelId(ctx echo.Context,
	req requests.GetModelDiscountByModelIdReq, resp responses.GetModelsDiscountResp) error {
	n.logger.Info("GetModelDiscountByModelId called", zap.Int64("modelId", req.ModelId))

	if req.ModelId <= 0 {
		n.logger.Warn("获取模型折扣时，模型ID无效")
		return protocol.Response(ctx,
			constants.ErrInvalidParams.AppendErrors(fmt.Errorf("模型ID不能为空")),
			nil)
	}

	session := n.dao.NewSession()
	defer session.Close()

	discount := &models.ModelsDiscount{ModelId: req.ModelId}
	if req.ModelProviderId > 0 {
		discount.ModelProviderId = req.ModelProviderId
	}

	ok, err := session.FindOne(discount)
	if err != nil {
		n.logger.Error("查询模型折扣失败", zap.Error(err), zap.Int64("modelId", req.ModelId))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	if !ok {
		n.logger.Warn("模型折扣不存在", zap.Int64("modelId", req.ModelId))
		return protocol.Response(ctx, constants.ErrNotDataSet, nil)
	}

	resp.Id = discount.Id
	resp.ModelId = discount.ModelId
	resp.ModelProviderId = discount.ModelProviderId
	resp.Type = discount.Type
	resp.Value = discount.Value
	resp.CreatedAt = discount.CreatedAt
	resp.LastUpdate = discount.LastUpdate

	return protocol.Response(ctx, nil, resp)
}
