package gorm

import (
	"errors"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/why444216978/gin-api/library/jaeger"
	"github.com/why444216978/go-util/orm"
)

type TestTable struct {
	Id      uint   `gorm:"column:id" json:"id"`
	GoodsId uint64 `gorm:"column:goods_id" json:"goods_id"`
	Name    string `gorm:"column:name" json:"name"`
}

func (TestTable) TableName() string {
	return "test"
}

func Test_before(t *testing.T) {
	convey.Convey("Test_before", t, func() {
		convey.Convey("Tracer nil", func() {
			db := orm.NewMemoryDB()
			if err := db.Migrator().CreateTable(&TestTable{}); err != nil {
				panic(err)
			}

			tracer := mocktracer.New()
			jaeger.Tracer = nil

			before(db)

			assert.Len(t, tracer.FinishedSpans(), 0)
		})
		convey.Convey("success", func() {
			db := orm.NewMemoryDB()
			if err := db.Migrator().CreateTable(&TestTable{}); err != nil {
				panic(err)
			}

			tracer := mocktracer.New()
			jaeger.Tracer = tracer

			before(db)

			assert.Len(t, tracer.FinishedSpans(), 0)
		})
	})
}

func Test_after(t *testing.T) {
	convey.Convey("Test_after", t, func() {
		convey.Convey("Tracer nil", func() {
			db := orm.NewMemoryDB()
			if err := db.Migrator().CreateTable(&TestTable{}); err != nil {
				panic(err)
			}

			tracer := mocktracer.New()
			jaeger.Tracer = nil

			after(db)

			assert.Len(t, tracer.FinishedSpans(), 0)
		})
		convey.Convey("!isExist", func() {
			db := orm.NewMemoryDB()
			if err := db.Migrator().CreateTable(&TestTable{}); err != nil {
				panic(err)
			}

			tracer := mocktracer.New()
			jaeger.Tracer = tracer

			after(db)

			assert.Len(t, tracer.FinishedSpans(), 0)
		})
		convey.Convey("span not opentracing.Span", func() {
			db := orm.NewMemoryDB()
			if err := db.Migrator().CreateTable(&TestTable{}); err != nil {
				panic(err)
			}

			tracer := mocktracer.New()
			jaeger.Tracer = tracer

			db = db.InstanceSet(gormSpanKey, 1)

			after(db)

			assert.Len(t, tracer.FinishedSpans(), 0)
		})
		convey.Convey("success and db err", func() {
			db := orm.NewMemoryDB()
			if err := db.Migrator().CreateTable(&TestTable{}); err != nil {
				panic(err)
			}

			tracer := mocktracer.New()
			jaeger.Tracer = tracer

			span, _ := opentracing.StartSpanFromContextWithTracer(db.Statement.Context, jaeger.Tracer, componentGorm)
			db = db.InstanceSet(gormSpanKey, span)

			after(db)

			assert.Len(t, tracer.FinishedSpans(), 1)
		})
		convey.Convey("success", func() {
			db := orm.NewMemoryDB()
			if err := db.Migrator().CreateTable(&TestTable{}); err != nil {
				panic(err)
			}

			tracer := mocktracer.New()
			jaeger.Tracer = tracer

			span, _ := opentracing.StartSpanFromContextWithTracer(db.Statement.Context, jaeger.Tracer, componentGorm)
			db = db.InstanceSet(gormSpanKey, span)
			db.Error = errors.New("")

			after(db)

			assert.Len(t, tracer.FinishedSpans(), 1)
		})
	})
}

func TestOpentracingPlugin_Name(t *testing.T) {
	convey.Convey("TestOpentracingPlugin_Name", t, func() {
		convey.Convey("success", func() {
			op := &opentracingPlugin{}
			assert.Equal(t, op.Name(), "opentracingPlugin")
		})
	})
}

func TestOpentracingPlugin_Initialize(t *testing.T) {
	convey.Convey("TestOpentracingPlugin_Name", t, func() {
		convey.Convey("success", func() {
			db := orm.NewMemoryDB()
			op := &opentracingPlugin{}
			assert.Equal(t, op.Initialize(db), nil)
		})
	})
}
