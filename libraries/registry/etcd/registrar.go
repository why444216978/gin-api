package etcd

import (
	"context"
	"time"

	"gin-api/libraries/registry"

	"github.com/coreos/etcd/clientv3"
)

// EtcdRegistrar
type EtcdRegistrar struct {
	serviceName   string
	addr          string
	cli           *clientv3.Client
	leaseID       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
	val           string
	lease         int64
	endpoints     []string
	dialTimeout   time.Duration
}

var _ registry.Registrar = (*EtcdRegistrar)(nil)

type RegistrarOption func(*EtcdRegistrar)

func WithRegistrarServiceName(serviceName string) RegistrarOption {
	return func(er *EtcdRegistrar) { er.serviceName = serviceName }
}

func WithRegistarAddr(addr string) RegistrarOption {
	return func(er *EtcdRegistrar) { er.addr = addr }
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
	r := &EtcdRegistrar{}

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

	r.key = "/" + r.serviceName + "/" + r.addr
	r.val = r.addr

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
