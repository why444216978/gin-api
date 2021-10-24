package job

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

type HandleFunc func(ctx context.Context) error

var Handlers = map[string]HandleFunc{}

func Handle(job string) {
	ctx := context.Background()

	log.Println("start by " + job)

	handle, ok := Handlers[job]
	if !ok {
		fmt.Println("job " + job + " not found")
		return
	}

	err := handle(ctx)

	time.Sleep(1 * time.Second)

	if err != nil {
		os.Exit(1)
	}
}
