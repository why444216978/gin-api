package xhop

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
)

//
// Generally speaking,for every http request, the hierarchy of the http rpc calls is a DAG(directed acycline graph).
//
// The X-Hop system is designed to accurately track the hierarchical http rpc call details. It use Hop, a binary encoded integer, to describe every module's hierarchy in the entire DAG from beginning/app to the farest backend end.
//
// In X-Hop system, a module is a http server instance that serves http requests and possibly send one or more rpc calls to other modules. A module must send correct X-Hop value in the request header X-Hop field when sending a http call and its downstream then receives the X-Hop from the request header. When a module receive a request without a X-Hop, it should set X-Hop to 1. The Hop value of a module is calcuted as follows:
//
// 1. Every positive bit in the X-Hop indicates a module.
// 2. A module's X-Hop stripping the leftmost positive bit indicates its parent's X-Hop.
//   a. The root module's (usually the outer GFW) Hop is always  1.
//   b. The count of positive bit in a module's X-Hop indicates the hierarchy of the entire http rpc calls for the specific module.
// 2. The ith http call's X-hop sending from a module M should add (100...0)b to the leftmost of M's X-Hop. There as (i - 1) 0s in total.
//
// For example, module A calls B and C, and B calls D and E, and C calls F, and D calls G at last.
// The DAG shows as follows:
//
//		  	               A
//		  	              / \
//		  	             B   C
//		  	            / \  |
//		  	           D   E F
//		  	          / \    |
//		  	         G   G2  H
//		  	        / \     /|\
//		  	       I   I2  J K L
//			      / \        |
//		         M   M2      N
//		        / \          |
//		       O   P         Q
//	              / \        |
//	             R   S       T
//	                / \      |
//	               U   V     W
//                     |    / \
//                     X   Y   Z
//
// Suppose Hop(i) indicates the module i's Hop, then:
// Hop(A) = (00000001)b
// Hop(B) = (00000011)b
// Hop(C) = (00000101)b
// Hop(D) = (00000111)b
// Hop(E) = (00001011)b
// Hop(F) = (00001101)b
// Hop(G) = (00001111)b
// Hop(G2) = (00010111)b
// Hop(H) = (00011101)b
// Hop(I) = (00011111)b
// Hop(I2) = (00101111)b
// Hop(J) = (00111101)b
// Hop(K) = (01011101)b
// Hop(L) = (10011101)b
// Hop(M) = (00111111)b
// Hop(M2) = (01011111)b
// Hop(N) = (11011101)b
// Hop(O) = (01111111)b
// Hop(P) = (10111111)b
// Hop(Q) = (00000001 11011101)b
// Hop(R) = (00000001 10111111)b
// Hop(S) = (00000010 10111111)b
// Hop(T) = (00000011 11011101)b
// Hop(U) = (00000110 10111111)b
// Hop(V) = (00001010 10111111)b
// Hop(W) = (00000111 11011101)b
// Hop(X) = (00011010 10111111)b
// Hop(Y) = (00001111 11011101)b

type XHop struct {
	buf []byte //use little endian for efficiency
	c   uint64 //current children rpc call counter
}

func New() *XHop {
	b := &XHop{
		buf: make([]byte, 1),
	}
	b.buf[0] = byte(1)

	return b
}

func (x *XHop) MarshalJSON() ([]byte, error) {
	if x == nil {
		return nil, errors.New("nil XHop")
	}

	return []byte(fmt.Sprintf("\"%s\"", x.Hex())), nil
}

func (x *XHop) IsRootXHop() bool {
	if x == nil {
		return false
	}

	return len(x.buf) == 1 && x.buf[0] == byte(1)
}

func NewFromHex(s string) (*XHop, error) {
	if buf, err := hex.DecodeString(s); err != nil {
		return nil, err
	} else {
		//change buf from BigEndian to LittleEndian
		//swap
		var (
			i int = 0
			j int = len(buf) - 1
			p byte
		)

		for i < j {
			p = buf[i]
			buf[i] = buf[j]
			buf[j] = p
			i++
			j--
		}
		//the byteorder of buf is already LittleEndian now.
		//the last byte, the most important byte, should not be zero.
		if buf[len(buf)-1] == 0 {
			return nil, fmt.Errorf("corrupted xhop, the last byte should not be zero:%s", s)
		}

		return &XHop{buf: buf}, nil
	}
}

func (x *XHop) Dup() *XHop {
	if x == nil {
		return New()
	}

	xhop := &XHop{
		buf: make([]byte, len(x.buf)),
	}
	copy(xhop.buf, x.buf)

	return xhop
}

func (x *XHop) Equal(x2 *XHop) bool {
	if x == nil || x2 == nil {
		return x == nil && x2 == nil
	} else {
		return bytes.Compare(x.buf, x2.buf) == 0 && x.c == x2.c
	}
}

//String should only be used for human reading.
//The func does not take performance into account.
//little endian
func (x *XHop) String() string {
	if x == nil {
		return ""
	}

	var (
		s = make([]byte, 0, 9*len(x.buf)-1) //8*N + N-1
		u uint8
	)
	for i := len(x.buf) - 1; i >= 0; i-- {
		u = 128
		for u > 0 {
			if u&x.buf[i] == u {
				s = append(s, '1')
			} else {
				s = append(s, '0')

			}
			u >>= 1
		}
		s = append(s, ' ') //add extra space for human reading
	}

	return string(s[:len(s)-1]) //strip the last extra space
}

func (x *XHop) Hex() string {
	if x == nil {
		return ""
	}

	//create a new buf in BigEndian
	var buf = make([]byte, len(x.buf))
	for i := 0; i < len(buf); i++ {
		buf[i] = x.buf[len(buf)-1-i]
	}

	return hex.EncodeToString(buf)
}

//https://leetcode.com/problems/counting-bits/#/description
func (x *XHop) Hierarchy() (n int) {
	if x == nil {
		return
	}

	for _, b := range []byte(x.buf) {
		n += byteBits(b)
	}

	return
}

// 1. strip the leftmost positive bit.
// 2. strip the leading zero bytes before 1.
//little endian
func (x *XHop) Parent() *XHop {
	if x == nil {
		return New()
	}

	//Root XHop's Parent is always itself.
	if x.IsRootXHop() {
		return x.Dup()
	}

	//attention: x.buf is in LittleEndian!
	xhop := x.Dup()
	var b = x.buf[len(x.buf)-1]
	xhop.buf[len(x.buf)-1] = b & ^(1 << (leftMostBitPos(b) - 1))

	//strip the left zero bytes
	n := len(xhop.buf) - 1
	for n > 0 {
		if xhop.buf[n] != 0 {
			break
		}
		n -= 1
	}

	xhop.buf = xhop.buf[:n+1]
	return xhop
}

//pading 100...0(seq 0) the x
func (x *XHop) Next() *XHop {
	if x == nil {
		return New().Next()
	}

	xhop := x.Dup()

	var (
		b  = x.buf[len(x.buf)-1]
		bp = uint64(leftMostBitPos(b)) //1~8
	)

	xhop.c = x.c

	if xhop.c >= 8-bp {
		xhop.c -= 8 - bp
		leading := make([]byte, xhop.c/8) //zero bytes
		leading = append(leading, 1<<(xhop.c%8))
		xhop.buf = append(xhop.buf, leading...)
	} else {
		xhop.buf[len(x.buf)-1] += (1 << xhop.c) << bp
	}

	//update counter
	xhop.c = 0
	x.c += 1

	return xhop
}

//pading 100...0(n 0 in total) the x
//different to Next(), NextN do not touch x.c
//used for Test_Smoke
//used for Test_Hierarchy
func (x *XHop) NextN(n uint64) *XHop {
	if x == nil {
		return New().NextN(n)
	}

	xhop := x.Dup()
	xhop.c = n
	return xhop.Next()
}

func byteBits(b byte) (c int) {
	n := uint8(b)
	for n > 0 {
		n &= (n - 1)
		c += 1
	}

	return
}

func leftMostBitPos(b byte) uint8 {
	//11111111
	//8...1

	var (
		p   uint8 = 128
		pos uint8 = 8
	)

	for p > 0 {
		if b&p == p {
			break
		}

		p >>= 1
		pos -= 1
	}

	return pos
}
