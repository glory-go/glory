# K8s label 服务发现支持

### 1. 服务注册发现原理
k8s将资源的所有字段保存在etcd分布式存储中，通过api server提供资源的修改接口。

glory框架目前的实现是：将当前服务名作为label key、pod真实ip（并非nodeip）经过转义后作为label value。

客户端服务发现：在当前命名空间下，筛选所有提供所需服务的pod真实ip，将拿到的ip列表做客户端负载均衡，进行调用

### 2. 使用方法
example：可见：example/k8s

运行服务首先保证本地存在k8s集群

- 框架配置k8s注册方式：config/glory.yaml

```yaml
registry:
  "registryKey":
    service: k8s
```
引入依赖
```go
import(
       _ "github.com/glory-go/glory/registry/k8s"
)

```

- 运行示例
构造namespace

```shell script
kubectl create namespace glory
```

分别在client和server下运行
```shell script
sudo sh build.sh
```
即可开启三个server和一个client，可通过docker dashboard 或者日志打印的方式，看到负载均衡调用的体现。

### 3. 接下来的工作
- 目前服务发现采用客户端隔五秒轮询一次刷新的方法。在高qps场景下上下线会造成大量失败调用。较为严谨的解决办法是通过监听注册中心的修改事件，实时刷新客户端本地缓存️。目前尚未支持。

    高并发，海量请求场景下的服务注册和发现是一个值的深度探讨的话题，关于大面积停机，流量摘除，优雅下线等等，其中涉及到的策略有很多，欢迎一起讨论和本框架此模块协同共建。
- 转义部分之后考虑使用base64实现