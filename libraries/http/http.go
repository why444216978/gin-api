package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gin-api/libraries/jaeger"
	"gin-api/libraries/logging"
	"gin-api/libraries/registry"

	load_balance "github.com/why444216978/load-balance"
)

type Response struct {
	HTTPCode int
	Response string
}

type RPC struct {
	logger *logging.RPCLogger
}

type Option func(r *RPC)

func WithLogger(logger *logging.RPCLogger) Option {
	return func(r *RPC) { r.logger = logger }
}

func New(opts ...Option) *RPC {
	r := &RPC{}
	for _, o := range opts {
		o(r)
	}

	return r
}

// Send is send HTTP request
func (r *RPC) Send(ctx context.Context, serviceName, method, uri string, header map[string]string, body io.Reader, timeout time.Duration) (ret *Response, err error) {
	node := &registry.ServiceNode{}
	_req := body
	ret = &Response{}

	defer func() {
		if r.logger == nil {
			return
		}
		fields := r.logger.Fields(ctx, serviceName, method, uri, header, _req, timeout, node, ret.Response, err)
		if err == nil {
			r.logger.Info("http rpc", fields...)
			return
		}
		r.logger.Error("http rpc errpr", fields...)
	}()

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
	req, err = http.NewRequestWithContext(ctx, method, fmt.Sprintf("http://%s:%d%s", node.Host, node.Port, uri), body)
	if err != nil {
		return
	}

	//设置请求header
	for k, v := range header {
		req.Header.Add(k, v)
	}

	if ctx.Err() != nil {
		return nil, err
	}

	//注入Jaeger
	logID := req.Header.Get(logging.LogHeader)
	jaeger.InjectHTTP(ctx, req, logID)

	//发送请求
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	ret.HTTPCode = resp.StatusCode
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http code is %d", resp.StatusCode)
		return
	}

	if b != nil {
		ret.Response = string(b)
	}

	return
}

func (r *RPC) loadBalance(serviceName string) (*registry.ServiceNode, error) {
	ser := registry.Services[serviceName]
	if ser == nil {
		return nil, errors.New("service is nil")
	}

	_nodes := ser.GetServices()
	l := len(_nodes)
	if l <= 0 {
		return nil, errors.New("service node empty")
	}

	nodes := make([]load_balance.Node, l)
	for k, v := range _nodes {
		nodes[k] = load_balance.Node{
			Node: fmt.Sprintf("%s:%d", v.Host, v.Port),
		}
	}

	load, err := load_balance.New(load_balance.BalanceType(ser.GetLoadBalance()))
	if err != nil {
		return nil, err
	}

	if err := load.InitNodeList(nodes); err != nil {
		return nil, err
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

// Send is send HTTP request
func Send(ctx context.Context, serviceName, method, uri string, header map[string]string, body io.Reader, timeout time.Duration) (ret *Response, err error) {
	node, err := loadBalance(serviceName)
	if err != nil {
		return
	}

	var req *http.Request

	client := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   timeout,
	}

	//构建req
	req, err = http.NewRequestWithContext(ctx, method, fmt.Sprintf("http://%s:%d%s", node.Host, node.Port, uri), body)
	if err != nil {
		return
	}

	//设置请求header
	for k, v := range header {
		req.Header.Add(k, v)
	}

	if ctx.Err() != nil {
		return nil, err
	}

	//注入Jaeger
	logID := req.Header.Get(logging.LogHeader)
	jaeger.InjectHTTP(ctx, req, logID)

	//发送请求
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	ret = &Response{
		HTTPCode: resp.StatusCode,
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http code is %d", resp.StatusCode)
		return
	}

	if b != nil {
		ret.Response = string(b)
	}

	return
}

func loadBalance(serviceName string) (*registry.ServiceNode, error) {
	ser := registry.Services[serviceName]
	if ser == nil {
		return nil, errors.New("service is nil")
	}

	_nodes := ser.GetServices()
	l := len(_nodes)
	if l <= 0 {
		return nil, errors.New("service node empty")
	}

	nodes := make([]load_balance.Node, l)
	for k, v := range _nodes {
		nodes[k] = load_balance.Node{
			Node: fmt.Sprintf("%s:%d", v.Host, v.Port),
		}
	}

	load, err := load_balance.New(load_balance.BalanceType(ser.GetLoadBalance()))
	if err != nil {
		return nil, err
	}

	if err := load.InitNodeList(nodes); err != nil {
		return nil, err
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
