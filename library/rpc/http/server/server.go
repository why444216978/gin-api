package server

import (
	"context"
	"net/http"
	"time"

	"github.com/why444216978/gin-api/library/rpc"
)

type Server struct {
	*http.Server
	ctx context.Context
}

var _ rpc.RPCServer = (*Server)(nil)

type Option func(s *Server)

func WithReadTimeout(timeout time.Duration) Option {
	return func(s *Server) { s.Server.ReadTimeout = timeout }
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(s *Server) { s.Server.WriteTimeout = timeout }
}

func New(ctx context.Context, addr string, handler http.Handler, opts ...Option) *Server {
	s := &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		ctx: ctx,
	}

	for _, o := range opts {
		o(s)
	}

	return s
}

func (s *Server) Start() (err error) {
	err = s.Server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return
}

func (s *Server) Close() (err error) {
	return s.Server.Shutdown(s.ctx)
}
