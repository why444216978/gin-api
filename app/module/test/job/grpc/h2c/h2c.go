package h2c

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/why444216978/gin-api/app/module/test/job/grpc/helloworld"
	client "github.com/why444216978/gin-api/client/grpc"
	serverGRPC "github.com/why444216978/gin-api/server/grpc"
	"github.com/why444216978/gin-api/server/grpc/h2c"
)

const (
	endpoint = ":8888"
)

type Server struct {
	pb.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: in.Name + " world"}, nil
}

func RegisterServer(s *grpc.Server) {
	pb.RegisterGreeterServer(s, &Server{})
}

func StartServer() {
	if err := h2c.NewH2C(endpoint,
		[]serverGRPC.Register{serverGRPC.NewRegister(RegisterServer, pb.RegisterGreeterHandlerFromEndpoint)},
	).Start(); err != nil {
		panic(err)
	}
}

func Start(ctx context.Context) (err error) {
	go StartServer()
	call()
	return
}

func call() {
	cc, err := client.Conn(context.Background(), endpoint)
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
}
