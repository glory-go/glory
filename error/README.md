#  分层错误码设计

glory协议引入了错误码和错误信息的概念，我将错误码使用了“按层分类”的策略。将“用户层error”、“框架层error”、“协议层error”、按照可扩展性分隔开。这样的设计一是为了增强框架的网络协议可插拔性，二是为了更容易定位问题。

在之前开发过程中，有阅读过一些框架的源码，他们将rpc请求server端用户实现代码返回的error与框架层的error并没有区分开，同时按照请求错误处理。最终返回的是框架的error报错。这样的设计是存在问题的，对于一些业务场景，用户代码返回的error是存在意义的，并且对于其他同时返回的response信息，应该按照正常逻辑返回给调用方，不应该因为抛出错误而阻断。

这也是我分层错误码设计解决的问题。框架能合理区分error是来自用户、框架、还是网络协议，按照层级来抛出error，保证问题的合理追溯。

在protocol/glory/error.go中，我有定义属于协议层的，适配与所使用协议的error，这些error是和协议绑定的，对于新协议的扩展，只要在协议代码中处理好来自上面框架层和用户层的代码，定义自己的errorCode，是可以优雅地横向扩展的。

```go
	// client error
	GloryErrorConnErr          = NewError(-1001, "conntion error")
	GloryErrorTimeoutErr       = NewError(-1002, "waiting for response time out")
	GloryErrorEmptyResponseErr = NewError(-1003, "get empty response")
```

- 下面是glory协议client端的一个协议层异常抛出场景:

protocol/glory/invoker.go

```go
rspPkg, ok := rspRawPkg.(*ResponsePackage)
if !ok {
 	log.Error("StreamInvoke:rspRawPkg assert not *ResponsePackage err")
  return rspChan, nil, 0, GloryErrorProtocol
}
```

当网络协议层出现来自server回包解包或断言错误，应当向上抛出GloryErrorProtocol 协议异常。

抛出的一场被框架层捕获到，向上返回给用户代码。

glory/client.go 

```go
err = invoker.Invoke(invokeCtx, &params)
// params store protocol error from server, err store error from client
// err is fatal than params.Error

return []reflect.Value{reflect.ValueOf(params.Out), reflect.ValueOf(&err).Elem()}
```

- 下面是glory协议server端一个协议层异常抛出场景

protocol/glory/protocol.go

```go
  if req.Header.Version != defaultVersion {
    log.Error("recv unsupportex version = ", req.Header.Version, " support version = ", defaultVersion)
    g.sendErrorResponse(req, gloryClient, GloryErrorVersion)
    continue
  }


func (g *GloryProtocol) sendErrorResponse(req *RequestPackage, client *gloryConnClient, err *Error) {
	rspPkg := NewResponsePackage(req.Header.TraceID, uint64(common.ErrorRspPkg), req.Header.Seq, err)
	rspPkg.sendToConn(client, g.gloryPkgHandler)
}

```

在检查协议版本号时，如果出现版本错误，则会返回错误包，错误码为`GloryErrorVersiuon`。

将在client端收到错误包，并返回给上层框架代码

protocol/glory/invoker.go

```go
recv, ok := recvPkg.(*ResponsePackage) // check package
if !ok {
  log.Error("Invoke: recvPkg assert *ResponsePackage err")
  return GloryErrorProtocol
}

if len(recv.Result) == 0 {
} else {
  param.Out = recv.Result[0]
}
param.Error = recv.Error // return to up level
```

- 下面是glory协议server端一个用户层错误的处理

框架层的default invoker在glory框架中的作用为调用最顶层用户server逻辑代码：

common/invoker_impl/default_invoker.go

```go
outValue := f.Call(valueList) // call user code logic

in.Out = outValue[0].Interface() // get output 

if outValue[len(outValue)-1].Interface() == nil { // check error
return nil
}
return outValue[len(outValue)-1].Interface().(error) // return error
```

而这段代码捕获到了用户自己手动抛出的error，并向下传递给协议层：

Protocol/glory/protocol.go

```go
err := invoker.Invoke(context.Background(), params)
// framework(user) level error cached! not protocol level error
if err != nil {
return params.Out, NewGloryErrorUserError(err)
}
```

对于glory协议server端的invoker，它拿到了来自上层的用户error，**选择使用NewGloryErrorUserError函数进行封装，再和用户想要的正确结果一起，通过网络层传递给client。**

glory/client.go

```go
err = invoker.Invoke(invokeCtx, &params)
// params store protocol error from server, err store error from client
// err is fatal than params.Error
if err != nil { // client protocol level error
	log.Error("rpc call result err = ", err)
} else { // server user level error
	log.Debugf("rpc call reply: %+v, err = %+v", params.Out, params.Error)
	err = params.Error 
}
```

可以看到，客户端的框架层拿到了来自底层的error，对于params参数里面的error，为从server端返回的用户级别错误，将不会影响框架正常运行，并且会通过err传递给客户端用户代码。

### 2.4 运行示例

- 对于一个用户级别的错误抛出，客户端调用后会看到日志：

```
get rsp = {SeqNum:1001 Value:payload string TimeStamp:2021-01-30 22:44:03.483 +0800 CST}, err = Code: -1004, Msg: user defined error
----timeCost =  212.590535ms
```

可看到既显示了需要的response，又打印出了用户定义的错误。

- 对于一个框架层面的异常，比如client未在注册中心找到对应provider。将会抛出框架层的错误和空返回值：

```
get rsp = {SeqNum:0 Value: TimeStamp:0001-01-01 00:00:00 +0000 UTC}, err = Code: -301, Msg: registry: can't found target provider
----timeCost =  88.341866ms
```

- 对于一个协议层的异常，比如协议版本号不一致，将会抛出协议层异常错误和空返回值

```
get rsp = {SeqNum:0 Value: TimeStamp:0001-01-01 00:00:00 +0000 UTC}, err = Code: -3002, Msg: glory version error
----timeCost =  131.86831ms
```

## 