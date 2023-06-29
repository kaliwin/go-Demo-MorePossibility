package main

import (
	"context"
	"fmt"
	"go-Demo-MorePossibility/MorePossibility/proto/BurpApi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	//dial, err := grpc.Dial("127.0.0.1:9523", grpc.WithTransportCredentials(insecure.NewCredentials()))
	maxMessageSize := 500 * 1024 * 1024

	dial, err := grpc.Dial("127.0.0.1:9525",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMessageSize)))
	if err != nil {
		panic(err)
	}

	//client := BurpApi.NewBurpServerClient(dial)
	//mirroring, err := client.RegisterRealTimeTrafficMirroring(context.Background(), &BurpApi.Str{Name: "di"})
	//if err != nil {
	//	panic(err)
	//}
	//
	//for true {
	//	recv, err := mirroring.Recv()
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println("接收到: " + recv.Req.Url)
	//
	//}

	// ====================   服务注册测试  ==========================  //

	client := BurpApi.NewBurpServerClient(dial)
	list, err := client.RegisterServerList(context.Background(), &BurpApi.ServiceRegisterRoutingList{
		ServiceList: []*BurpApi.ServiceRegisterRouting{
			&BurpApi.ServiceRegisterRouting{
				ServerType:  BurpApi.ServerTypeName_INTRUDER_GENERATE,
				Name:        "payloadGenerateTest",
				GrpcAddress: "127.0.0.1:9000",
			},
		},
	})

	if err != nil {
		panic(err)
	}

	if list.GetBoole() {
		fmt.Println("注册成功")
	} else {
		fmt.Println("注册失败")
	}
}
