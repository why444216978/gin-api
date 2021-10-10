package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"sync"

	"github.com/why444216978/gin-api/libraries/registry"

	"github.com/why444216978/go-util/validate"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

//EtcdDiscovery 服务发现
type EtcdDiscovery struct {
	config   *registry.DiscoveryConfig
	cli      *clientv3.Client
	nodeList map[string]*registry.ServiceNode //node list
	lock     sync.Mutex
	decode   registry.Decode
}

var _ registry.Discovery = (*EtcdDiscovery)(nil)

type DiscoverOption func(*EtcdDiscovery)

func WithDiscoverConfig(config *registry.DiscoveryConfig) DiscoverOption {
	if err := validate.ValidateCamel(config); err != nil {
		panic(err)
	}
	return func(ed *EtcdDiscovery) { ed.config = config }
}

func WithDiscoverClient(cli *clientv3.Client) DiscoverOption {
	return func(er *EtcdDiscovery) { er.cli = cli }
}

// NewDiscovery
func NewDiscovery(opts ...DiscoverOption) (*EtcdDiscovery, error) {
	ed := &EtcdDiscovery{
		nodeList: make(map[string]*registry.ServiceNode),
		decode:   JSONDecode,
	}

	for _, o := range opts {
		o(ed)
	}

	return ed, nil
}

// WatchService
func (s *EtcdDiscovery) WatchService(ctx context.Context) error {
	if s.config.Type == registry.TypeHostPort {
		s.SetServiceList(s.config.ServiceName, &registry.ServiceNode{
			Host: s.config.Host,
			Port: s.config.Port,
		})
		return nil
	}

	if s.config.Type == registry.TypeHostDomain {
		host, err := net.ResolveIPAddr("ip", "localhost")
		if err != nil {
			panic(err)
		}
		s.SetServiceList(s.config.ServiceName, &registry.ServiceNode{
			Host: host.IP.String(),
			Port: 80,
		})
	}

	if s.cli == nil {
		return errors.New("cli is nil")
	}

	//根据前缀获取现有的key
	resp, err := s.cli.Get(ctx, s.config.ServiceName, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kv := range resp.Kvs {
		key := string(kv.Key)
		val := string(kv.Value)

		node, err := s.decode(val)
		if err != nil {
			log.Println("service:", s.config.ServiceName, " put key:", key, "val:", val, "err:", err.Error())
			continue
		}
		s.SetServiceList(key, node)
	}

	//监视前缀，修改变更的服务节点
	go s.watcher(ctx)

	return nil
}

// watcher
func (s *EtcdDiscovery) watcher(ctx context.Context) {
	rch := s.cli.Watch(ctx, s.config.ServiceName, clientv3.WithPrefix())
	log.Printf("watching prefix:%s now...", s.config.ServiceName)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			key := string(ev.Kv.Key)
			val := string(ev.Kv.Value)

			switch ev.Type {
			case mvccpb.PUT:
				node, err := s.decode(val)
				if err != nil {
					log.Println("service:", s.config.ServiceName, " put key:", key, "val:", val, "err:", err.Error())
					return
				}
				s.SetServiceList(key, node)
				log.Println("service", s.config.ServiceName, " put key:", key, "val:", val)
			case mvccpb.DELETE:
				s.DelServiceList(key)
				log.Println("service:", s.config.ServiceName, " del key:", key)
			}
		}
	}
}

// SetServiceList
func (s *EtcdDiscovery) SetServiceList(key string, node *registry.ServiceNode) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.nodeList[key] = node
}

// DelServiceList
func (s *EtcdDiscovery) DelServiceList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.nodeList, key)
}

// GetServices
func (s *EtcdDiscovery) GetServices() []*registry.ServiceNode {
	s.lock.Lock()
	defer s.lock.Unlock()
	nodes := make([]*registry.ServiceNode, 0)

	for _, node := range s.nodeList {
		nodes = append(nodes, node)
	}
	return nodes
}

// GetLoadBalance
func (s *EtcdDiscovery) GetLoadBalance() string {
	return s.config.LoadBalance
}

// Close
func (s *EtcdDiscovery) Close() error {
	if s.cli == nil {
		return nil
	}
	return s.cli.Close()
}

func JSONDecode(val string) (*registry.ServiceNode, error) {
	node := &registry.ServiceNode{}
	err := json.Unmarshal([]byte(val), node)
	if err != nil {
		return nil, errors.New("Unmarshal val " + err.Error())
	}

	return node, nil
}
