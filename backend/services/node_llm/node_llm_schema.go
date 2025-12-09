package node_llm

// ModelConfig 模型配置根结构
type ModelConfig struct {
	NodeId      string      `json:"node_id"`
	NodeAddress string      `json:"node_address"`
	Models      []ModelInfo `json:"models"`
}

// ModelInfo 模型信息结构
type ModelInfo struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	APIVersion  string     `json:"api_version"`
	DeployName  string     `json:"deploy_name"`
	InputPrice  int        `json:"input_price"`  // 销售价
	OutputPrice int        `json:"output_price"` // 销售价
	CachePrice  int        `json:"cache_price"`  // 销售价
	Providers   []Provider `json:"providers"`
	Status      string     `json:"status"`     // 模型状态
	Address     string     `json:"address"`    // 模型地址
	ApiStyles   []string   `json:"api_styles"` // 支持的API风格列表

	// 上报信息
	Metrics    ModelMetrics `json:"metrics"`     // 模型指标
	ExpireTime int64        `json:"expire_time"` // 模型过期时间
	KeepLive   int64        `json:"keep_live"`   // 模型最后上报时间
}

type NodeKeepLiveInfo struct {
	ModelId    int64        `json:"model_id"`
	Metrics    ModelMetrics `json:"metrics"`
	ExpireTime int64        `json:"expire_time"` // 模型过期时间
	KeepLive   int64        `json:"keep_live"`   // 模型最后上报时间
}

// ModelMetrics 模型指标
type ModelMetrics struct {
	Latency     float64 `json:"latency"`      // 平均延迟(秒)
	HealthScore float64 `json:"health_score"` // 健康评分 0-100
}

// Provider 提供商信息
type Provider struct {
	ID          int64    `json:"id"`
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Endpoint    string   `json:"endpoint"`
	APIType     string   `json:"api_type"`
	APIKeys     []string `json:"api_keys"`
	ModelName   string   `json:"model_name"`
	InputPrice  int      `json:"input_price"`  // 成本价
	OutputPrice int      `json:"output_price"` // 成本价
	CachePrice  int      `json:"cache_price"`  // 成本价
	Quota       int64    `json:"quota"`        // 并发限制
}

// ProviderType 提供商类型枚举
type ProviderType string

const (
	ProviderTypePrimary ProviderType = "primary"
	ProviderTypeBackup  ProviderType = "backup"
)

// APIType API类型枚举
type APIType string

const (
	APITypeOpenAI    APIType = "openai"
	APITypeAnthropic APIType = "anthropic"
	APITypeGoogle    APIType = "google"
	APITypeLocal     APIType = "local"
)

// BalanceMode 平衡模式枚举
type BalanceMode int

const (
	BalanceModeRoundRobin BalanceMode = iota // 0: 轮询
	BalanceModeRandom                        // 1: 随机
	BalanceModeWeighted                      // 2: 加权
	BalanceModeFailover                      // 3: 故障转移
)

// UserSelectModelIds 用户选择的模型IDS
type UserSelectModelIds struct {
	NodeId   string   `json:"node_id"`
	ModelIds []string `json:"model_ids"` // 节点ID到模型ID列表的映射
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

// LLMUsageReport 大模型使用报告
type LLMUsageReport struct {
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

// NodeInfo节点信息
type NodeInfo struct {
	Id           int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string `json:"name" gorm:"type:varchar(100);not null;uniqueIndex"`
	NodeUserId   int64  `json:"node_user_id"`
	CreatedAt    int64  `json:"created_at"`
	LastUpdateAt int64  `json:"lastupdate_at"`
	Domain       string `json:"domain"`
	AccessKey    string `json:"access_key"`
	SecretKey    string `json:"secret_key"`
	CompanyId    int64  `json:"company_id"`
}
