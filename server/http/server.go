package http

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/why444216978/gin-api/app/response"
	"github.com/why444216978/gin-api/server"
)

type Server struct {
	*http.Server
	ctx                context.Context
	middlewares        []gin.HandlerFunc
	registerRouterFunc RegisterRouter
	pprofTurn          bool
}

var _ server.RPCServer = (*Server)(nil)

type RegisterRouter func(server *gin.Engine)

type Option func(s *Server)

func WithReadTimeout(timeout time.Duration) Option {
	return func(s *Server) { s.Server.ReadTimeout = timeout }
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(s *Server) { s.Server.WriteTimeout = timeout }
}

func WithMiddlewares(middlewares ...gin.HandlerFunc) Option {
	return func(s *Server) { s.middlewares = middlewares }
}

func WithRegisterRouter(f RegisterRouter) Option {
	return func(s *Server) { s.registerRouterFunc = f }
}

func WithPprof(pprofTurn bool) Option {
	return func(s *Server) { s.pprofTurn = pprofTurn }
}

func New(ctx context.Context, addr string, opts ...Option) *Server {
	s := &Server{
		Server: &http.Server{
			Addr: addr,
		},
		ctx: ctx,
	}

	for _, o := range opts {
		o(s)
	}

	s.Handler = s.initHandler()

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

func (s *Server) initHandler() *gin.Engine {
	server := gin.New()

	s.startPprof(server)

	server.Use(s.middlewares...)

	if s.registerRouterFunc != nil {
		s.registerRouterFunc(server)
	}

	server.NoRoute(func(c *gin.Context) {
		response.Response(c, response.CodeUriNotFound, nil, "")
		c.AbortWithStatus(http.StatusNotFound)
	})

	return server
}

func (s *Server) startPprof(server *gin.Engine) {
	if !s.pprofTurn {
		return
	}

	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	pprof.Register(server)
}
