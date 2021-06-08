# <center>GOLORY 协议介绍和使用</center>


[toc]

## 1. glory

### 1.1 glory协议字段

我希望我的glory框架需要拥有自己的服务治理能力：服务发现、负载均衡、默认glory协议、hessian2序列化工具支持、流式rpc支持。

因此我参考jsonrpc 2.0规范 设计了glory协议

请求包、回包、错误码
```go
// RequestPackage is sent from client to server
type RequestPackage struct {
	Header     *Header       // Header is glory header
	MethodName string        // MethodName is method to invoke
	Params     []interface{} // Params is param list
}

// ResponsePackage is sent from server to client
type ResponsePackage struct {
	Header *Header       // Header is glory header
	Result []interface{} // Result is response result list
	Error  *Error        // Error defined the response status
}

// Error is a field of rsp pkg
type Error struct {
	Code int32  // Code shows glory error code
	Msg  string // Msg shows error message
}

```

针对glory协议头，它可以在日后的开发和完善中进行扩展，目前根据我的设计，拥有如下字段：

```go
// Header is glory network header field,
type Header struct {
	Version    string // Version is glory version
	TraceID    string // TraceID is used to get all link way trace usage
	PkgType    uint64 // PkgType define the package type
	Seq        uint64 // Seq is the request seq number
	ChanOffset uint8 // ChanOffset is used when stream channel send data, define which target Chan sends
}
```

其中，trace ID用于之后请求链路追踪，标识一次rpc调用的id。

MethodName表示请求调用的方法名。

PkgType目前包含以下类型：

```go
const (
    PingPkg           PkgType = 0
    PongPkg           PkgType = 1
    NormalRequestPkg  PkgType = 2
    NormalResponsePkg PkgType = 3
    StreamRequestPkg  PkgType = 4
    StreamSendPkg     PkgType = 5
    StreamRecvPkg     PkgType = 6
    StreamReadyPkg    PkgType = 7
    ErrorRspPkg       PkgType = 8
)
```

用于标识当前RPC是流式RPC还是普通rpc，以及是请求包还是返回包。

seq是请求序列号，用于pending请求，针对单个链接的多次请求服务进行区分，并且是目前支持的负载均衡算法roundrobin的依赖字段。

ChanOffset字段为流式RPC设计，会在之后用于流式RPC的channe查找。

以上glory协议基础字段可以基本保证rpc请求的完成。

### 1.2 Hessian2序列化

hessian2序列化极大程度保证了go-java的互通性。并且hessian2-go被应用于dubbo-go开源框架，并且保证和dubbo-java的连通性。

但是由于已有的hessian2框架并不支持流式RPC调用，这也是我本次毕业设计最大的贡献：在hessian2和glory协议的基础之上，实现流式调用。

就目前而言，我希望从unary（普通）RPC调用开始实现，之后再实现流式RPC接口。

Hessian2用法可见：https://github.com/apache/dubbo-go-hessian2

针对需要序列化的任何go struct，只需要定义好如下ID：

```go
func (Object) JavaClassName() string {
	return "com.company.Circular"
}
```

再将其结构以POJO的格式注册到hessian2框架上：

`hessian.RegisterPOJO(&GloryHeader{})`

即可直接通过Encode  API进行序列化和反序列化。无需预定义和编译IDL，并方便兼容java服务，十分方便。

我选择hessian2而不是pb，还有一个重要原因是它支持动态生成Interface的注册和序列化，使得预编译型的框架支持称为可能：即可以通过框架编译出来的二进制文件和用户定义的配置文件，直接实现服务的调用，无需编写代码。

很荣幸，这个feature是我提出并完成的：https://github.com/apache/dubbo-go-hessian2/pull/243

## 2. server端开启普通rpc服务过程实现

### 2.1 配置文件定义暴露接口

glory.yaml

```yaml
provider:
  "gloryService": # 服务名
    protocol: glory # 协议名
    registry_key: registryKey # 注册中心名 （可选）
    service_id: GoOnline-IDE-gloryService # 用于注册的服务发现ID（可选）
    port: 8080 # 本地暴露端口
```

定义了如上配置，即可在服务main函数中，引入框架，进行provider端服务的实例话。

### 2.2 服务实例化、和注册

1. 开发者可以在代码中， 定义需要暴露服务的业务代码（示例）：

```go
type GloryProvider struct {
}

func (g GloryProvider) SayHello(ctx context.Context, req *ReqBody, str2 string) (*RspBody, error) {
	log.Info("req = ", *req, "+", str2)
	fmt.Println(time.Now().String())
	return &RspBody{
		Value:     req.Value,
		SeqNum:    req.SeqNum + 1,
		TimeStamp: time.Now(),
	}, nil
}

```

可以看到，GloryProvider实现了一个业务函数，传输过程包括ReqBody、RspBody两个自定义参数。

2. 为了保证正确序列化，需要将两个参数注册到hessian2序列化框架内

```go
func init() {
	hessian.RegisterPOJO(&ReqBody{})
	hessian.RegisterPOJO(&RspBody{})
}
```

3. 服务注册

   ```go
   func main() {
      gloryServer := glory.NewServer()
      gloryService := service.NewGloryService("gloryService", &GloryProvider{}) // 服务配置读入
      gloryServer.RegisterService(gloryService) // 业务代码实例化
      gloryServer.Run() // 服务启动
   }
   ```

4. 实例化invoker的构造

   invoker接口在常见开源框架内被应用的十分广泛，也是他们对于抽象化中间件实现的基础。

   invoker往往提供Invoker()方法，可以被用于封装网络逻辑、封装代理逻辑、封装集群策略、重试逻辑、链式调用等等。我设计的glory框架在服务注册的过程中，根据用户代码，构造出了实例化invoker，封装好用户业务逻辑代码，并通过框架暴露。

   在service/glory_service.go的gloryService: Run()函数中，可以看到通过用户代码构造invoker的起点。

   他构造了一个DefaultInvoker类型的invoker

   ```go
   type DefaultInvoker struct {
   	fValMap  map[string]reflect.Value
   	fTypeMap map[string]reflect.Type
   
   	// streamMethod store methodName which is stream
   	streamMethod map[string]bool
   
   	// chanValMap store alive channel of stream function
   	chanValMap map[string][]reflect.Value
   }
   ```

   这个结构可以抽象任何Provider实例化函数，将其method通过反射注册到invoker内部。

   ```go
   for i := 0; i < nums; i++ {
   		invoker.setFunc(typ.Method(i).Name, val.Method(i), typ.Method(i).Type)
   	}
   ```

   最终得到invoker 作为server端provider的函数代理。

### 2.3 服务开启

通过上述介绍，拿到了封装业务代码的invoker，下面将其通过glory协议暴露给外部调用。

在glory/glory_protocol.go的Export函数中，可以看到glory协议暴露上述invoker 的细节

1. 接受TCPconn

2. 开启协程接受请求

3. 开启协程处理请求，并返回

   其中，会针对特定pkg，选择特定的handler进行处理

   ```go
   switch common.PkgType(pkg.Header.PkgType) {
     case common.NormalRequestPkg:
     go g.handleNormalRequest(conn, pkg)
     case common.StreamRequestPkg:
     go g.handleStreamRequest(ctx, conn, pkg)
     case common.StreamSendPkg:
     go g.handleStreamSendPkg(pkg)
     default:
     log.Error("recv unsupported pkg.Header.PkgType = ", pkg.Header.PkgType)
   }
   ```

   对于普通rpc调用，可看到如下处理：

   ```go
   func (g *GloryProtocol) handleNormalRequest(conn *net.TCPConn, pkg *GloryPackage) {
   	err := callNormalInvoker(g.invoker, pkg)
   	if err != nil {
   		log.Error("callNormalInvoker err = ", err)
   		//todo error rsp
   		return
   	}
   	rspPkg := newGloryPackage(pkg.Header.TraceID, pkg.Header.MethodName, uint64(common.NormalResponsePkg), pkg.Header.Seq)
   	rspPkg.Out = pkg.Out
   	if err := rspPkg.sendToConn(conn, g.gloryPkgHandler); err != nil {
   		log.Error("handleNormalRequest: rspPkg.sendToConn(conn, g.gloryPkgHandler) err:", err)
   	}
   }
   ```

   通过调用上述封装好用户业务逻辑的invoker代理，传入glory 协议pkg，获取返回参数，写入了pkg.Out。

   通过pkg.Out构造返回glory协议包，请求序列号不变

   使用序列化工具，将返回包序列化，并写回链接。

   至此，server端可以成功处理一次glory协议请求。

## 3. client端调用普通rpc服务过程

### 3.1 配置文件定义主调接口

```yaml
consumer :
  "gloryClient":
    registry_key: registryKey # 注册中心名（可选）
    service_id: GoOnline-IDE-gloryService # 服务发现ID（可选）
    server_address: 30.225.19.225:8080 # 被调地址和端口 （可选）
    protocol: glory # 协议名
```

在main函数中即可获取client的配置

`glory.NewGloryClient(context.Background(), "gloryClient", &gloryClient)`

### 3.2 定义接口桩和传输结构体

定义接口桩

```go
type GloryClient struct {
	SayHello func(ctx context.Context, req *ReqBody, str2 string) (*RspBody, error)
}
```

和服务端保持一致，定义函数接口即可

注册传输结构体

```go
func init() {
	hessian.RegisterPOJO(&ReqBody{})
	hessian.RegisterPOJO(&RspBody{})
}
```

### 3.3 RPC调用

在main函数中，通过框架根据上述接口桩实例化好调用逻辑，进而直接调用即可。

```go
glory.NewGloryClient(context.Background(), "gloryClient", &gloryClient)
	// test method
rsp, err := gloryClient.SayHello(context.Background(), &ReqBody{
				SeqNum: 1000,
				Value:  "payload string",
				ID:     "24234",
			}, "req2"))
```

在上述实例化客户端接口桩的部分，有值得探讨的细节，具体可见glory/client.go NewGloryClient()函数

1. 配置读入

2. 服务发现

3. 根据服务发现的地址列表，生成对应的所有代理invoker

   ```go
   func newGloryInvoker(targetAddress string) (*GloryInvoker, error) {
   	var conn *net.TCPConn
   	addr, err := net.ResolveTCPAddr("tcp", targetAddress)
   	if err != nil {
   		log.Error("new glory invoker:net.ResolveTCPAddr failed with err = ", err, " address = ", targetAddress)
   		return nil, err
   	}
   	log.Debugf("dail addr = ", targetAddress)
   	conn, err = net.DialTCP("tcp", nil, addr)
   	if err != nil {
   		log.Error("new glory invoker:net.DialTCP failed with err = ", err, " address = ", addr)
   		return nil, err
   	}
   	newInvoker := &GloryInvoker{
   		addr:          common.NewAddress(targetAddress),
   		targetAddress: targetAddress,
   		conn:          conn,
   		handler:       NewGloryPkgHandler(),
   		pendingMap:    sync.Map{},
   	}
   	go newInvoker.startRspListening()
   	return newInvoker, nil
   }
   ```

   在上述代码中，可见发起了tcp链接，并且开启返回监听。

4. 负载均衡封装了选择代理invoker的逻辑

   根据负载均衡策略， 将上述代理invoker列表中选择一个实例。并且发起调用

   每次调用都会选择一次，保证流量的均衡化

5. Glory通过负载均衡invoker来调用远程的网络细节封装入代理函数，通过reflect写入接口桩。

## 4. 客户端负载均衡策略

在上述函数的第三步，即为负载均衡机制。我设计的负载均衡逻辑被封装为特定负载均衡策略的invoker内部

以roundRobin负载均衡策略为例

```go
func (rri *RoundRobinInvoker) Invoke(ctx context.Context, in *common.Params) error {
	in.Seq = uint64(rri.seq.Inc())
	return rri.invokerList[in.Seq%uint64(len(rri.invokerList))].Invoke(ctx, in)
}
```

拿到框架通用的Param参数，设置序列号，为了保证每次调用的序列号不同，采用线程安全的自增策略。

通过取模操作，选择代理invokers列表之一的invoker进行远程调用。



远程调用的过程被封装到代理invoker内：

```go
func (gi *GloryInvoker) Invoke(ctx context.Context, param *common.Params) error {
   gloryPkg := newGloryPackage("", param.MethodName, uint64(common.NormalRequestPkg), param.Seq)
   gloryPkg.Ins = param.Ins
   rspChannel := make(chan interface{})
   gi.pendingMap.Store(param.Seq, rspChannel)
   if err := gloryPkg.sendToConn(gi.conn, gi.handler); err != nil {
      log.Error("Invoke: gloryPkg.sendToConn(gi.conn, gi.handler) err = ", err)
      return err
   }
   recvPkg := <-rspChannel
   param.Out = recvPkg.(*GloryPackage).Out[0]
   close(rspChannel)
   return nil
}
```

可见，其将框架通用的param结构体转化为了glory协议包，通过序列化和pending写入，向服务端发起了一次请求。



## 5. 超时机制
   
   ### 5.1 超时机制实现
   
   在请求的过程中，可能因为网络延迟、server端下线等原因，出现请求超时的情况，针对这一情况，需要在客户端增加超时机制
   
   通过在请求发起后通过time.After建立超时触发器，使用select实现阻塞
   
   protocol/glory/invoker.go
   
   ```go
   	gi.pendingMap.Store(param.Seq, rspChannel)
   // close channel and delete from pending if necessary
   	defer func() {
   		close(rspChannel)
   		gi.pendingMap.Delete(param.Seq)
   	}()
   // call server
   	if err := gloryPkg.sendToConn(gi.gloryConnClient, gi.handler); err != nil {
   		log.Error("Invoke: gloryPkg.sendToConn(gi.conn, gi.handler) err = ", err)
   		return GloryErrorConnErr
   	}
   // start timeout caller 
   	timeoutcaller := time.After(time.Duration(gi.timeout) * time.Millisecond)
   
   	var recvPkg interface{}
   // pending
   	select {
   	case <-timeoutcaller: // timeout 
   		log.Error(`"invoke timeout"`)
   		return GloryErrorTimeoutErr // reurn timeout err
   	case recvPkg = <-rspChannel: // recv rsp
   	}
   ```
   
   可看到在一次rpc调用时，可以针对超时请求进行识别和处理。
   
   ### 5.2 用户配置
   
   Config/client_config.go
   
   ```go
   type NetworkConfig struct {
   	Timeout int `yaml:"timeout" default:"3000"`
   }
   ```
   
   用户可通过配置timeout，实现特定超时的设定。
   
   ### 5.3 超时触发示例
   
   ```
   get rsp = {SeqNum:0 Value: TimeStamp:0001-01-01 00:00:00 +0000 UTC}, err = Code: -1002, Msg: waiting for response time out
   ----timeCost =  3.197914652s
   ```
   
   可看到以网络协议级别错误的形式展现。
   
## 6 高并发支持-粘包处理策略
      
为了保证高并发场景下不出现线程安全问题，我将临界区资源例如：请求序列号seq、pending队列，都使用了sync库。
### 6.1 高并发情景的数据粘包

我希望我的glory框架为稳定的，高并发线程安全的框架，但通过实验，我发现在高并发场景下，我在之前报告中所介绍的将请求write进tcp链接，再read出来的时候，会出现多个数据包粘在一起的情况，导致解包失败。

参考grpc的设计，我选择在经过序列化的data到基础之上，增加帧头（4 byte）保存data length的方式，写入tcp链接，在服务端用相似的思路进行解包。

### 6.2 实现思路

在glory协议的两端，增加打解包中间件，从而实现在序列化data之后的封装。

protocol/glory/transfer.go

```go
// ReadFrame can split data frame by glory header
// length | data
// 4 byte | data
func (cc *gloryConnClient) ReadFrame() ([]byte, error) {
	data := <-cc.frameQueue
	if data == nil {
		return nil, errors.New("read on closed client")
	}
	return data, nil
}

func (cc *gloryConnClient) WriteFrame(data []byte) (int, error) {
	return cc.conn.Write(data2Frame(data)) // data2Frame middleware function
}

func frame2Datas(fm []byte) [][]byte {
	result := make([][]byte, 0)
	for len(fm) > 0 {
		h := fm[:4]
		length := binary.BigEndian.Uint32(h)
		result = append(result, fm[4:4+length])
		fm = fm[4+length:]
	}
	return result
}

func data2Frame(data []byte) []byte {
	length := uint32(len(data))
	fm := make([]byte, 4+length)
	binary.BigEndian.PutUint32(fm, length)
	copy(fm[4:], data)
	return fm
}
```

从而可实现高并发场景的调用。