package rpc_client

import (
	"context"
	"github.com/g-Halo/go-server/pkg/pb"
	"google.golang.org/grpc"
)

var (
	AuthClient pb.AuthClient
	LogicClient pb.LogicClient
	WsClient pb.WsConnClient
)

func InitAuthClient(addr string) {
	conn, err := grpc.DialContext(context.TODO(), addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	AuthClient = pb.NewAuthClient(conn)
}

func InitLogicClient(addr string) {
	conn, err := grpc.DialContext(context.TODO(), addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	LogicClient = pb.NewLogicClient(conn)
}

func InitWsClient(addr string) {
	conn, err := grpc.DialContext(context.TODO(), addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	WsClient = pb.NewWsConnClient(conn)
}