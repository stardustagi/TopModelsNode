package responses

import "github.com/stardustagi/TopModelsNode/models"

// CreateModelsDiscountResp 创建模型折扣响应
type CreateModelsDiscountResp struct {
	Id              int64  `json:"id"`
	ModelId         int64  `json:"model_id"`
	ModelProviderId int64  `json:"model_provider_id"`
	Type            string `json:"type"`
	Value           int64  `json:"value"`
	CreatedAt       int64  `json:"created_at"`
}

// GetModelsDiscountResp 获取模型折扣响应
type GetModelsDiscountResp struct {
	Id              int64  `json:"id"`
	ModelId         int64  `json:"model_id"`
	ModelProviderId int64  `json:"model_provider_id"`
	Type            string `json:"type"`
	Value           int64  `json:"value"`
	CreatedAt       int64  `json:"created_at"`
	LastUpdate      int64  `json:"last_update"`
}

// ListModelsDiscountResp 获取模型折扣列表响应
type ListModelsDiscountResp struct {
	Discounts []models.ModelsDiscount `json:"discounts"`
	Total     int                     `json:"total"`
}

// UpdateModelsDiscountResp 更新模型折扣响应
type UpdateModelsDiscountResp struct {
	Id              int64  `json:"id"`
	ModelId         int64  `json:"model_id"`
	ModelProviderId int64  `json:"model_provider_id"`
	Type            string `json:"type"`
	Value           int64  `json:"value"`
	LastUpdate      int64  `json:"last_update"`
}
