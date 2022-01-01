package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/why444216978/gin-api/library/registry"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

//EtcdDiscovery 服务发现
type EtcdDiscovery struct {
	ctx         context.Context
	serviceName string
	cli         *clientv3.Client
	nodeList    map[string]*registry.Node //node list
	lock        sync.RWMutex
	decode      registry.Decode
}

var _ registry.Discovery = (*EtcdDiscovery)(nil)

type DiscoverOption func(*EtcdDiscovery)

func WithContext(ctx context.Context) DiscoverOption {
	return func(ed *EtcdDiscovery) { ed.ctx = ctx }
}

func WithServierName(serviceName string) DiscoverOption {
	return func(ed *EtcdDiscovery) { ed.serviceName = serviceName }
}

func WithDiscoverClient(cli *clientv3.Client) DiscoverOption {
	return func(ed *EtcdDiscovery) { ed.cli = cli }
}

// NewDiscovery
func NewDiscovery(opts ...DiscoverOption) (registry.Discovery, error) {
	ed := &EtcdDiscovery{
		nodeList: make(map[string]*registry.Node),
		decode:   JSONDecode,
	}

	for _, o := range opts {
		o(ed)
	}

	if ed.serviceName == "" {
		return nil, errors.New("serviceName is nil")
	}

	if ed.cli == nil {
		return nil, errors.New("cli is nil")
	}

	if ed.ctx == nil {
		ed.ctx = context.Background()
	}

	if err := ed.init(); err != nil {
		return nil, err
	}

	return ed, nil
}

// WatchService
func (s *EtcdDiscovery) init() error {
	//根据前缀获取现有的key
	resp, err := s.cli.Get(s.ctx, s.serviceName, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kv := range resp.Kvs {
		key := string(kv.Key)
		val := string(kv.Value)

		node, err := s.decode(val)
		if err != nil {
			log.Println("service:", s.serviceName, " put key:", key, "val:", val, "err:", err.Error())
			continue
		}
		s.SetServiceList(key, node)
	}

	//监视前缀，修改变更的服务节点
	go s.watcher()

	return nil
}

// watcher
func (s *EtcdDiscovery) watcher() {
	rch := s.cli.Watch(s.ctx, s.serviceName, clientv3.WithPrefix())
	log.Printf("watching prefix:%s now...", s.serviceName)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			key := string(ev.Kv.Key)
			val := string(ev.Kv.Value)

			switch ev.Type {
			case mvccpb.PUT:
				node, err := s.decode(val)
				if err != nil {
					log.Println("service:", s.serviceName, " put key:", key, "val:", val, "err:", err.Error())
					return
				}
				s.SetServiceList(key, node)
				log.Println("service", s.serviceName, " put key:", key, "val:", val)
			case mvccpb.DELETE:
				s.DelServiceList(key)
				log.Println("service:", s.serviceName, " del key:", key)
			}
		}
	}
}

// SetServiceList
func (s *EtcdDiscovery) SetServiceList(key string, node *registry.Node) {
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

// GetNodes
func (s *EtcdDiscovery) GetNodes() []*registry.Node {
	s.lock.RLock()
	defer s.lock.RUnlock()
	nodes := make([]*registry.Node, 0)

	for _, node := range s.nodeList {
		nodes = append(nodes, node)
	}
	return nodes
}

// Close
func (s *EtcdDiscovery) Close() error {
	if s.cli == nil {
		return nil
	}
	return s.cli.Close()
}

func JSONDecode(val string) (*registry.Node, error) {
	node := &registry.Node{}
	err := json.Unmarshal([]byte(val), node)
	if err != nil {
		return nil, errors.New("Unmarshal val " + err.Error())
	}

	return node, nil
}
