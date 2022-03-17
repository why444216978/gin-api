package gorm

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracing_log "github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"

	"github.com/why444216978/gin-api/library/jaeger"
	"github.com/why444216978/go-util/assert"
)

// gorm hook
const (
	componentGorm      = "Gorm"
	gormSpanKey        = "gorm_span"
	callBackBeforeName = "opentracing:before"
	callBackAfterName  = "opentracing:after"
)

// before gorm before execute action do something
func before(db *gorm.DB) {
	if assert.IsNil(jaeger.Tracer) {
		return
	}
	span, _ := opentracing.StartSpanFromContextWithTracer(db.Statement.Context, jaeger.Tracer, componentGorm)
	db.InstanceSet(gormSpanKey, span)
	return
}

// after gorm after execute action do something
func after(db *gorm.DB) {
	if assert.IsNil(jaeger.Tracer) {
		return
	}
	_span, isExist := db.InstanceGet(gormSpanKey)
	if !isExist {
		return
	}
	span, ok := _span.(opentracing.Span)
	if !ok {
		return
	}
	defer span.Finish()

	jaeger.SetCommonTag(db.Statement.Context, span)

	if db.Error != nil {
		span.LogFields(opentracing_log.Error(db.Error))
		span.SetTag(string(ext.Error), true)
	}
	span.LogFields(opentracing_log.String("SQL", db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)))

	return
}

type opentracingPlugin struct{}

var GormTrace gorm.Plugin = &opentracingPlugin{}

func (op *opentracingPlugin) Name() string {
	return "opentracingPlugin"
}

func (op *opentracingPlugin) Initialize(db *gorm.DB) (err error) {
	// create
	if err = db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before); err != nil {
		return err
	}
	if err = db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after); err != nil {
		return err
	}

	// query
	if err = db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before); err != nil {
		return err
	}
	if err = db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after); err != nil {
		return err
	}

	// delete
	if err = db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before); err != nil {
		return err
	}
	if err = db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after); err != nil {
		return err
	}

	// update
	if err = db.Callback().Update().Before("gorm:before_update").Register(callBackBeforeName, before); err != nil {
		return err
	}
	if err = db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after); err != nil {
		return err
	}

	// row
	if err = db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before); err != nil {
		return err
	}
	if err = db.Callback().Row().After("gorm:row").Register(callBackAfterName, after); err != nil {
		return err
	}

	// raw
	if err = db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before); err != nil {
		return err
	}
	if err = db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after); err != nil {
		return err
	}

	// associations
	if err = db.Callback().Raw().Before("gorm:save_before_associations").Register(callBackBeforeName, before); err != nil {
		return err
	}
	if err = db.Callback().Update().After("gorm:save_after_associations").Register(callBackAfterName, after); err != nil {
		return err
	}
	return nil
}
