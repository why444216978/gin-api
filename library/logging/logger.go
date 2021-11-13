package logging

import (
	"context"
	"errors"
	"io"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/why444216978/go-util/conversion"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
	level       zapcore.Level
	callSkip    int
	module      string
	serviceName string
}

type Config struct {
	InfoFile  string
	ErrorFile string
	Level     string
}

type Option func(l *Logger)

func WithCallerSkip(skip int) Option {
	return func(l *Logger) { l.callSkip = skip }
}

func WithModule(module string) Option {
	return func(l *Logger) { l.module = module }
}

func WithServiceName(serviceName string) Option {
	return func(l *Logger) { l.serviceName = serviceName }
}

func NewLogger(cfg *Config, opts ...Option) (l *Logger, err error) {
	level, err := zapLevel(cfg.Level)
	if err != nil {
		return
	}

	l = &Logger{
		level: level,
	}

	for _, o := range opts {
		o(l)
	}

	encoder := l.formatEncoder()

	infoEnabler := l.infoEnabler()
	errorEnabler := l.errorEnabler()

	infoWriter, err := l.getWriter(cfg.InfoFile)
	if err != nil {
		return
	}
	errorWriter, err := l.getWriter(cfg.ErrorFile)
	if err != nil {
		return
	}

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoEnabler),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorEnabler),
	)

	if l.callSkip == 0 {
		l.callSkip = 1
	}

	fields := make([]zapcore.Field, 0)

	if l.module != "" {
		fields = append(fields, zap.String(Module, l.module))
	}

	if l.serviceName != "" {
		fields = append(fields, zap.String(SericeName, l.serviceName))
	}

	l.Logger = zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(errorEnabler),
		zap.AddCallerSkip(l.callSkip),
		zap.Fields(fields...),
	)

	return
}

func (l *Logger) infoEnabler() zap.LevelEnablerFunc {
	return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if lvl < l.level {
			return false
		}
		return lvl <= zapcore.InfoLevel
	})
}

func (l *Logger) errorEnabler() zap.LevelEnablerFunc {
	return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if lvl < l.level {
			return false
		}
		return lvl >= zapcore.WarnLevel
	})
}

func (l *Logger) formatEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		TimeKey:       "time",
		CallerKey:     "file",
		FunctionKey:   "func",
		StacktraceKey: "stack",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
}

func (l *Logger) getWriter(filename string) (io.Writer, error) {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		strings.Replace(filename, ".log", "", -1)+"-%Y%m%d%H.log", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		return nil, err
	}

	return hook, nil
}

func (l *Logger) GetLevel() zapcore.Level {
	return l.level
}

func zapLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug", "DEBUG":
		return zapcore.DebugLevel, nil
	case "info", "INFO", "":
		return zapcore.InfoLevel, nil
	case "warn", "WARN":
		return zapcore.WarnLevel, nil
	case "error", "ERROR":
		return zapcore.ErrorLevel, nil
	case "dpanic", "DPANIC":
		return zapcore.DPanicLevel, nil
	case "panic", "PANIC":
		return zapcore.PanicLevel, nil
	case "fatal", "FATAL":
		return zapcore.FatalLevel, nil
	default:
		return 0, errors.New("error level:" + level)
	}
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, append(fields, l.extractFields(ctx)...)...)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Info(msg, append(fields, l.extractFields(ctx)...)...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, append(fields, l.extractFields(ctx)...)...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Error(msg, append(fields, l.extractFields(ctx)...)...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, append(fields, l.extractFields(ctx)...)...)
}

func (l *Logger) extractFields(ctx context.Context) []zap.Field {
	fieldsMap, _ := conversion.StructToMap(ValueHTTPFields(ctx))

	fields := make([]zap.Field, len(fieldsMap))

	i := 0
	for k, v := range fieldsMap {
		fields[i] = zap.Reflect(k, v)
		i = i + 1
	}

	return fields
}
