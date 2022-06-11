package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/why444216978/gin-api/app/config"
	pb "github.com/why444216978/gin-api/app/module/test/service/grpc/helloworld"
	client "github.com/why444216978/gin-api/client/grpc"
)

func Start(ctx context.Context) (err error) {
	call()
	return
}

func call() {
	cc, err := client.Conn(context.Background(), fmt.Sprintf(":%d", config.App.AppPort))
	if err != nil {
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewGreeterClient(cc)

	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "why"})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(reply)
	}

	return
}
