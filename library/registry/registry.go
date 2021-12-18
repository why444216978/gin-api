package registry

import (
	"context"
)

const (
	TypeRegistry   uint8 = 1
	TypeHostPort   uint8 = 2
	TypeHostDomain uint8 = 3
)

type ServiceNode struct {
	Host      string
	Port      int
	CaCrt     []byte
	ClientPem []byte
	ClientKey []byte
}

type RegistryConfig struct {
	Lease int64
}

// Registrar is service registrar
type Registrar interface {
	Register(ctx context.Context) error
	DeRegister(ctx context.Context) error
}

type DiscoveryConfig struct {
	ServiceName string `validate:"required"`
	Type        uint8  `validate:"required,oneof=1 2"`
	Host        string `validate:"required"`
	Port        int    `validate:"required"`
	LoadBalance string `validate:"required,oneof=random round_robin"`
	CaCrt       string
	ClientPem   string
	ClientKey   string
}

// Discovery is service discovery
type Discovery interface {
	WatchService(ctx context.Context) error
	SetServiceList(key string, val *ServiceNode)
	DelServiceList(key string)
	GetServices() []*ServiceNode
	GetLoadBalance() string
	Close() error
}

// Encode func is encode service node info
type Encode func(node *ServiceNode) (string, error)

// Decode func is decode service node info
type Decode func(val string) (*ServiceNode, error)

var Services = make(map[string]Discovery)
