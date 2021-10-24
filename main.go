package main

import (
	"log"
	"syscall"

	"github.com/why444216978/gin-api/app"
)

func main() {
	log.Printf("Actual pid is %d", syscall.Getpid())

	app.Init()
	app.Start()
}
