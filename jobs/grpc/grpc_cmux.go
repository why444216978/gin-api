package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "gin-api/jobs/grpc/helloworld"

	server "gin-api/libraries/grpc"

	"google.golang.org/grpc"
)

const (
	endpoint = ":8888"
)

type Server struct {
	*server.Server
	pb.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: in.Name + " world"}, nil
}

func startHTTP(ctx context.Context, s *server.Server) {
	if err := s.InitGateway(ctx); err != nil {
		panic(err)
	}
	if err := pb.RegisterGreeterHandler(ctx, s.ServerMux, s.GRPCConn); err != nil {
		panic(err)
	}
	s.StartGateway()
}

func startGRPC(ctx context.Context, s *server.Server) {
	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, new(Server))
	if err := grpcServer.Serve(s.GRPCListener); err != nil {
		panic(err)
	}
}

func client() {
	cc, err := newClientConn(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewGreeterClient(cc)

	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "why"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}

func newClientConn(target string) (*grpc.ClientConn, error) {
	cc, err := grpc.Dial(
		target,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	return cc, nil
}

func GrpcCmux() (err error) {
	conn, err := net.Listen("tcp", endpoint)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	s := server.New(endpoint, conn,
		server.WithGRPCStartFunc(startGRPC),
		server.WithHTTPStartFunc(startHTTP))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		panic(s.Start())
		wg.Done()
	}()

	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		client()
	}

	wg.Wait()

	return
}
