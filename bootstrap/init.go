package bootstrap

import (
	"flag"
	"log"
	"syscall"

	"github.com/why444216978/gin-api/library/app"
)

func Init(load func() error) (err error) {
	flag.Parse()

	log.Printf("Actual pid is %d", syscall.Getpid())

	if err = app.InitApp(); err != nil {
		return
	}

	if err = load(); err != nil {
		return
	}

	return
}
