## 统一化配置的实现

### 1.1 配置文件选型

配置文件调研：

java框架常见配置文件：xml

nodejs框架常见配置文件：json

golang框架常见配置文件：json（go-zero），yaml（dubbo-go、k8s）

经过调研，使用yaml文件作为golang框架的配置文件应用较为广泛，并且go语言有解析yaml文件的专用库，yaml配置文件的编写和修改较为方便，支持多种常见结构和数据类型的定义，配置文件内容较为简洁。

最终选择yaml格式作为glory框架的配置文件格式，就在下文称本框架的配置文件为glory.yaml。

### 1.2 根据环境变量对配置文件的选择

这个是在本阶段最后，开发人员提出的需求，要求使用多个配置文件，通过环境变量的方式，对多种环境：dev、release、test环境进行区分。

由于考虑到产品兼容性问题，将在之后针对这部分内容进行修改。

思路为：配置加载之前，首先从环境变量中读入：GLORY_ENV 

- 如果存在环境变量，则读入glory_$GLROY_ENV.yaml配置文件
- 如果不存在，默认读入glory.yaml文件作为配置

从而达到针对多环境的动态兼容。

### 1.3 配置读入灵活性设计

- 设计思路

由于GoOnline是基于k8s管理的微服务，在之前没有框架的时候，都是使用环境变量保存配置，为了兼容已有服务，我决定使用特定字段对同一级别的配置进行标注，被标注的配置会事先从配置文件中读入配置信息，再以读到的字段为key逐一尝试从环境变量读入，如果存在配置则替换。

这样的策略可以保证和GoOnline已有服务配置的兼容。

- 实现方式

  1. 配置中增加`config_source`字段，供开发者选择是否触发上述兼容逻辑。

     以/config/service_config.go中的ServiceConfig结构体为例：

     ```go
     type ServiceConfig struct {
     	ConfigSource string `yaml:"config_source"` # mark
     	Protocol     string `yaml:"protocol"`
     	RegistryKey  string `yaml:"registry_key"`
     	Port         int    `yaml:"port"`
     	ServiceID    string `yaml:"service_id"`
     }
     ```

     service_config配置是一个抽象的配置，里面针对provider端，可以配置开启服务的协议`protocol`， 选择注册中心时使用的`registry_key`注册中心名， 本地暴露的端口`port`，以及服务注册发现使用的唯一ID：`service_id`。

     除此之外，mark标注的字段起到了上述作用，在我的框架逻辑中，一旦发现这个字段为env, 则触发环境变量尝试读入逻辑。

     判断逻辑写在了/tools/toosl.go ReadFromEnvIfNeed函数中

     ```go
     func ReadFromEnvIfNeed(rawConfig interface{}) error {
     	rawVal := reflect.ValueOf(rawConfig)
     	if rawVal.Kind() != reflect.Ptr { //保证是指针
     		log.Println("ReadFromEnvIfNeed param rawConfig should be pointer to struct")
     		return errors.New("input rawConfig kind invalid")
     	}
     	val := rawVal.Elem()
     	typ := reflect.TypeOf(rawConfig).Elem()
     	kd := val.Kind() //获取到a对应的类别
     	if kd != reflect.Struct {
     		log.Println("ReadAllConfFromEnv expect struct but got ", kd)
     		return errors.New("ReadAllConfFromEnv expect struct")
     	}
     
     	num := val.NumField()
     	for i := 0; i < num; i++ {
     		tagVal := typ.Field(i).Tag.Get("yaml")
     		if tagVal == "config_source" {
     			val, ok := val.Field(i).Interface().(string)
     			if !ok {
     				//不符合要求，直接跳过
     				log.Println("yaml tag 'config_source' type should be string")
     				continue
     			}
     			if val == "env" {
     				//符合环境变量读入要求
     				readAllConfFromEnv(rawConfig)
     				return nil
     			}
     		}
     	}
     	return nil
     }
     ```

     可以看到，我应用了反射，针对特定字段进行判断，再针对需要环境读入的config调用readAllConfFromEnv函数

     ```go
     // readAllConfFromEnv 提供从环境变量拉取的服务
     // 对于传入的rawConfig结构所有yaml标签不为config_source 的字段，尝试从环境变量中拉取，如果拉取到则替换原有值
     func readAllConfFromEnv(rawConfig interface{}) {
     	val := reflect.ValueOf(rawConfig) //获取reflect.Type类型
     	typ := reflect.TypeOf(rawConfig)
     	//获取到该结构体有几个字段
     	num := val.Elem().NumField()
     	//遍历结构体的所有字段
     	for i := 0; i < num; i++ {
     		//fmt.Printf("Field %d:值=%v\n", i, val.Elem().Field(i))
     		tagVal := typ.Elem().Field(i).Tag.Get("yaml")
     		//如果该字段有tag标签就显示，否则就不显示
     		if envValue := GetEnv(val.Elem().Field(i).String()); tagVal != "config_source" && envValue != "" {
     			// 非config字段，并且环境变量里面有定义
     			val.Elem().Field(i).SetString(envValue)
     		}
     	}
     }
     
     ```

     根据上述实现，可以达到环境变量读入检查的兼容性需求。

### 1.4 中心化配置管理（设计）

在未来，当GoOnline有足够的资源维护一个配置中心例如nacos，可以将启动时根据服务ID从配置中心拉取特定环境配置的逻辑写入框架。