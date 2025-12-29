package models

type UserConsumeDetailText struct {
	ID           int64 `xorm:"pk autoincr comment('主键，自增')" json:"id"`                    // 主键，自增
	ConsumdId    int64 `xorm:"consume_id comment('消费记录id')" json:"consume_id"`            // 消费记录id
	InputTokens  int64 `xorm:"bigint default 0 comment('输入token数')" json:"input_tokens"`  // 输入token数
	OutputTokens int64 `xorm:"bigint default 0 comment('输出token数')" json:"output_tokens"` // 输出token数
	CacheTokens  int64 `xorm:"bigint default 0 comment('缓存token数')" json:"cache_tokens"`  // 缓存token数
	InputPrice   int   `xorm:"int default 0 comment('输入token价格')" json:"input_price"`     // 输入token价格
	OutputPrice  int   `xorm:"int default 0 comment('输出token价格')" json:"output_price"`    // 输出token价格
	CachePrice   int   `xorm:"int default 0 comment('缓存token价格')" json:"cache_price"`     // 缓存token价格
	CreatedAt    int64 `xorm:"created_at comment('创建时间')" json:"created"`                 // 创建时间
}

func (UserConsumeDetailText) TableName() string {
	return "user_consume_detail_text"
}

// UserConsumeDetailImage 图片消费明细，可能是多张
type UserConsumeDetailImage struct {
	ID        int64  `xorm:"pk autoincr comment('主键，自增')" json:"id"`
	ConsumeId int64  `xorm:"consume_id comment('消费记录id')" json:"consume_id"`
	Quality   string `xorm:"varchar(64) comment('Quality')" json:"quality"`
	Size      string `xorm:"varchar(64) comment('Size')" json:"size"`
	CreatedAt int64  `xorm:"created_at comment('创建时间')" json:"created"`
}

func (UserConsumeDetailImage) TableName() string {
	return "user_consume_detail_image"
}

type UserConsumeDetailVideo struct {
	ID        int64   `xorm:"pk autoincr comment('主键，自增')" json:"id"`         // 主键，自增
	ConsumdId int64   `xorm:"consume_id comment('消费记录id')" json:"consume_id"` // 消费记录id
	Seconds   float64 `xorm:"double comment('Seconds')" json:"Seconds"`       // Seconds
	Size      string  `xorm:"varchar(64) comment('Size')" json:"size"`        // Size
	CreatedAt int64   `xorm:"created_at comment('创建时间')" json:"created"`      // 创建时间
}

func (UserConsumeDetailVideo) TableName() string {
	return "user_consume_detail_video"
}
