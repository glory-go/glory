Glory
===================================

Glory框架为一款Go语言的轻量级RPC框架，您可以使用它快速开发你的服务实例。如果您希望在**微服务场景下**使用**gRPC**进行网络通信，那么Glory会使您的开发、运维工作量减轻不少。

欢迎访问Glory主页： [glory-go.github.io](https://glory-go.github.io/introduction)

示例仓库：[github.com/glory-go/glory-demo](https://github.com/glory-go/glory-demo)

Glory提供以下能力：

- 通信协议：Glory框架提供gRPC（client端和server端）、HTTP（server端）、Websocket（server端）**脚手架**，你可以通过**几行配置和几行代码**快速开启多个gRPC或HTTP服务。

- 配置：Glory框架提供**统一化的配置服务**，您只需要在main文件同级目录config文件夹下定义glory.yaml，在配置文件内按照约定格式写入配置信息，在引入框架后执行时，框架会自动读入配置文件，并开启所需服务。

  您也可以选择从nacos **配置中心拉取**当前服务所需配置。

- 日志：Glory框架提供**日志支持**，您可以在配置文件中定义自己需要的日志记录方式。支持命令行、文件、远程（基于elastic、阿里云sls）的日志收集方式。

- 链路追踪：glory框架提供适配于 **gRPC 的链路追踪**服务，你可以选择将服务调用链路上报至本地jaeger或阿里云链路追踪平台进行监控和错误追溯。

- 数据上报：glory框架提供**基于Promethus的数据上报**服务，你可以在配置文件中定义自己需要的数据上报方式，同时支持基于promethus-pushgateway的推模式数据上报。

- 第三方工具常用sdk支持：glory框架提供mysql、redis、mongodb、rabbitmq等常见工具的sdk封装，开发者可以在配置中引入服务，使用框架提供的sdk进行快速开发。

- 服务治理：Glory框架提供基于K8s、Nacos的服务发现机制，可以在**k8s集群中自动进行Glory-gRPC服务实例的注册、发现和负载均衡。**

如果您觉得不错的话麻烦留下一颗星星⭐