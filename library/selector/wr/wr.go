// wr is Weighted random
package wr

import (
	"math/rand"
	"sort"
	"sync"

	"github.com/why444216978/gin-api/library/selector"
)

type Node struct {
	lock       sync.RWMutex
	address    string
	weight     int
	meta       selector.Meta
	statistics selector.Statistics
}

var _ selector.Node = (*Node)(nil)

func NewNode(address string, weight int, meta selector.Meta) *Node {
	return &Node{
		address:    address,
		weight:     weight,
		meta:       meta,
		statistics: selector.Statistics{},
	}
}

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

func (n *Node) Weight() int {
	return n.weight
}

func (n *Node) incrSuccess() {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.statistics.Success = n.statistics.Success + 1
}

func (n *Node) incrFail() {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.statistics.Fail = n.statistics.Fail + 1
}

type nodeOffset struct {
	Address     string
	Weight      int
	OffsetStart int
	OffsetEnd   int
}

type Selector struct {
	lock        sync.RWMutex
	nodeCount   int
	nodes       map[string]*Node
	list        []*Node
	offsetList  []nodeOffset
	sameWeight  bool
	totalWeight int
	serviceName string
}

var _ selector.Selector = (*Selector)(nil)

type SelectorOption func(*Selector)

func WithServiceName(name string) SelectorOption {
	return func(s *Selector) { s.serviceName = name }
}

func NewSelector(opts ...SelectorOption) *Selector {
	s := &Selector{
		nodes:      make(map[string]*Node),
		list:       make([]*Node, 0),
		offsetList: make([]nodeOffset, 0),
	}

	for _, o := range opts {
		o(s)
	}

	return s
}

func (s *Selector) ServiceName() string {
	return s.serviceName
}

func (s *Selector) AddNode(node selector.Node) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	var (
		offsetStart = 0
		offsetEnd   = s.totalWeight + node.Weight()
	)
	if s.nodeCount > 0 {
		offsetStart = s.totalWeight + 1
	}

	offset := nodeOffset{
		Address:     node.Address(),
		Weight:      node.Weight(),
		OffsetStart: offsetStart,
		OffsetEnd:   offsetEnd,
	}

	wrNode := s.node2WRNode(node)

	s.totalWeight = offsetEnd
	s.nodes[node.Address()] = wrNode
	s.list = append(s.list, wrNode)
	s.offsetList = append(s.offsetList, offset)
	s.nodeCount = s.nodeCount + 1

	s.sortOffset()
	s.checkSameWeight()

	return
}

func (s *Selector) DeleteNode(node selector.Node) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.nodeCount = s.nodeCount - 1

	delete(s.nodes, node.Address())

	for idx, n := range s.list {
		if n.Address() != node.Address() {
			continue
		}
		new := make([]*Node, len(s.list)-1)
		new = append(s.list[:idx], s.list[idx+1:]...)
		s.list = new
	}

	for idx, n := range s.offsetList {
		if n.Address != node.Address() {
			continue
		}
		s.totalWeight = s.totalWeight - node.Weight()
		new := make([]nodeOffset, len(s.offsetList)-1)
		new = append(s.offsetList[:idx], s.offsetList[idx+1:]...)
		s.offsetList = new
	}

	s.sortOffset()
	s.checkSameWeight()

	return
}

func (s *Selector) GetNodes() (nodes []selector.Node, err error) {
	nodes = make([]selector.Node, 0)
	for _, n := range s.list {
		nodes = append(nodes, n)
	}
	return
}

func (s *Selector) Select() (node selector.Node, err error) {
	if s.sameWeight {
		idx := rand.Intn(s.nodeCount)
		node = s.list[idx]
		return
	}

	idx := rand.Intn(s.totalWeight + 1)
	for _, n := range s.offsetList {
		if idx >= n.OffsetStart && idx <= n.OffsetEnd {
			node = s.nodes[n.Address]
			break
		}
	}

	return
}

func (s *Selector) AfterHandle(info selector.HandleInfo) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	node := s.nodes[info.Node.Address()]

	if info.Err != nil {
		node.incrFail()
		return
	}
	node.incrSuccess()

	return
}

func (s *Selector) checkSameWeight() {
	s.sameWeight = true

	var last int
	for _, n := range s.list {
		cur := int(n.weight)
		if last == 0 {
			last = cur
			continue
		}
		if last == cur {
			last = cur
			continue
		}
		s.sameWeight = false
		return
	}
}

func (s *Selector) sortOffset() {
	sort.Slice(s.offsetList, func(i, j int) bool {
		return s.offsetList[i].Weight > s.offsetList[j].Weight
	})
}

func (s *Selector) node2WRNode(node selector.Node) *Node {
	return &Node{
		address:    node.Address(),
		weight:     node.Weight(),
		meta:       node.Meta(),
		statistics: node.Statistics(),
	}
}
