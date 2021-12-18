package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/bytebufferpool"
	load_balance "github.com/why444216978/load-balance"

	logging_rpc "github.com/why444216978/gin-api/library/logging/rpc"
	"github.com/why444216978/gin-api/library/registry"
	"github.com/why444216978/gin-api/library/rpc/codec"
	timeoutLib "github.com/why444216978/gin-api/library/timeout"
)

type Response struct {
	HTTPCode int
	Resp     []byte
	Response string
}

type RPC struct {
	codec         codec.Codec
	logger        *logging_rpc.RPCLogger
	beforePlugins []BeforeRequestPlugin
	afterPlugins  []AfterRequestPlugin
}

type Option func(r *RPC)

func WithCodec(c codec.Codec) Option {
	return func(r *RPC) { r.codec = c }
}

func WithLogger(logger *logging_rpc.RPCLogger) Option {
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
func (r *RPC) Send(ctx context.Context, serviceName, method, uri string, header http.Header, reqData interface{}, timeout time.Duration) (ret Response, err error) {
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
		fields := logging_rpc.RPCLogFields{
			ServiceName: serviceName,
			Header:      header,
			Method:      method,
			URI:         uri,
			Request:     reqByte,
			Response:    ret.Response,
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

	node, err = r.loadBalance(serviceName)
	if err != nil {
		return
	}

	var req *http.Request

	client := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   timeout,
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

	ret.Resp = b.Bytes()
	ret.Response = string(b.Bytes())

	return
}

func (r *RPC) loadBalance(serviceName string) (*registry.ServiceNode, error) {
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

	nodes := make([]load_balance.Node, l)
	for k, v := range _nodes {
		nodes[k] = load_balance.Node{
			Node: fmt.Sprintf("%s:%d", v.Host, v.Port),
		}
	}

	load, err := load_balance.New(load_balance.BalanceType(ser.GetLoadBalance()))
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
