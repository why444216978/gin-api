package app

import (
	"time"

	"github.com/why444216978/gin-api/library/config"
)

var app struct {
	AppName        string
	AppPort        int
	Pprof          bool
	IsDebug        bool
	ContextTimeout int
	ConnectTimeout int
	WriteTimeout   int
	ReadTimeout    int
}

func InitApp() (err error) {
	return config.ReadConfig("app", "toml", &app)
}

func Name() string {
	return app.AppName
}

func Port() int {
	return app.AppPort
}

func Pprof() bool {
	return app.Pprof
}

func Debug() bool {
	return app.IsDebug
}

func ContextTimeout() time.Duration {
	if app.ConnectTimeout == 0 {
		return time.Duration(100000) * time.Millisecond
	}
	return time.Duration(app.ContextTimeout) * time.Millisecond
}

func ConnectTimeout() time.Duration {
	if app.ConnectTimeout == 0 {
		return time.Duration(100000) * time.Millisecond
	}
	return time.Duration(app.ConnectTimeout) * time.Millisecond
}

func WriteTimeout() time.Duration {
	if app.WriteTimeout == 0 {
		return time.Duration(100000) * time.Millisecond
	}
	return time.Duration(app.WriteTimeout) * time.Millisecond
}

func ReadTimeout() time.Duration {
	if app.ReadTimeout == 0 {
		return time.Duration(100000) * time.Millisecond
	}
	return time.Duration(app.ReadTimeout) * time.Millisecond
}
