package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	opentracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"github.com/why444216978/go-util/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	"github.com/why444216978/gin-api/library/jaeger"
	"github.com/why444216978/gin-api/library/logger"
	"github.com/why444216978/gin-api/server/grpc/middleware/log"
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

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

type DialOption struct{}

type DialOptionFunc func(*DialOption)

func NewDialOption(opts ...DialOptionFunc) []grpc.DialOption {
	return []grpc.DialOption{
		// TODO
		// grpc.WithResolvers(resolver),
		grpc.WithTimeout(10 * time.Second),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(kacp),
		grpc.WithChainUnaryInterceptor(
			log.UnaryClientInterceptor(),
			otgrpc.OpenTracingClientInterceptor(
				opentracing.GlobalTracer(),
				otgrpc.SpanDecorator(func(span opentracing.Span, method string, req, resp interface{}, err error) {
					if assert.IsNil(span) {
						return
					}

					bs, _ := json.Marshal(req)
					jaeger.SetRequest(span, string(bs))

					if err != nil {
						span.LogFields(opentracingLog.Error(err))
					}
				}),
			),
		),
	}
}

type ServerOption struct {
	logger logger.Logger
}

type ServerOptionFunc func(*ServerOption)

func ServerOptionLogger(l logger.Logger) ServerOptionFunc {
	return func(o *ServerOption) { o.logger = l }
}

func NewServerOption(opts ...ServerOptionFunc) []grpc.ServerOption {
	opt := &ServerOption{}
	for _, o := range opts {
		o(opt)
	}

	interceptors := []grpc.UnaryServerInterceptor{
		otgrpc.OpenTracingServerInterceptor(
			opentracing.GlobalTracer(),
			otgrpc.SpanDecorator(func(span opentracing.Span, method string, req, resp interface{}, err error) {
				if assert.IsNil(span) {
					return
				}

				bs, _ := json.Marshal(resp)
				jaeger.SetResponse(span, string(bs))

				if err != nil {
					span.LogFields(opentracingLog.Error(err))
				}
			})),
		grpc_recovery.UnaryServerInterceptor(
			grpc_recovery.WithRecoveryHandlerContext(func(ctx context.Context, p interface{}) (err error) {
				err = errors.WithStack(fmt.Errorf("%v", p))
				return status.Errorf(codes.Internal, "%+v", err)
			})),
	}
	if !assert.IsNil(opt.logger) {
		interceptors = append(interceptors, log.UnaryServerInterceptor(opt.logger))
	}

	return []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
		grpc.ChainUnaryInterceptor(interceptors...),
	}
}

type CallOption struct{}

type CallOptionFunc func(*CallOption)

func NewCallOption(opts ...CallOption) []grpc.CallOption {
	return []grpc.CallOption{}
}
