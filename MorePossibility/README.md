



# GO-Demo-MorePossibility



官方文档 : https://grpc.io/docs/languages/go/quickstart/ 



## 先前条件



需要下载protoc 工具用于解析proto生成出对应语言的文件

请参考官方文档下载  (国内请求会很慢)

```
 go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
 go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

将burpApi.proto 放到一个文件夹



![image-20230628142809640](README.assets/image-20230628142809640.png)



运行该命令就会在当前目录创建名为BurpApi的文件夹

包含用于proto序列化反序列化的文件和一个Grpc服务的文件   直接用就行

```
 protoc --proto_path=.  --go_out=.  --go-grpc_out=. burpApi.proto 
```

![image-20230628142659701](README.assets/image-20230628142659701.png)



文件夹中包含这两个文件就行   第一次创建需要引入依赖IDE会自动完成

![image-20230628142846822](README.assets/image-20230628142846822.png)









### proxy

####  代理请求处理器

```
当burp代理模块收到请求时便会调用该函数 该函数可用于

1、丢弃请求	(丢弃请求)
2、修改请求	(修改请求内容 路由地址 比如将baidu.com 的包发到 jd.com)
3、拦截请求 （在拦截面板中显示 给用户二次确认）
4、添加注解  (添加注解和颜色用于标识特别的请求)
```

服务签名

```
// 代理请求处理器 提供请求修改拦截
service ProxyRequestHandler{
  rpc ProxyHandleRequestReceived(httpReqGroup) returns(ProxyRequestAction){} // 代理请求处理
}
```

入参

```
// 单个请求组 用于包含请求信息
message httpReqGroup{
  httpReqData httpReqData = 1;
  annotationsText annotationsText = 2 ;
}

// http 单个请求 包含头部和体
message httpReqData {
  bytes Data = 1;                       //请求包数据
  string Url = 2;                       //请求url
  int64 BodyIndex = 3;                  //请求体开始下标,如若请求体下标与总长度一致就是没有请求体
  string httpVersion = 4;               // http版本
  httpReqService httpReqService = 5;    // 目标信息
}
```

回参

需要做的动作设置为true

除非丢弃请求和继续、修改请求和拦截即便没有改也要设置httpReqGroup 需要注意有请求体要设置Content-Length

```
// 收到代理请求的操作
message ProxyRequestAction{
  bool continue = 1; // 继续不做任何处理
  bool drop = 2 ;  // 丢弃请求
  bool isReviseReq = 3;  // 修改请求
  bool isIntercept = 4 ;  // 是否拦截   无论是否修改都要返回数据也就是httpReqData 字段
  httpReqGroup httpReqGroup = 5;
}
```



###### 场景演示



首先需要组合一个基本实现    (可以理解为继承但是go没有继承概念只有组合 ，组合之后重写函数即可)

这个结构体是burpApi.proto 生成的包含服务的基本实现我们重写函数就行

技巧: 	使用U加你要组合的服务名ProxyRequestHandler 之后IDE会给提示

![image-20230628144651321](README.assets/image-20230628144651321.png)



Cttrl + 点击 到结构体声明的地方

![image-20230628145015311](README.assets/image-20230628145015311.png)



将这个函数复制过来

![image-20230628145135438](README.assets/image-20230628145135438.png)



复制出来为以下

![image-20230628145210940](README.assets/image-20230628145210940.png)

修改一下 这样我们便可在ProxyHandleRequestReceived 函数中编写我们的具体代码了

![image-20230628145251070](README.assets/image-20230628145251070.png)





###### 添加注解

检查请求url中包含cyvk添加注释并设置颜色

Continue: true 继续就是当前请求不变 你可以添加注解也可以不添加注解类型IsInfo: true 就是添加

```go
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
```

启动服务

```go
func main() {
	listen, err := net.Listen("tcp", ":9000") // 监听端口
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()                                        // 创建一个grpc服务 默认是明文传输
	BurpApi.RegisterProxyRequestHandlerServer(server, proxyReqDemo{}) // 注册服务 将server放入还有实现这个服务接口的结构体放入
	// 有命令规范 BurpApi.Register开头之后便是你服务的名称以Server结尾
	fmt.Println("启动服务")
	err = server.Serve(listen)

	if err != nil {
		panic(err)
	}
}
```

服务启动后在插件这注册服务 选择Proxy请求处理器 Grpc地址  名字随便写

![image-20230628151008485](README.assets/image-20230628151008485.png)



符合我们的预期

![image-20230628151348450](README.assets/image-20230628151348450.png)



###### 丢弃请求

将请求的数据[]byte 转为request对象 在通过host判断是否包含baidu.com 包含就丢弃

注意动作是在返回的参数中 你返回Drop: true 就会丢弃 , Continue: true 就是继续

```go
// ProxyHandleRequestReceived 丢弃请求
func (proxyReqDropDemo) ProxyHandleRequestReceived(c context.Context, httpReqGroup *BurpApi.HttpReqGroup) (*BurpApi.ProxyRequestAction, error) {
	reader := bytes.NewReader(httpReqGroup.GetHttpReqData().GetData())
	request, err := http.ReadRequest(bufio.NewReader(reader))    //将字节流转为http请求实例
	if err != nil {
		log.Println(err)
		return &BurpApi.ProxyRequestAction{Continue: true}, nil  // 报错退出
	}
	if strings.Contains(request.Host, "baidu.com") {    // 看host是否包含baidu.com
		return &BurpApi.ProxyRequestAction{Drop: true}, nil   // 丢弃该请求
	}
	return &BurpApi.ProxyRequestAction{Continue: true}, nil // 不包含直接返回Continue: true 即可
}
```

编写完成后要重启服务 , burp插件不用重启

![image-20230628152745952](README.assets/image-20230628152745952.png)

百度的包被丢弃了 

实例化了ReadRequest 你可以修改请求从go程序发出

![image-20230628152853074](README.assets/image-20230628152853074.png)



###### 修改请求

url中有nacos的就添加一个默认token     ( 这里是在做演示 我知道头是Authorization: Bearer ！！！)

accessToken: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJuYWNvcyIsInRpbWUiOiIxNzc5NDcwMjAwIn0.7ig_ZTKgi7HTckMxdLK3yJm0yACNFsxzRHN4JQddwv4

```go
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
				AnnotationsText: nil,  // 想要注解就给字段
			},
		}, nil
	}
	return &BurpApi.ProxyRequestAction{Continue: true}, nil // 不包含直接返回Continue: true 即可
}
```

注意这个点 你如果修改了原始请求就会出现这个可以切换的按钮

![image-20230628160929323](README.assets/image-20230628160929323.png)



添加了我们想要的字段

![image-20230628161029046](README.assets/image-20230628161029046.png)



url没有nacos的就不会做修改也就没有哪个切换按钮了

![image-20230628161217991](README.assets/image-20230628161217991.png)

你可以在logger这看到实际发出的请求

![image-20230628161059428](README.assets/image-20230628161059428.png)

###### 拦截请求

拦截post请求并且url包含login 的请求



```go
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
```



可以看到即便你的拦截开关是关的也能拦截到 可以用于精准拦截

![image-20230628164611927](README.assets/image-20230628164611927.png)



#### 代理响应处理器



？？？？  漏了就漏了吧 累了不想写了





#### 	总结



```
要实现代理请求 你需要实现服务中的接口函数并且监听端口开启Grpc服务 然后在插件中填写Proxy请求处理器 Grpc的地址 和名称注册即可
修改自身服务不用重启插件服最多抛出几个异常不会影响burp等你服务改完在重新启动便会恢复正常

该功能用于精确拦截、丢弃无用请求、修改请求 
```



### http编辑框

#### 代理响应处理器

```
当proxy收到响应时会调用该函数 与请求一样

1、设置注解
2、丢弃响应
3、修改响应
4、拦截响应
```



服务签名

```
// 代理响应处理器 提供对相应的处理
service ProxyResponseHandler{
  rpc ProxyHandleResponseReceived(httpReqAndRes)returns(ProxyResponseAction){}  // 代理响应处理
}
```

入参

```
一组http数据包含请求和响应 以及注解


// 一组http请求 (请求和响应)
message httpReqAndRes {
  httpReqData req = 1;
  httpResData res = 2;
  httpInfo info = 3;
  annotationsText annotationsText = 4;
}
```

回参

```
与请求一样 

// 收到代理响应的操作
message ProxyResponseAction{
  bool continue = 1; // 继续不做任何处理
  bool drop = 2 ;  // 丢弃响应
  bool isReviseRes = 3;  // 修改响应
  bool isIntercept = 4 ;  // 是否拦截  无论是否修改都要返回数据也就是httpReqData 字段
  httpResGroup httpResGroup = 5;
}
```





###### 场景演示



与请求一样组合基本实现结构体

![image-20230628171351336](README.assets/image-20230628171351336.png)

编写业务代码即可

![image-20230628171832624](README.assets/image-20230628171832624.png)



###### 添加注解

响应状态码是404的就添加注解

```go
// ProxyResDemo 代理响应Demo
type ProxyResDemo struct {
	BurpApi.UnimplementedProxyResponseHandlerServer
}
// ProxyHandleResponseReceived 代理响应处理器 注解测试
func (ProxyResDemo) ProxyHandleResponseReceived(c context.Context, hr *BurpApi.HttpReqAndRes) (*BurpApi.ProxyResponseAction, error) {
	if hr.GetRes().GetStatusCode() == 404 {   // 状态码等于404
		return &BurpApi.ProxyResponseAction{
			Continue:    true,
			Drop:        false,
			IsReviseRes: false,
			IsIntercept: false,
			HttpResGroup: &BurpApi.HttpResGroup{AnnotationsText: &BurpApi.AnnotationsText{
				IsInfo: true,      // 添加注解
				Notes:  "404啦",
				Color:  BurpApi.HighlightColor_GREEN,
			}},
		}, nil
	}
	return &BurpApi.ProxyResponseAction{Continue: true}, nil
}
func main() {
	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	BurpApi.RegisterProxyResponseHandlerServer(server, ProxyResDemo{})
	fmt.Println("服务启动")
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
```

注册服务

![image-20230628172529372](README.assets/image-20230628172529372.png)

和预期一样

![image-20230628172738664](README.assets/image-20230628172738664.png)



###### 丢弃响应

丢弃响应体有 page 的

```go
// ProxyHandleResponseReceived 代理响应处理器 注解测试
func (ProxyResDropDemo) ProxyHandleResponseReceived(c context.Context, hr *BurpApi.HttpReqAndRes) (*BurpApi.ProxyResponseAction, error) {

	if strings.Contains(string(hr.Res.Data), "page") {
		fmt.Println("丢弃")
		return &BurpApi.ProxyResponseAction{Drop: true}, nil
	}
	return &BurpApi.ProxyResponseAction{Continue: true}, nil
}	
```

![image-20230628181222870](README.assets/image-20230628181222870.png)



后续与请求一致就不写了





#### http请求编辑框

​	

```
点击数据时可以有个组件用于渲染信息
```

就是这个组件

![image-20230629095538622](README.assets/image-20230629095538622.png)



服务签名

```
需要实现一个服务两个函数  首先burp会调用IsReqHttpEditFor 返回值是个布尔类型用于判断是否要渲染显示true的话就会调用ReqHttpEdit 返回的是字节数组用于渲染
要注意的是  IsReqHttpEditFor 会频繁的反复的调用 因此IsReqHttpEditFor只是判断是否要渲染不要执行压力过大的程序


// 提供请求编辑框
service HttpReqEditBoxAssist{
  rpc ReqHttpEdit (HttpEditBoxData) returns (ByteData){}    // 返回数据渲染
  rpc IsReqHttpEditFor(HttpEditBoxData) returns (Boole){}   // 是否要渲染
}
```

入参



```
两个函数入参一致包含的是编辑器名称和一组数据


// http编辑框数据, 名字
message HttpEditBoxData{
  string name = 1; // 名字 编辑器名称
  httpReqAndRes httpReqAndResData = 2; // 一组http数据 请求 / 响应
  string info = 3 ; // 预留字段
}
```



回参



```
IsReqHttpEditFor 要返回一个布尔类型 用于表示后续 ReqHttpEdit 函数是否要执行

ReqHttpEdit 返回的是字节数组用于渲染
```



###### 场景演示





###### 渲染关键信息



使用正则匹配目标中的敏感信息 

只要url中有nacos的就执行匹配

```react
(?:"|')(((?:[a-zA-Z]{1,10}://|//)[^"'/]{1,}\.[a-zA-Z]{2,}[^"']{0,})|((?:/|\.\./|\./)[^"'><,;|*()(%%$^/\\\[\]][^"'><,;|()]{1,})|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{1,}\.(?:[a-zA-Z]{1,4}|action)(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{3,}(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-]{1,}\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)(?:[\?|#][^"|']{0,}|)))(?:"|')
```



```go
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
	newData := strings.Join(allString, "\n") // 将字符串切片中的元素连接成一个字符串 每个元素用\n隔开
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
```



服务注册

![image-20230629105434578](README.assets/image-20230629105434578.png)



只要url有nacos的就会有reqTest 名字就是你注册服务的名字

![image-20230629105405861](README.assets/image-20230629105405861.png)







#### http响应编辑框



```
响应和请求是一样的无非是一个在左边的请求框一个在右键的响应框  这里就不演示了
```



#### 总结



```
编辑框用于格式化数据 提取关键数据 渲染十分有用 
还能用于加解密处理 假如请求体中的数据是加密的而你已经知道了解密的方式那你就可以用这个功能将明文渲染 这样可以更直观的分析流量
```



### http流量处理器



```
http流量处理器用于处理burp全局流量 
与proxy 请求 / 响应 处理器一样能够修改请求和响应 包括代理流量 重放器流量 其他插件流量 都能处理

因为burp中注册http处理器要实现的接口就是两个函数 所以我也把他放在一个服务里不像代理那样请求和响应分开

如果你只是对请求或响应做处理 另一个也要实现不做处理就返回继续就行
```



#### 场景演示





服务签名

```
要实现一个服务两个函数 即便你用不到也要实现


// http流量处理器 请求和响应放一起
service HttpFlowHandler{
  rpc HttpHandleRequestReceived(httpReqGroup) returns(HttpRequestAction){}
  rpc HttpHandleResponseReceived(httpReqAndRes) returns(HttpResponseAction){}
}

```



入参



```

请求处理器只要请求组  这里的流量不能丢弃只能改或者不改  如果你想要在这丢弃可以将路由指定一个不存在的ip端口
响应处理器有一组请求 一样的这个包你不能丢弃 只能改或者不改

注意: 修改请求体要自己更新Content-Length 插件和burp不会在对你返回的请求做任何处理 如果你不处理Content-Length服务器会认为是个畸形的请求会被丢弃
计算body的字节长度更新就行 

// 单个请求组 用于包含 不同服务请求中的附带信息
message httpReqGroup{
  httpReqData httpReqData = 1;
  annotationsText annotationsText = 2 ;
}



// 一组http请求 (请求和响应)
message httpReqAndRes {
  httpReqData req = 1;
  httpResData res = 2;
  httpInfo info = 3;
  annotationsText annotationsText = 4;
}
```



回参

```

只能继续或者修改   


// http请求操作
message HttpRequestAction{
  bool continue = 1; // 继续不做任何处理
  bool isReviseReq = 3;  // 修改请求
  httpReqGroup httpReqGroup = 4;
}

// http响应操作
message HttpResponseAction{
  bool continue = 1; // 继续不做任何处理
  bool isReviseRes = 3;  // 修改请求
  httpResGroup httpResGroup = 4;
}
```



#### http请求处理器



###### 修改请求体



将POST请求 并且url 包含login的请求体进行修改

我们要把请求体中的admin替换为cyvk  之后添加注解

http的注解可以在代理和logger中显示

![image-20230629112224629](README.assets/image-20230629112224629.png)





```go
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

		hrg.AnnotationsText = &BurpApi.AnnotationsText{  // 设置注解
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
```



服务注册

![image-20230629130834548](README.assets/image-20230629130834548.png)



请求看上去是admin实际发出的是cyvk

![image-20230629130856611](README.assets/image-20230629130856611.png)



在Logger中才能看到实际丢出的是改过的包 并且有注解

![image-20230629130945649](README.assets/image-20230629130945649.png)



代理也是会被处理   所有burp的流量都能被处理  但是代理这没有切换按钮只有在Logger中能看到实际发出的请求和实际接收到的响应

![image-20230629131037104](README.assets/image-20230629131037104.png)



可以看到Content-Length 被正确处理了

![image-20230629131341476](README.assets/image-20230629131341476.png)



######  请求重定向



将POST请求 并且url 包含login的请求体进行修改

我们要把请求体中的admin替换为cyvk  之后添加注解

http的注解可以在代理和logger中显示   之后修改路由点将包发给百度



就是把包发给百度

![image-20230629142341622](README.assets/image-20230629142341622.png)



只要添加服务器信息就行 入参中的HttpReqGroup 请求组有这些信息

![image-20230629142731188](README.assets/image-20230629142731188.png)



重定向到了百度

![image-20230629142714586](README.assets/image-20230629142714586.png)



###### 请求重定向 无中生响应



场景: 比如一个APP你能抓到包但是他有些功能是内网ip你无法访问但是你知道要给他返回什么响应就可以将原本无法访问到地址改为可以访问的之后在修改响应即可

http://8.8.8.8/       这是一个不存在的域名或者无法访问的ip       我们要修改他的路由让他有数据

想要丢弃就设置一个访问不了的地址



```go
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
```



正常这是无法访问的 我们把go写的Grpc服务开起来在插件里注册好

![image-20230629144328342](README.assets/image-20230629144328342.png)

 

这样就可以了

![image-20230629144505166](README.assets/image-20230629144505166.png)



![image-20230629144614498](README.assets/image-20230629144614498.png)





#### http 响应处理器



```
响应与请求一样 要么继续要么改数据 
```



服务签名

```
// http流量处理器 请求和响应放一起
service HttpFlowHandler{
  rpc HttpHandleRequestReceived(httpReqGroup) returns(HttpRequestAction){}
  rpc HttpHandleResponseReceived(httpReqAndRes) returns(HttpResponseAction){}
}
```



###### 修改响应状态码



大致逻辑其实都是一样的这里就演示修改404为200状态码



```go
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
		response, err := http.ReadResponse(reader, nil)   // 将字节流解析为响应体对象
		if err != nil {
			panic(err)
		}
		response.StatusCode = 200     // 设置状态码
		response.Status = ""
		var buffer bytes.Buffer
		_ = response.Write(&buffer)     

		res.Data = buffer.Bytes()    // 将响应对象转回字节数组

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
```

![image-20230629151027137](README.assets/image-20230629151027137.png)



#### 总结

```
流量处理器最好用于重放器这处理     ????  
```




### 上下文菜单提供程序

```
也就是右键菜单 用于对特定请求响应包做特殊的事情
```



两个服务首先是注册上下文提供程序 burp会调用GetConTextMenuItems 该函数获取到菜单用于渲染

之后点击的事件由ContextMenuItemsProvider 处理

这样说比较难以理解我们来看 逻辑说明

```
// 获取上下文菜单项服务器  发送菜单的名字
service GetConTextMenuItemsServer {
  rpc GetConTextMenuItems(Str)returns(MenuInfo){}
}



// 右键菜单项执行程序
service ContextMenuItemsProvider{
  rpc MenuItemsProvider(ContextMenuItems) returns(MenuItemsReturn){}  // 菜单项处理
}
```



入参

```
入参是你注册的名称

//名称 String类型
message Str{
  string name = 1;
}
```

回参

```
回参用于渲染  tarGet是你Grpc的地址
menu 是一个菜单 他可以包含多个菜单和菜单项    菜单项是你点击的事件

// 一个上下文菜单项信息
message MenuInfo{
  string tarGet = 2; // 服务地址往后菜单项都用这个地址服务
  Menu menu = 3;   // 菜单项
}

// 菜单 菜单下可包含多个菜单和菜单项
message Menu{
  string name = 1;
  repeated Menu menuList = 2;
  repeated MenuItem menuItemList = 3;
}
```



#### 逻辑说明

菜单就是下面还有菜单和菜单项           菜单项是真正可以点击的

可以看到 knife  是一个菜单 他下面有很多菜单项和菜单   insert Payload For ALL 这个菜单下面有两个菜单项 是可以点击的

![image-20230629152618993](README.assets/image-20230629152618993.png)





演示

```
// ConTextMenuItemsDemo 上下文提供程序
type ConTextMenuItemsDemo struct {
	BurpApi.UnimplementedGetConTextMenuItemsServerServer
}

// GetConTextMenuItems 获取菜单项
func (ConTextMenuItemsDemo) GetConTextMenuItems(c context.Context, str *BurpApi.Str) (*BurpApi.MenuInfo, error) {

	return &BurpApi.MenuInfo{
		TarGet: "127.0.0.1:9000", // 地址
		Menu: &BurpApi.Menu{ // 一个菜单  一个菜单可以有多个菜单和菜单项
			Name:     "第一个菜单",
			MenuList: nil, // 菜单列表
			MenuItemList: []*BurpApi.MenuItem{ // 菜单项列表
				&BurpApi.MenuItem{Name: "第一个菜单项呀"},
				&BurpApi.MenuItem{Name: "第二个菜单项呀"},
			},
		},
	}, nil
}
```



注册服务  burp 就会发起请求 获取菜单用于渲染

![image-20230629153506956](README.assets/image-20230629153506956.png)



可以自行对比

![image-20230629153615899](README.assets/image-20230629153615899.png)



较为复杂的结构 你理解不了就用最上面那个多建立几个菜单项就行



```go
// ConTextMenuItemsDemo 上下文提供程序
type ConTextMenuItemsDemo struct {
	BurpApi.UnimplementedGetConTextMenuItemsServerServer
}

// GetConTextMenuItems 获取菜单项
func (ConTextMenuItemsDemo) GetConTextMenuItems(c context.Context, str *BurpApi.Str) (*BurpApi.MenuInfo, error) {

	return &BurpApi.MenuInfo{
		TarGet: "127.0.0.1:9000", // 地址
		Menu: &BurpApi.Menu{ // 一个菜单  一个菜单可以有多个菜单和菜单项
			Name: "第一个菜单",
			MenuList: []*BurpApi.Menu{
				&BurpApi.Menu{
					Name:     "第一下的第一个菜单",
					MenuList: nil,
					MenuItemList: []*BurpApi.MenuItem{
						&BurpApi.MenuItem{Name: "第一个菜单项"},
					},
				},
				&BurpApi.Menu{
					Name:         "第一下的第二个菜单",
					MenuList:     nil,
					MenuItemList: nil,
				},
			}, // 菜单列表
			MenuItemList: []*BurpApi.MenuItem{ // 菜单项列表
				&BurpApi.MenuItem{Name: "第一个菜单项呀"},
				&BurpApi.MenuItem{Name: "第二个菜单项呀"},
			},
		},
	}, nil
}
```



![image-20230629154216937](README.assets/image-20230629154216937.png)





#### 菜单项点击事件处理



渲染后 点击事件会调用MenuItemsProvider grpc地址就是你获取菜单时填写的



服务签名

```
// 右键菜单项执行程序
service ContextMenuItemsProvider{
  rpc MenuItemsProvider(ContextMenuItems) returns(MenuItemsReturn){}  // 菜单项处理
}
```



入参

当你点击时burp会发起rpc调用  

名称就是你菜单项的名称 ， 选中数据是你右键时有没有划到数据 然后就是一组请求和响应



```go
// 上下文处理器 , 提供 一组请求响应 和选中的数据
message ContextMenuItems{
  string  name = 1;                  //名称   点击的菜单项名称
  bool isSelect = 2;                 // 是否有选中的数据
  HttpSource selectSource = 3;       // 选中来源 只有选中数据才能知道从哪点的
  httpReqAndRes httpReqAndRes = 4;   // 一组请求和响应
  bytes selectData = 5;              // 用户选中的数据
  SubscriptOffsets selectOffsets = 6; // 选中的下标偏移量
}
```



没有的话isSelect就为false

![image-20230629154846058](README.assets/image-20230629154846058.png)



#### 修改选中的数据



将选中的消息替换为Log4j2的测试载荷

首先要先实现GetConTextMenuItems 来获取菜单项渲染

然后再实现菜单项处理来处理点击事件

${jndi:ldap://127.0.0.1:1389/xxx}



```go
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
				&BurpApi.MenuItem{Name: "log4j2"},
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

```





![image-20230629164615861](README.assets/image-20230629164615861.png)





![image-20230629164606541](README.assets/image-20230629164606541.png)





#### 修改请求包



对请求体进行base64编码



```go
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
		body := reqData[index:]        // 截取请求体
		base64String := base64.StdEncoding.EncodeToString(body)   // 进行base64编码

		newReq := append(reqData[:index], []byte(base64String)...) // 组装新的请求

		return &BurpApi.MenuItemsReturn{
			IsContinue:     false,
			IsReviseSelect: false,
			IsReviseReq:    true,      // 修改请求
			IsReviseRes:    false,
			ReqData:        newReq,   // 新的请求数据
			ResData:        nil,
			SelectDate:     nil,
		}, nil
	}
	fmt.Println("没有请求体")
	return &BurpApi.MenuItemsReturn{IsContinue: true}, nil
}
```





我不好截图 这是原来的右键点击bash64



![image-20230629170314854](README.assets/image-20230629170314854.png)



点击后

![image-20230629170345341](README.assets/image-20230629170345341.png)



符合预期

![image-20230629170357049](README.assets/image-20230629170357049.png)







目前修改响应是不能用的但是官方的API文档中存在这种函数先留着等待官网后续更新





### 迭代器

```
迭代器 常用与用户名密码穷举目录穷举和载荷穷举

一共两个功能 一个是迭代处理器 用于对载荷做处理 一般常用于对载荷的加密编码

还有一个是生成器 用于生成用于迭代的载荷

```







#### 迭代处理器

服务签名

```
// 迭代处理器
service IntruderPayloadProcessorServer {
  // 迭代器载荷处理
  rpc IntruderPayloadProcessor(ByteData)returns(ByteData){}
}
```



入参

```
// 字节数据
message ByteData {
  bytes byteData = 1;
}
```

回参

```
// 字节数据
message ByteData {
  bytes byteData = 1;
}
```





###### 对载荷进行base64编码



```go
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
```



注册服务 名字不要中文不要特殊符号最好就是大小写字符的组合就行

![image-20230629172022237](README.assets/image-20230629172022237.png)



选择插入点

![image-20230629172015944](README.assets/image-20230629172015944.png)

选择 处理器 名字就是你注册服务的名字



![image-20230629172157742](README.assets/image-20230629172157742.png)



符合预期



![image-20230629172303243](README.assets/image-20230629172303243.png)





#### 迭代生成器

```
生成器就是 载荷的生成器 比如你可以通过host域名和之前收集到的信息 生成一个针对性字典

入参有整个http请求和插入点列表 
burp会并发不断调用你的函数直到 IsEnd=true 才能停下
```



服务签名

```
// 迭代器 生成器
service IntruderPayloadGeneratorServer{
  // 迭代器载荷生成, 服务端流 流终止就为生成结束
  rpc IntruderPayloadGeneratorProvider(IntruderGeneratorData) returns (PayloadGeneratorResult){}
}
```



入参

```
//迭代生成器 数据
message IntruderGeneratorData{
  bytes contentData = 1; // http 请求模板
  repeated  InsertionPointOffsets  insertionPointOffsets = 2;  // 插入点列表
  bytes IntruderInsertionPoint = 3;  //插入点基值
}
```



回参

```
// 载荷生成器结果
message PayloadGeneratorResult{
  bytes byteData = 1; // 生成的载荷
  bool IsEnd = 2;     // 是否结束 true 为结束
}
```





###### 字典生成构造



将域名拼接在字典前面即可



```go
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
			IsEnd:    true,   // 结束
		}, nil
	}
	reqData := igd.GetContentData()
	request, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(reqData))) // 解析请求
	if err != nil {
		log.Println(err)
		return &BurpApi.PayloadGeneratorResult{
			ByteData: nil,
			IsEnd:    true,  // 出错结束
		}, nil
	}
	host := request.Host

	str := host + dictionary[i.index]   // host拼接字典
	i.index++
	return &BurpApi.PayloadGeneratorResult{
		ByteData: []byte(str),
		IsEnd:    false,
	}, nil
}
```



服务注册



![image-20230629174522239](README.assets/image-20230629174522239.png)



下拉选 选中这个

![image-20230629174550399](README.assets/image-20230629174550399.png)



点击选择你的载荷生成器

![image-20230629174617428](README.assets/image-20230629174617428.png)

符合预期 



![image-20230629174703469](README.assets/image-20230629174703469.png)







#### 总结

```
生成器和处理器可以一起工作

生成器可以用于动态构造载荷 目录、用户名、密码等 目前只能注册一个生成器获取处理器 多个只能换端口如果有人使用后面我在改吧

```







### 实时流量镜像



```
该功能是MorePossibility 最核心的功能也是最初设计的初心, MorePossibility将为史诗庇护所提供数据来源 

该服务会将burp的所有流量以一组http请求的方式发送给调用方 用于记录burp流量或用于异步处理分析 

```



服务签名

```
// 实时流量镜像
service RealTimeTrafficMirroring{
  /*
      实时流量传输
      burp将主动建立连接通过客户端流进行实时流量镜像
    */
  rpc RealTimeTrafficMirroring(stream httpReqAndRes) returns (Str){}
}
```



注册服务后会建立客户端流 burp作为客户端 





```
package main

import (
	"fmt"
	"go-Demo-MorePossibility/MorePossibility/proto/BurpApi"
	"google.golang.org/grpc"
	"net"
)

// RealTimeTrafficMirroringDemo 实时流量镜像
type RealTimeTrafficMirroringDemo struct {
	BurpApi.UnimplementedRealTimeTrafficMirroringServer
}

// RealTimeTrafficMirroring 实时流量镜像
func (RealTimeTrafficMirroringDemo) RealTimeTrafficMirroring(e BurpApi.RealTimeTrafficMirroring_RealTimeTrafficMirroringServer) error {
	fmt.Println("[+] 建立实时流量通道")
	for true {
		recv, err := e.Recv()
		if err != nil {
			panic(err)
		}
		req := recv.GetReq()
		res := recv.GetRes()
		fmt.Println(string(req.GetData())) // 打印请求
		fmt.Println()
		fmt.Println(string(res.GetData())) // 打印响应
	}
	return nil
}

func main() {
	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	BurpApi.RegisterRealTimeTrafficMirroringServer(server, RealTimeTrafficMirroringDemo{})
	fmt.Println("服务启动")
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}

```





所有流量都会转发出来





![image-20230629190246704](README.assets/image-20230629190246704.png)
