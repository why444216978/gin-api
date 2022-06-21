package service

import (
	"github.com/why444216978/gin-api/library/logger"
	"github.com/why444216978/gin-api/library/logger/zap"
)

type Config struct {
	InfoFile  string
	ErrorFile string
	Level     string
}

type ServiceLogger struct {
	*zap.ZapLogger
	config *Config
}

func NewServiceLogger(serviceName string, config *Config) (*ServiceLogger, error) {
	infoWriter, errWriter, err := logger.RotateWriter(config.InfoFile, config.ErrorFile)
	if err != nil {
		return nil, err
	}

	l, err := zap.NewLogger(
		zap.WithModule(logger.ModuleHTTP),
		zap.WithServiceName(serviceName),
		zap.WithInfoWriter(infoWriter),
		zap.WithErrorWriter(errWriter),
		zap.WithLevel(config.Level),
	)
	if err != nil {
		return nil, err
	}

	return &ServiceLogger{ZapLogger: l, config: config}, nil
}
