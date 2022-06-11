package grpc

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/reflection"
)

type (
	RegisterGRPC func(s *grpc.Server)
	RegisterMux  func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
)

type Register struct {
	RegisterGRPC RegisterGRPC
	RegisterMux  RegisterMux
}

func NewRegister(registerGRPC RegisterGRPC, registerMux RegisterMux) Register {
	return Register{
		RegisterGRPC: registerGRPC,
		RegisterMux:  registerMux,
	}
}

func RegisterTools(s *grpc.Server) {
	reflection.Register(s)
	service.RegisterChannelzServiceToServer(s)
}
