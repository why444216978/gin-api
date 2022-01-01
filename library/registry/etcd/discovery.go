package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/why444216978/gin-api/library/registry"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const defaultRefreshDuration = time.Second * 10

//EtcdDiscovery 服务发现
type EtcdDiscovery struct {
	ctx             context.Context
	cli             *clientv3.Client
	nodeList        map[string]*registry.Node
	lock            sync.RWMutex
	decode          registry.Decode
	ticker          *time.Ticker
	refreshDuration time.Duration
	serviceName     string
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

func WithRefreshDuration(d int) DiscoverOption {
	return func(ed *EtcdDiscovery) { ed.refreshDuration = time.Duration(d) * time.Second }
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
	if s.ticker != nil {
		s.ticker.Stop()
	}
	if s.cli != nil {
		s.cli.Close()
	}
	return nil
}

// WatchService
func (s *EtcdDiscovery) init() error {
	//set all nodes
	s.setNodes()

	//start etcd watcher
	go s.watcher()

	//start refresh ticker
	go s.refresh()

	return nil
}

// loadKVs
func (s *EtcdDiscovery) loadKVs() (kvs []*mvccpb.KeyValue) {
	resp, err := s.cli.Get(s.ctx, s.serviceName, clientv3.WithPrefix())
	if err != nil {
		s.logErr("get by prefix", s.serviceName, "", err)
		return
	}
	kvs = resp.Kvs
	return
}

// watcher
func (s *EtcdDiscovery) watcher() {
	rch := s.cli.Watch(s.ctx, s.serviceName, clientv3.WithPrefix())
	s.log("Watch", "")
	for wresp := range rch {
		for _, ev := range wresp.Events {
			key := string(ev.Kv.Key)
			val := string(ev.Kv.Value)

			switch ev.Type {
			case mvccpb.PUT:
				node, err := s.decode(val)
				if err != nil {
					s.logErr("decode val", key, val, err)
					return
				}
				s.setNode(key, node)
				s.log("mvccpb.PUT", key)
			case mvccpb.DELETE:
				s.delServiceList(key)
				s.log("mvccpb.DELETE", key)
			}
		}
	}
}

//refresh
func (s *EtcdDiscovery) refresh() {
	if s.refreshDuration == 0 {
		s.refreshDuration = defaultRefreshDuration
	}

	s.ticker = time.NewTicker(s.refreshDuration)
	for range s.ticker.C {
		s.setNodes()
		s.log("refresh", "all")
	}
}

// setNodes
func (s *EtcdDiscovery) setNodes() {
	nodeList := make(map[string]*registry.Node)
	kvs := s.loadKVs()
	for _, kv := range kvs {
		key := string(kv.Key)
		val := string(kv.Value)

		node, err := s.decode(val)
		if err != nil {
			s.logErr("decode val", key, val, err)
			continue
		}
		nodeList[key] = node
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	s.nodeList = nodeList
}

// setNode
func (s *EtcdDiscovery) setNode(key string, node *registry.Node) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.nodeList[key] = node
}

// delServiceList
func (s *EtcdDiscovery) delServiceList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.nodeList, key)
}

func (s *EtcdDiscovery) logErr(action, key, val string, err error) {
	log.Printf("[%s]: action:%s, err:%s, service:%s, key:%s, val:%s\n", time.Now().Format("2006-01-02 15:04:05"), action, s.serviceName, key, val, err.Error())
}

func (s *EtcdDiscovery) log(action, key string) {
	log.Printf("[%s]: [action:%s, service:%s, key:%s]\n", time.Now().Format("2006-01-02 15:04:05"), action, s.serviceName, key)
}

func JSONDecode(val string) (*registry.Node, error) {
	node := &registry.Node{}
	err := json.Unmarshal([]byte(val), node)
	if err != nil {
		return nil, errors.New("Unmarshal val " + err.Error())
	}

	return node, nil
}
