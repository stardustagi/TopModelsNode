package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type StatusReport struct {
	Id               uint64    `json:"id" xorm:"'id' not null pk autoincr comment('主键ID') UNSIGNED BIGINT(20)"`
	TraceId          string    `json:"trace_id" xorm:"'trace_id' not null default '' comment('跟踪ID') index VARCHAR(64)"`
	NodeAddr         string    `json:"node_addr" xorm:"'node_addr' not null default '' comment('LLM代理地址') VARCHAR(128)"`
	Model            string    `json:"model" xorm:"'model' not null default '' comment('模型名字') index VARCHAR(64)"`
	ModelId          int       `json:"model_id" xorm:"'model_id' not null default 0 comment('模型ID（计费使用）') INT(11)"`
	ActualModel      string    `json:"actual_model" xorm:"'actual_model' not null default '' comment('实际使用的模型') VARCHAR(64)"`
	Provider         string    `json:"provider" xorm:"'provider' not null default '' comment('虚拟provider') VARCHAR(64)"`
	ActualProvider   string    `json:"actual_provider" xorm:"'actual_provider' not null default '' comment('实际服务商') VARCHAR(64)"`
	ActualProviderId string    `json:"actual_provider_id" xorm:"'actual_provider_id' not null default '' comment('实际服务商ID') VARCHAR(64)"`
	CallerKey        string    `json:"caller_key" xorm:"'caller_key' not null default '' comment('客户端key') index VARCHAR(128)"`
	Stream           int       `json:"stream" xorm:"'stream' not null default 0 comment('是否流式访问：0-否，1-是') TINYINT(1)"`
	ReportType       string    `json:"report_type" xorm:"'report_type' not null default '' comment('报告类型：text/image/video') VARCHAR(16)"`
	TokensPerSec     int       `json:"tokens_per_sec" xorm:"'tokens_per_sec' not null default 0 comment('每秒输出token') INT(11)"`
	Latency          string    `json:"latency" xorm:"'latency' not null default 0.0000 comment('请求延迟（秒）') DECIMAL(10,4)"`
	Step             string    `json:"step" xorm:"'step' not null default '' comment('调用环节：call_llm_agent/check_user_balance/select_provider/send_llm_request/send_llm_completed/llm_agent_done/user_agent_done') index VARCHAR(32)"`
	StatusCode       string    `json:"status_code" xorm:"'status_code' not null default '' comment('状态码（非空为失败）') VARCHAR(16)"`
	StatusMessage    string    `json:"status_message" xorm:"'status_message' not null default '' comment('状态消息（状态码非空时有值）') VARCHAR(512)"`
	CreatedAt        time.Time `json:"created_at" xorm:"'created_at' not null default current_timestamp() comment('请求时间') index DATETIME"`
	UserId           int64     `json:"user_id" xorm:"'user_id' BIGINT(20)"`
}

func (o *StatusReport) TableName() string {
	return "status_report"
}

func (o *StatusReport) GetSliceName(slice string, num uint32) string {
	var hash uint32
	for _, c := range slice {
		hash = hash*31 + uint32(c)
	}
	shardIndex := hash % num
	return fmt.Sprintf("status_report_%d", shardIndex)
}

func (o *StatusReport) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("status_report_%d%02d", t.Year(), t.Month())
}

func (o *StatusReport) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("status_report_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *StatusReport) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *StatusReport) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *StatusReport) PrimaryKey() interface{} {
	return o.Id
}
