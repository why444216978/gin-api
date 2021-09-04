package jobs

import (
	"context"
	"fmt"
	"net/http"
	"time"

	lib_http "gin-api/libraries/http"
)

func Registry() (err error) {
	for {
		select {
		case <-time.Tick(1 * time.Second):
			ret, err := lib_http.Send(context.Background(), "gin-api", http.MethodGet, "/ping", nil, nil, time.Second)
			fmt.Println(ret)
			fmt.Println(err)
		}
	}

	return
}
