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
	"strings"
	"time"

	"github.com/valyala/bytebufferpool"
	loadBalance "github.com/why444216978/load-balance"

	loggingRPC "github.com/why444216978/gin-api/library/logging/rpc"
	"github.com/why444216978/gin-api/library/registry"
	"github.com/why444216978/gin-api/library/rpc/codec"
	timeoutLib "github.com/why444216978/gin-api/library/timeout"
)

type Response struct {
	HTTPCode int
	Response []byte
}

type RPC struct {
	codec         codec.Codec
	logger        *loggingRPC.RPCLogger
	beforePlugins []BeforeRequestPlugin
	afterPlugins  []AfterRequestPlugin
}

type Option func(r *RPC)

func WithCodec(c codec.Codec) Option {
	return func(r *RPC) { r.codec = c }
}

func WithLogger(logger *loggingRPC.RPCLogger) Option {
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

	if r.codec == nil {
		panic("codec is nil")
	}

	return r
}

// Send is send HTTP request
func (r *RPC) Send(ctx context.Context, serviceName, method, uri string, header http.Header, timeout time.Duration, reqData interface{}, respData interface{}) (ret Response, err error) {
	var (
		reqByte []byte
		cost    int64
		node    = &registry.ServiceNode{}
	)

	if header == nil {
		header = http.Header{}
	}

	defer func() {
		if r.logger == nil {
			return
		}
		fields := loggingRPC.RPCLogFields{
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

	var req *http.Request

	client, err := r.getClient(serviceName)
	if err != nil {
		return
	}

	//构建req
	req, err = http.NewRequestWithContext(ctx, method, fmt.Sprintf("http://%s:%d%s", node.Host, node.Port, uri), bytes.NewReader(reqByte))
	if err != nil {
		return
	}

	//超时传递
	remain, err := timeoutLib.CalcRemainTimeout(ctx)
	if err != nil {
		return
	}
	header.Set(timeoutLib.TimeoutKey, strconv.FormatInt(remain, 10))

	//设置请求header
	req.Header = header

	//请求结束前插件
	for _, plugin := range r.beforePlugins {
		_ = plugin.Handle(ctx, req)
	}

	//请求开始时间
	start := time.Now()

	//判断是否cancel
	if err = ctx.Err(); err != nil {
		return
	}

	//发送请求
	resp, err := client.Do(req)

	//请求结束后插件
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

func (r *RPC) getClient(serviceName string) (client *http.Client, err error) {
	node, err := r.pick(serviceName)
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
			pool.AppendCertsFromPEM(node.CaCrt)
			cliCrt, err := tls.X509KeyPair(node.ClientPem, node.ClientKey)
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

	return &http.Client{Transport: tp}, nil
}

func (r *RPC) pick(serviceName string) (*registry.ServiceNode, error) {
	node := &registry.ServiceNode{}

	ser := registry.Services[serviceName]
	if ser == nil {
		return node, errors.New("service is nil")
	}

	_nodes := ser.GetServices()
	l := len(_nodes)
	if l <= 0 {
		return node, errors.New("service node empty")
	}
	//TODO 如果只有一个直接返回

	nodes := make([]loadBalance.Node, l)
	for k, v := range _nodes {
		nodes[k] = loadBalance.Node{
			Node: fmt.Sprintf("%s:%d", v.Host, v.Port),
		}
	}

	//TODO 初始化 services 的时候设置selector，无需重复New
	load, err := loadBalance.New(loadBalance.BalanceType(ser.GetLoadBalance()))
	if err != nil {
		return node, err
	}

	if err := load.InitNodeList(nodes); err != nil {
		return node, err
	}

	target := load.GetNodeAddress()
	arr := strings.Split(target, ":")
	host := arr[0]
	port, _ := strconv.Atoi(arr[1])

	return &registry.ServiceNode{
		Host: host,
		Port: port,
	}, nil
}
