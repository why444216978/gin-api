// dwrr is Dynamic Weighted Round Robin
package dwrr

import (
	"math"
	"sync"

	"github.com/why444216978/gin-api/library/selector"
)

const defaultStep float64 = 0.1

type Node struct {
	lock          sync.RWMutex
	address       string
	weight        float64
	currentWeight float64
	meta          selector.Meta
	statistics    selector.Statistics
}

var _ selector.Node = (*Node)(nil)

func (n *Node) Address() string {
	return n.address
}

func (n *Node) Meta() selector.Meta {
	return n.meta
}

func (n *Node) Statistics() selector.Statistics {
	n.lock.RLock()
	defer n.lock.RUnlock()
	return n.statistics
}

func (n *Node) Weight() float64 {
	return n.weight
}

type Selector struct {
	lock        sync.RWMutex
	nodeCount   int
	nodes       map[string]*Node
	list        []*Node
	step        float64
	serviceName string
}

var _ selector.Selector = (*Selector)(nil)

type SelectorOption func(*Selector)

func WithServiceName(name string) SelectorOption {
	return func(s *Selector) { s.serviceName = name }
}

func WithStep(step float64) SelectorOption {
	return func(s *Selector) { s.step = step }
}

func NewSelector(opts ...SelectorOption) *Selector {
	s := &Selector{
		nodes: make(map[string]*Node),
		list:  make([]*Node, 0),
	}

	for _, o := range opts {
		o(s)
	}

	if s.step <= 0 {
		s.step = defaultStep
	}

	return s
}

func (s *Selector) incrWeight(n *Node) {
	n.lock.Lock()
	defer n.lock.Unlock()

	n.currentWeight = math.Trunc(n.currentWeight*(1+s.step)*1e2+0.5) * 1e-2

	if n.currentWeight > n.weight {
		n.currentWeight = n.weight
	}

	return
}

func (s *Selector) decreaseWeight(n *Node) {
	n.lock.Lock()
	defer n.lock.Unlock()

	n.currentWeight = math.Trunc(n.currentWeight*(1-s.step)*1e2+0.5) * 1e-2

	return
}
