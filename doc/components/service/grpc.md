# GRPC服务

GRPC服务(Service)提供用户初始化并注册多个grpc服务的能力，框架基于用户配置完成了grpc服务的初始化，结合Service的`Run`能力，帮助用户降低重复的代码编写负担

## 配置定义

参考`components/service/grpc/config.go`中的配置定义，一个可供参考的完整配置定义如下：

```
config_center:
    xxx
service:
    grpc:
        addr: ":8080"
```

## 中间件

由于grpc初始化后不再提供注册中间件的能力，因此我们提供了`WithOptions()`方法，来帮助用户在使用时注入中间件。**请确保中间件的注册发生在config.Init()方法之前**，否则后续新增的中间件将会被忽略