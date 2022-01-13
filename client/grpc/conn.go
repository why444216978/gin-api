package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

func Conn(ctx context.Context, target string) (cc *grpc.ClientConn, err error) {
	//TODO resolver
	cc, err = grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
	if err != nil {
		return
	}

	return
}
