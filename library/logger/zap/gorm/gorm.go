package gorm

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/why444216978/gin-api/library/logger"
	zapLogger "github.com/why444216978/gin-api/library/logger/zap"
)

// GormConfig is used to parse configuration file
// logger should be controlled with Options
type GormConfig struct {
	ServiceName               string
	SlowThreshold             int
	InfoFile                  string
	ErrorFile                 string
	Level                     int
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

type GormLogger struct {
	*zap.Logger
	config                    *GormConfig
	LogLevel                  gormLogger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

type GormOption func(gl *GormLogger)

var _ gormLogger.Interface = (*GormLogger)(nil)

func NewGorm(config *GormConfig, opts ...GormOption) (gl *GormLogger, err error) {
	gl = &GormLogger{
		config:                    config,
		LogLevel:                  gormLogger.LogLevel(config.Level),
		SlowThreshold:             time.Duration(config.SlowThreshold) * time.Millisecond,
		SkipCallerLookup:          config.SkipCallerLookup,
		IgnoreRecordNotFoundError: config.IgnoreRecordNotFoundError,
	}

	for _, o := range opts {
		o(gl)
	}

	level := zap.InfoLevel.String()
	switch gl.LogLevel {
	case gormLogger.Silent:
		level = zapcore.FatalLevel.String()
	case gormLogger.Error:
		level = zapcore.ErrorLevel.String()
	case gormLogger.Warn:
		level = zapcore.WarnLevel.String()
	case gormLogger.Info:
		level = zapcore.InfoLevel.String()
	}

	infoWriter, errWriter, err := logger.RotateWriter(config.InfoFile, config.ErrorFile)
	if err != nil {
		return
	}

	l, err := zapLogger.NewLogger(
		zapLogger.WithModule(logger.ModuleMySQL),
		zapLogger.WithServiceName(config.ServiceName),
		zapLogger.WithInfoWriter(infoWriter),
		zapLogger.WithErrorWriter(errWriter),
		zapLogger.WithLevel(level),
	)
	if err != nil {
		return
	}
	gl.Logger = l.Logger

	gormLogger.Default = gl

	return
}

func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return &GormLogger{
		Logger:                    l.Logger,
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		SkipCallerLookup:          l.SkipCallerLookup,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
	}
}

func (l *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {}

func (l *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {}

func (l *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}

	fields := logger.ValueFields(ctx)

	elapsed := time.Since(begin)

	sql, rows := fc()
	sqlSlice := strings.Split(sql, " ")
	api := ""
	if len(sqlSlice) > 1 {
		api = sqlSlice[0]
	}

	f := []zapcore.Field{
		zap.String(logger.LogID, logger.ValueLogID(ctx)),
		zap.String(logger.TraceID, logger.ValueTraceID(ctx)),
		zap.Int64(logger.Cost, elapsed.Milliseconds()),
		zap.String(logger.Request, sql),
		zap.Int64(logger.Response, rows),
		zap.String(logger.API, api),
		zap.Reflect(logger.ClientIP, logger.Find(logger.ServerIP, fields)),
		zap.Reflect(logger.ClientPort, logger.Find(logger.ServerPort, fields)),
	}

	switch {
	case err != nil && l.LogLevel >= gormLogger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		l.logger().Error(err.Error(), f...)
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormLogger.Warn:
		l.logger().Warn("warn", f...)
	case l.LogLevel >= gormLogger.Info:
		l.logger().Info("info", f...)
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
			return l.Logger.WithOptions(zap.AddCallerSkip(i - 2))
		}
	}
	return l.Logger
}

func (l *GormLogger) Close() error {
	return l.Logger.Sync()
}
