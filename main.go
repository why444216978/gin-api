package main

import (
	"fmt"
	"gin-api/app_const"
	"gin-api/bootstrap"
	"gin-api/routers"
	"net/http"
	"strconv"
	"time"
)

func main() {
	bootstrap.Bootstrap()

	router := routers.InitRouter()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", strconv.Itoa(app_const.SERVICE_PORT)),
		Handler:      router,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	server.ListenAndServe()

	// endless.DefaultReadTimeOut = 3 * time.Second
	// endless.DefaultWriteTimeOut = 3 * time.Second
	// server := endless.NewServer(fmt.Sprintf(":%s", strconv.Itoa(app_const.SERVICE_PORT)), router)
	// server.BeforeBegin = func(add string) {
	// 	log.Printf("Actual pid is %d", syscall.Getpid())
	// }
	// err := server.ListenAndServe()
	// if err != nil {
	// 	log.Printf("Server err: %v", err)
	// }
}
