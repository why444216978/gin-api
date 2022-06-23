package selector

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	TypeWR   = "wr"
	TypeWrr  = "wrr"
	TypeDwrr = "dwrr"
	TypeP2C  = "p2c"
	TypeICMP = "icmp"
)

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

type NewNodeFunc func(host string, port, weight int, meta Meta) Node

type HandleInfo struct {
	Node Node
	Err  error
}

type Selector interface {
	ServiceName() string
	AddNode(node Node) (err error)
	DeleteNode(host string, port int) (err error)
	GetNodes() (nodes []Node, err error)
	GetNode(host string, port int) (node Node, ok bool)
	Select() (node Node, err error)
	AfterHandle(address string, err error)
}

func GenerateAddress(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

func ExtractAddress(address string) (string, int) {
	arr := strings.Split(address, ":")
	if len(arr) != 2 {
		return "", 0
	}

	port, _ := strconv.Atoi(arr[1])
	return arr[0], port
}
