## 3. http-server filter 中间件的设计

http服务的启动过程与上述grpc十分类似，其中亮点是，根据业务需求，http-server封装了mux，在提供基础服务：路由、统一化http参数获取、自动打解包等等等基础上，增加了链式请求过滤器的支持。

### 3.1 链式http请求过滤的实现

可在http/filter.go中看到针对过滤器和业务处理函数接口的定义

```go
// HandleFunc 业务处理函数接口
type HandleFunc func(controller *GRegisterController) (err error)

// Filter 过滤器（拦截器），根据dispatch处理流程进行上下文拦截处理
type Filter func(controller *GRegisterController, f HandleFunc) (err error)

// Chain 链式过滤器
type Chain []Filter
```

框架会为http服务分配一个filter链，即上述的Chain结构，

```go
// http处理函数
func getGloryHttpHandler(handler func(*GRegisterController) error, req, rsp interface{}, filters []Filter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		retPkg := rspImpPackage

		// recovery
		...

		// 创建针对当前接口的过滤器chain
		chain := Chain{}
		chain.AddFilter(filters) // 注册过滤器

		tRegisterController := GRegisterController{
			Ctx: r.Context(),
			R:   r,
		}
		// 处理 req
		// 处理 rsp
			// 执行业务函数
		if err := chain.Handle(&tRegisterController, handler); err != nil {
			retPkg.SetErrorPkg(w, err)
			return
		}
		// 最终回包
    ...
		return
	}
}

```

可见，在上述函数中，将用户传入的自定义filter过滤函数，放在chain内形成filter链，再由chain调用Handle 函数，逐个执行链内所有filter，通过所有过滤的请求才最终被执行业务逻辑，否则按照filter逻辑返回。

http/filter.go:  Handle函数内的实现

```go
//多个Filter,递归执行
	lastI := n - 1
	return func(controller *GRegisterController, f HandleFunc) error {

		var (
			chainFunc HandleFunc
			curI      int
		)
		chainFunc = func(controller *GRegisterController) error {
			if curI == lastI {
				return f(controller)
			}
			curI++
			err := (*fc)[curI](controller, chainFunc)
			curI--
			return err
		}
		return (*fc)[0](controller, chainFunc)
	}(controller, f)
```

### 3.2 链式http请求过滤的使用

```go
// 测试用filter 如果input字段为-2则报错
func myFilter2(controller *ghttp.GRegisterController, f ghttp.HandleFunc) error {
	req, ok := controller.Req.(*gloryHttpReq)
	if !ok {
		log.Error("req type err")
		return errors.New("req type err")
	}

	if req.Input[0] == -2 {
		log.Error("filting because input == -2")
		return errors.New("filting because input == -2")
	}
	err := f(controller)
	return err
}

func main() {
	gloryServer := glory.NewServer()
	httpService := service.NewHttpService("httpDemo")
	httpService.RegisterRouter("/testwithfilter/{hello}/{hello2}", testHandler, &gloryHttpReq{}, &gloryHttpRsp{}, "POST", myFilter1, myFilter2)
	gloryServer.RegisterService(httpService)
	gloryServer.Run()
}
```

可按照此方法，将，myfilter1，myfilter2等过滤器注册在httpService上，从而针对请求进行过滤。

在GoOnline项目中，此filter被广泛应用于auth部分token鉴权。

