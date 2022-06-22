package service

import (
	"context"
	"errors"
	"net"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/why444216978/go-util/assert"
	utilDir "github.com/why444216978/go-util/dir"
	"github.com/why444216978/go-util/validate"

	"github.com/why444216978/gin-api/library/config"
	"github.com/why444216978/gin-api/library/etcd"
	"github.com/why444216978/gin-api/library/registry"
	registryEtcd "github.com/why444216978/gin-api/library/registry/etcd"
	"github.com/why444216978/gin-api/library/selector"
	"github.com/why444216978/gin-api/library/selector/wr"
	"github.com/why444216978/gin-api/library/servicer"
)

func LoadGlobPattern(path, suffix string, etcd *etcd.Etcd) (err error) {
	var (
		dir   string
		files []string
	)

	if dir, err = config.Dir(); err != nil {
		return
	}

	if files, err = filepath.Glob(filepath.Join(dir, path, "*."+suffix)); err != nil {
		return
	}

	var discover registry.Discovery
	info := utilDir.FileInfo{}
	cfg := &Config{}
	for _, f := range files {
		if info, err = utilDir.GetPathInfo(f); err != nil {
			return
		}
		if err = config.ReadConfig(filepath.Join("services", info.BaseNoExt), info.ExtNoSpot, cfg); err != nil {
			return
		}

		if cfg.Type == servicer.TypeRegistry {
			if assert.IsNil(etcd) {
				return errors.New("LoadGlobPattern etcd nil")
			}
			opts := []registryEtcd.DiscoverOption{
				registryEtcd.WithServierName(cfg.ServiceName),
				registryEtcd.WithRefreshDuration(cfg.RefreshSecond),
				registryEtcd.WithDiscoverClient(etcd.Client),
			}
			if discover, err = registryEtcd.NewDiscovery(opts...); err != nil {
				return
			}
		}

		if err = LoadService(cfg, WithDiscovery(discover)); err != nil {
			return
		}
	}

	return
}

func LoadService(config *Config, opts ...Option) (err error) {
	s, err := NewService(config, opts...)
	if err != nil {
		return
	}

	servicer.SetServicer(s)

	return nil
}

type Config struct {
	ServiceName   string `validate:"required"`
	Type          uint8  `validate:"required,oneof=1 2"`
	Host          string `validate:"required"`
	Port          int    `validate:"required"`
	Selector      string `validate:"required,oneof=wr"` // TODO 后续支持其它
	CaCrt         string
	ClientPem     string
	ClientKey     string
	RefreshSecond int
}

type Service struct {
	sync.RWMutex
	selector        selector.Selector
	selectorNewNode selector.NewNodeFunc
	adjusting       int32
	updateTime      time.Time
	discovery       registry.Discovery
	caCrt           []byte
	clientPem       []byte
	clientKey       []byte
	config          *Config
}

type Option func(*Service)

func WithDiscovery(discovery registry.Discovery) Option {
	return func(s *Service) { s.discovery = discovery }
}

func NewService(config *Config, opts ...Option) (*Service, error) {
	s := &Service{
		adjusting: 0,
		config:    config,
		caCrt:     []byte(config.CaCrt),
		clientPem: []byte(config.ClientPem),
		clientKey: []byte(config.ClientKey),
	}

	for _, o := range opts {
		o(s)
	}

	if err := validate.ValidateCamel(config); err != nil {
		return nil, err
	}

	if err := s.initSelector(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Service) Name() string {
	return s.config.ServiceName
}

func (s *Service) Pick(ctx context.Context) (node *servicer.Node, err error) {
	node = &servicer.Node{
		Port: s.config.Port,
	}

	if s.config.Type == servicer.TypeIPPort {
		node.Host = s.config.Host
		return
	}

	if s.config.Type == servicer.TypeDomain {
		var host *net.IPAddr
		host, err = net.ResolveIPAddr("ip", s.config.Host)
		if err != nil {
			return
		}
		node.Host = host.IP.String()
		return
	}

	s.adjustSelectorNode()

	target, err := s.selector.Select()
	if err != nil {
		return
	}

	node.Host, node.Port = selector.ExtractAddress(target.Address())

	return
}

func (s *Service) initSelector() (err error) {
	if s.config.Type != servicer.TypeRegistry {
		return nil
	}

	if assert.IsNil(s.discovery) {
		return errors.New("discovery is nil")
	}

	switch s.config.Selector {
	case selector.TypeWR:
		s.selector = wr.NewSelector(wr.WithServiceName(s.config.ServiceName))
		s.selectorNewNode = wr.NewNode
	}

	s.adjustSelectorNode()

	return nil
}

func (s *Service) adjustSelectorNode() {
	if s.discovery.GetUpdateTime().Before(s.updateTime) {
		return
	}

	if !atomic.CompareAndSwapInt32(&s.adjusting, 0, 1) {
		return
	}

	s.Lock()
	defer s.Unlock()

	var (
		address     string
		host        string
		port        int
		nowNodes    = s.discovery.GetNodes()
		nowMap      = make(map[string]struct{})
		selectorMap = make(map[string]selector.Node)
	)

	// selector add new nodes
	for _, n := range nowNodes {
		host = n.Host
		port = n.Port
		address = selector.GenerateAddress(host, port)
		node := s.selectorNewNode(host, port, n.Weight, selector.Meta{})

		nowMap[address] = struct{}{}
		selectorMap[address] = node

		_ = s.selector.AddNode(node)
	}

	// selector delete non-existent nodes
	selectorNodes, _ := s.selector.GetNodes()
	for _, n := range selectorNodes {
		if _, ok := nowMap[n.Address()]; ok {
			continue
		}
		host, port = selector.ExtractAddress(n.Address())
		_ = s.selector.DeleteNode(host, port)
	}

	s.updateTime = time.Now()
	atomic.StoreInt32(&s.adjusting, 0)
}

func (s *Service) Done(ctx context.Context, node *servicer.Node, err error) error {
	if assert.IsNil(s.selector) {
		return errors.New("selector is nil")
	}
	s.selector.AfterHandle(selector.GenerateAddress(node.Host, node.Port), err)
	return nil
}

func (s *Service) GetCaCrt() []byte {
	return s.caCrt
}

func (s *Service) GetClientPem() []byte {
	return s.clientPem
}

func (s *Service) GetClientKey() []byte {
	return s.clientKey
}
