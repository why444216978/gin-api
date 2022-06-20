package reliablequeue

import "time"

const (
	RecordStatusUnsuccess uint8 = 0
	RecordStatusSuccess   uint8 = 1
)

// ReliableMqMessage 消息实体表
type ReliableMqMessage struct {
	Id         uint64    `gorm:"column:id" json:"id"`                   // 主键id
	CreateUser string    `gorm:"column:create_user" json:"create_user"` // 创建方标识
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"` // 创建时间
	UpdateUser string    `gorm:"column:update_user" json:"update_user"` // 更新方标识
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"` // 更新时间
	Version    uint      `gorm:"column:version" json:"version"`         // 版本号
	IsDel      uint8     `gorm:"column:is_del" json:"is_del"`           // 0-未删除，1-已删除
	Scene      string    `gorm:"column:scene" json:"scene"`             // 唯一消息scene
	SceneDesc  string    `gorm:"column:scene_desc" json:"scene_desc"`   // 描述信息
}

func (ReliableMqMessage) TableName() string {
	return "reliable_mq_message"
}

// ReliableMqMessageDistribute 消息分发关联表
type ReliableMqMessageDistribute struct {
	Id          uint64    `gorm:"column:id" json:"id"`                     // 主键id
	CreateUser  string    `gorm:"column:create_user" json:"create_user"`   // 创建方标识
	CreateTime  time.Time `gorm:"column:create_time" json:"create_time"`   // 创建时间
	UpdateUser  string    `gorm:"column:update_user" json:"update_user"`   // 更新方标识
	UpdateTime  time.Time `gorm:"column:update_time" json:"update_time"`   // 更新时间
	Version     uint      `gorm:"column:version" json:"version"`           // 版本号
	IsDel       uint8     `gorm:"column:is_del" json:"is_del"`             // 0-未删除，1-已删除
	MessageId   uint64    `gorm:"column:message_id" json:"message_id"`     // 关联message表主键id
	Scene       string    `gorm:"column:scene" json:"scene"`               // 关联message表scene
	ServiceName string    `gorm:"column:service_name" json:"service_name"` // service_name
	Uri         string    `gorm:"column:uri" json:"uri"`                   // uri
	Method      string    `gorm:"column:method" json:"method"`             // http method
}

func (ReliableMqMessageDistribute) TableName() string {
	return "reliable_mq_message_distribute"
}

// ReliableMqMessageRecord 业务侧可靠消息表
type ReliableMqMessageRecord struct {
	Id                  uint64    `gorm:"column:id" json:"id"`                                       // 主键id
	CreateUser          string    `gorm:"column:create_user" json:"create_user"`                     // 创建方标识
	CreateTime          time.Time `gorm:"column:create_time" json:"create_time"`                     // 创建时间
	UpdateUser          string    `gorm:"column:update_user" json:"update_user"`                     // 更新方标识
	UpdateTime          time.Time `gorm:"column:update_time" json:"update_time"`                     // 更新时间
	Version             uint      `gorm:"column:version" json:"version"`                             // 版本号
	IsDel               uint8     `gorm:"column:is_del" json:"is_del"`                               // 0-未删除，1-已删除
	MessageId           uint64    `gorm:"column:message_id" json:"message_id"`                       // 关联message表主键id
	MessageDistributeId uint64    `gorm:"column:message_distribute_id" json:"message_distribute_id"` // 关联reliable_mq_message_distribute表主键id
	LogId               string    `gorm:"column:log_id" json:"log_id"`                               // 消息产生的log_id
	Uuid                string    `gorm:"column:uuid" json:"uuid"`                                   // 消息唯一id
	ServiceName         string    `gorm:"column:service_name" json:"service_name"`                   // service_name
	Uri                 string    `gorm:"column:uri" json:"uri"`                                     // uri
	Method              string    `gorm:"column:method" json:"method"`                               // http method
	Body                string    `gorm:"column:body" json:"body"`                                   // 请求body
	Delay               int64     `gorm:"column:delay" json:"delay"`                                 // 重试间隔
	RetryTime           time.Time `gorm:"column:retry_time" json:"retry_time"`                       // 最后一次重试时间
	NextTime            time.Time `gorm:"column:next_time" json:"next_time"`                         // 下次重试时间
	IsSuccess           uint8     `gorm:"column:is_success" json:"is_success"`                       // 是否消费成功
}

func (ReliableMqMessageRecord) TableName() string {
	return "reliable_mq_message_record"
}
