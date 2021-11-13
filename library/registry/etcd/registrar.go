package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/why444216978/gin-api/library/registry"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdRegistrar
type EtcdRegistrar struct {
	serviceName   string
	host          string
	port          int
	cli           *clientv3.Client
	leaseID       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
	val           string
	lease         int64
	encode        registry.Encode
}

var _ registry.Registrar = (*EtcdRegistrar)(nil)

type RegistrarOption func(*EtcdRegistrar)

func WithRegistrarClient(cli *clientv3.Client) RegistrarOption {
	return func(er *EtcdRegistrar) { er.cli = cli }
}

func WithRegistrarServiceName(serviceName string) RegistrarOption {
	return func(er *EtcdRegistrar) { er.serviceName = serviceName }
}

func WithRegistarHost(host string) RegistrarOption {
	return func(er *EtcdRegistrar) { er.host = host }
}

func WithRegistarPort(port int) RegistrarOption {
	return func(er *EtcdRegistrar) { er.port = port }
}

func WithRegistrarLease(lease int64) RegistrarOption {
	return func(er *EtcdRegistrar) { er.lease = lease }
}

// NewRegistry
func NewRegistry(opts ...RegistrarOption) (*EtcdRegistrar, error) {
	var err error

	r := &EtcdRegistrar{
		encode: JSONEncode,
	}

	for _, o := range opts {
		o(r)
	}

	r.key = fmt.Sprintf("%s.%s.%d", r.serviceName, r.host, r.port)

	if r.val, err = r.encode(&registry.ServiceNode{
		Host: r.host,
		Port: r.port,
	}); err != nil {
		return nil, err
	}

	return r, nil
}

func (s *EtcdRegistrar) Register(ctx context.Context) error {
	if s.cli == nil {
		return errors.New("cli is nil")
	}

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
	log.Println("deregister success")
	return s.cli.Close()
}

func JSONEncode(node *registry.ServiceNode) (string, error) {
	val, err := json.Marshal(node)
	if err != nil {
		return "", errors.New("marshal node " + err.Error())
	}

	return string(val), nil
}
