package registry

import (
	"context"
	"time"
)

type RegistryConfig struct {
	Endpoints string
	Lease     int64
}

// Registrar is service registrar
type Registrar interface {
	Register(ctx context.Context) error
	DeRegister(ctx context.Context) error
}

type DiscoveryConfig struct {
	ServiceName string
	Endpoints   string
	DialTimeout time.Duration
}

// Discovery is service discovery
type Discovery interface {
	WatchService(ctx context.Context, serviceName string) error
	SetServiceList(key, val string)
	DelServiceList(key string)
	GetServices() []string
	Close() error
}
