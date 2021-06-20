package jaeger

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracerLog "github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

//gorm hook
const (
	gormSpanKey        = "gorm_span"
	opentracingSpanKey = "gorm"
	callBackBeforeName = "opentracing:before"
	callBackAfterName  = "opentracing:after"
)

//before gorm before execute action do something
func before(db *gorm.DB) {
	if Tracer == nil {
		panic(ErrNotJaeger)
	}
	span, _ := opentracing.StartSpanFromContextWithTracer(db.Statement.Context, Tracer, opentracingSpanKey)
	db.InstanceSet(gormSpanKey, span)
	return
}

//after gorm after execute action do something
func after(db *gorm.DB) {
	if Tracer == nil {
		panic(ErrNotJaeger)
	}
	_span, isExist := db.InstanceGet(gormSpanKey)
	if !isExist {
		return
	}
	span, ok := _span.(opentracing.Span)
	if !ok {
		return
	}
	setGormTag(db.Statement.Context, span)
	defer span.Finish()
	if db.Error != nil {
		span.LogFields(
			tracerLog.Error(db.Error),
		)
		span.SetTag(string(ext.Error), true)
	}
	span.LogFields(tracerLog.String("SQL", db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)))

	return
}

type OpentracingPlugin struct{}

var GormTrace gorm.Plugin = &OpentracingPlugin{}

func (op *OpentracingPlugin) Name() string {
	return "opentracingPlugin"
}

//Initialize 初始化GormHook
func (op *OpentracingPlugin) Initialize(db *gorm.DB) (err error) {
	//create
	if err = db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before); err != nil {
		return err
	}
	if err = db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after); err != nil {
		return err
	}
	//query
	if err = db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before); err != nil {
		return err
	}
	if err = db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after); err != nil {
		return err
	}
	//delete
	if err = db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before); err != nil {
		return err
	}
	if err = db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after); err != nil {
		return err
	}
	//gorm:begin_transaction
	if err = db.Callback().Update().Before("gorm:begin_transaction").Register(callBackBeforeName, before); err != nil {
		return err
	}
	if err = db.Callback().Update().After("gorm:commit_or_rollback_transaction").Register(callBackAfterName, after); err != nil {
		return err
	}
	if err = db.Callback().Raw().After("gorm:after_transaction").Register(callBackAfterName, after); err != nil {
		return err
	}
	//row raw
	if err = db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before); err != nil {
		return err
	}
	if err = db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before); err != nil {
		return err
	}
	if err = db.Callback().Row().After("gorm:row").Register(callBackAfterName, after); err != nil {
		return err
	}
	if err = db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after); err != nil {
		return err
	}

	//update
	if err = db.Callback().Raw().Before("gorm:before_update").Register(callBackBeforeName, before); err != nil {
		return err
	}
	if err = db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after); err != nil {
		return err
	}
	if err = db.Callback().Raw().Before("gorm:save_before_associations").Register(callBackBeforeName, before); err != nil {
		return err
	}

	if err = db.Callback().Update().After("gorm:save_after_associations").Register(callBackAfterName, after); err != nil {
		return err
	}

	return nil
}

func setGormTag(ctx context.Context, span opentracing.Span) {
	setTag(ctx, span)
	span.SetTag(string(ext.Component), operationTypeGorm)
}
