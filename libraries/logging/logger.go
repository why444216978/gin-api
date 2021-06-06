package logging

import (
	"io"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	encoder zapcore.Encoder
	logger  *zap.Logger
}

type Config struct {
	InfoFile  string
	ErrorFile string
	Level     string
}

func NewLogger(cfg Config) (logger *Logger) {
	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	// 实现两个判断日志等级的interface
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if lvl.String() < cfg.Level {
			return false
		}
		return lvl <= zapcore.InfoLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if lvl.String() < cfg.Level {
			return false
		}
		return lvl >= zapcore.WarnLevel
	})

	logger = &Logger{
		encoder: encoder,
	}

	// 获取 info、error日志文件的io.Writer 抽象 getWriter() 在下方实现
	infoWriter := logger.getWriter(cfg.InfoFile)
	errorWriter := logger.getWriter(cfg.ErrorFile)

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
	)

	// 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数
	logger.logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return
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
