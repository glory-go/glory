##  grpc/http服务的封装

服务开启大致流程：

glory服务端启动的大体流程为：根据上述/config/service_config.go中字段的描述，我们将服务端可以抽象为一个一个Service，这些Service在main函数执行时候，根据init阶段读入的配置文件和用户实现的代码，生成好对应的service实例，注册在中心化的Server上，再由gloryServer逐一启动各个实力化Service，从而使得server和每个Service之间实现解耦，增强框架的扩展性和可插拔性。

服务客户端启动过程比服务端简单，只需要在init的时候，由框读入配置，在用户函数中直接调用框架提供的接口来根据配置初始化client，再经过必要的封住后在代码中调用client即可。

由于同一协议的服务端和客户端是相对的，所以在本部分，将按照框架启动的先后顺序，穿插介绍客户端和服务端启动的操作。

### 2.1 配置的读入

- 客户端

对于一个glory服务，如果存在主调其他服务的客户端，在yaml配置文件中应当存在如下配置例子：

```yaml
consumer:
  "grpc-helloworld-demo": # 客户端Key
    registry_key: registryKey # 注册中心标识
    config_source: file # 配置文件读入方式
    service_id: GoOnline-IDE-gloryService # 注册服务名
    protocol: grpc # 协议名
```

可看到，例子描述了一个grpc客户端，目前无需了解其他字段的意义。只需要关心协议名和客户端Key即可。协议名标注了本客户端所使用的协议，用于之后协议抽象时进行区分。客户端Key则用于在开发者书写的代码中，使用这个Key找到它对应的配置。

将配置读入后，框架将配置保存在格式化的结构体内。

- 服务端

同理，针对服务端也是如此

```yaml
provider:
  "gloryGrpcService": # 服务端key
    protocol: grpc
    registry_key: registryKey 
    service_id: GoOnline-IDE-gloryService
    port: 8080 # 暴露端口
```

将配置读入框架后，开始执行main函数。

### 2.2 客户端/服务端的加载

以grpc服务举例

- 客户端grpc服务的启动

  ```go
  func main() {
  	client := grpc.NewGrpcClient("dev", "grpc-helloworld-demo")
  	greeterClient := NewGreeterClient(client.GetConn())
  	...
  }
  ```

  在用户代码main函数中，调用框架提供的借口，将特定环境、特定客户端key的服务建立起来，拿到grpc协议的client。NewGrpcClient中封装了针对根据配置对官方grpcclient的调用，从而拿到封装好的官方客户端的client。

  之后，根据特定协议的需求，进一步使用pb文件的NewGreeterClient函数，引入我们建立好链接的client，拿到封装好grpc协议的实例化client。

  在后面，就可以根据client直接调用业务函数了。

- 服务端grpc服务的启动

  ```go
  func main() {
  	gloryServer := glory.NewServer()
  	gloryService := service.NewGrpcService("gloryGrpcService")
  }
  ```

  首先生成上面提到的gloryServer，然后再新建一个抽象的service，在新建的过程中用到了之前提到的服务端Key，根据这个key框架可以找到目标配置，从而按照协议初始化一个服务端Service。

### 2.3 server端grpc服务的实例化、注册到服务上

服务端业务代码的实现：

```go
// server is used to implement helloworld.GreeterServer.
type server struct {
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &HelloReply{Message: "Hello " + in.GetName()}, nil
}
```

由开发者撰写业务函数代码，封装到server 结构体内

```go
func main() {
  ...
	RegisterGreeterServer(gloryService.GetGrpcServer(), &server{})
	gloryServer.RegisterService(gloryService)
	gloryServer.Run()
}
```

拿到了gloryService之后，按照grpc官方库的要求，在pb上注册好server链接，并传入实现结构。

再将Service注册于server，最终启动框架server。启动时，会逐个启动所有注册在server上的servive。