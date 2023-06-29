package main

import (
	"context"
	"fmt"
	"go-Demo-MorePossibility/MorePossibility/proto/BurpApi"
	"google.golang.org/grpc"
	"log"
	"net"
	"regexp"
	"strings"
)

type ReqEditBoxDemo struct {
	BurpApi.UnimplementedHttpReqEditBoxAssistServer
}

// ReqHttpEdit 处理渲染信息
func (ReqEditBoxDemo) ReqHttpEdit(c context.Context, heb *BurpApi.HttpEditBoxData) (*BurpApi.ByteData, error) {
	compile, err := regexp.Compile("(?:\"|')(((?:[a-zA-Z]{1,10}://|//)[^\"'/]{1,}\\.[a-zA-Z]{2,}[^\"']{0,})|((?:/|\\.\\./|\\./)[^\"'><,;|*()(%%$^/\\\\\\[\\]][^\"'><,;|()]{1,})|([a-zA-Z0-9_\\-/]{1,}/[a-zA-Z0-9_\\-/]{1,}\\.(?:[a-zA-Z]{1,4}|action)(?:[\\?|#][^\"|']{0,}|))|([a-zA-Z0-9_\\-/]{1,}/[a-zA-Z0-9_\\-/]{3,}(?:[\\?|#][^\"|']{0,}|))|([a-zA-Z0-9_\\-]{1,}\\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)(?:[\\?|#][^\"|']{0,}|)))(?:\"|')")
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	// 正则匹配最多命中100个
	allString := compile.FindAllString(string(heb.GetHttpReqAndResData().GetRes().GetData()), 100)
	newData := strings.Join(allString, "\n")                 // 将字符串切片中的元素连接成一个字符串 每个元素用\n隔开
	return &BurpApi.ByteData{ByteData: []byte(newData)}, nil // 返回命中的信息
}

// IsReqHttpEditFor 是否要渲染
func (ReqEditBoxDemo) IsReqHttpEditFor(c context.Context, heb *BurpApi.HttpEditBoxData) (*BurpApi.Boole, error) {
	// 只要url包含nacos 就返回true
	if strings.Contains(heb.GetHttpReqAndResData().GetReq().GetUrl(), "nacos") {
		return &BurpApi.Boole{Boole: true}, nil
	}
	return &BurpApi.Boole{Boole: false}, nil
}

func main() {

	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	BurpApi.RegisterHttpReqEditBoxAssistServer(server, ReqEditBoxDemo{})

	fmt.Println("服务启动")
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}

}
