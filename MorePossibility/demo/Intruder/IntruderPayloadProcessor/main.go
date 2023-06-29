package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"go-Demo-MorePossibility/MorePossibility/proto/BurpApi"
	"google.golang.org/grpc"
	"net"
)

type IntruderDemo struct {
	BurpApi.UnimplementedIntruderPayloadProcessorServerServer
}

// IntruderPayloadProcessor 迭代处理器 将载荷进行base64编码
func (IntruderDemo) IntruderPayloadProcessor(c context.Context, byteS *BurpApi.ByteData) (*BurpApi.ByteData, error) {
	byteData := byteS.GetByteData() // 获取载荷
	fmt.Println("原始载荷: " + string(byteData))
	encodeToString := base64.StdEncoding.EncodeToString(byteData)   // 进行base64 编码
	return &BurpApi.ByteData{ByteData: []byte(encodeToString)}, nil //返回数据
}

func main() {
	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	BurpApi.RegisterIntruderPayloadProcessorServerServer(server, IntruderDemo{})

	fmt.Println("服务启动")
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}

}
