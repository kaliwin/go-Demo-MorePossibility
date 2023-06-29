package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"go-Demo-MorePossibility/MorePossibility/proto/BurpApi"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

// HttpFlowHandler http请求流量处理
type httpReqFlowDemo struct {
	BurpApi.UnimplementedHttpFlowHandlerServer
}

// HttpHandleRequestReceived 请求处理
func (httpReqFlowDemo) HttpHandleRequestReceived(c context.Context, hrg *BurpApi.HttpReqGroup) (*BurpApi.HttpRequestAction, error) {
	reader := bytes.NewReader(hrg.GetHttpReqData().GetData()) // 请求的字节数组
	request, err := http.ReadRequest(bufio.NewReader(reader)) // 转为请求实例
	if err != nil {
		log.Println(err)
		return &BurpApi.HttpRequestAction{Continue: true}, nil // 出错打印错误信息后返回继续
	}
	url := hrg.GetHttpReqData().GetUrl()                                             // 获取URL
	if strings.Contains(url, "login") && strings.EqualFold(request.Method, "POST") { // 请求方法是POST
		all, err := io.ReadAll(request.Body) // 读出请求体
		if err != nil {
			log.Println(err)
			return &BurpApi.HttpRequestAction{Continue: true}, nil
		}
		replace := strings.Replace(string(all), "admin", "cyvk", 1) // 替换admin为cyvk

		newReader := bytes.NewReader([]byte(replace))
		request.Body = io.NopCloser(newReader) // 重新设置body

		fmt.Println("修改后的请求体: " + replace)
		fmt.Println(len(replace))
		var buffer bytes.Buffer
		request.ContentLength = int64(len(replace)) // 设置Content-Length长度就是[]body的长度
		_ = request.Write(&buffer)                  // 将请求对象读成字节流
		data := buffer.Bytes()                      // 将请求对象读成字节流
		hrg.HttpReqData.Data = data                 // 直接修改字段

		hrg.AnnotationsText = &BurpApi.AnnotationsText{ // 设置注解
			IsInfo: true,
			Notes:  "修改请求",
			Color:  BurpApi.HighlightColor_RED,
		}
		fmt.Println(string(data)) // 打印修改后的请求
		return &BurpApi.HttpRequestAction{
			Continue:     false,
			IsReviseReq:  true, // 修改请求
			HttpReqGroup: hrg,  // 返回修改后的请求
		}, nil
	}
	return &BurpApi.HttpRequestAction{Continue: true}, nil
}

// HttpHandleResponseReceived 响应处理  不处理返回Continue: true
func (httpReqFlowDemo) HttpHandleResponseReceived(c context.Context, httpReqAndRes *BurpApi.HttpReqAndRes) (*BurpApi.HttpResponseAction, error) {
	fmt.Println("响应")
	return &BurpApi.HttpResponseAction{Continue: true}, nil
}

// HttpFlowHandler http修改请求理由
type httpReqFlowRoutingDemo struct {
	BurpApi.UnimplementedHttpFlowHandlerServer
}

// HttpHandleRequestReceived 请求处理 修改请求理由
func (httpReqFlowRoutingDemo) HttpHandleRequestReceived(c context.Context, hrg *BurpApi.HttpReqGroup) (*BurpApi.HttpRequestAction, error) {
	reader := bytes.NewReader(hrg.GetHttpReqData().GetData()) // 请求的字节数组
	request, err := http.ReadRequest(bufio.NewReader(reader)) // 转为请求实例
	if err != nil {
		log.Println(err)
		return &BurpApi.HttpRequestAction{Continue: true}, nil // 出错打印错误信息后返回继续
	}
	url := hrg.GetHttpReqData().GetUrl()                                             // 获取URL
	if strings.Contains(url, "login") && strings.EqualFold(request.Method, "POST") { // 请求方法是POST
		all, err := io.ReadAll(request.Body) // 读出请求体
		if err != nil {
			log.Println(err)
			return &BurpApi.HttpRequestAction{Continue: true}, nil
		}
		replace := strings.Replace(string(all), "admin", "cyvk", 1) // 替换admin为cyvk

		newReader := bytes.NewReader([]byte(replace))
		request.Body = io.NopCloser(newReader) // 重新设置body

		fmt.Println("修改后的请求体: " + replace)
		fmt.Println(len(replace))
		var buffer bytes.Buffer
		request.ContentLength = int64(len(replace)) // 设置Content-Length长度就是[]body的长度
		_ = request.Write(&buffer)                  // 将请求对象读成字节流
		data := buffer.Bytes()                      // 将请求对象读成字节流
		hrg.HttpReqData.Data = data                 // 直接修改字段

		hrg.HttpReqData.HttpReqService = &BurpApi.HttpReqService{ // 请求的服务器信息
			Ip:     "www.baidu.com",
			Port:   443,
			Secure: true, // 是否用TLS 加密传输
		}

		hrg.AnnotationsText = &BurpApi.AnnotationsText{ // 设置注解
			IsInfo: true,
			Notes:  "修改请求路由到百度",
			Color:  BurpApi.HighlightColor_RED,
		}
		fmt.Println(string(data)) // 打印修改后的请求
		return &BurpApi.HttpRequestAction{
			Continue:     false,
			IsReviseReq:  true, // 修改请求
			HttpReqGroup: hrg,  // 返回修改后的请求
		}, nil
	}
	return &BurpApi.HttpRequestAction{Continue: true}, nil
}

// HttpHandleResponseReceived 响应处理  不处理返回Continue: true
func (httpReqFlowRoutingDemo) HttpHandleResponseReceived(c context.Context, httpReqAndRes *BurpApi.HttpReqAndRes) (*BurpApi.HttpResponseAction, error) {

	fmt.Println("响应")
	return &BurpApi.HttpResponseAction{Continue: true}, nil
}

// HttpFlowHandler http修改请求路由创造响应
type httpReqFlowCreateResponseDemo struct {
	BurpApi.UnimplementedHttpFlowHandlerServer
}

// HttpHandleRequestReceived 请求处理 http修改请求路由创造响应
func (httpReqFlowCreateResponseDemo) HttpHandleRequestReceived(c context.Context, hrg *BurpApi.HttpReqGroup) (*BurpApi.HttpRequestAction, error) {
	fmt.Println("进入")
	if hrg.GetHttpReqData().GetHttpReqService().GetIp() == "8.8.8.8" { // 不存在的域名
		hrg.HttpReqData.HttpReqService = &BurpApi.HttpReqService{ // 修改原路由
			Ip:     "127.0.0.1",
			Port:   8848,
			Secure: false,
		}
		return &BurpApi.HttpRequestAction{ // 修改请求
			Continue:     false,
			IsReviseReq:  true,
			HttpReqGroup: hrg,
		}, nil
	}
	return &BurpApi.HttpRequestAction{
		Continue:     true,
		IsReviseReq:  false,
		HttpReqGroup: nil,
	}, nil
}

// HttpHandleResponseReceived 响应处理  不处理返回Continue: true
func (httpReqFlowCreateResponseDemo) HttpHandleResponseReceived(c context.Context, httpReqAndRes *BurpApi.HttpReqAndRes) (*BurpApi.HttpResponseAction, error) {

	fmt.Println("响应")
	return &BurpApi.HttpResponseAction{Continue: true}, nil
}

func main() {

	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	//BurpApi.RegisterHttpFlowHandlerServer(server, httpReqFlowDemo{})

	//BurpApi.RegisterHttpFlowHandlerServer(server, httpReqFlowRoutingDemo{})

	BurpApi.RegisterHttpFlowHandlerServer(server, httpReqFlowCreateResponseDemo{})
	fmt.Println("服务启动")
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
