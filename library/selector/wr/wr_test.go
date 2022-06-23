// wr is Weighted random
package wr

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/why444216978/gin-api/library/selector"
)

func TestNewNode(t *testing.T) {
	convey.Convey("TestNewNode", t, func() {
		convey.Convey("success", func() {
			ip := "127.0.0.1"
			port := 80
			weight := 10
			meta := selector.Meta{}
			node := NewNode(ip, port, weight, meta)
			assert.Equal(t, node.Address(), selector.GenerateAddress(ip, port))
			assert.Equal(t, node.Weight(), weight)
			assert.Equal(t, node.Meta(), meta)
		})
	})
}

func TestNode_Address(t *testing.T) {
}

func TestNode_Meta(t *testing.T) {
}

func TestNode_Statistics(t *testing.T) {
}

func TestNode_Weight(t *testing.T) {
}

func TestNode_incrSuccess(t *testing.T) {
}

func TestNode_incrFail(t *testing.T) {
}

func TestWithServiceName(t *testing.T) {
}

func TestNewSelector(t *testing.T) {
}

func TestSelector_ServiceName(t *testing.T) {
	convey.Convey("TestSelector_ServiceName", t, func() {
		convey.Convey("success", func() {
			s := NewSelector(WithServiceName("test_service"))
			serviceName := s.ServiceName()
			assert.Equal(t, serviceName, "test_service")
		})
	})
}

func TestSelector_AddNode(t *testing.T) {
}

func TestSelector_DeleteNode(t *testing.T) {
}

func TestSelector_GetNodes(t *testing.T) {
}

func TestSelector_Select(t *testing.T) {
}

func TestSelector_AfterHandle(t *testing.T) {
}

func TestSelector_checkSameWeight(t *testing.T) {
}

func TestSelector_sortOffset(t *testing.T) {
}

func TestSelector_node2WRNode(t *testing.T) {
}

func TestWR(t *testing.T) {
	convey.Convey("TestWR", t, func() {
		convey.Convey("testNoDeleteHandle same weight", func() {
			nodes := []*Node{
				{
					address: "127.0.0.1:80",
				},
				{
					address: "127.0.0.2:80",
				},
				{
					address: "127.0.0.3:80",
				},
			}
			res := testNoDeleteHandle(t, nodes)
			fmt.Println("\ntestNoDeleteHandle same weight")
			for _, n := range res {
				fmt.Println(n.Address(), ":", n.Statistics())
			}
		})
		convey.Convey("testNoDeleteHandle diff weight", func() {
			nodes := []*Node{
				{
					address: "127.0.0.1:80",
					weight:  2,
				},
				{
					address: "127.0.0.2:80",
					weight:  2,
				},
				{
					address: "127.0.0.3:80",
					weight:  1,
				},
			}
			res := testNoDeleteHandle(t, nodes)
			fmt.Println("\ntestNoDeleteHandle diff weight")
			for _, n := range res {
				fmt.Println(n.Address(), ":", n.Statistics())
			}
		})
		convey.Convey("testDeleteHandle same weight", func() {
			nodes := []*Node{
				{
					address: "127.0.0.1:80",
				},
				{
					address: "127.0.0.2:80",
				},
				{
					address: "127.0.0.3:80",
				},
			}
			res := testDeleteHandle(t, nodes)
			fmt.Println("\ntestDeleteHandle same weight")
			for _, n := range res {
				fmt.Println(n.Address(), ":", n.Statistics())
			}
		})
		convey.Convey("testDeleteHandle diff weight", func() {
			nodes := []*Node{
				{
					address: "127.0.0.1:80",
					weight:  2,
				},
				{
					address: "127.0.0.2:80",
					weight:  2,
				},
				{
					address: "127.0.0.3:80",
					weight:  1,
				},
			}
			res := testDeleteHandle(t, nodes)
			fmt.Println("\ntestDeleteHandle diff weight")
			for _, n := range res {
				fmt.Println(n.Address(), ":", n.Statistics())
			}
		})
	})
}

func testNoDeleteHandle(t *testing.T, nodes []*Node) []selector.Node {
	s := NewSelector(WithServiceName("test_service"))

	for _, node := range nodes {
		s.AddNode(node)
	}

	i := 1
	for {
		if i > 10000 {
			break
		}
		node, _ := s.Select()

		random := rand.Intn(100)
		err := errors.New("error")
		if random != 0 {
			err = nil
		}
		s.AfterHandle(node.Address(), err)
		i++
	}

	res, _ := s.GetNodes()
	return res
}

func testDeleteHandle(t *testing.T, nodes []*Node) []selector.Node {
	s := NewSelector(WithServiceName("test_service"))

	for _, node := range nodes {
		s.AddNode(node)
	}

	i := 1
	for {
		if i > 9000 {
			break
		}
		node, _ := s.Select()

		random := rand.Intn(100)
		err := errors.New("error")
		if random != 0 {
			err = nil
		}

		s.AfterHandle(node.Address(), err)
		i++
	}

	del := nodes[2]
	host, port := selector.ExtractAddress(del.address)
	_ = s.DeleteNode(host, port)
	i = 1
	for {
		if i > 1000 {
			break
		}
		node, _ := s.Select()

		random := rand.Intn(10)
		err := errors.New("error")
		if random != 0 {
			err = nil
		}

		assert.Equal(t, node.Address() != del.address, true)
		s.AfterHandle(node.Address(), err)
		i++
	}

	res, _ := s.GetNodes()
	return res
}
