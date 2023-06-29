package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"go-Demo-MorePossibility/MorePossibility/proto/BurpApi"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

// intruderPayloadGeneratorDemo 迭代生成器
type intruderPayloadGeneratorDemo struct {
	BurpApi.UnimplementedIntruderPayloadGeneratorServerServer
	index int
}

var dictionary = []string{"123456", "admin", "sa", "admin123", "9521cc"}

// IntruderPayloadGeneratorProvider 生成器
func (i *intruderPayloadGeneratorDemo) IntruderPayloadGeneratorProvider(c context.Context, igd *BurpApi.IntruderGeneratorData) (*BurpApi.PayloadGeneratorResult, error) {
	if len(dictionary) <= i.index { // 下标是否还有可取元素
		i.index = 0
		return &BurpApi.PayloadGeneratorResult{
			ByteData: nil,
			IsEnd:    true, // 结束
		}, nil
	}
	reqData := igd.GetContentData()
	request, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(reqData))) // 解析请求
	if err != nil {
		log.Println(err)
		return &BurpApi.PayloadGeneratorResult{
			ByteData: nil,
			IsEnd:    true, // 出错结束
		}, nil
	}
	host := request.Host

	str := host + dictionary[i.index] // host拼接字典
	i.index++
	return &BurpApi.PayloadGeneratorResult{
		ByteData: []byte(str),
		IsEnd:    false,
	}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	BurpApi.RegisterIntruderPayloadGeneratorServerServer(server, &intruderPayloadGeneratorDemo{})

	fmt.Println("服务启动")
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
