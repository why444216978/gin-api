package main

import (
	"log"
	"syscall"

	"github.com/why444216978/gin-api/app/bootstrap"
)

func main() {
	log.Printf("Actual pid is %d", syscall.Getpid())

	bootstrap.Start()
}
