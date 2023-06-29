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
	"strings"
)

// proxyReqDemo 代理请求设置注解测试
type proxyReqDemo struct {
	BurpApi.UnimplementedProxyRequestHandlerServer
}

// ProxyHandleRequestReceived 添加注解
func (proxyReqDemo) ProxyHandleRequestReceived(c context.Context, httpReqGroup *BurpApi.HttpReqGroup) (*BurpApi.ProxyRequestAction, error) {
	if strings.Contains(httpReqGroup.GetHttpReqData().GetUrl(), "cyvk") { // 判断url中是否包含cyvk
		return &BurpApi.ProxyRequestAction{ // 返回信息
			Continue: true, // 继续
			HttpReqGroup: &BurpApi.HttpReqGroup{ // 请求组 因为注解要在里面设置
				HttpReqData: nil, // http请求 因为我们继续了没有要修改请求就给nil
				AnnotationsText: &BurpApi.AnnotationsText{ // 注解
					IsInfo: true,                          // true 表示有注解
					Notes:  "url有cyvk",                    // 注解文本
					Color:  BurpApi.HighlightColor_YELLOW, // 颜色
				},
			},
		}, nil
	}
	return &BurpApi.ProxyRequestAction{Continue: true}, nil // 不包含直接返回Continue: true 即可
}

// proxyReqDropDemo 丢弃请求
type proxyReqDropDemo struct {
	BurpApi.UnimplementedProxyRequestHandlerServer
}

// ProxyHandleRequestReceived 丢弃请求
func (proxyReqDropDemo) ProxyHandleRequestReceived(c context.Context, httpReqGroup *BurpApi.HttpReqGroup) (*BurpApi.ProxyRequestAction, error) {

	reader := bytes.NewReader(httpReqGroup.GetHttpReqData().GetData())

	request, err := http.ReadRequest(bufio.NewReader(reader)) //将字节流转为http请求实例
	if err != nil {
		log.Println(err)
		return &BurpApi.ProxyRequestAction{Continue: true}, nil // 报错退出
	}

	if strings.Contains(request.Host, "baidu.com") { // 看host是否包含baidu.com
		return &BurpApi.ProxyRequestAction{Drop: true}, nil // 丢弃该请求
	}
	return &BurpApi.ProxyRequestAction{Continue: true}, nil // 不包含直接返回Continue: true 即可
}

// proxyReqReviseDemo 修改请求
type proxyReqReviseDemo struct {
	BurpApi.UnimplementedProxyRequestHandlerServer
}

// ProxyHandleRequestReceived 修改请求
func (proxyReqReviseDemo) ProxyHandleRequestReceived(c context.Context, httpReqGroup *BurpApi.HttpReqGroup) (*BurpApi.ProxyRequestAction, error) {

	if strings.Contains(httpReqGroup.GetHttpReqData().GetUrl(), "nacos") { // 包含nacos
		data := httpReqGroup.GetHttpReqData().GetData()
		request, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(data))) // 将字节数组转为请求实例
		if err != nil {
			log.Println(err)
			return &BurpApi.ProxyRequestAction{Continue: true}, nil // 不包含直接返回Continue: true 即可
		}
		// 设置请求头
		request.Header.Set("accessToken", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJuYWNvcyIsInRpbWUiOiIxNzc5NDcwMjAwIn0.7ig_ZTKgi7HTckMxdLK3yJm0yACNFsxzRHN4JQddwv4")

		// 将修改后的请求转换回[]byte
		var buffer bytes.Buffer
		_ = request.Write(&buffer)
		modifiedRequestBytes := buffer.Bytes()

		reqData := httpReqGroup.GetHttpReqData() // 获取原始的请求实例
		reqData.Data = modifiedRequestBytes      // 修改原来的数据

		return &BurpApi.ProxyRequestAction{ //
			Continue:    false,
			Drop:        false,
			IsReviseReq: true, // 修改请求
			IsIntercept: false,
			HttpReqGroup: &BurpApi.HttpReqGroup{
				HttpReqData: reqData, // 返回新数据  如果你需要更复杂的操作比如将请求发给其他服务器或强制使用ssl 需要自己创造HttpReqData填写你要的字段因为
				// 这里只是改请求数据其他参数不变就把原来的拿来改下数据就可以了
				AnnotationsText: nil,
			},
		}, nil
	}
	return &BurpApi.ProxyRequestAction{Continue: true}, nil // 不包含直接返回Continue: true 即可
}

// proxyReqReviseDemo 修改请求
type proxyReqInterceptDemo struct {
	BurpApi.UnimplementedProxyRequestHandlerServer
}

// ProxyHandleRequestReceived 拦截请求
func (proxyReqInterceptDemo) ProxyHandleRequestReceived(c context.Context, httpReqGroup *BurpApi.HttpReqGroup) (*BurpApi.ProxyRequestAction, error) {
	reader := bytes.NewReader(httpReqGroup.GetHttpReqData().GetData())
	request, err := http.ReadRequest(bufio.NewReader(reader)) // 将[]byte 转为请求实例 便于操作
	if err != nil {
		log.Println(err)
		return &BurpApi.ProxyRequestAction{Continue: true}, nil
	}
	// 方法是POST并且url包含login
	if strings.EqualFold(request.Method, "POST") && strings.Contains(httpReqGroup.GetHttpReqData().GetUrl(), "login") {
		return &BurpApi.ProxyRequestAction{
			Continue:     false,
			Drop:         false,
			IsReviseReq:  false,
			IsIntercept:  true,
			HttpReqGroup: httpReqGroup, // 你可以修改请求 不改就把他返回就行
		}, nil
	}
	return &BurpApi.ProxyRequestAction{Continue: true}, nil // 不包含直接返回Continue: true 即可
}

func main() {
	listen, err := net.Listen("tcp", ":9000") // 监听端口
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer() // 创建一个grpc服务 默认是明文传输
	//BurpApi.RegisterProxyRequestHandlerServer(server, proxyReqDemo{}) 		// 添加注解
	//BurpApi.RegisterProxyRequestHandlerServer(server, proxyReqDropDemo{}) 	// 丢弃请求
	//BurpApi.RegisterProxyRequestHandlerServer(server, proxyReqReviseDemo{})   // 修改请求
	BurpApi.RegisterProxyRequestHandlerServer(server, proxyReqInterceptDemo{}) // 拦截请求
	// 有命令规范 BurpApi.Register开头之后便是你服务的名称以Server结尾
	fmt.Println("启动服务")
	err = server.Serve(listen)

	if err != nil {
		panic(err)
	}
}
