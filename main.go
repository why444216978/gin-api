package main

import (
	"fmt"
	"gin-api/app_const"
	"gin-api/bootstrap"
	"gin-api/libraries/endless"
	"gin-api/routers"
	"log"
	"strconv"
)

func main() {
	bootstrap.Bootstrap()

	router := routers.InitRouter()

	// server := &http.Server{
	// 	Addr:         fmt.Sprintf(":%s", strconv.Itoa(app_const.SERVICE_PORT)),
	// 	Handler:      router,
	// 	ReadTimeout:  3 * time.Second,
	// 	WriteTimeout: 3 * time.Second,
	// }

	// server.ListenAndServe()

	// endless.DefaultReadTimeOut = 3 * time.Second
	// endless.DefaultWriteTimeOut = 3 * time.Second
	err := endless.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(app_const.SERVICE_PORT)), router)
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
