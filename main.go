package main

import (
	"gin-api/bootstrap"
	"log"
	"syscall"
)

func main() {
	log.Printf("Actual pid is %d", syscall.Getpid())

	bootstrap.Bootstrap()
	bootstrap.StartApp()
}
