package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"gin-api/libraries/registry"
	"log"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

//EtcdDiscovery 服务发现
type EtcdDiscovery struct {
	serviceName string
	cli         *clientv3.Client
	nodeList    map[string]string //node list
	lock        sync.Mutex
	endpoints   []string
	dialTimeout time.Duration
	decode      registry.Decode
}

var _ registry.Discovery = (*EtcdDiscovery)(nil)

type DiscoverOption func(*EtcdDiscovery)

func WithDiscoverServiceName(serviceName string) DiscoverOption {
	return func(ed *EtcdDiscovery) { ed.serviceName = serviceName }
}

func WithDiscoverEndpoints(endpoints []string) DiscoverOption {
	return func(ed *EtcdDiscovery) { ed.endpoints = endpoints }
}

func WithDiscoverDialTimeout(duration time.Duration) DiscoverOption {
	return func(ed *EtcdDiscovery) { ed.dialTimeout = duration * time.Second }
}

// NewDiscovery
func NewDiscovery(opts ...DiscoverOption) (*EtcdDiscovery, error) {
	var err error
	ed := &EtcdDiscovery{
		nodeList: make(map[string]string),
		decode:   JSONDecode,
	}

	for _, o := range opts {
		o(ed)
	}

	ed.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   ed.endpoints,
		DialTimeout: ed.dialTimeout,
	})
	if err != nil {
		return nil, err
	}

	return ed, nil
}

// WatchService
func (s *EtcdDiscovery) WatchService(ctx context.Context) error {
	//根据前缀获取现有的key
	resp, err := s.cli.Get(ctx, s.serviceName, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, ev := range resp.Kvs {
		s.SetServiceList(string(ev.Key), string(ev.Value))
	}

	//监视前缀，修改变更的服务节点
	go s.watcher(ctx)

	return nil
}

// watcher
func (s *EtcdDiscovery) watcher(ctx context.Context) {
	rch := s.cli.Watch(ctx, s.serviceName, clientv3.WithPrefix())
	log.Printf("watching prefix:%s now...", s.serviceName)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT: //修改或者新增
				s.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE: //删除
				s.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

// SetServiceList
func (s *EtcdDiscovery) SetServiceList(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.nodeList[key] = string(val)
	log.Println("put key :", key, "val:", val)
}

// DelServiceList
func (s *EtcdDiscovery) DelServiceList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.nodeList, key)
	log.Println("del key:", key)
}

// GetServices
func (s *EtcdDiscovery) GetServices() []*registry.ServiceNode {
	s.lock.Lock()
	defer s.lock.Unlock()
	addrs := make([]*registry.ServiceNode, 0)

	for _, v := range s.nodeList {
		addr, err := s.decode(v)
		if err != nil {
			continue
		}
		addrs = append(addrs, addr)
	}
	return addrs
}

// Close
func (s *EtcdDiscovery) Close() error {
	if s.cli == nil {
		return nil
	}
	return s.cli.Close()
}

func JSONDecode(val interface{}) (*registry.ServiceNode, error) {
	var (
		ok        bool
		stringVal string
		node      = &registry.ServiceNode{}
	)

	if stringVal, ok = val.(string); !ok {
		return nil, errors.New("val is not string")
	}

	err := json.Unmarshal([]byte(stringVal), node)
	if err != nil {
		return nil, errors.New("Unmarshal val " + err.Error())
	}

	return node, nil
}
