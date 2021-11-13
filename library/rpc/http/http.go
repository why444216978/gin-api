package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	jaeger_http "github.com/why444216978/gin-api/library/jaeger/http"
	"github.com/why444216978/gin-api/library/logging"
	logging_rpc "github.com/why444216978/gin-api/library/logging/rpc"
	"github.com/why444216978/gin-api/library/registry"

	load_balance "github.com/why444216978/load-balance"
)

type Response struct {
	HTTPCode int
	Response string
}

type RPC struct {
	logger *logging_rpc.RPCLogger
}

type Option func(r *RPC)

func WithLogger(logger *logging_rpc.RPCLogger) Option {
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
func (r *RPC) Send(ctx context.Context, serviceName, method, uri string, header http.Header, body io.Reader, timeout time.Duration) (ret *Response, err error) {
	var cost int64

	node := &registry.ServiceNode{}

	var buf []byte
	if body != nil {
		buf, _ = ioutil.ReadAll(body)
		body = ioutil.NopCloser(bytes.NewBuffer(buf))
	}

	ret = &Response{}

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
			Request:     buf,
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
	req.Header = header

	if ctx.Err() != nil {
		return nil, err
	}

	//注入Jaeger
	logID := req.Header.Get(logging.LogHeader)
	jaeger_http.InjectHTTP(ctx, req, logID)

	//发送请求
	start := time.Now()
	resp, err := client.Do(req)
	cost = time.Since(start).Milliseconds()
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
