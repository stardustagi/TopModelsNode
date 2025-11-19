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

type NodeRegisterReq struct {
	Mail        string `json:"mail" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
	AccessToken string `json:"access_token"`
	Once        string `json:"once"`
}

// TokenUsage 记录 token 使用情况
type TokenUsage struct {
	InputTokens     int     `json:"input_tokens"`
	OutputTokens    int     `json:"output_tokens"`
	CachedTokens    int     `json:"cached_tokens"`
	ReasoningTokens int     `json:"reasoning_tokens"`
	TokensPerSec    int     `json:"tokens_per_sec"`
	Latency         float64 `json:"latency"`
}
type UsageReport struct {
	ID               string     `json:"id"`
	NodeId           string     `json:"node_id"`
	Model            string     `json:"model"`
	ModelID          string     `json:"model_id"`     // 模型id（计费使用）
	ActualModel      string     `json:"actual_model"` // 实际使用的模型
	Provider         string     `json:"provider"`
	ActualProvider   string     `json:"actual_provider"`    // 实际服务商
	ActualProviderId string     `json:"actual_provider_id"` // 实际服务商id
	Caller           string     `json:"caller"`
	CallerKey        string     `json:"caller_key"`
	ClientVersion    string     `json:"client_version,omitempty"`
	TokenUsage       TokenUsage `json:"token_usage"`
	AgentVersion     string     `json:"agent_version,omitempty"`
	Stream           bool       `json:"stream"`
	CreatedAt        int64      `json:"created_at"`
	IsPrivate        int        `json:"is_private"` // 是否私有模型
}

type NodeReportUsageReq struct {
	NodeId string `json:"node_id" validate:"required"`
	Report []UsageReport
}

type NodeUnRegisterReq struct {
	Mail   string `json:"mail" validate:"required,email"`
	NodeId string `json:"node_id" validate:"required"`
}
