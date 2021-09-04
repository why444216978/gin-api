package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gin-api/libraries/registry"

	"github.com/coreos/etcd/clientv3"
)

// EtcdRegistrar
type EtcdRegistrar struct {
	serviceName   string
	host          string
	port          int
	addr          string
	endpoints     []string
	cli           *clientv3.Client
	dialTimeout   time.Duration
	leaseID       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
	val           string
	lease         int64
	encode        registry.Encode
}

var _ registry.Registrar = (*EtcdRegistrar)(nil)

type RegistrarOption func(*EtcdRegistrar)

func WithRegistrarServiceName(serviceName string) RegistrarOption {
	return func(er *EtcdRegistrar) { er.serviceName = serviceName }
}

func WithRegistarHost(host string) RegistrarOption {
	return func(er *EtcdRegistrar) { er.host = host }
}

func WithRegistarPort(port int) RegistrarOption {
	return func(er *EtcdRegistrar) { er.port = port }
}

func WithRegistrarEndpoints(endpoints []string) RegistrarOption {
	return func(er *EtcdRegistrar) { er.endpoints = endpoints }
}

func WithRegistrarLease(lease int64) RegistrarOption {
	return func(er *EtcdRegistrar) { er.lease = lease }
}

// NewRegistry
func NewRegistry(opts ...RegistrarOption) (registry.Registrar, error) {
	var err error
	r := &EtcdRegistrar{
		encode: JSONEncode,
	}

	for _, o := range opts {
		o(r)
	}

	r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.endpoints,
		DialTimeout: r.dialTimeout,
	})
	if err != nil {
		return nil, err
	}

	r.key = r.serviceName + "." + r.addr

	val, err := r.encode(&registry.ServiceNode{
		Host: r.host,
		Port: r.port,
	})
	if err != nil {
		return nil, errors.New("marshal node " + err.Error())
	}
	var ok bool
	if r.val, ok = val.(string); !ok {
		return nil, errors.New("assert val fail")
	}

	return r, nil
}

func (s *EtcdRegistrar) Register(ctx context.Context) error {
	//申请租约设置时间keepalive
	if err := s.putKeyWithRegistrarLease(ctx, s.lease); err != nil {
		return err
	}

	//监听续租相应chan
	go s.listenLeaseRespChan()

	return nil
}

// putKeyWithRegistrarLease
func (s *EtcdRegistrar) putKeyWithRegistrarLease(ctx context.Context, lease int64) error {
	//设置租约时间
	resp, err := s.cli.Grant(ctx, lease)
	if err != nil {
		return err
	}
	//注册服务并绑定租约
	_, err = s.cli.Put(ctx, s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	//设置续租 定期发送需求请求
	leaseRespChan, err := s.cli.KeepAlive(context.Background(), resp.ID)

	if err != nil {
		return err
	}
	s.leaseID = resp.ID
	s.keepAliveChan = leaseRespChan
	return nil
}

// listenLeaseRespChan
func (s *EtcdRegistrar) listenLeaseRespChan() {
	for leaseKeepResp := range s.keepAliveChan {
		_ = leaseKeepResp
		// log.Println("续租：", leaseKeepResp)
	}
}

// Close
func (s *EtcdRegistrar) DeRegister(ctx context.Context) error {
	//撤销租约
	if _, err := s.cli.Revoke(ctx, s.leaseID); err != nil {
		return err
	}
	// log.Println("续租结束")
	return s.cli.Close()
}

func JSONEncode(node *registry.ServiceNode) (interface{}, error) {
	val, err := json.Marshal(node)
	if err != nil {
		return nil, errors.New("marshal node " + err.Error())
	}

	return string(val), nil
}
