package rpc_client

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/g-Halo/go-server/pkg/pb"
	"google.golang.org/grpc"
)

var (
	AuthClient  pb.AuthClient
	LogicClient pb.LogicClient
	WsClient    pb.WsConnClient
	CometConn   net.Conn
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

func ConnectToComet(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Print("commet connect not work")
		return
	}

	var zeroTime time.Time
	// 不超时
	conn.SetDeadline(zeroTime)
	conn.Write([]byte("  g-halo"))
	CometConn = conn
}
