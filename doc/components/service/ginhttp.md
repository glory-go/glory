# 基于Gin实现的http服务(ginhttp)

该服务(Service)提供用户初始化并注册多个gin http服务的能力，框架基于用户配置完成了gin服务的初始化，结合Service的`Run`能力，帮助用户降低重复的代码编写负担

## 配置定义

参考`components/service/ginhttp/config.go`中的配置定义，一个可供参考的完整配置定义如下：

```
config_center:
    xxx
service:
    ginhttp:
        addr: ":8080"
        read_timeout: 5
        write_timeout: 10
```