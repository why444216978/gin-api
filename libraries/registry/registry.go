package registry

import (
	"context"
)

type ServiceNode struct {
	Host string
	Port int
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
	ServiceName string
}

// Discovery is service discovery
type Discovery interface {
	WatchService(ctx context.Context) error
	SetServiceList(key, val string)
	DelServiceList(key string)
	GetServices() []*ServiceNode
	Close() error
}

// Encode func is encode service node info
type Encode func(node *ServiceNode) (string, error)

// Decode func is decode service node info
type Decode func(val string) (*ServiceNode, error)
