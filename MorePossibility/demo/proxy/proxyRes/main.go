package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"go-Demo-MorePossibility/MorePossibility/proto/BurpApi"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"strings"
)

// ProxyResDemo 代理响应Demo
type ProxyResDemo struct {
	BurpApi.UnimplementedProxyResponseHandlerServer
}

// ProxyHandleResponseReceived 代理响应处理器 注解测试
func (ProxyResDemo) ProxyHandleResponseReceived(c context.Context, hr *BurpApi.HttpReqAndRes) (*BurpApi.ProxyResponseAction, error) {
	if hr.GetRes().GetStatusCode() == 404 { // 状态码等于404
		return &BurpApi.ProxyResponseAction{
			Continue:    true,
			Drop:        false,
			IsReviseRes: false,
			IsIntercept: false,
			HttpResGroup: &BurpApi.HttpResGroup{AnnotationsText: &BurpApi.AnnotationsText{
				IsInfo: true, // 添加注解
				Notes:  "404啦",
				Color:  BurpApi.HighlightColor_GREEN,
			}},
		}, nil
	}
	return &BurpApi.ProxyResponseAction{Continue: true}, nil
}

// ProxyResDropDemo 代理响应Demo
type ProxyResDropDemo struct {
	BurpApi.UnimplementedProxyResponseHandlerServer
}

// ProxyHandleResponseReceived 代理响应处理器 注解测试
func (ProxyResDropDemo) ProxyHandleResponseReceived(c context.Context, hr *BurpApi.HttpReqAndRes) (*BurpApi.ProxyResponseAction, error) {
	fmt.Println(string(hr.GetReq().GetData()))

	reader := bytes.NewReader(hr.GetReq().GetData())
	request, err := http.ReadRequest(bufio.NewReader(reader))

	if err != nil {
		panic(err)
	}

	//request.

	fmt.Println(request.URL)
	if strings.Contains(string(hr.Res.Data), "page") {
		fmt.Println("丢弃")
		return &BurpApi.ProxyResponseAction{Drop: true}, nil
	}
	return &BurpApi.ProxyResponseAction{Continue: true}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	//BurpApi.RegisterProxyResponseHandlerServer(server, ProxyResDemo{})  // 添加注解
	BurpApi.RegisterProxyResponseHandlerServer(server, ProxyResDropDemo{}) // 丢弃响应
	fmt.Println("服务启动")
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
