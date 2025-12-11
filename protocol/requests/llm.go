package requests

type ModelsInfo struct {
	Id          int64  `json:"id"`
	Name        string `json:"name" validate:"required"`
	ApiVersion  string `json:"api_version" validate:"required"`
	DeployName  string `json:"deploy_name" validate:"required"`
	InputPrice  int    `json:"input_price" validate:"required"`
	OutputPrice int    `json:"output_price" validate:"required"`
	CachePrice  int    `json:"cache_price" validate:"required"`
	Status      string `json:"status"`
	IsPrivate   int    `json:"is_private"`
	OwnerId     int64  `json:"owner_id"`
	Address     string `json:"address"`
	ApiStyles   string `json:"api_styles" validate:"required"`
}

type UpsertModelsInfoRequest struct {
	ModelsInfos []ModelsInfo `json:"models_infos"`
}

type UpsetNodeInfoRequest struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Domain    string `json:"domain"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	CompanyId int64  `json:"company_id"`
}

type ModelsProviderInfo struct {
	Id          int64  `json:"id"`
	Type        string `json:"type" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Endpoint    string `json:"endpoint" validate:"required"`
	ApiType     string `json:"api_type" validate:"required"`
	ModelName   string `json:"model_name" validate:"required"`
	InputPrice  int    `json:"input_price"`
	OutputPrice int    `json:"output_price"`
	CachePrice  int    `json:"cache_price"`
	ApiKeys     string `json:"api_keys" validate:"required"`
}

type UpsertModelsProviderInfoRequest struct {
	ModelsProviderInfo []ModelsProviderInfo `json:"models_provider_info"`
}

type MapModelsProviderInfoToNodeRequest struct {
	NodeId      int64   `json:"node_id" validate:"required"`
	ModelId     int64   `json:"model_id" validate:"required"`
	ProviderIds []int64 `json:"provider_ids" validate:"required"`
}

type UnMapModelsProviderInfoToNodeRequest struct {
	NodeId int64   `json:"node_id" validate:"required"`
	MapIds []int64 `json:"map_ids" validate:"required"`
}
type ListNodeInfoRequest struct {
	PageInfo PageReq `json:"page_info"`
}

type NodeLoginReq struct {
	Mail        string `json:"mail" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
	AccessToken string `json:"access_token" validate:"required,min=6"`
	Name        string `json:"name"`
	Once        string `json:"once" validate:"required"`
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

type ReportType string

const (
	TextReportType  ReportType = "text"
	ImageReportType ReportType = "image"
	VideoReportType ReportType = "video"
)

type UsageReport struct {
	ID               string     `json:"id"`
	NodeId           int64      `json:"node_id"`
	Model            string     `json:"model"`
	ModelID          int64      `json:"model_id"`     // 模型id（计费使用）
	ActualModel      string     `json:"actual_model"` // 实际使用的模型
	Provider         string     `json:"provider"`
	ActualProvider   string     `json:"actual_provider"`    // 实际服务商
	ActualProviderId string     `json:"actual_provider_id"` // 实际服务商id
	Caller           string     `json:"caller"`
	CallerKey        string     `json:"caller_key"`
	ClientVersion    string     `json:"client_version,omitempty"`
	AgentVersion     string     `json:"agent_version,omitempty"`
	Stream           bool       `json:"stream"`
	ReportType       ReportType `json:"report_type"`
	TokenUsage       any        `json:"token_usage"`

	CreatedAt int64 `json:"created_at"`
	IsPrivate int   `json:"is_private"` // 是否私有模型
}

type NodeReportUsageReq struct {
	NodeId int64 `json:"node_id" validate:"required,gt=0"`
	Report []UsageReport
}

type NodeUnRegisterReq struct {
	Mail string `json:"mail" validate:"required,email"`
	Name string `json:"name" validate:"required"`
}

type ModelMetrics struct {
	Latency     float64 `json:"latency"`      // 平均延迟(秒)
	HealthScore float64 `json:"health_score"` // 健康评分 0-100
}

type NodeKeepLiveModelInfo struct {
	ModelId    int64        `json:"model_id"`
	Metrics    ModelMetrics `json:"metrics"`
	ExpireTime int64        `json:"expire_time"` // 模型过期时间
	KeepLive   int64        `json:"keep_live"`   // 模型最后上报时间
}

type NodeKeepLiveReq struct {
	NodeId   int64                   `json:"node_id" validate:"required,gt=0"`
	NodeName string                  `json:"node_name" validate:"required"`
	Info     []NodeKeepLiveModelInfo `json:"info"`
}
