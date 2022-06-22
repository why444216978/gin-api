package h2c

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"

	"github.com/why444216978/gin-api/library/logger"
	"github.com/why444216978/gin-api/server"
	serverGRPC "github.com/why444216978/gin-api/server/grpc"
)

type Option struct {
	logger logger.Logger
}

type OptionFunc func(*Option)

func WithLogger(l logger.Logger) OptionFunc {
	return func(s *Option) { s.logger = l }
}

type H2CServer struct {
	*Option
	*grpc.Server
	ctx        context.Context
	endpoint   string
	registers  []serverGRPC.Register
	httpServer *http.Server
}

var _ server.Server = (*H2CServer)(nil)

func NewH2C(endpoint string, registers []serverGRPC.Register, opts ...OptionFunc) *H2CServer {
	if len(registers) < 1 {
		panic("len(registers) < 1")
	}

	option := &Option{}
	for _, o := range opts {
		o(option)
	}

	s := &H2CServer{
		Option:    option,
		ctx:       context.Background(),
		endpoint:  endpoint,
		registers: registers,
	}

	return s
}

func (s *H2CServer) Start() (err error) {
	grpcServer := grpc.NewServer(serverGRPC.NewServerOption(serverGRPC.ServerOptionLogger(s.logger))...)

	mux := http.NewServeMux()
	gwmux := runtime.NewServeMux()
	mux.Handle("/", gwmux)

	for _, r := range s.registers {
		if r.RegisterGRPC == nil {
			return errors.New("r.RegisterGRPC nil")
		}

		r.RegisterGRPC(grpcServer)
		if err = r.RegisterMux(s.ctx, gwmux, s.endpoint, serverGRPC.NewDialOption()); err != nil {
			return
		}
	}

	serverGRPC.RegisterTools(grpcServer)

	s.Server = grpcServer

	s.httpServer = &http.Server{
		Addr: s.endpoint,
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				grpcServer.ServeHTTP(w, r)
			} else {
				mux.ServeHTTP(w, r)
			}
		}), &http2.Server{}),
	}

	return s.httpServer.ListenAndServe()
}

func (s *H2CServer) Close() (err error) {
	s.GracefulStop()
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	return s.httpServer.Shutdown(ctx)
}
