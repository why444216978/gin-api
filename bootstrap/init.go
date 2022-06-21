package bootstrap

import (
	"github.com/why444216978/gin-api/library/app"
)

func Init(load func() error) (err error) {
	if err = app.InitApp(); err != nil {
		return
	}

	if err = load(); err != nil {
		return
	}

	return
}
