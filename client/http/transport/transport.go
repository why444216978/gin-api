package transport

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/why444216978/go-util/assert"

	client "github.com/why444216978/gin-api/client/http"
	"github.com/why444216978/gin-api/library/logger"
	loggerRPC "github.com/why444216978/gin-api/library/logger/zap/rpc"
	"github.com/why444216978/gin-api/library/servicer"
	timeoutLib "github.com/why444216978/gin-api/server/http/middleware/timeout"
)

type RPC struct {
	logger        logger.Logger
	beforePlugins []client.BeforeRequestPlugin
	afterPlugins  []client.AfterRequestPlugin
}

type Option func(r *RPC)

func WithLogger(logger *loggerRPC.RPCLogger) Option {
	return func(r *RPC) { r.logger = logger }
}

func WithBeforePlugins(plugins ...client.BeforeRequestPlugin) Option {
	return func(r *RPC) { r.beforePlugins = plugins }
}

func WithAfterPlugins(plugins ...client.AfterRequestPlugin) Option {
	return func(r *RPC) { r.afterPlugins = plugins }
}

func New(opts ...Option) *RPC {
	r := &RPC{}
	for _, o := range opts {
		o(r)
	}

	return r
}

// Send is send HTTP request
func (r *RPC) Send(ctx context.Context, serviceName string, request client.Request, response *client.Response) (err error) {
	var (
		cost int64
		node = &servicer.Node{}
	)

	if response == nil {
		return errors.New("response is nil")
	}

	if assert.IsNil(request.Codec) {
		return errors.New("request.Codec is nil")
	}

	if assert.IsNil(response.Codec) {
		return errors.New("request.Codec is nil")
	}

	if request.Header == nil {
		request.Header = http.Header{}
	}

	defer func() {
		if r.logger == nil {
			return
		}
		fields := []logger.Field{
			logger.Reflect(logger.ServiceName, serviceName),
			logger.Reflect(logger.Header, request.Header),
			logger.Reflect(logger.Method, request.Method),
			logger.Reflect(logger.API, request.URI),
			logger.Reflect(logger.Request, request.Body),
			logger.Reflect(logger.Response, response.Body),
			logger.Reflect(logger.ServerIP, node.Host),
			logger.Reflect(logger.ServerPort, node.Port),
			logger.Reflect(logger.Code, response.HTTPCode),
			logger.Reflect(logger.Cost, cost),
			logger.Reflect(logger.Timeout, request.Timeout),
		}
		if err == nil {
			r.logger.Info(ctx, "rpc success", fields...)
			return
		}
		r.logger.Error(ctx, err.Error(), fields...)
	}()

	reqReader, err := request.Codec.Encode(request.Body)
	if err != nil {
		return
	}

	var (
		client *http.Client
		req    *http.Request
	)

	service, ok := servicer.GetServicer(serviceName)
	if !ok {
		err = errors.New("service is nil")
		return
	}

	client, node, err = r.getClient(ctx, serviceName, service)
	if err != nil {
		return
	}

	// 构建req
	url := fmt.Sprintf("http://%s:%d%s", node.Host, node.Port, request.URI)
	req, err = http.NewRequestWithContext(ctx, request.Method, url, reqReader)
	if err != nil {
		return
	}

	// 超时传递
	remain, err := timeoutLib.CalcRemainTimeout(ctx)
	if err != nil {
		return
	}
	request.Header.Set(timeoutLib.TimeoutKey, strconv.FormatInt(remain, 10))

	// 设置请求header
	req.Header = request.Header

	// 请求结束前插件
	for _, plugin := range r.beforePlugins {
		_ = plugin.Handle(ctx, req)
	}

	// 请求开始时间
	start := time.Now()

	// 判断是否cancel
	if err = ctx.Err(); err != nil {
		return
	}

	// 发送请求
	resp, err := client.Do(req)

	_ = service.Done(ctx, node, err)

	// 请求结束后插件
	for _, plugin := range r.afterPlugins {
		_ = plugin.Handle(ctx, req, resp)
	}

	cost = time.Since(start).Milliseconds()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	response.HTTPCode = resp.StatusCode
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http code is %d", resp.StatusCode)
		return
	}

	err = response.Codec.Decode(resp.Body, response.Body)

	return
}

func (r *RPC) getClient(ctx context.Context, serviceName string, service servicer.Servicer) (client *http.Client, node *servicer.Node, err error) {
	node, err = service.Pick(ctx)
	if err != nil {
		return
	}

	address := fmt.Sprintf("%s:%d", node.Host, node.Port)

	tp := &http.Transport{
		MaxIdleConnsPerHost: 30,
		MaxConnsPerHost:     30,
		IdleConnTimeout:     time.Minute,
		DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
			conn, err := (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext(context.TODO(), "tcp", address)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
		DialTLSContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
			pool := x509.NewCertPool()
			pool.AppendCertsFromPEM(service.GetCaCrt())
			cliCrt, err := tls.X509KeyPair(service.GetClientPem(), service.GetClientKey())
			if err != nil {
				err = errors.New("server pem error " + err.Error())
				return nil, err
			}

			conn, err := (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext(context.TODO(), "tcp", address)
			if err != nil {
				return nil, err
			}

			return tls.Client(conn, &tls.Config{
				RootCAs:      pool,
				Certificates: []tls.Certificate{cliCrt},
				ServerName:   serviceName,
			}), err
		},
	}
	client = &http.Client{Transport: tp}

	return
}
