package rpc

type RPCServer interface {
	Start() error
	Close() error
}
