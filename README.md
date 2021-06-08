# Glory —— 一款轻量级微服务框架
**请移步至[Glory/glory](https://github.com/glory-go/glory)获取最新版本**
## 1. release note

## 2. 功能介绍

### 2.1 总体介绍

glory框架为一款轻量级微服务框架，你可以使用它快速开发你的服务端。

- 在协议方面：glory框架提供grpc（client端和server端）、http（server端）脚手架，你可以通过几行配置和几行代码快速开启一个grpc或http服务。

  glory框架还提供具备服务治理能力的**glory协议**，通过glory协议暴露的服务，在整个rpc过程中，框架可提供**服务注册、服务发现、集群策略、负载均衡**等服务治理能力，并且支持基于glory协议的**流式rpc**。

- 在配置方面：glory框架提供**统一化的配置服务**，你只需要在main文件同级目录config文件夹下定义glory.yml，在配置文件内按照约定格式写入配置信息，在引入框架后执行时，框架会自动读入配置文件，并开启所需服务。

  你也可以选择从阿里云 nacos **配置中心拉取**当前服务所需配置。

- 日志：glory框架提供**日志服务**，你可以在配置文件中定义自己需要的日志记录方式。支持命令行、文件、远程（基于elastic、阿里云sls）的日志收集。

- 链路追踪：glory框架提供适配于 **grpc 的链路追踪**服务，你可以选择将服务调用链路上报至阿里云链路追踪平台进行监控和错误追溯。

- 数据上报：glory框架提供**基于promethus的数据上报**服务，你可以在配置文件中定义自己需要的数据上报方式，支持基于promethus-pushgateway的推式数据上报，以及传统拉式数据上报。

- 第三方工具常用sdk支持：glory框架提供mysql、redis、mongodb、oss、rabbitmq等常见工具的sdk封装，开发者可以在配置中引入服务，使用框架提供的sdk进行快速开发。

### 2.2 各模块详细介绍

//这部分待不断补充

- [统一配置服务](./config/README.md)

    - 配置中心拉取

- [单个特定协议Service启动](./service/README.md)
    - 开启grpc服务
    
      - grpc-server
      - rpc-client
    
    - 开启http服务
    
      - http-server
        - [http filter链实现](./http/README.md)
    - Triple (dubbo3) 协议和网络模型接入

- [使用glory协议实现RPC](./protocol/glory/README.md)

  - 服务注册发现
    - nacos
    - [k8s](./registry/k8s/README.md)
    - redis
  - 负载均衡
    - round_robin
    - random（后续支持）
  - 集群策略

- 日志

    - 阿里云sls
    - elastic

- [数据上报(基于 prometheus)](./metrics/README.md)

- 链路追踪

    - [grpc - 阿里云链路追踪收集](./filter/filter_impl/README.md)

- 数据库

  - redis
  - mysql
  - mongodb
  
- [分层错误码](./error/README.md)

- oss 对象存储

  - qiniu sdk

- 消息

  - rabbitmq

  

## 3. quick start 

### 3.1 手把手带你开启一次grpc调用

##### 3.1.1 server端

1. 定义IDL (接口描述语言）：

   新建sever文件夹，server文件夹下新建helloworld.proto

   ```protobuf
   syntax = "proto3";
   
   package main;
   
   // The greeting service definition.
   service Greeter {
     // Sends a greeting
     rpc SayHello (HelloRequest) returns (HelloReply) {}
   }
   
   // The request message containing the user's name.
   message HelloRequest {
     string name = 1;
   }
   
   // The response message containing the greetings
   message HelloReply {
     string message = 1;
   }
   ```

   执行

   `$ protoc --go_out=plugins=grpc:. *.proto`

   同级目录下生成helloworld.pb.go文件

2. 撰写配置文件

   server/config/glory.yaml

   ```yml
   org_name: ide
   server_name: grpc-demo-server
   log : # 日志配置
     "console-log":
       log_type: console # 命令行输出日志
       level: debug # 日志等级
   
   provider:
     "gloryGrpcService": # service 名称可自定义
       protocol: grpc # 通过grpc暴露
       service_id: GoOnline-IDE-gloryService # service 唯一ID，用于服务注册
       port: 8080 # 暴露端口
   ```

3. main.go文件

   ```go
   package main
   
   import (
   	"context"
   
       // 开启框架服务必须引入
   	"github.com/glory-go/glory/glory"
       // 使用日志组件必须引入
   	"github.com/glory-go/glory/log"
       // 注册service必须引入
   	"github.com/glory-go/glory/service"
   )
   
   // server is used to implement helloworld.GreeterServer.
   type server struct {
   }
   
   // SayHello implements helloworld.GreeterServer
   func (s *server) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
   	log.Info("Received: %v", in.GetName())
   	return &HelloReply{Message: "Hello " + in.GetName()}, nil
   }
   
   func main() {
   	gloryServer := glory.NewServer()
   	gloryService := service.NewGrpcService("gloryGrpcService")
   	RegisterGreeterServer(gloryService.GetGrpcServer(), &server{})
   	gloryServer.RegisterService(gloryService)
   	gloryServer.Run()
   }
   
   
   ```

4. 拉取依赖

   `$ go mod init glory-grpc-server-demo `

   `$ export GOPROXY="http://goproxy.io"`

   `$ export GOPRIVATE="git.go-online.org.cn"`

   `$ go get`

   可能会要求输入密码，因为框架属于私有仓库  

   go get成功后如果ide报错，尝试重启goland

   如果还是爆红，尝试在goland-setting-go-gomodules-environment 配置上面两个环境变量。

5. 运行服务

   `$ go run .`

   可看到控制台输出：

   ```text
   $ go run .
   org_name: ide
   server_name: grpc-demo-server
   yamlFile: log : # 日志配置
     "console-log":
       log_type: console # 命令行输出日志
       level: debug # 日志等级
   
   provider:
     "gloryGrpcService": # service 名称可自定义
       protocol: grpc # 通过grpc暴露
       service_id: GoOnline-IDE-gloryService # service 唯一ID，用于服务注册
       port: 8080 # 暴露端口
   grpc start listening on :8080

   ```
   
   代表grpc服务启动成功

#### 3.1.2 client端

1. server同级目录下建立client文件夹，定义IDL (接口描述语言）helloworld.proto，文件内容与server端完全一致。相同方法编译生成.pb.go文件。

2. 撰写配置文件

   client/config/glory.yaml

   ```yml
   org_name: ide
   server_name: grpc-demo-client
   log : # 日志配置
     "console-log":
       log_type: console # 命令行输出日志
       level: debug # 日志等级
   
   consumer :
     "grpc-helloworld-demo":
       service_id: GoOnline-IDE-gloryService
       server_address: 127.0.0.1:8080
       protocol: grpc
   ```

3. main.go

   ```go
   package main
   
   import (
   	"context"
   
   	// grpc客户端需引入
   	"github.com/glory-go/glory/grpc"
   	// 框架日志组件引入
   	"github.com/glory-go/glory/log"
   )
   
   func main() {
   	// 从配置生成grpc客户端，与配置中serviceName对应
   	client := grpc.NewGrpcClient("grpc-helloworld-demo")
       
       // 与协议文件结合，拿到greeterClient
   	greeterClient := NewGreeterClient(client.GetConn())
       
       // 发起rpc调用，传递参数
   	reply, err := greeterClient.SayHello(context.Background(), &HelloRequest{
   		Name: "grpcDemo",
   	})
   	if err != nil {
   		panic(err)
   	}
       // 打印结果
   	log.Infof("reply = %+v", reply)
   }
   
   ```

4. 和server端完全一样，拉取依赖

   `$ go mod init glory-grpc-client-demo`

   `$ export GOPROXY="http://goproxy.io"`

   `$ export GOPRIVATE="git.go-online.org.cn"`

   `$ go get`

   可能会要求输入密码，因为框架属于私有仓库  

   go get成功后如果ide报错，尝试重启goland

   如果还是爆红，尝试在goland-setting-go-gomodules-environment 配置上面两个环境变量。

5. 运行服务

   `$ go run .`

   可看到控制台输出：

   ```text
   $ go run .
   org_name: ide
   server_name: grpc-demo-client
   yamlFile: log : # 日志配置
     "console-log":
       log_type: console # 命令行输出日志
       level: debug # 日志等级
   
   consumer :
     "grpc-helloworld-demo":
       service_id: GoOnline-IDE-gloryService
       server_address: 127.0.0.1:8080
       protocol: grpc
   2020-12-13T22:33:16.330+0800    info    grpc/client.go:16       [[{8080 127.0.0.1}]]
   
   2020-12-13T22:33:16.438+0800    info    client/main.go:19       reply = [[message:"Hello grpcDemo" ]]

   ```
   
   rpc调用成功

### 3.2 运行一个简易http-server

config/glory.yaml

```yaml
org_name: ide
server_name: http-demo-server
log : # 日志配置
  "console-log":
    log_type: console # 命令行输出日志
    level: debug # 日志等级

provider:
  "httpDemo":
    protocol: http
    service_id: GoOnline-IDE-gloryService
    port: 8080
```



main.go

```go
package main

import (
	"github.com/glory-go/glory/glory"
	ghttp "github.com/glory-go/glory/http"
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/service"
)

// 定义 request 结构体，支持 validate 标签校验，具体语法参考：https://godoc.org/github.com/go-playground/validator
type gloryHttpReq struct {
	Input    []int  `schema:"input" validate:"required"` // query参数使用schema 标签
	BodyStr  string `json:"body_str"`                    // body 参数使用json标签
	BodyStr2 string `json:"body_str_2"`                  // body 参数使用json标签

}

// 定义 response 结构体
type gloryHttpRsp struct {
	Output int `schema:"output"`
}

// 自定义业务逻辑处理 handler
func testHandler(controller *ghttp.GRegisterController) error {
	req := controller.Req.(*gloryHttpReq)
	rsp := controller.Rsp.(*gloryHttpRsp)
	log.Info("req = ", *req)                                                                     // 打印query和body参数
	log.Info("hello = ", controller.VarsMap["hello"], "hello2 = ", controller.VarsMap["hello2"]) // 打印path内变量
	rsp.Output = req.Input[0] + 1
	return nil
}

func main() {
	gloryServer := glory.NewServer()
	// 与 yaml文件中的key保持一致
	httpService := service.NewHttpService("httpDemo")
	// 注册http服务注册：path、method、handler、bodySturcture、filter...
	httpService.RegisterRouter("/testwithfilter/{hello}/{hello2}", testHandler, &gloryHttpReq{}, &gloryHttpRsp{}, "POST")
	// 注册service到glory服务
	gloryServer.RegisterService(httpService)
	// 开启glory server
	gloryServer.Run()
	// 使用postman测试
}

```
开启服务并测试

![1607870851747](./img/1607870851747.png)

