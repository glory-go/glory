package glory

import (
	"context"
	"reflect"

	"github.com/glory-go/glory/common"
	"github.com/glory-go/glory/config"
	_ "github.com/glory-go/glory/load_balance/round_robin"
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/plugin"
	"github.com/glory-go/glory/registry"
)

func NewClient(ctx context.Context, clientName string, clientService interface{}) {
	var targetAddress []common.Address
	// get glory client config
	gloryClientConfig, ok := config.GlobalServerConf.ClientConfig[clientName]
	if !ok {
		panic("serviceName " + clientName + " in your source code not found in config file!")
	}

	var registryProtoc registry.Registry

	// get target server address, if not direct link, choose target glory_registry to get target address
	if gloryClientConfig.ServerAddress != "" {
		// direct link
		targetAddress = append(targetAddress, *common.NewAddress(gloryClientConfig.ServerAddress))
	} else if gloryClientConfig.RegistryKey != "" {
		// get target from glory_registry
		registryConfig := config.GlobalServerConf.RegistryConfig[gloryClientConfig.RegistryKey]
		registryProtoc = plugin.GetRegistry(registryConfig)
		targetAddress = registryProtoc.Refer(gloryClientConfig.ServiceID)
	} else {
		panic("no target address found! ")
	}

	serviceList := make([]interface{}, 1)
	serviceList = append(serviceList, clientService)

	// use glory_registry protocol to handle target address to invoker
	protoc := plugin.GetProtocol(gloryClientConfig.Protocol, gloryClientConfig.NetworkConfig, serviceList...)
	invokers := make([]common.Invoker, 0)
	for _, addr := range targetAddress {
		invokers = append(invokers, protoc.Refer(addr))
	}

	// use load balance to handler multiple invokers to one invoker
	// todo loadBanlanceType from config
	loadBalance := plugin.GetLoadBalance("round_robin")
	invoker := loadBalance.Select(invokers, registryProtoc, gloryClientConfig.ServiceID, protoc)

	// 传入方法名和返回值类型，返回值类型最后一位一定是error，返回值长度只能是一个或者两个
	funcProxyFactory := func(methodName string, outs []reflect.Type) func(in []reflect.Value) []reflect.Value {
		// 代理函数入参 in
		return func(in []reflect.Value) []reflect.Value {
			var (
				err    error
				inIArr []interface{}
				reply  reflect.Value
				params common.Params
			)

			if len(outs) > 2 {
				log.Error()
			}

			if len(outs) == 2 { // 返回值为 (结构体, error)
				if outs[0].Kind() == reflect.Ptr {
					reply = reflect.New(outs[0].Elem())
				} else {
					reply = reflect.New(outs[0])
				}
			}

			start := 0
			end := len(in)
			invokeCtx := context.Background()
			if end > 0 && in[0].Type().String() == "context.Context" { // 如果存在传入的ctx
				if !in[0].IsNil() {
					invokeCtx = in[0].Interface().(context.Context)
				}
				start += 1
				if len(outs) == 1 && in[end-1].Type().Kind() == reflect.Ptr {
					end -= 1
					reply = in[len(in)-1]
				}
			}

			if end-start <= 0 {
				inIArr = []interface{}{}
			} else if v, ok := in[start].Interface().([]interface{}); ok && end-start == 1 {
				inIArr = v
			} else {
				inIArr = make([]interface{}, end-start)
				index := 0
				for i := start; i < end; i++ {
					inIArr[index] = in[i].Interface()
					index++
				}
			}

			params.MethodName = methodName
			params.Ins = inIArr
			params.Out = reply.Interface()

			err = invoker.Invoke(invokeCtx, &params)
			// params store protocol error from server, err store error from client
			// err is fatal than params.Error
			if err != nil {
				log.Error("rpc call result err = ", err)
			} else {
				log.Debugf("rpc call reply: %+v, %+v, err = %+v", params.Out, reply.Interface(), params.Error)
				err = params.Error
			}

			if len(outs) == 1 {
				return []reflect.Value{reflect.ValueOf(&err).Elem()}
			}
			if len(outs) == 2 && outs[0].Kind() != reflect.Ptr {
				return []reflect.Value{reflect.ValueOf(params.Out), reflect.ValueOf(&err).Elem()}
			}
			return []reflect.Value{reflect.ValueOf(params.Out), reflect.ValueOf(&err).Elem()}
		}
	}

	// 传入方法名和返回值类型，返回值类型最后一位一定是error, 倒数第二位为rspchan，返回值长度最少为2
	streamFuncProxyFactory := func(methodName string, outs []reflect.Type) func(in []reflect.Value) []reflect.Value {
		// 代理函数入参 in
		return func(in []reflect.Value) []reflect.Value {
			log.Debug("here outLen = ", len(outs))
			var (
				err    error
				inIArr []interface{}
				params common.Params
			)
			if len(outs) < 2 {
				panic("too few func out param, out param must be (chan, chan..., err)")
			}

			start := 0
			end := len(in)
			//invokeCtx := context.Background()
			if end > 0 && in[0].Type().String() == "context.Context" { // 如果存在传入的ctx
				if !in[0].IsNil() {
					//invokeCtx = in[0].Interface().(context.Context)
				}
				start += 1
			}

			if end-start <= 0 {
				inIArr = []interface{}{}
			} else if v, ok := in[start].Interface().([]interface{}); ok && end-start == 1 {
				inIArr = v
			} else {
				inIArr = make([]interface{}, end-start)
				index := 0
				for i := start; i < end; i++ {
					inIArr[index] = in[i].Interface()
					index++
				}
			}

			params.MethodName = methodName
			params.Ins = inIArr

			// creat out Stream to recv streamRsp package
			outStream, addr, seq, err := invoker.StreamInvoke(ctx, &params, outs[len(outs)-2])
			// params store protocol error from server, err store error from client
			// err is fatal than params.Error
			if err != nil {
				log.Error("rpc call result err = ", err)
				// just panic!!!
				panic(err)
			} else {
				log.Debug("rpc call reply: %+v, err: %=v", params.Out, params.Error)
				err = params.Error
			}

			// create many send chans to send streamReq package
			chanListValue := []reflect.Value{}
			for j := 0; j < len(outs)-2; j++ {
				if outs[j].Kind() != reflect.Chan {
					log.Error("stream mod must have return format:(chan recv, chan rsp, err), don't add none chan param")
					continue
				}
				tempChanValue := reflect.MakeChan(outs[j], 0)
				chanListValue = append(chanListValue, tempChanValue)

				// recv stream from client goroutine:
				go func(ctx context.Context, chanValue reflect.Value, offset int, seq uint64) {
					valueChan := make(chan reflect.Value)
					stopChan := make(chan interface{}, 1)
					for {
						// chanValue -> chan interface{}

						go func() {
							for {
								select {
								case <-stopChan:
									return
								default:
									recv, ok := chanValue.Recv()
									if !ok {
										log.Error(" chanValue.Recv() err ")
										break
									}
									valueChan <- recv
								}
							}
						}()
						select {
						case <-ctx.Done():
							stopChan <- struct{}{}
							log.Debugf("GRACEFUL SHUTDOWN")
							return
						case recv := <-valueChan:
							param := common.Params{
								MethodName: methodName,
								ChanOffset: uint8(offset),
								Value:      recv.Interface(),
								Seq:        seq,
								Addr:       addr,
							}
							if err := invoker.StreamRecv(&param); err != nil {
								log.Error("stream recv error = ", err)
							}
						}
					}
				}(ctx, tempChanValue, j, seq)
			}
			returnValues := chanListValue
			returnValues = append(returnValues, outStream, reflect.ValueOf(&err).Elem())
			return returnValues
		}
	}

	val := reflect.ValueOf(clientService).Elem()
	typ := reflect.TypeOf(clientService).Elem()
	nums := val.NumField()
	for i := 0; i < nums; i++ { //为每一个需要的函数创造一个函数代理，
		funcType := typ.Field(i)
		funcValue := val.Field(i)
		if funcValue.Kind() == reflect.Func && funcValue.IsValid() && funcValue.CanSet() {
			outPutNums := funcType.Type.NumOut()
			outTypes := make([]reflect.Type, 0)
			for j := 0; j < outPutNums; j++ {
				outTypes = append(outTypes, funcType.Type.Out(j))
			}
			if outTypes[0].Kind() == reflect.Chan {
				// stream function
				val.Field(i).Set(reflect.MakeFunc(val.Field(i).Type(), streamFuncProxyFactory(typ.Field(i).Name, outTypes)))
				continue
			}
			// normal function
			val.Field(i).Set(reflect.MakeFunc(val.Field(i).Type(), funcProxyFactory(typ.Field(i).Name, outTypes)))
		}
	}
}
