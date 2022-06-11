package cmux

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	"github.com/why444216978/go-util/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/why444216978/gin-api/server"
	serverGRPC "github.com/why444216978/gin-api/server/grpc"
)

type (
	RegisterHTTP func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
)

type Option struct {
	registerHTTP RegisterHTTP
	httpServer   *http.Server
}

type OptionFunc func(*Option)

func WithHTTP(httpServer *http.Server, registerHTTP RegisterHTTP) OptionFunc {
	return func(s *Option) {
		s.httpServer = httpServer
		s.registerHTTP = registerHTTP
	}
}

type CMUXServer struct {
	*Option
	ctx       context.Context
	endpoint  string
	registers []serverGRPC.Register
	tcpMux    cmux.CMux
}

var _ server.Server = (*CMUXServer)(nil)

func NewCMUX(endpoint string, registers []serverGRPC.Register, opts ...OptionFunc) *CMUXServer {
	if len(registers) < 1 {
		panic("len(registers) < 1")
	}

	option := &Option{}
	for _, o := range opts {
		o(option)
	}

	s := &CMUXServer{
		Option:    option,
		ctx:       context.Background(),
		registers: registers,
		endpoint:  endpoint,
	}

	return s
}

func (s *CMUXServer) Start() (err error) {
	listener, err := net.Listen("tcp", s.endpoint)
	if err != nil {
		return
	}
	s.tcpMux = cmux.New(listener)

	go s.startGRPC()
	go s.startHTTP()

	return s.tcpMux.Serve()
}

func (s *CMUXServer) startGRPC() {
	grpcServer := grpc.NewServer(serverGRPC.NewServerOption()...)

	for _, r := range s.registers {
		if r.RegisterGRPC == nil {
			panic("r.RegisterGRPC == nil")
		}
		r.RegisterGRPC(grpcServer)
	}

	serverGRPC.RegisterTools(grpcServer)

	listener := s.tcpMux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func (s *CMUXServer) startHTTP() (err error) {
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

	grpcConn, err := grpc.DialContext(s.ctx, s.endpoint, serverGRPC.NewDialOption()...)
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

func (s *CMUXServer) Close() (err error) {
	s.tcpMux.Close()
	return
}
