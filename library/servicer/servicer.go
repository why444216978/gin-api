package servicer

import (
	"context"
	"sync"
)

const (
	TypeRegistry uint8 = 1
	TypeIPPort   uint8 = 2
	TypeDomain   uint8 = 3
)

type Node struct {
	Host string
	Port int
}

type DoneInfo struct {
	Node *Node
	Err  error
}

var (
	lock      sync.RWMutex
	Servicers = make(map[string]Servicer)
)

func SetServicer(s Servicer) {
	lock.Lock()
	defer lock.Unlock()
	Servicers[s.Name()] = s
}

func DelServicer(s Servicer) {
	lock.Lock()
	defer lock.Unlock()
	delete(Servicers, s.Name())
}

func GetServicer(serviceName string) (Servicer, bool) {
	s, has := Servicers[serviceName]
	return s, has
}

type Servicer interface {
	Name() string
	Pick(ctx context.Context) (*Node, error)
	Done(ctx context.Context, node *Node, err error) error
	GetCaCrt() []byte
	GetClientPem() []byte
	GetClientKey() []byte
}
