package logging

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"time"

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

func New(zapLogger *zap.Logger) *GormLogger {
	return &GormLogger{
		ZapLogger:                 zapLogger,
		LogLevel:                  logger.Warn,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: false,
	}
}

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

	l, err := NewLogger(&Config{
		InfoFile:  cfg.InfoFile,
		ErrorFile: cfg.ErrorFile,
		Level:     zapLever,
	})
	if err != nil {
		return
	}
	gl.ZapLogger = l.GetLogger()

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

func (l *GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < logger.Info {
		return
	}
	l.logger().Sugar().Debugf(str, args...)
}

func (l *GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < logger.Warn {
		return
	}
	l.logger().Sugar().Warnf(str, args...)
}

func (l *GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < logger.Error {
		return
	}
	l.logger().Sugar().Errorf(str, args...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		l.logger().Error("trace", zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		l.logger().Warn("trace", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.LogLevel >= logger.Info:
		sql, rows := fc()
		l.logger().Info("trace", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
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
