package main

import (
	"fmt"
	"go-Demo-MorePossibility/MorePossibility/proto/BurpApi"
	"google.golang.org/grpc"
	"net"
)

// RealTimeTrafficMirroringDemo 实时流量镜像
type RealTimeTrafficMirroringDemo struct {
	BurpApi.UnimplementedRealTimeTrafficMirroringServer
}

// RealTimeTrafficMirroring 实时流量镜像
func (RealTimeTrafficMirroringDemo) RealTimeTrafficMirroring(e BurpApi.RealTimeTrafficMirroring_RealTimeTrafficMirroringServer) error {
	fmt.Println("[+] 建立实时流量通道")
	for true {
		recv, err := e.Recv()
		if err != nil {
			panic(err)
		}
		req := recv.GetReq()
		res := recv.GetRes()
		fmt.Println(string(req.GetData())) // 打印请求
		fmt.Println()
		fmt.Println(string(res.GetData())) // 打印响应
	}
	return nil
}

func main() {
	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	BurpApi.RegisterRealTimeTrafficMirroringServer(server, RealTimeTrafficMirroringDemo{})
	fmt.Println("服务启动")
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
