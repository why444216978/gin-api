package logging

import (
	"errors"
	"io"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
	level  zapcore.Level
}

type Config struct {
	InfoFile  string
	ErrorFile string
	Level     string
}

type Option func(l *Logger)

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

	// 实现两个判断日志等级的interface
	infoLevel := l.infoEnabler()
	errorLevel := l.errorEnabler()

	// 获取 info、error日志文件的io.Writer 抽象 getWriter() 在下方实现
	infoWriter := l.getWriter(cfg.InfoFile)
	errorWriter := l.getWriter(cfg.ErrorFile)

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
	)

	// 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数
	l.logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

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
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "time",
		CallerKey:   "file",
		FunctionKey: "func",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
}

func (l *Logger) getWriter(filename string) io.Writer {
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
		panic(err)
	}
	return hook
}

func (l *Logger) GetLogger() *zap.Logger {
	return l.logger
}

func (l *Logger) GetLevel() zapcore.Level {
	return l.level
}

func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	data := l.withFields(fields)
	l.logger.Debug(msg, data...)
}

func (l *Logger) Info(msg string, fields map[string]interface{}) {
	data := l.withFields(fields)
	l.logger.Info(msg, data...)
}

func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	data := l.withFields(fields)
	l.logger.Warn(msg, data...)
}

func (l *Logger) Error(msg string, fields map[string]interface{}) {
	data := l.withFields(fields)
	l.logger.Error(msg, data...)
}

func (l *Logger) Panic(msg string, fields map[string]interface{}) {
	data := l.withFields(fields)
	l.logger.Panic(msg, data...)
}

func (l *Logger) Fatal(msg string, fields map[string]interface{}) {
	data := l.withFields(fields)
	l.logger.Fatal(msg, data...)
}

func (l *Logger) withFields(fields map[string]interface{}) []zapcore.Field {
	ret := make([]zapcore.Field, 0)
	for k, v := range fields {
		ret = append(ret, zap.Reflect(k, v))
	}
	return ret
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
