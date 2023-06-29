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
)

type httpFlowResDemo struct {
	BurpApi.UnimplementedHttpFlowHandlerServer
}

// HttpHandleRequestReceived 不用也要实现这个接口返回 Continue: true
func (httpFlowResDemo) HttpHandleRequestReceived(c context.Context, hrg *BurpApi.HttpReqGroup) (*BurpApi.HttpRequestAction, error) {
	return &BurpApi.HttpRequestAction{Continue: true}, nil
}

// HttpHandleResponseReceived http流量响应演示
func (httpFlowResDemo) HttpHandleResponseReceived(c context.Context, httpReqAndRes *BurpApi.HttpReqAndRes) (*BurpApi.HttpResponseAction, error) {
	fmt.Println("流量处理")
	if httpReqAndRes.GetRes().GetStatusCode() == 404 { // 响应状态码是404
		fmt.Println("修改响应")
		res := httpReqAndRes.GetRes() // 获取响应
		reader := bufio.NewReader(bytes.NewReader(res.GetData()))
		response, err := http.ReadResponse(reader, nil) // 将字节流解析为响应体对象
		if err != nil {
			panic(err)
		}
		response.StatusCode = 200 // 设置状态码
		response.Status = ""
		var buffer bytes.Buffer
		_ = response.Write(&buffer)

		res.Data = buffer.Bytes() // 将响应对象转回字节数组

		return &BurpApi.HttpResponseAction{ // 返回数据
			Continue:    false,
			IsReviseRes: true, // 修改响应
			HttpResGroup: &BurpApi.HttpResGroup{
				HttpResData:     res,
				AnnotationsText: nil,
			},
		}, nil
	}
	return &BurpApi.HttpResponseAction{Continue: true}, nil
}

func main() {

	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	BurpApi.RegisterHttpFlowHandlerServer(server, httpFlowResDemo{})

	fmt.Println("服务启动")
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
