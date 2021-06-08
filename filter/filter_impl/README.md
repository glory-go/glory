### 调用链追踪和RPC Filter设计

#### 1. Trace RPC追踪

在当今企业的生产环境中，经常会出现几十跳甚至上百跳的RPC调用链路，如果没有trace机制，一旦生产环境中出现问题，将会非常难以追溯，从而容易造成较大的生产事故。

在trace的设计中，一般会在协议头部放置traceid，用来查找单次调用的每一跳，在每一个trace节点，都进行相应数据的上报，在数据收集平台即可查看到每一跳的服务以及状态，从而最快速度地定位问题。

在glory框架中，结合GoOnline的项目落地情况，grpc 协议被广泛地应用。我选择在glory框架中针对 grpc 协议的暴露增加链路trace 追踪。并使用行业内通用解决方案jaeger进行数据收集和展示。最终可以实现将glory框架的调用链路展示在阿里云平台上。

在此基础之上，我基于grpc interceptor 增加了针对rpc调用的Filter调用链，可以通过Filter接口的实现，为调用链中增加用户需要的过滤器。

#### 2. Glory 框架的 Grpc Trace 实现

- 配置文件

```yaml
filter:
  "grpc_filter":
    filter_name: "jaeger"
    aliyun_token_1: xxxxxxx
    aliyun_token_2: xxxxxxx
```

只需要在glory.yaml 配置文件中，增加上述配置，即可将框架的grpc 调用引入trace filter。并将追溯结果上报至阿里云平台。

- grpc 客户端和服务端 filter 的实现。
    1. 读取配置中的Filter key，尝试获取其对应的已经实现好的interceptor，

 ```go
// getDialOption 根据filter返回对应的DialOption
func getOptionFromFilter(filterKeys []string) []grpc.DialOption {
	intercepter, err := intercepter_impl.NewDefaultGRPCIntercepter(filterKeys)
	if err != nil {
		panic(err)
	}

	return []grpc.DialOption{
		grpc.WithUnaryInterceptor(intercepter.ClientIntercepterHandle),
	}
}
 ```

2. 将所有Filterkey的实现（在本例子中只包括jaeger上报） 封装至 Chain filter 中，作为一个Filter抽象结构

```go
func NewDefaultGRPCIntercepter(filterKeys []string) (filter.Intercepter, error) {
    // 获取封装好所有filter 的chain，该chain本质上也是一个Filter 的实现。
	chainFilter, err := plugin.GetFilter("chain", &config.FilterConfig{
		SubFiltersKey: filterKeys, // filter Keys
	})
	if err != nil {
		log.Errorf("new default grpc intercepter error = %v\n", err)
		return nil, err
	}
    // 将当前filter 封装至 grpc interceptor接口中，用于直接交给grpc进行操作。
	return &defaultIntercepter{
		chainFilter: chainFilter,
	}, nil
}
```

3. Jaeger Filter 的构造过程，根据配置构造出tracer，实现ServeHandle和ClientHandle两个接口函数。在实现的过程中将tracer的逻辑实现出来即可。

   具体Filter Chain的构造过程实现将在下一部分展示。

    - 构造tracer

   ```go
   conf, err := jaegerFilter.setup(filterConfig)
   // ...
   jaegerFilter.tracer, _, err = jaegerFilter.newJaegerTracer(conf)
   ```

    - 实现ClientHandle接口的主要trace逻辑：

   ```go
   // 通过构造好的tracer 生成trace span
   span := j.tracer.StartSpan(
   		method,
   		opentracing.ChildOf(parentCtx),
   		opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
   		ext.SpanKindRPCClient,
   	)
   defer span.Finish()
   // 获取jaeger 上下文
   sc, ok := span.Context().(jaeger.SpanContext)
   /*
   ...
   */
   mdWriter := MDReaderWriter{md}
   err := j.tracer.Inject(span.Context(), opentracing.TextMap, mdWriter)
   // 将trace 相关数据写入context.Context 结构，例如traceid。
   newCtx := metadata.NewOutgoingContext(ctx, md)
   
   // 保证当前Filter存在下游filter，否则无法完成Filter链的调用
   if j.next == nil {
       log.Errorf("err = %v", filter.ErrNotSetNextFilter)
       return filter.ErrNotSetNextFilter
   }
   // 调用next Filter，传入本Filter 生成好的context上下文，完成client端trace逻辑
   err = j.next.ClientHandle(newCtx, method, req, reply, cc, invoker, opts...)
   ```

   可看到，在当前Filter中，只关心属于自己的trace 实现逻辑。由于Filter 链中所有Filter都实现了对应接口，所以完成trace逻辑后，直接调用下游。

    - 实现ServerHandle接口的主要trace逻辑：

   ```go
   spanContext, err := j.tracer.Extract(opentracing.TextMap, MDReaderWriter{md})
   // 确保服务端拿到的上下文存在trace数据
   if err != nil && err != opentracing.ErrSpanContextNotFound {
       log.Errorf("JaegerError: extract from metadata err: %v", err)
   } else {
       // 生成新的span结构
       span := j.tracer.StartSpan(
           info.FullMethod,
           ext.RPCServerOption(spanContext),
           opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
           ext.SpanKindRPCServer,
       )
      /*
      ...
      */
       // 追加当前server端trace逻辑
       ctx = opentracing.ContextWithSpan(ctx, span)
   }
   // 确保下游Filter调用链是完整的
   if j.next == nil {
       log.Errorf("err = %v", filter.ErrNotSetNextFilter)
       return nil, filter.ErrNotSetNextFilter
   }
   // 将新构造的context 传入下游
   return j.next.ServerHandle(ctx, req, info)
   ```

   Server端和client端同理，也是在原有上下文的基础上增加当前节点的trace逻辑，例如traceid的上报和链路上下文字段的追加，之后使用新的上下文。

   因此，如果引入glory框架的grpc服务，并且引入了对应的filter 实现，即可构造出一个filter 链，使得每次rpc调用的请求结构和返回结构，都需要经过一层层filter的执行，而在单个filter执行的过程中，会针对自己关心的部分，对请求和返回结构进行特定操作（比如上报tracceid、记录成功与否），而在用户（开发者）的角度，是无感的。

   而如果需要增加新的Filter需求，例如通过Prometheus上报的形式，记录所有请求的耗时，也可以实现对应的Filter，引入filter链即可，具有很强的可配置和可插拔性。

- 结果展示

  我开启三个RPC服务，其中一个服务jaeger-client作为客户端，会在一次RPC请求中通过调用链依次访问到另外三个服务：jaeger-server、jaeger-subserver、jaeger-subsubserver，从而来观察trace链路和耗时情况。

  经过服务的开启和调用，可以在数据平台上收集到如下链路展示。



![](../../img/trace.jpg)
可注意到，当调用经过jaeger-server时，出现了很大的耗时（从server被调用到发起调用经历了很长时间）。则表明该服务具有较大的耗时，考虑瓶颈优化

### 3 Glory 的 grpc Filter 调用链设计

- Filter 接口

  参考grpc提供的Interceptor结构，我定义的Glory-grpc Filter 接口为：

  ```go
  // GRPCFilter is the normal grpc filter interface
  type GRPCFilter interface {
      // 用于设置下一个Filter，链的最后一个Filter无需设置。
  	SetNext(filter GRPCFilter)
  
      // grpc服务端需要链式调用的函数
  	ServerHandle(ctx context.Context,
  		req interface{},
  		info *grpc.UnaryServerInfo) (resp interface{}, err error)
  
      // grpc客户端需要链式调用的寒素
  	ClientHandle(ctx context.Context,
  		method string,
  		req, reply interface{},
  		cc *grpc.ClientConn,
  		invoker grpc.UnaryInvoker,
  		opts ...grpc.CallOption) error
  }
  
  ```

  有了上述接口的约束，即可根据特定Filter在RPC过程中关心的点，针对入参/返回值进行操作，如日志打印、数据上报、失败率统计、请求延迟等。

- Filter Chain 的实现

  glory框架中的 Filter Chain是继承了Filter 接口。可以传入通过一组filter的配置进行初始化，初始化后将获得一整条Filter 链。

  初始化的过程大致分为两步：

    1. 通过配置获取所有Filter 实例
    2. 将所有Filter 实例串联起来形成调用链。

```go
// NewChainFilter never setup failed. If it's subfilter setup with error, chain filter will jump it and setup next.
// The worst result is there is no sub filter successfully setup, and this chain is useless.
func NewChainFilter(filterConfig *config.FilterConfig) (filter.GRPCFilter, error) {
	filterChain := make([]filter.GRPCFilter, 0)
	for _, k := range filterConfig.SubFiltersKey {
		filterConfig, ok := config.GlobalServerConf.FilterConfigMap[k]
		if !ok {
			log.Warnf("filter key %s not defined in config file's filters block, this filter not loaded!", k)
			continue
		}
        // 根据当前Filter配置，从插件中生成当前实例化Filter
		tempFilter, err := plugin.GetFilter(filterConfig.FilterName, filterConfig)
		if err != nil {
			log.Warnf("filter key %s setup error = %v", err)
			continue
		}
        // 加入当前链
		filterChain = append(filterChain, tempFilter)
	}

    // 将链条串起来
	for i, f := range filterChain {
		if i == len(filterChain)-1 {
			break
		}
		f.SetNext(filterChain[i+1])
	}

	return &ChainFilter{
		chain: filterChain,
	}, nil
}
```

- grpc接口Interceptor的适配

  由于grpc与glory的filter接口并不兼容，我在glory框架中，定义了拥有grpc 提供的两个Interceptor 函数的接口如下：

  ```
  // Intercepter is the grpc invocation api, and is the entrance of glory filter
  type Intercepter interface {
  	ServerIntercepterHandle(ctx context.Context,
  		req interface{},
  		info *grpc.UnaryServerInfo,
  		handler grpc.UnaryHandler) (resp interface{}, err error)
  
  	ClientIntercepterHandle(ctx context.Context,
  		method string,
  		req, reply interface{},
  		cc *grpc.ClientConn,
  		invoker grpc.UnaryInvoker,
  		opts ...grpc.CallOption) error
  }
  
  
  ```

  可见，grpc的两个Interceptor 函数的入参中并不包含Interceptor函数本身，因此无法形成调用链扩展。这是我我在此基础之上引入Filter 的原因。有了Filter 和Interceptor接口，我只需要实现好一个包含glory框架Filter链的Interceptor结构，然后将其注入grpc的客户端和服务端。

  默认grpc Interceptor的实现

  ```go
  // defaultIntercepter is default grpc intercepter, it calls chain_filter and it self is the final filter of grpc invocation filter procedure
  type defaultIntercepter struct {
  	chainFilter filter.GRPCFilter
  	handler     grpc.UnaryHandler
  }
  
  ```

  defaultIntercepter 同时实现了Interceptor和Filter接口。他有一个自己的Filter链，以及负责真正执行业务代码的handler。在这它的调用过程中，会首先将请求通过Filter链执行，再将链的最后一个Filter指向自己，请求的最终会运行到自己实现Filter的函数中，完成整个调用。

  以Client端调用为例：

  ```go
  func (j *defaultIntercepter) ClientIntercepterHandle(ctx context.Context,
  	method string,
  	req, reply interface{},
  	cc *grpc.ClientConn,
  	invoker grpc.UnaryInvoker,
  	opts ...grpc.CallOption) error {
  
      // 将自己作为Chain的next Filter
  	j.chainFilter.SetNext(j)
      
      // 将请求传入chain，开始依次调用Filter 链，最终调用到自身实现的函数中，结束调用。
  	return j.chainFilter.ClientHandle(ctx, method, req, reply, cc, invoker, opts...)
  }
  ```



defaultIntercepter 的最后一次执行（由自己作为Filter链的最后一环）的实现：

  ```go
  func (j *defaultIntercepter) ServerHandle(ctx context.Context,
     req interface{},
     info *grpc.UnaryServerInfo) (resp interface{}, err error) {
     return j.handler(ctx, req)
  }
  ```

只需要执行handler，完成rpc调用即可，不需要关心其他逻辑。