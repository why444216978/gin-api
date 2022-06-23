package bootstrap

import (
	"log"
	"syscall"

	"github.com/why444216978/gin-api/library/app"
	"github.com/why444216978/gin-api/library/config"
)

func Init(env string, load func() error) (err error) {
	log.Printf("Actual pid is %d", syscall.Getpid())

	config.Init(env)

	if err = app.InitApp(); err != nil {
		return
	}

	if err = load(); err != nil {
		return
	}

	return
}
