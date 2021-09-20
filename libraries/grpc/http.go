package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

type Server struct {
	endpoint     string
	listener     net.Listener
	HTTPListener net.Listener
	GRPCListener net.Listener
	HTTPServer   *http.Server
	GRPCConn     *grpc.ClientConn
	startHTTP    startFunc
	startGRPC    startFunc
	ServerMux    *runtime.ServeMux
	router       *http.ServeMux
}

type startFunc func(ctx context.Context, s *Server)

type Option func(*Server)

func WithHTTPStartFunc(startHTTP startFunc) Option {
	return func(s *Server) { s.startHTTP = startHTTP }
}

func WithGRPCStartFunc(startGRPC startFunc) Option {
	return func(s *Server) { s.startGRPC = startGRPC }
}

// New returns a Server instance
func New(endpoint string, l net.Listener, opts ...Option) *Server {
	s := &Server{
		endpoint: endpoint,
		listener: l,
	}

	for _, o := range opts {
		o(s)
	}

	return s
}

func (s *Server) Start() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tcpMux := cmux.New(s.listener)

	s.GRPCListener = tcpMux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	s.HTTPListener = tcpMux.Match(cmux.HTTP1Fast())

	go s.startGRPC(ctx, s)
	go s.startHTTP(ctx, s)

	return tcpMux.Serve()
}

func (s *Server) InitGateway(ctx context.Context) error {
	var err error

	s.router = http.NewServeMux()

	s.GRPCConn, err = grpc.Dial(s.endpoint, []grpc.DialOption{
		grpc.WithTimeout(10 * time.Second),
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}...)
	if err != nil {
		return fmt.Errorf("Fail to dial: %v", err)
	}

	s.ServerMux = runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:  true,
				UseEnumNumbers: true,
			},
		}),
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
	)

	return nil
}

func (s *Server) StartGateway() {
	s.router.Handle("/", s.ServerMux)

	s.HTTPServer = &http.Server{
		Addr:         s.endpoint,
		Handler:      s.router,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		IdleTimeout:  time.Second,
	}

	if err := s.HTTPServer.Serve(s.HTTPListener); err != nil {
		panic(err)
	}
}
