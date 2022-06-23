package reliablequeue

import (
	"context"
	"math"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"github.com/why444216978/go-util/assert"
	"github.com/why444216978/go-util/conversion"
	"github.com/why444216978/go-util/orm"
	"github.com/why444216978/go-util/snowflake"

	"github.com/why444216978/gin-api/library/queue"
)

var (
	defaultFirstDelaySecond         time.Duration = 10
	defaultRetryDelaySecondMultiple int64         = 2
)

type ReliableQueue struct {
	opt   *Option
	queue queue.Queue
}

type Option struct {
	FirstDelaySecond         time.Duration // 发布时退避秒数
	RetryDelaySecondMultiple int64         // 重试退避倍数
}

type ReliableQueueOption func(*Option)

func defaultOption() *Option {
	return &Option{
		FirstDelaySecond:         defaultFirstDelaySecond,
		RetryDelaySecondMultiple: defaultRetryDelaySecondMultiple,
	}
}

func WithFirstDelaySecond(t time.Duration) ReliableQueueOption {
	return func(o *Option) {
		o.FirstDelaySecond = t
	}
}

func WithRetryDelaySecondMultiple(i int64) ReliableQueueOption {
	return func(o *Option) {
		o.RetryDelaySecondMultiple = i
	}
}

func NewReliableQueue(q queue.Queue, opts ...ReliableQueueOption) (*ReliableQueue, error) {
	if assert.IsNil(q) {
		return nil, errors.New("Queue is nil")
	}

	opt := defaultOption()
	for _, o := range opts {
		o(opt)
	}

	rq := &ReliableQueue{
		queue: q,
		opt:   opt,
	}

	return rq, nil
}

// PublishParams 发布消息方法参数
type PublishParams struct {
	LogID string
	Scene string
	Data  map[string]interface{}
}

// Publish 发布消息，注意此方法在本地事务最后一步调用，会自动提交事务
func (rq *ReliableQueue) Publish(ctx context.Context, tx *gorm.DB, msg PublishParams) (err error) {
	if tx == nil {
		return errors.New("tx is nil")
	}

	// 根据scene获取对应分发列表
	distributeList, err := rq.getDistributeList(ctx, tx, msg)
	if err != nil {
		return
	}

	// 生成消息
	record := ReliableMqMessageRecord{}
	records := []ReliableMqMessageRecord{}
	for _, v := range distributeList {
		record, err = rq.generateMessage(ctx, tx, v, msg)
		if err != nil {
			tx.Rollback()
			return
		}
		records = append(records, record)
	}

	// 本地事务提交
	if err = tx.Commit().Error; err != nil {
		err = errors.Wrap(err, "reliablequeue Publish tx.Commit err")
		return
	}

	// 并行分发消息
	return rq.publish(ctx, records)
}

// generateMessage 生成消息表记录
func (rq *ReliableQueue) generateMessage(ctx context.Context, tx *gorm.DB, messageDistribute ReliableMqMessageDistribute, msg PublishParams) (record ReliableMqMessageRecord, err error) {
	uuid := snowflake.Generate().String()

	msg.Data["uuid"] = uuid
	b, err := conversion.JsonEncode(msg.Data)
	if err != nil {
		err = errors.Wrap(err, "reliablequeue generateMessage conversion.JsonEncode err")
		return
	}

	record = ReliableMqMessageRecord{
		MessageId:           messageDistribute.MessageId,
		MessageDistributeId: messageDistribute.Id,
		LogId:               msg.LogID,
		Uuid:                uuid,
		ServiceName:         messageDistribute.ServiceName,
		Uri:                 messageDistribute.Uri,
		Method:              messageDistribute.Method,
		Body:                string(b),
		NextTime:            time.Now().Add(time.Second * rq.opt.FirstDelaySecond),
	}
	if _, err = orm.Insert(ctx, tx, &record); err != nil {
		err = errors.Wrap(err, "reliablequeue generateMessage orm.Insert err")
		return
	}

	return
}

// getDistributeList 获得分发列表
func (rq *ReliableQueue) getDistributeList(ctx context.Context, tx *gorm.DB, msg PublishParams) (distributeList []ReliableMqMessageDistribute, err error) {
	where := map[string]interface{}{"scene": msg.Scene}
	tx = tx.WithContext(ctx).Select("*")
	if err = orm.WithWhere(ctx, tx, where).Find(&distributeList).Error; err != nil {
		err = errors.Wrap(err, "reliablequeue getDistributeList err")
		return
	}
	return
}

// Retry 消费失败退避重试
func (rq *ReliableQueue) Retry(ctx context.Context, tx *gorm.DB, record ReliableMqMessageRecord) (err error) {
	delay := record.Delay * rq.opt.RetryDelaySecondMultiple
	if delay >= math.MaxInt64 || delay < 1 {
		delay = math.MaxInt64
	}

	t := time.Now()
	where := map[string]interface{}{"uuid": record.Uuid}
	update := map[string]interface{}{
		"delay":       delay,
		"retry_time":  t,
		"next_time":   t.Add(time.Second * time.Duration(delay)),
		"update_time": t,
	}
	if _, err = orm.Update(ctx, tx, &ReliableMqMessageRecord{}, where, update); err != nil {
		err = errors.Wrap(err, "reliablequeue Retry orm.Update err")
		return
	}

	return
}

// Republish 重新发布，一般用离线任务调用
func (rq *ReliableQueue) Republish(ctx context.Context, tx *gorm.DB) (err error) {
	records, err := rq.getUnsuccessRecords(ctx, tx)
	if err != nil {
		return
	}
	return rq.publish(ctx, records)
}

// getUnsuccessRecords 获得未成功的记录列表
func (rq *ReliableQueue) getUnsuccessRecords(ctx context.Context, tx *gorm.DB) (records []ReliableMqMessageRecord, err error) {
	where := map[string]interface{}{
		orm.FormatEq("is_success"): RecordStatusUnsuccess,
		orm.FormatLt("next_time"):  time.Now(),
	}
	tx = tx.WithContext(ctx).Select("*")
	if err = orm.WithWhere(ctx, tx, where).Find(&records).Error; err != nil {
		err = errors.Wrap(err, "reliablequeue getUnsuccessRecords err")
		return
	}
	return
}

// publish 分发消息
func (rq *ReliableQueue) publish(ctx context.Context, records []ReliableMqMessageRecord) error {
	g, _ := errgroup.WithContext(ctx)
	for _, v := range records {
		record := v
		g.Go(func() (err error) {
			return rq.queue.Produce(ctx, record)
		})
	}
	return g.Wait()
}

// SetSuccess 消费完成后标记成功
func (rq *ReliableQueue) SetSuccess(ctx context.Context, tx *gorm.DB, record ReliableMqMessageRecord) (err error) {
	where := map[string]interface{}{"uuid": record.Uuid}
	update := map[string]interface{}{"is_success": RecordStatusSuccess}
	if _, err = orm.Update(ctx, tx, &ReliableMqMessageRecord{}, where, update); err != nil {
		err = errors.Wrap(err, "reliablequeue SetSuccess orm.Update err")
		return
	}

	return
}
