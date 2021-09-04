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

func WithDiscoverClient(cli *clientv3.Client) DiscoverOption {
	return func(er *EtcdDiscovery) { er.cli = cli }
}

func WithDiscoverServiceName(serviceName string) DiscoverOption {
	return func(ed *EtcdDiscovery) { ed.serviceName = serviceName }
}

// NewDiscovery
func NewDiscovery(opts ...DiscoverOption) (*EtcdDiscovery, error) {
	ed := &EtcdDiscovery{
		nodeList: make(map[string]string),
		decode:   JSONDecode,
	}

	for _, o := range opts {
		o(ed)
	}

	return ed, nil
}

// WatchService
func (s *EtcdDiscovery) WatchService(ctx context.Context) error {
	if s.cli == nil {
		return errors.New("cli is nil")
	}

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

func JSONDecode(val string) (*registry.ServiceNode, error) {
	node := &registry.ServiceNode{}
	err := json.Unmarshal([]byte(val), node)
	if err != nil {
		return nil, errors.New("Unmarshal val " + err.Error())
	}

	return node, nil
}
