package grpc

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	"github.com/why444216978/go-util/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/why444216978/gin-api/server"
)

var kaep = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
	PermitWithoutStream: true,            // Allow pings even when there are no active streams
}

var kasp = keepalive.ServerParameters{
	MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
	MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
	MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
	Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
	Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
}

var _ server.Server = (*Server)(nil)

type Server struct {
	ctx          context.Context
	endpoint     string
	registerHTTP RegisterHTTP
	registerGRPC RegisterGRPC
	httpServer   *http.Server
	tcpMux       cmux.CMux
}

type (
	RegisterGRPC func(s *grpc.Server)
	RegisterHTTP func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
)

type Option func(*Server)

func WithEndpoint(endpoint string) Option {
	return func(s *Server) { s.endpoint = endpoint }
}

func WithRegisterGRPCFunc(registerGRPC RegisterGRPC) Option {
	return func(s *Server) { s.registerGRPC = registerGRPC }
}

func WithHTTP(httpServer *http.Server, registerHTTP RegisterHTTP) Option {
	return func(s *Server) {
		s.httpServer = httpServer
		s.registerHTTP = registerHTTP
	}
}

// New returns a Server
func New(opts ...Option) *Server {
	s := &Server{
		ctx: context.Background(),
	}

	for _, o := range opts {
		o(s)
	}

	if assert.IsNil(s.registerGRPC) {
		panic("registerGRPC is nil")
	}

	return s
}

func (s *Server) Start() (err error) {
	listener, err := net.Listen("tcp", s.endpoint)
	if err != nil {
		return
	}
	s.tcpMux = cmux.New(listener)

	go s.startGRPC()
	go s.startHTTP()

	return s.tcpMux.Serve()
}

func (s *Server) startGRPC() {
	grpcServer := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))
	s.registerGRPC(grpcServer)
	listener := s.tcpMux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func (s *Server) startHTTP() (err error) {
	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	if assert.IsNil(s.registerHTTP) {
		return
	}

	if assert.IsNil(s.httpServer) {
		panic("httpServer is nil")
	}

	grpcConn, err := grpc.DialContext(s.ctx, s.endpoint, []grpc.DialOption{
		grpc.WithTimeout(10 * time.Second),
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}...)
	if err != nil {
		return
	}

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:  true,
				UseEnumNumbers: true,
			},
		}),
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
	)
	if err = s.registerHTTP(s.ctx, mux, grpcConn); err != nil {
		return
	}

	router := http.NewServeMux()
	router.Handle("/", mux)
	s.httpServer.Addr = s.endpoint
	s.httpServer.Handler = router
	listener := s.tcpMux.Match(cmux.HTTP1Fast())
	if err = s.httpServer.Serve(listener); err != nil {
		return
	}

	return
}

func (s *Server) Close() (err error) {
	s.tcpMux.Close()
	return
}
