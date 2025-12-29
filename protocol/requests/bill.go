package requests

// CreateModelsDiscountReq 创建模型折扣请求
type CreateModelsDiscountReq struct {
	ModelId         int64  `json:"model_id" validate:"required"`
	ModelProviderId int64  `json:"model_provider_id"`
	Type            string `json:"type" validate:"required"` // percent, money, time
	Value           int64  `json:"value" validate:"required"`
}

// GetModelsDiscountReq 获取模型折扣请求
type GetModelsDiscountReq struct {
	Id int64 `json:"id" validate:"required"`
}

// ListModelsDiscountReq 获取模型折扣列表请求
type ListModelsDiscountReq struct {
	ModelId         int64   `json:"model_id"`
	ModelProviderId int64   `json:"model_provider_id"`
	Type            string  `json:"type"`
	PageInfo        PageReq `json:"page_info"`
}

// UpdateModelsDiscountReq 更新模型折扣请求
type UpdateModelsDiscountReq struct {
	Id              int64  `json:"id" validate:"required"`
	ModelId         int64  `json:"model_id"`
	ModelProviderId int64  `json:"model_provider_id"`
	Type            string `json:"type"`
	Value           int64  `json:"value"`
}

// DeleteModelsDiscountReq 删除模型折扣请求
type DeleteModelsDiscountReq struct {
	Id int64 `json:"id" validate:"required"`
}

// GetModelDiscountByModelIdReq 根据模型ID获取折扣请求
type GetModelDiscountByModelIdReq struct {
	ModelId         int64 `json:"model_id" validate:"required"`
	ModelProviderId int64 `json:"model_provider_id"`
}
