package autowire

var allWrapperAutowires = make(map[string]WrapperAutowire)

func Load() error {
	// 1. impl autowire
	allWrapperAutowires = GetAllWrapperAutowires()

	// impl all struct
	for _, aw := range allWrapperAutowires {
		for sdID := range aw.GetAllStructDescribers() {
			if aw.IsSingleton() {
				if _, err := aw.Impl(sdID, nil); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func Impl(autowireType, structDescriberID string, param interface{}) (interface{}, error) {
	for _, wrapperAutowire := range allWrapperAutowires {
		if wrapperAutowire.TagKey() == autowireType {
			return wrapperAutowire.Impl(structDescriberID, param)
		}
	}
	return nil, nil
}

//func impl(p *autowire.StructDescriber) interface{} {
//	// 1. check if already impled
//	tempInterfaceId := util.GetIdByInterfaceAndImplPtr(p.Interface, p.ImplStructPtr)
//	if impledPtr, ok := implCompletedMap[tempInterfaceId]; ok {
//		// if already impleted, return
//		return impledPtr
//	}
//	defer func() { // assure the impl procedure of one service run once
//		if r := recover(); r != nil {
//			log.Println("recover panic = %s", r)
//		}
//		implCompletedMap[tempInterfaceId] = p.ImplStructPtr
//	}()
//
//	// 2. check type
//	valueOf := reflect.ValueOf(p.ImplStructPtr)
//	valueOfElem := valueOf.Elem()
//	typeOf := valueOfElem.Type()
//	if typeOf.Kind() != reflect.Struct {
//		panic("invalid struct ptr")
//	}
//
//	// 3. tag inject
//	numField := valueOfElem.NumField()
//	for i := 0; i < numField; i++ {
//		t := typeOf.Field(i)
//		var impledPtr interface{}
//		tagKey := ""
//		tagValue := ""
//		for _, aw := range allAutowires{
//			if val, ok := t.Tag.Lookup(aw.TagKey()); ok{
//				fieldTypeName := t.Type.Name()
//				impledPtr = aw.Impl(allEntries[aw.TagKey() +"-"+ util.GetIdByNamePair(fieldTypeName, val)], allAutowires)
//				tagKey = aw.TagKey()
//				tagValue = val
//				break // only one tag is support
//			}
//		}
//		//
//		//
//		//if svcImplStructName := t.Tag.Get("service"); svcImplStructName != "" {
//		//	// get impled sub local service
//		//	fieldTypeName := t.Type.Name()
//		//	if fieldTypeName == "" { // autowire struct ptr
//		//		fieldTypeName = svcImplStructName
//		//	}
//		//	impledPtr = impl(singletonsMap[util.GetIdByNamePair(fieldTypeName, svcImplStructName)])
//		//	tagKey = "service"
//		//	tagValue = svcImplStructName
//		//} else if grpcClientName := t.Tag.Get("grpc"); grpcClientName != "" {
//		//	// `service:"grpc"` means auto wire grpc client
//		//	impledPtr = implGRPC(t.Type.Name(), grpcClientName, t.Tag.Get("interceptorsKey"))
//		//	tagKey = "grpc"
//		//	tagValue = grpcClientName
//		//} else if configPath := t.Tag.Get("config"); configPath != "" {
//		//	// XXX string `config:"mysvc.config"` means auto wire config
//		//	stringConfig := &StringConfig{}
//		//	err := config.LoadConfigPathValue(stringConfig)
//		//	if err != nil {
//		//		log.Println("error load string from config path", configPath, " with error ", err)
//		//		continue
//		//	}
//		//	tagKey = "config"
//		//	tagValue = configPath
//		//
//		//	subService := valueOfElem.Field(i)
//		//	if !(subService.Kind() == reflect.String && subService.IsValid() && subService.CanSet()) {
//		//		err := perrors.Errorf("Failed to autowire interface %s's confige. It's field %s with tag '%s:\"%s\"', please check if the field is exported",
//		//			util.GetStructName(p.Interface), t.Type.Name(), tagKey, tagValue)
//		//		panic(err)
//		//	}
//		//	subService.Set(reflect.ValueOf(stringConfig.string))
//		//	continue
//		//}
//		if tagKey == "" && tagValue == "" {
//			continue
//		}
//		// set field
//		subService := valueOfElem.Field(i)
//		if !(subService.IsValid() && subService.CanSet()) {
//			err := perrors.Errorf("Failed to autowire interface %s's impl %s service. It's field %s with tag '%s:\"%s\"', please check if the field is exported",
//				util.GetStructName(p.Interface), util.GetStructName(p.ImplStructPtr), t.Type.Name(), tagKey, tagValue)
//			panic(err)
//		}
//		subService.Set(reflect.ValueOf(impledPtr))
//	}
//
//	// 4. load config
//	if err := config.LoadConfigPathValue(p.Config); err != nil{
//		log.Println("load config error = ", err)
//	}
//
//	// 5. run construct function
//	if p.ConstructFunc != nil{
//		if err := p.ConstructFunc(p.Config, p.ImplStructPtr); err != nil{
//			log.Println("error construct impl ", util.GetStructName(p.Interface), util.GetStructName(p.ImplStructPtr), " error = ", err)
//			return p.ImplStructPtr
//		}
//	}
//
//	// 6. monkey
//	if os.Getenv("GOARCH") == "amd64" || runtime.GOARCH == "amd64" {
//		// only service, only amd64 mod can inject monkey function
//		implMonkey(p.ImplStructPtr, tempInterfaceId)
//	}
//	return p.ImplStructPtr
//}
