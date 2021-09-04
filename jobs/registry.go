package jobs

import (
	"context"
	"fmt"
	"gin-api/libraries/jaeger"
	"gin-api/resource"
	"log"
	"net/http"
	"time"
)

func Registry() (err error) {
	ser := resource.Services["gin-api"]
	defer ser.Close()
	err = ser.WatchService(context.Background())
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-time.Tick(1 * time.Second):
			node := ser.GetServices()
			if len(node) <= 0 {
				log.Println("node empty")
				continue
			}
			sendUrl := fmt.Sprintf("http://%s:%d/ping", node[0].Host, node[0].Port)

			ret, err := jaeger.JaegerSend(context.Background(), http.MethodGet, sendUrl, nil, nil, time.Second)
			fmt.Println(ret)
			fmt.Println(err)
		}
	}

	return
}
