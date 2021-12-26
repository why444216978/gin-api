package selector

type Statistics struct {
	Success uint64
	Fail    uint64
}

type Meta struct{}

type Node interface {
	Address() string
	Weight() int
	Meta() Meta
	Statistics() Statistics
}

type HandleInfo struct {
	Node Node
	Err  error
}
type Selector interface {
	ServiceName() string
	AddNode(node Node) (err error)
	DeleteNode(node Node) (err error)
	GetNodes() (nodes []Node, err error)
	Select() (node Node, err error)
	AfterHandle(info HandleInfo) (err error)
}
