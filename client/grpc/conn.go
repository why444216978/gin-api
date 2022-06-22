package client

import (
	"context"

	"google.golang.org/grpc"

	serverGRPC "github.com/why444216978/gin-api/server/grpc"
)

func Conn(ctx context.Context, target string) (cc *grpc.ClientConn, err error) {
	// TODO resolver
	cc, err = grpc.DialContext(ctx, target, serverGRPC.NewDialOption()...)
	if err != nil {
		return
	}

	return
}
