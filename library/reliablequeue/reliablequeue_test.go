package reliablequeue

import (
	"context"
	"math"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/why444216978/go-util/orm"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/gorm"

	"github.com/why444216978/gin-api/library/queue"
)

type Queue struct{}

func (*Queue) Produce(ctx context.Context, msg interface{}, opts ...queue.ProduceOptionFunc) error {
	return nil
}

func (*Queue) Consume(consumer queue.Consumer) {}

func createTable() (db *gorm.DB, err error) {
	db = orm.NewMemoryDB()
	if err = db.Migrator().CreateTable(&ReliableMqMessage{}); err != nil {
		return
	}
	if err = db.Migrator().CreateTable(&ReliableMqMessageDistribute{}); err != nil {
		return
	}
	if err = db.Migrator().CreateTable(&ReliableMqMessageRecord{}); err != nil {
		return
	}
	return
}

func Test_defaultOption(t *testing.T) {
	Convey("Test_defaultOption", t, func() {
		Convey("success", func() {
			opt := defaultOption()
			assert.Equal(t, opt.FirstDelaySecond, defaultFirstDelaySecond)
			assert.Equal(t, opt.RetryDelaySecondMultiple, defaultRetryDelaySecondMultiple)
		})
	})
}

func TestNewReliableQueue(t *testing.T) {
	Convey("TestNewReliableQueue", t, func() {
		Convey("success", func() {
			rq, err := NewReliableQueue(&Queue{}, WithFirstDelaySecond(time.Second*10), WithRetryDelaySecondMultiple(int64(10)))
			assert.Equal(t, err, nil)
			assert.Equal(t, rq != nil, true)
		})
		Convey("Queue is nil", func() {
			rq, err := NewReliableQueue(nil)
			assert.Equal(t, err.Error(), "Queue is nil")
			assert.Equal(t, rq, nil)
		})
	})
}

func TestReliableQueue_Publish(t *testing.T) {
	Convey("TestReliableQueue_Publish", t, func() {
		Convey("success", func() {
			db, err := createTable()
			assert.Equal(t, err, nil)
			defer orm.CloseMemoryDB(db)

			distrubute := &ReliableMqMessageDistribute{
				Scene:       "test",
				MessageId:   1,
				ServiceName: "test",
				Uri:         "/test",
				Method:      "POST",
			}
			err = db.Create(distrubute).Error
			assert.Equal(t, err, nil)

			rq, err := NewReliableQueue(&Queue{})
			assert.Equal(t, err, nil)

			msg := PublishParams{
				LogID: "logId",
				Scene: "test",
				Data:  map[string]interface{}{"a": "a"},
			}
			err = rq.Publish(context.Background(), db.Begin(), msg)
			assert.Equal(t, err, nil)

			records := []ReliableMqMessageRecord{}
			err = db.Select("*").Where("message_distribute_id", distrubute.Id).Find(&records).Error
			assert.Equal(t, err, nil)
			assert.Equal(t, len(records), 1)
		})
		Convey("tx is nil", func() {
			rq, err := NewReliableQueue(&Queue{})
			assert.Equal(t, err, nil)

			msg := PublishParams{
				LogID: "logId",
				Scene: "test",
				Data:  map[string]interface{}{"a": "a"},
			}
			err = rq.Publish(context.Background(), nil, msg)
			assert.Equal(t, err.Error(), "tx is nil")
		})
	})
}

func TestReliableQueue_generateMessage(t *testing.T) {
	Convey("TestReliableQueue_generateMessage", t, func() {
		Convey("success", func() {
			db, err := createTable()
			assert.Equal(t, err, nil)
			defer orm.CloseMemoryDB(db)

			rq, err := NewReliableQueue(&Queue{})
			assert.Equal(t, err, nil)

			distrubute := ReliableMqMessageDistribute{
				Scene:       "test",
				MessageId:   1,
				ServiceName: "test",
				Uri:         "/test",
				Method:      "POST",
			}
			msg := PublishParams{
				LogID: "logId",
				Scene: "test",
				Data:  map[string]interface{}{"a": "a"},
			}
			record, err := rq.generateMessage(context.Background(), db, distrubute, msg)
			assert.Equal(t, err, nil)
			assert.Equal(t, record.MessageId, distrubute.MessageId)
			assert.Equal(t, record.MessageDistributeId, distrubute.Id)
			assert.Equal(t, record.LogId, msg.LogID)
			assert.Equal(t, record.ServiceName, distrubute.ServiceName)
			assert.Equal(t, record.Uri, distrubute.Uri)
			assert.Equal(t, record.Method, distrubute.Method)
		})
	})
}

func TestReliableQueue_getDistributeList(t *testing.T) {
	Convey("TestReliableQueue_getDistributeList", t, func() {
		Convey("success", func() {
			db, err := createTable()
			assert.Equal(t, err, nil)
			defer orm.CloseMemoryDB(db)

			distrubute := &ReliableMqMessageDistribute{
				Scene:       "test",
				MessageId:   1,
				ServiceName: "test",
				Uri:         "/test",
				Method:      "POST",
			}
			err = db.Create(distrubute).Error
			assert.Equal(t, err, nil)

			rq, err := NewReliableQueue(&Queue{})
			assert.Equal(t, err, nil)

			msg := PublishParams{
				LogID: "logId",
				Scene: "test",
				Data:  map[string]interface{}{"a": "a"},
			}
			distributeList, err := rq.getDistributeList(context.Background(), db, msg)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(distributeList), 1)
			for _, v := range distributeList {
				assert.Equal(t, v.Scene, distrubute.Scene)
			}
		})
	})
}

func TestReliableQueue_Retry(t *testing.T) {
	Convey("TestReliableQueue_Retry", t, func() {
		Convey("success", func() {
			db, err := createTable()
			assert.Equal(t, err, nil)
			defer orm.CloseMemoryDB(db)

			record := &ReliableMqMessageRecord{
				Uuid:  "uuid",
				Delay: 60,
			}
			err = db.Create(record).Error
			assert.Equal(t, err, nil)

			rq, err := NewReliableQueue(&Queue{})
			assert.Equal(t, err, nil)

			err = rq.Retry(context.Background(), db, *record)
			assert.Equal(t, err, nil)

			target := &ReliableMqMessageRecord{}
			err = db.Select("*").Where("uuid", record.Uuid).First(target).Error
			assert.Equal(t, err, nil)
			assert.Equal(t, target.Delay, record.Delay*rq.opt.RetryDelaySecondMultiple)
		})
		Convey("delay >= math.MaxUint64", func() {
			db, err := createTable()
			assert.Equal(t, err, nil)
			defer orm.CloseMemoryDB(db)

			record := &ReliableMqMessageRecord{
				Uuid:  "uuid",
				Delay: math.MaxInt64,
			}
			err = db.Create(record).Error
			assert.Equal(t, err, nil)

			rq, err := NewReliableQueue(&Queue{})
			assert.Equal(t, err, nil)

			err = rq.Retry(context.Background(), db, *record)
			assert.Equal(t, err, nil)

			target := &ReliableMqMessageRecord{}
			err = db.Select("*").Where("uuid", record.Uuid).First(target).Error
			assert.Equal(t, err, nil)
			assert.Equal(t, target.Delay, int64(math.MaxInt64))
		})
	})
}

func TestReliableQueue_Republish(t *testing.T) {
	Convey("TestReliableQueue_Republish", t, func() {
		Convey("success", func() {
			db, err := createTable()
			assert.Equal(t, err, nil)
			defer orm.CloseMemoryDB(db)

			record := &ReliableMqMessageRecord{
				Uuid:      "uuid",
				Delay:     60,
				NextTime:  time.Now().Add(-time.Minute),
				IsSuccess: RecordStatusUnsuccess,
			}
			err = db.Create(record).Error
			assert.Equal(t, err, nil)

			rq, err := NewReliableQueue(&Queue{})
			assert.Equal(t, err, nil)

			err = rq.Republish(context.Background(), db)
			assert.Equal(t, err, nil)
		})
	})
}

func TestReliableQueue_getUnsuccessRecords(t *testing.T) {
	Convey("TestReliableQueue_getUnsuccessRecords", t, func() {
		Convey("success", func() {
			db, err := createTable()
			assert.Equal(t, err, nil)
			defer orm.CloseMemoryDB(db)

			record := &ReliableMqMessageRecord{
				Uuid:      "uuid",
				NextTime:  time.Now().Add(-time.Minute),
				IsSuccess: RecordStatusUnsuccess,
			}
			err = db.Create(record).Error
			assert.Equal(t, err, nil)

			rq, err := NewReliableQueue(&Queue{})
			assert.Equal(t, err, nil)

			records, err := rq.getUnsuccessRecords(context.Background(), db)
			assert.Equal(t, err, nil)

			for _, v := range records {
				assert.Equal(t, v.IsSuccess, RecordStatusUnsuccess)
			}
		})
	})
}

func TestReliableQueue_publish(t *testing.T) {
	Convey("TestReliableQueue_publish", t, func() {
		Convey("success", func() {
			db, err := createTable()
			assert.Equal(t, err, nil)
			defer orm.CloseMemoryDB(db)

			rq, err := NewReliableQueue(&Queue{})
			assert.Equal(t, err, nil)

			records := []ReliableMqMessageRecord{
				{
					Uuid:      "uuid",
					Delay:     60,
					IsSuccess: RecordStatusUnsuccess,
				},
			}
			err = rq.publish(context.Background(), records)
			assert.Equal(t, err, nil)
		})
	})
}

func TestReliableQueue_SetSuccess(t *testing.T) {
	Convey("TestReliableQueue_SetSuccess", t, func() {
		Convey("success", func() {
			db, err := createTable()
			assert.Equal(t, err, nil)
			defer orm.CloseMemoryDB(db)

			record := &ReliableMqMessageRecord{
				Uuid:      "uuid",
				IsSuccess: RecordStatusUnsuccess,
			}
			err = db.Create(record).Error
			assert.Equal(t, err, nil)

			rq, err := NewReliableQueue(&Queue{})
			assert.Equal(t, err, nil)

			err = rq.SetSuccess(context.Background(), db, *record)
			assert.Equal(t, err, nil)

			target := &ReliableMqMessageRecord{}
			err = db.Select("*").Where("uuid", record.Uuid).First(target).Error
			assert.Equal(t, err, nil)
			assert.Equal(t, target.IsSuccess, RecordStatusSuccess)
		})
	})
}
