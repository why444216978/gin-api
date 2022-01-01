package servicer

import (
	"context"
	"errors"
	"net"

	"github.com/why444216978/go-util/validate"

	"github.com/why444216978/gin-api/library/registry"
	"github.com/why444216978/gin-api/library/selector"
	"github.com/why444216978/gin-api/library/selector/wr"
)

const (
	TypeRegistry uint8 = 1
	TypeIPPort   uint8 = 2
	TypeDomain   uint8 = 3
)

var Servicers = make(map[string]Servicer)

type Node struct {
	Host string
	Port int
}

type DoneInfo struct {
	Node *Node
	Err  error
}

type Servicer interface {
	Pick(ctx context.Context) (*Node, error)
	Done(ctx context.Context, node *Node, err error) error
	GetCaCrt() []byte
	GetClientPem() []byte
	GetClientKey() []byte
}

type Config struct {
	ServiceName   string `validate:"required"`
	Type          uint8  `validate:"required,oneof=1 2"`
	Host          string `validate:"required"`
	Port          int    `validate:"required"`
	Selector      string `validate:"required,oneof=wr"` //TODO 后续支持其它
	CaCrt         string
	ClientPem     string
	ClientKey     string
	RefreshSecond int
}

type Service struct {
	selector        selector.Selector
	selectorNewNode selector.NewNodeFunc
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

func LoadService(config *Config, opts ...Option) error {
	s := &Service{
		config:    config,
		caCrt:     []byte(config.CaCrt),
		clientPem: []byte(config.ClientPem),
		clientKey: []byte(config.ClientKey),
	}

	for _, o := range opts {
		o(s)
	}

	if err := validate.ValidateCamel(config); err != nil {
		return err
	}

	s.initSelector()

	Servicers[config.ServiceName] = s

	return nil
}

func (s *Service) Pick(ctx context.Context) (node *Node, err error) {
	node = &Node{
		Port: s.config.Port,
	}

	if s.config.Type == TypeIPPort {
		node.Host = s.config.Host
		return
	}

	if s.config.Type == TypeDomain {
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
	if s.config.Type != TypeRegistry {
		return nil
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
	var (
		address     string
		host        string
		port        int
		nowNodes    = s.discovery.GetNodes()
		nowMap      = make(map[string]struct{})
		selectorMap = make(map[string]selector.Node)
	)

	for _, n := range nowNodes {
		host = n.Host
		port = n.Port
		address = selector.GenerateAddress(host, port)
		node := s.selectorNewNode(host, port, n.Weight, selector.Meta{})

		nowMap[address] = struct{}{}
		selectorMap[address] = node

		_, ok := s.selector.GetNode(host, port)
		if ok {
			continue
		}
		_ = s.selector.AddNode(node)
	}

	selectorNodes, _ := s.selector.GetNodes()
	for _, n := range selectorNodes {
		if _, ok := nowMap[n.Address()]; ok {
			continue
		}
		host, port = selector.ExtractAddress(n.Address())
		_ = s.selector.DeleteNode(host, port)
	}
}

func (s *Service) Done(ctx context.Context, node *Node, err error) error {
	if s.selector == nil {
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
