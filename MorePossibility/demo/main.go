package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"go-Demo-MorePossibility/MorePossibility/proto/BurpApi"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

type ServerTest struct {
	BurpApi.UnimplementedHttpFlowHandlerServer           // http 流量处理
	BurpApi.UnimplementedGetConTextMenuItemsServerServer // 获取菜单项
	BurpApi.UnimplementedContextMenuItemsProviderServer  // 菜单项处理

}

// HttpHandleRequestReceived 请求
func (ServerTest) HttpHandleRequestReceived(c context.Context, reqGroup *BurpApi.HttpReqGroup) (*BurpApi.HttpRequestAction, error) {
	start := time.Now()
	fmt.Println("http流量处理器 请求调用")

	fmt.Println("注解: " + reqGroup.GetAnnotationsText().GetNotes())

	reqData := reqGroup.HttpReqData

	fmt.Println(fmt.Sprintf("ip :%s 端口 : %d   是否用安全协议 %t", reqData.HttpReqService.Ip, reqData.HttpReqService.Port, reqData.HttpReqService.Secure))

	//reqGroup

	ret := &BurpApi.HttpRequestAction{
		Continue:    true,
		IsReviseReq: false,
		HttpReqGroup: &BurpApi.HttpReqGroup{
			HttpReqData: &BurpApi.HttpReqData{
				Data: nil,
				HttpReqService: &BurpApi.HttpReqService{
					Ip:     "",
					Port:   0,
					Secure: false,
				},
			},
			AnnotationsText: &BurpApi.AnnotationsText{
				IsInfo: true,
				Notes:  "请求注解",
				Color:  BurpApi.HighlightColor_GREEN,
			},
		},
	}
	elapsed := time.Since(start) // 计算经过的时间

	fmt.Printf("The function took %v to run.\n", elapsed)
	return ret, nil
}

// GetConTextMenuItems 获取菜单项
func (ServerTest) GetConTextMenuItems(context.Context, *BurpApi.Str) (*BurpApi.MenuInfo, error) {

	//item := BurpApi.MenuItem{Name: "sd"}

	//items := []BurpApi.MenuItem{item}
	fmt.Println("获取菜单项")
	items := []*BurpApi.MenuItem{{Name: "第一个处理"}, {Name: "第二个处理"}}

	ret := &BurpApi.MenuInfo{
		TarGet: "127.0.0.1:9000",
		Menu: &BurpApi.Menu{ // 菜单项
			Name:         "第一个处理程序",
			MenuList:     nil,
			MenuItemList: items,
		},
	}

	return ret, nil
}

// MenuItemsProvider 菜单项处理
func (ServerTest) MenuItemsProvider(c context.Context, mi *BurpApi.ContextMenuItems) (*BurpApi.MenuItemsReturn, error) {

	fmt.Println("接收到右键处理请求: " + mi.GetName())

	source := mi.GetSelectSource()

	switch *source.Enum() {
	case BurpApi.HttpSource_Request:
		fmt.Println("从请求点来的")

	case BurpApi.HttpSource_Response:
		fmt.Println("从响应点来的")

	}

	if mi.GetIsSelect() { // 用户是否有选中的数据

		fmt.Println("选中: " + string(mi.GetSelectData()))
	}

	newReq := append(mi.GetHttpReqAndRes().Req.Data, []byte("go 菜单项处理")...)

	newRes := append(mi.GetHttpReqAndRes().Res.Data, []byte("go 响应测试")...)

	return &BurpApi.MenuItemsReturn{
		IsReviseRes: true,
		ReqData:     newReq,
		ResData:     newRes,
	}, nil
}
func (ProxyHandleRequest) ProxyHandleRequestReceived(c context.Context, h *BurpApi.HttpReqGroup) (*BurpApi.ProxyRequestAction, error) {
	fmt.Println("代理请求处理")

	return &BurpApi.ProxyRequestAction{
		Continue:    true,
		Drop:        false,
		IsReviseReq: false,
		IsIntercept: false,
		HttpReqGroup: &BurpApi.HttpReqGroup{
			HttpReqData: nil,
			AnnotationsText: &BurpApi.AnnotationsText{
				IsInfo: true,
				Notes:  "注解测试",
				Color:  BurpApi.HighlightColor_BLUE,
			},
		},
	}, nil
}

type ProxyHandleRequest struct {
	BurpApi.UnimplementedProxyRequestHandlerServer
	BurpApi.UnimplementedProxyResponseHandlerServer
	BurpApi.UnimplementedHttpResEditBoxAssistServer
}

func (ProxyHandleRequest) ProxyHandleResponseReceived(c context.Context, reqAndres *BurpApi.HttpReqAndRes) (*BurpApi.ProxyResponseAction, error) {

	resData := string(reqAndres.Res.Data)
	compile, err := regexp.Compile("(?:\"|')(((?:[a-zA-Z]{1,10}://|//)[^\"'/]{1,}\\.[a-zA-Z]{2,}[^\"']{0,})|((?:/|\\.\\./|\\./)[^\"'><,;|*()(%%$^/\\\\\\[\\]][^\"'><,;|()]{1,})|([a-zA-Z0-9_\\-/]{1,}/[a-zA-Z0-9_\\-/]{1,}\\.(?:[a-zA-Z]{1,4}|action)(?:[\\?|#][^\"|']{0,}|))|([a-zA-Z0-9_\\-/]{1,}/[a-zA-Z0-9_\\-/]{3,}(?:[\\?|#][^\"|']{0,}|))|([a-zA-Z0-9_\\-]{1,}\\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)(?:[\\?|#][^\"|']{0,}|)))(?:\"|')")
	if err != nil {
		log.Println(err)
		return nil, nil
	}

	if compile.MatchString(resData) {
		return &BurpApi.ProxyResponseAction{
			Continue:    true,
			Drop:        false,
			IsReviseRes: false,
			IsIntercept: false,
			HttpResGroup: &BurpApi.HttpResGroup{
				HttpResData: nil,
				AnnotationsText: &BurpApi.AnnotationsText{
					IsInfo: true,
					Notes:  "命中正则规则",
					Color:  BurpApi.HighlightColor_GREEN,
				},
			},
		}, nil
	}

	return &BurpApi.ProxyResponseAction{
		Continue:     true,
		Drop:         false,
		IsReviseRes:  false,
		IsIntercept:  false,
		HttpResGroup: nil,
	}, nil
}

func (ProxyHandleRequest) ResHttpEdit(c context.Context, heb *BurpApi.HttpEditBoxData) (*BurpApi.ByteData, error) {

	data := string(heb.HttpReqAndResData.Res.Data)
	compile, err := regexp.Compile("(?:\"|')(((?:[a-zA-Z]{1,10}://|//)[^\"'/]{1,}\\.[a-zA-Z]{2,}[^\"']{0,})|((?:/|\\.\\./|\\./)[^\"'><,;|*()(%%$^/\\\\\\[\\]][^\"'><,;|()]{1,})|([a-zA-Z0-9_\\-/]{1,}/[a-zA-Z0-9_\\-/]{1,}\\.(?:[a-zA-Z]{1,4}|action)(?:[\\?|#][^\"|']{0,}|))|([a-zA-Z0-9_\\-/]{1,}/[a-zA-Z0-9_\\-/]{3,}(?:[\\?|#][^\"|']{0,}|))|([a-zA-Z0-9_\\-]{1,}\\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)(?:[\\?|#][^\"|']{0,}|)))(?:\"|')")
	if err != nil {
		log.Println(err)
		return nil, nil
	}

	allString := compile.FindAllString(data, 100)
	var str string
	for _, v := range allString {
		str += v + "\n"

	}

	return &BurpApi.ByteData{ByteData: []byte(str)}, nil
}
func (ProxyHandleRequest) IsResHttpEditFor(c context.Context, heb *BurpApi.HttpEditBoxData) (*BurpApi.Boole, error) {

	return &BurpApi.Boole{Boole: true}, nil
}

// Demo 测试汇总
type Demo struct {
	// 请求编辑框测试
	BurpApi.UnimplementedHttpReqEditBoxAssistServer
	// 响应编辑框测试
	BurpApi.UnimplementedHttpResEditBoxAssistServer
	// http 流量测试
	BurpApi.UnimplementedHttpFlowHandlerServer

	// 菜单项提供程序
	BurpApi.UnimplementedGetConTextMenuItemsServerServer

	// 菜单项处理程序
	BurpApi.UnimplementedContextMenuItemsProviderServer

	// 实时流量 镜像
	BurpApi.UnimplementedRealTimeTrafficMirroringServer

	// 代理请求处理器
	BurpApi.UnimplementedProxyRequestHandlerServer

	// 代理响应处理器
	BurpApi.UnimplementedProxyResponseHandlerServer

	// 迭代处理器
	BurpApi.UnimplementedIntruderPayloadProcessorServerServer

	// 迭代生成器
	BurpApi.UnimplementedIntruderPayloadGeneratorServerServer

	fileStringS []string
}

var i = 0

func (d Demo) IntruderPayloadGeneratorProvider(c context.Context, ig *BurpApi.IntruderGeneratorData) (*BurpApi.PayloadGeneratorResult, error) {
	fmt.Println(fmt.Sprintf("迭代生成器 : %d", i))
	if len(d.fileStringS) <= i {
		i = 0
		return &BurpApi.PayloadGeneratorResult{
			ByteData: []byte("end"),
			IsEnd:    true,
		}, nil
	} else {

		data := []byte(d.fileStringS[i])
		i++
		return &BurpApi.PayloadGeneratorResult{
			ByteData: data,
			IsEnd:    false,
		}, nil
	}

}

// IntruderPayloadProcessor 迭代处理器
func (Demo) IntruderPayloadProcessor(c context.Context, byteData *BurpApi.ByteData) (*BurpApi.ByteData, error) {
	fmt.Println(fmt.Sprintf("迭代处理: %s", string(byteData.GetByteData())))
	toString := base64.StdEncoding.EncodeToString(byteData.GetByteData())
	return &BurpApi.ByteData{ByteData: []byte(toString)}, nil
}

// ProxyHandleResponseReceived 代理响应测试
func (Demo) ProxyHandleResponseReceived(c context.Context, hr *BurpApi.HttpReqAndRes) (*BurpApi.ProxyResponseAction, error) {

	if hr.GetAnnotationsText().GetIsInfo() {
		fmt.Println("有注解: " + hr.GetAnnotationsText().GetNotes())
	} else {
		return &BurpApi.ProxyResponseAction{
			Continue:    true,
			Drop:        false,
			IsReviseRes: false,
			IsIntercept: false,
			HttpResGroup: &BurpApi.HttpResGroup{
				HttpResData: nil,
				AnnotationsText: &BurpApi.AnnotationsText{
					IsInfo: true,
					Notes:  "代理响应注解测试",
					Color:  BurpApi.HighlightColor_MAGENTA,
				},
			},
		}, nil
	}

	return &BurpApi.ProxyResponseAction{
		Continue:     true,
		Drop:         false,
		IsReviseRes:  false,
		IsIntercept:  false,
		HttpResGroup: nil,
	}, nil
}

// ProxyHandleRequestReceived 代理请求处理器
func (Demo) ProxyHandleRequestReceived(c context.Context, hr *BurpApi.HttpReqGroup) (*BurpApi.ProxyRequestAction, error) {

	if strings.Contains(hr.GetHttpReqData().GetUrl(), "cyvk") {
		return &BurpApi.ProxyRequestAction{
			Continue:    true,
			Drop:        false,
			IsReviseReq: false,
			IsIntercept: false,
			HttpReqGroup: &BurpApi.HttpReqGroup{
				HttpReqData: nil,
				AnnotationsText: &BurpApi.AnnotationsText{
					IsInfo: true,
					Notes:  "代理请求命中",
					Color:  BurpApi.HighlightColor_GRAY,
				},
			},
		}, nil
	}

	return &BurpApi.ProxyRequestAction{
		Continue:     true,
		Drop:         false,
		IsReviseReq:  false,
		IsIntercept:  false,
		HttpReqGroup: nil,
	}, nil
}

func (Demo) RealTimeTrafficMirroring(server BurpApi.RealTimeTrafficMirroring_RealTimeTrafficMirroringServer) error {
	for true {

		recv, err := server.Recv()
		if err != nil {
			if err == io.EOF {
				_ = server.SendAndClose(&BurpApi.Str{Name: "结束"})
			}
		}

		fmt.Println("接收到: " + recv.GetReq().GetUrl())

		if recv.GetAnnotationsText().GetIsInfo() {
			fmt.Println("有注解 " + recv.GetAnnotationsText().GetNotes())
		}

	}
	fmt.Println("结束")
	return nil
}

// GetConTextMenuItems 菜单项提供程序
func (Demo) GetConTextMenuItems(c context.Context, str *BurpApi.Str) (*BurpApi.MenuInfo, error) {
	fmt.Println("demo " + str.GetName())
	return &BurpApi.MenuInfo{
		TarGet: "127.0.0.1:9000",
		Menu: &BurpApi.Menu{
			Name: "ggi",
			MenuList: []*BurpApi.Menu{
				&BurpApi.Menu{
					Name:     "mu",
					MenuList: nil,
					MenuItemList: []*BurpApi.MenuItem{
						&BurpApi.MenuItem{Name: "mu1"},
					},
				},
			},
			MenuItemList: []*BurpApi.MenuItem{
				&BurpApi.MenuItem{Name: "svi"},
			},
		},
	}, nil
}

func (Demo) MenuItemsProvider(c context.Context, cm *BurpApi.ContextMenuItems) (*BurpApi.MenuItemsReturn, error) {

	fmt.Println("点击来源 " + cm.GetName())
	if cm.GetHttpReqAndRes().GetAnnotationsText().GetIsInfo() {
		fmt.Println("点击事件有注解  " + cm.GetHttpReqAndRes().GetAnnotationsText().GetNotes())
	}

	return &BurpApi.MenuItemsReturn{
		IsContinue:     false,
		IsReviseSelect: true,
		IsReviseReq:    false,
		IsReviseRes:    false,
		ReqData:        nil,
		ResData:        nil,
		SelectDate:     []byte("选中替换测试"),
	}, nil
}

// HttpHandleRequestReceived 请求流量处理
func (Demo) HttpHandleRequestReceived(c context.Context, hr *BurpApi.HttpReqGroup) (*BurpApi.HttpRequestAction, error) {
	fmt.Println("接收到请求 " + hr.GetHttpReqData().GetUrl())

	if hr.GetAnnotationsText().GetIsInfo() {
		fmt.Println("有注解 " + hr.GetAnnotationsText().GetNotes())
	}

	return &BurpApi.HttpRequestAction{
		Continue:    true,
		IsReviseReq: false,
		HttpReqGroup: &BurpApi.HttpReqGroup{
			HttpReqData: nil,
			AnnotationsText: &BurpApi.AnnotationsText{
				IsInfo: true,
				Notes:  "请求流量注解测试",
				Color:  BurpApi.HighlightColor_YELLOW,
			},
		},
	}, nil
}

// HttpHandleResponseReceived 响应流量处理
func (Demo) HttpHandleResponseReceived(c context.Context, hr *BurpApi.HttpReqAndRes) (*BurpApi.HttpResponseAction, error) {

	if hr.GetAnnotationsText().GetIsInfo() {
		fmt.Println("有注解 " + hr.GetAnnotationsText().GetNotes())
	}
	fmt.Println(hr.GetRes().GetStatusCode())
	return &BurpApi.HttpResponseAction{
		Continue:     true,
		IsReviseRes:  false,
		HttpResGroup: nil,
	}, nil
}

func (Demo) ResHttpEdit(c context.Context, he *BurpApi.HttpEditBoxData) (*BurpApi.ByteData, error) {
	fmt.Println("请求编辑框 : 渲染调用 ==== ")
	if he.GetHttpReqAndResData().GetAnnotationsText().GetIsInfo() {
		log.Println("请求编辑框 : 有注解 " + he.GetHttpReqAndResData().GetAnnotationsText().GetNotes())
	}
	bytes := append(he.GetHttpReqAndResData().GetReq().GetData(), []byte("响应编辑框: 渲染测试")...)

	return &BurpApi.ByteData{ByteData: bytes}, nil
}

func (Demo) IsResHttpEditFor(c context.Context, he *BurpApi.HttpEditBoxData) (*BurpApi.Boole, error) {
	fmt.Println("响应编辑框 : 是否渲染调用")
	if strings.Contains(he.GetHttpReqAndResData().GetReq().GetUrl(), "cyvk") {
		return &BurpApi.Boole{Boole: true}, nil
	}

	return &BurpApi.Boole{Boole: false}, nil
}

func (Demo) ReqHttpEdit(c context.Context, he *BurpApi.HttpEditBoxData) (*BurpApi.ByteData, error) {
	fmt.Println("请求编辑框 : 渲染调用 ==== ")
	if he.GetHttpReqAndResData().GetAnnotationsText().GetIsInfo() {
		log.Println("请求编辑框 : 有注解 " + he.GetHttpReqAndResData().GetAnnotationsText().GetNotes())
	}
	bytes := append(he.GetHttpReqAndResData().GetRes().GetData(), []byte("请求编辑框 渲染测试")...)

	return &BurpApi.ByteData{ByteData: bytes}, nil
}
func (Demo) IsReqHttpEditFor(c context.Context, he *BurpApi.HttpEditBoxData) (*BurpApi.Boole, error) {
	fmt.Println("请求编辑框 : 是否渲染调用")
	if strings.Contains(he.GetHttpReqAndResData().GetReq().GetUrl(), "nacos") {
		return &BurpApi.Boole{Boole: true}, nil
	}

	return &BurpApi.Boole{Boole: false}, nil
}

func main() {

	dial, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()

	file, err := os.ReadFile("/root/tmp/Top100.txt")
	if err != nil {
		panic(err)
	}

	split := strings.Split(string(file), "\n")

	demo := Demo{fileStringS: split}

	BurpApi.RegisterHttpReqEditBoxAssistServer(server, demo) // 请求编辑框
	BurpApi.RegisterHttpResEditBoxAssistServer(server, demo) // 响应编辑框

	BurpApi.RegisterHttpFlowHandlerServer(server, demo) // http流处理器

	BurpApi.RegisterGetConTextMenuItemsServerServer(server, demo) // 菜单项提供程序

	BurpApi.RegisterContextMenuItemsProviderServer(server, demo) // 菜单项处理程序

	BurpApi.RegisterRealTimeTrafficMirroringServer(server, demo) // 实时流量镜像 舍弃

	BurpApi.RegisterProxyRequestHandlerServer(server, demo) // 代理请求处理

	BurpApi.RegisterProxyResponseHandlerServer(server, demo) // 代理响应处理

	BurpApi.RegisterIntruderPayloadProcessorServerServer(server, demo) // 迭代处理器

	BurpApi.RegisterIntruderPayloadGeneratorServerServer(server, demo) // 迭代生成器

	fmt.Println("服务启动")
	err1 := server.Serve(dial)
	if err1 != nil {
		panic(err1)
	}

}
