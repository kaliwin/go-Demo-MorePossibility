package main

import (
	"context"
	"fmt"
	"go-Demo-MorePossibility/MorePossibility/proto/BurpApi" // burp生成地址
	"google.golang.org/grpc"
	"net"
	"strings"
)

// Demo 测试
type Demo struct {
	BurpApi.UnimplementedProxyRequestHandlerServer // 继承默认实现
}

// ProxyHandleRequestReceived 重写代理请求处理
func (Demo) ProxyHandleRequestReceived(c context.Context, httpReq *BurpApi.HttpReqGroup) (*BurpApi.ProxyRequestAction, error) {
	url := httpReq.GetHttpReqData().GetUrl() // 获取url
	if strings.Contains(url, "cyvk") {       // 是否包含cyvk
		return &BurpApi.ProxyRequestAction{
			Drop: true, // 为true 就是丢弃该请求
		}, nil
	}
	if strings.Contains(url, "test") { // 是否包含test
		return &BurpApi.ProxyRequestAction{
			Continue: true, // 继续
			HttpReqGroup: &BurpApi.HttpReqGroup{
				AnnotationsText: &BurpApi.AnnotationsText{
					IsInfo: true,                          // 为true 表示修改注解
					Notes:  "注解测试",                        // 文本
					Color:  BurpApi.HighlightColor_YELLOW, // 颜色
				},
			},
		}, nil
	}
	return &BurpApi.ProxyRequestAction{
		Continue: true, // 继续没有任何别的动作
	}, nil
}
func main() {
	listen, err := net.Listen("tcp", ":9000") // 监听端口
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()                                // 创建默认服务 明文传输
	BurpApi.RegisterProxyRequestHandlerServer(server, Demo{}) // 注册服务实现相关接口即可
	fmt.Println("启动服务")
	err = server.Serve(listen) // 启动服务
	if err != nil {
		panic(err)
	}
}
