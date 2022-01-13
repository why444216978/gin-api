package grpc

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"

	pb "github.com/why444216978/gin-api/app/module/test/job/grpc/helloworld"
	client "github.com/why444216978/gin-api/client/grpc"
	server "github.com/why444216978/gin-api/server/grpc"
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
	httpServer := &http.Server{
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		IdleTimeout:  time.Second,
	}
	err := server.New(
		server.WithEndpoint(endpoint),
		server.WithRegisterGRPCFunc(RegisterServer),
		server.WithHTTP(httpServer, pb.RegisterGreeterHandler),
	).Start()
	if err != nil {
		panic(err)
	}
}

func GrpcCmux(ctx context.Context) (err error) {
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

	return
}
