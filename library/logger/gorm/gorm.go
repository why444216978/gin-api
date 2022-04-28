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
	Config                    *GormConfig
	ZapLogger                 *zap.Logger
	LogLevel                  gormLogger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

type GormOption func(gl *GormLogger)

var _ gormLogger.Interface = (*GormLogger)(nil)

func NewGorm(cfg *GormConfig, opts ...GormOption) (gl *GormLogger, err error) {
	gl = &GormLogger{
		Config:                    cfg,
		LogLevel:                  gormLogger.LogLevel(cfg.Level),
		SlowThreshold:             time.Duration(cfg.SlowThreshold) * time.Millisecond,
		SkipCallerLookup:          cfg.SkipCallerLookup,
		IgnoreRecordNotFoundError: cfg.IgnoreRecordNotFoundError,
	}

	for _, o := range opts {
		o(gl)
	}

	zapLever := zap.InfoLevel.String()
	switch gl.LogLevel {
	case gormLogger.Silent:
		zapLever = zapcore.FatalLevel.String()
	case gormLogger.Error:
		zapLever = zapcore.ErrorLevel.String()
	case gormLogger.Warn:
		zapLever = zapcore.WarnLevel.String()
	case gormLogger.Info:
		zapLever = zapcore.InfoLevel.String()
	}

	infoWriter, errWriter, err := logger.RotateWriter(cfg.InfoFile, cfg.ErrorFile)
	if err != nil {
		return
	}

	l, err := logger.NewLogger(&logger.Config{
		InfoFile:  cfg.InfoFile,
		ErrorFile: cfg.ErrorFile,
		Level:     zapLever,
	},
		logger.WithModule(logger.ModuleMySQL),
		logger.WithServiceName(cfg.ServiceName),
		logger.WithInfoWriter(infoWriter),
		logger.WithErrorWriter(errWriter),
	)
	if err != nil {
		return
	}
	gl.ZapLogger = l.Logger

	gormLogger.Default = gl

	return
}

func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return &GormLogger{
		ZapLogger:                 l.ZapLogger,
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

	logFields := logger.ValueHTTPFields(ctx)

	elapsed := time.Since(begin)

	sql, rows := fc()
	sqlSlice := strings.Split(sql, " ")
	api := ""
	if len(sqlSlice) > 1 {
		api = sqlSlice[0]
	}

	fields := []zapcore.Field{
		zap.String(logger.LogID, logger.ValueLogID(ctx)),
		zap.String(logger.TraceID, logger.ValueTraceID(ctx)),
		zap.Int64(logger.Cost, elapsed.Milliseconds()),
		zap.String(logger.Request, sql),
		zap.Int64(logger.Response, rows),
		zap.String(logger.API, api),
		zap.String(logger.ClientIP, logFields.ServerIP),
		zap.Int(logger.ClientPort, logFields.ServerPort),
	}

	switch {
	case err != nil && l.LogLevel >= gormLogger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		l.logger().Error(err.Error(), fields...)
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormLogger.Warn:
		l.logger().Warn("warn", fields...)
	case l.LogLevel >= gormLogger.Info:
		l.logger().Info("info", fields...)
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
