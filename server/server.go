package server

type RPCServer interface {
	Start() error
	Close() error
}
