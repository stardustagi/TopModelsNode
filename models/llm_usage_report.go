package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type LlmUsageReport struct {
	Id                        string  `json:"id" xorm:"'id' not null pk VARCHAR(255)"`
	NodeId                    int64   `json:"node_id" xorm:"'node_id' BIGINT(12)"`
	Model                     string  `json:"model" xorm:"'model' VARCHAR(255)"`
	ModelId                   int64   `json:"model_id" xorm:"'model_id' BIGINT(12)"`
	ActualModel               string  `json:"actual_model" xorm:"'actual_model' VARCHAR(255)"`
	Provider                  string  `json:"provider" xorm:"'provider' VARCHAR(255)"`
	ActualProvider            string  `json:"actual_provider" xorm:"'actual_provider' VARCHAR(255)"`
	Caller                    string  `json:"caller" xorm:"'caller' VARCHAR(255)"`
	CallerKey                 string  `json:"caller_key" xorm:"'caller_key' VARCHAR(255)"`
	ClientVersion             string  `json:"client_version" xorm:"'client_version' VARCHAR(255)"`
	TokenUsageInputTokens     int     `json:"token_usage_input_tokens" xorm:"'token_usage_input_tokens' INT(11)"`
	TokenUsageOutputTokens    int     `json:"token_usage_output_tokens" xorm:"'token_usage_output_tokens' INT(11)"`
	TokenUsageCachedTokens    int     `json:"token_usage_cached_tokens" xorm:"'token_usage_cached_tokens' INT(11)"`
	TokenUsageReasoningTokens int     `json:"token_usage_reasoning_tokens" xorm:"'token_usage_reasoning_tokens' INT(11)"`
	TokenUsageTokensPerSec    int     `json:"token_usage_tokens_per_sec" xorm:"'token_usage_tokens_per_sec' INT(11)"`
	TokenUsageLatency         float64 `json:"token_usage_latency" xorm:"'token_usage_latency' DOUBLE"`
	AgentVersion              string  `json:"agent_version" xorm:"'agent_version' VARCHAR(255)"`
	Stream                    int     `json:"stream" xorm:"'stream' not null default 0 TINYINT(1)"`
	UpdatedAt                 int64   `json:"updated_at" xorm:"'updated_at' not null BIGINT(12)"`
}

func (o *LlmUsageReport) TableName() string {
	return "llm_usage_report"
}

func (o *LlmUsageReport) GetSliceName(slice string) string {
	return fmt.Sprintf("llm_usage_report_%s", slice)
}

func (o *LlmUsageReport) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("llm_usage_report_%d%02d", t.Year(), t.Month())
}

func (o *LlmUsageReport) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("llm_usage_report_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *LlmUsageReport) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *LlmUsageReport) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *LlmUsageReport) PrimaryKey() interface{} {
	return o.Id
}
