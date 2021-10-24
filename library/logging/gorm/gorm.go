package logging

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"time"

	"github.com/why444216978/gin-api/library/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormConfig struct {
	SlowThreshold             int
	InfoFile                  string
	ErrorFile                 string
	Level                     int
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

type GormLogger struct {
	ZapLogger                 *zap.Logger
	LogLevel                  logger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

type GormOption func(gl *GormLogger)

var _ logger.Interface = (*GormLogger)(nil)

func NewGorm(cfg *GormConfig, opts ...GormOption) (gl *GormLogger, err error) {
	gl = &GormLogger{
		LogLevel:                  logger.LogLevel(cfg.Level),
		SlowThreshold:             time.Duration(cfg.SlowThreshold) * time.Millisecond,
		SkipCallerLookup:          cfg.SkipCallerLookup,
		IgnoreRecordNotFoundError: cfg.IgnoreRecordNotFoundError,
	}

	zapLever := zap.InfoLevel.String()
	switch gl.LogLevel {
	case logger.Silent:
		zapLever = zapcore.FatalLevel.String()
	case logger.Error:
		zapLever = zapcore.ErrorLevel.String()
	case logger.Warn:
		zapLever = zapcore.WarnLevel.String()
	case logger.Info:
		zapLever = zapcore.InfoLevel.String()
	}

	for _, o := range opts {
		o(gl)
	}

	l, err := logging.NewLogger(&logging.Config{
		InfoFile:  cfg.InfoFile,
		ErrorFile: cfg.ErrorFile,
		Level:     zapLever,
	}, logging.WithModule(logging.ModuleMySQL))
	if err != nil {
		return
	}
	gl.ZapLogger = l.Logger

	logger.Default = gl

	return
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &GormLogger{
		ZapLogger:                 l.ZapLogger,
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		SkipCallerLookup:          l.SkipCallerLookup,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
	}
}

func (l *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel < logger.Info {
		return
	}
	l.logger().Info(msg, zap.Reflect("data", args))

}

func (l *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel < logger.Warn {
		return
	}
	l.logger().Warn(msg, zap.Reflect("data", args))
}

func (l *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel < logger.Error {
		return
	}
	l.logger().Error(msg, zap.Reflect("data", args))
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}

	elapsed := time.Since(begin)

	sql, rows := fc()
	fields := []zapcore.Field{
		zap.Error(err),
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
		zap.String("sql", sql),
		zap.String(logging.LogID,
			logging.ValueTraceID(ctx)),
		zap.String(logging.TraceID,
			logging.ValueLogID(ctx)),
	}

	switch {
	case err != nil && l.LogLevel >= logger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		l.logger().Error("trace", fields...)
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= logger.Warn:
		l.logger().Warn("trace", fields...)
	case l.LogLevel >= logger.Info:
		l.logger().Info("trace", fields...)
	}
}

func (l *GormLogger) logger() *zap.Logger {
	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.Contains(file, "gorm.io"):
		case strings.Contains(file, "go-util/orm/orm.go"):
		default:
			return l.ZapLogger.WithOptions(zap.AddCallerSkip(i - 2))
		}
	}
	return l.ZapLogger
}
