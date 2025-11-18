package requests

type ModelsInfo struct {
	Name        string `json:"name" validate:"required"`
	ApiVersion  string `json:"api_version" validate:"required"`
	DeployName  string `json:"deploy_name" validate:"required"`
	InputPrice  int    `json:"input_price"`
	OutputPrice int    `json:"output_price"`
	CachePrice  int    `json:"cache_price"`
	Status      string `json:"status"`
	IsPrivate   int    `json:"is_private"`
	OwnerId     int64  `json:"owner_id"`
	Address     string `json:"address"`
}

type AddModelsInfoRequest struct {
	ModelsInfos []ModelsInfo `json:"models_infos"`
}

type ModelsProviderInfo struct {
	Type        string `json:"type" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Endpoint    string `json:"endpoint" validate:"required"`
	ApiType     string `json:"api_type"`
	ModelName   string `json:"model_name"`
	InputPrice  int    `json:"input_price"`
	OutputPrice int    `json:"output_price"`
	CachePrice  int    `json:"cache_price"`
	ApiKeys     string `json:"api_keys"`
}

type AddModelsProviderInfoRequest struct {
	ModelsProviderInfo []ModelsProviderInfo `json:"models_provider_info"`
}

type MapModelsProviderInfoToNodeRequest struct {
	NodeId      int64   `json:"node_id" validate:"required"`
	ModelId     int64   `json:"model_id" validate:"required"`
	ProviderIds []int64 `json:"provider_ids" validate:"required"`
}

type UpsetNodeInfoRequest struct {
	NodeUserId int64  `json:"node_user_id" validate:"required"`
	Code       string `json:"code"`
	Domain     string `json:"domain"`
}

type ListNodeInfoRequest struct {
	NodeUserId int64   `json:"node_user_id"`
	PageInfo   PageReq `json:"page_info"`
}
