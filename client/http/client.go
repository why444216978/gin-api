package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/valyala/bytebufferpool"
	"github.com/why444216978/go-util/assert"

	"github.com/why444216978/gin-api/client/codec"
	loggerRPC "github.com/why444216978/gin-api/library/logger/rpc"
	"github.com/why444216978/gin-api/library/servicer"
	timeoutLib "github.com/why444216978/gin-api/server/http/middleware/timeout"
)

type Response struct {
	HTTPCode int
	Response []byte
}

type RPC struct {
	codec         codec.Codec
	logger        *loggerRPC.RPCLogger
	beforePlugins []BeforeRequestPlugin
	afterPlugins  []AfterRequestPlugin
}

type Option func(r *RPC)

func WithCodec(c codec.Codec) Option {
	return func(r *RPC) { r.codec = c }
}

func WithLogger(logger *loggerRPC.RPCLogger) Option {
	return func(r *RPC) { r.logger = logger }
}

func WithBeforePlugins(plugins ...BeforeRequestPlugin) Option {
	return func(r *RPC) { r.beforePlugins = plugins }
}

func WithAfterPlugins(plugins ...AfterRequestPlugin) Option {
	return func(r *RPC) { r.afterPlugins = plugins }
}

func New(opts ...Option) *RPC {
	r := &RPC{}
	for _, o := range opts {
		o(r)
	}

	if assert.IsNil(r.codec) {
		panic("codec is nil")
	}

	return r
}

// Send is send HTTP request
func (r *RPC) Send(ctx context.Context, serviceName, method, uri string, header http.Header, timeout time.Duration, reqData interface{}, respData interface{}) (ret Response, err error) {
	var (
		reqByte []byte
		cost    int64
		node    = &servicer.Node{}
	)

	if header == nil {
		header = http.Header{}
	}

	defer func() {
		if r.logger == nil {
			return
		}
		fields := loggerRPC.RPCLogFields{
			ServiceName: serviceName,
			Header:      header,
			Method:      method,
			URI:         uri,
			Request:     reqData,
			Response:    respData,
			ServerIP:    node.Host,
			ServerPort:  node.Port,
			HTTPCode:    ret.HTTPCode,
			Cost:        cost,
			Timeout:     timeout,
		}
		if err == nil {
			r.logger.Info(ctx, "rpc success", fields)
			return
		}
		r.logger.Error(ctx, err.Error(), fields)
	}()

	reqByte, err = r.codec.Encode(reqData)
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
	req, err = http.NewRequestWithContext(ctx, method, fmt.Sprintf("http://%s:%d%s", node.Host, node.Port, uri), bytes.NewReader(reqByte))
	if err != nil {
		return
	}

	// 超时传递
	remain, err := timeoutLib.CalcRemainTimeout(ctx)
	if err != nil {
		return
	}
	header.Set(timeoutLib.TimeoutKey, strconv.FormatInt(remain, 10))

	// 设置请求header
	req.Header = header

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

	ret.HTTPCode = resp.StatusCode
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http code is %d", resp.StatusCode)
		return
	}

	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return
	}

	l := b.Len()
	if l < 0 {
		return
	}

	ret.Response = b.Bytes()

	err = r.codec.Decode(ret.Response, respData)

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
