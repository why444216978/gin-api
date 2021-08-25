package jobs

import (
	"fmt"
	"log"
	"os"
	"time"
)

type handleFunc func() error

var handlers = map[string]handleFunc{}

func Handle(job string) {
	log.Println("start by " + job)

	handle, ok := handlers[job]
	if !ok {
		fmt.Println("job " + job + " not found")
		return
	}

	err := handle()

	time.Sleep(1 * time.Second)

	if err != nil {
		os.Exit(1)
	}
}
