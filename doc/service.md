# 服务(Service)

服务是一类特殊的组件(Component)，开发者编写的项目通过该类组件对外提供服务。我们所熟悉的Http、各种RPC以及消息队列的消费者都是服务的一员。Service将这些组件进行抽象，允许社区开发者开发出适合自己的服务实现，也便于框架使用者更方便地注册和使用这些外部服务

## 注册服务

glory提供了服务的注册能力，调用`service.GetService().RegisterService()`即可完成服务的注册

## 运行服务

已注册的服务将自动跟随配置的加载完成初始化，但初始化的服务并没有真正启动，用户需调用`service.Run()`运行服务，该方法为阻塞调用，当所有服务都退出时，该方法才会退出

service不会做panic保护，这意味着任何一个服务panic都会导致整个进程的退出