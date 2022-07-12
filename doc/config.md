# 配置中心

配置中心提供了一组抽象的接口定义，从而支持开发者按需接入不同的配置来源

## 基础操作

### 配置文件默认位置及格式

配置文件默认的读取位置为`./glory/glory.${env}.yaml`，其中`env`为环境名称，将从`GLORY_ENV`环境变量中读取。用户可通过调用`config.ChangeDefaultConfigPath()`来指定新的文件位置，支持yaml后缀的文件，且`${env}`会被默认加入到读取配置文件的路径中

若启动时找不到带环境后缀的配置文件时，将默认读取不带后缀的配置文件。当所有配置文件均无法找到时，启动过程会被终止，用户可根据打印的panic信息进行定位和排查

## 配置源接入

配置源的接口信息定义在`config/interface.go`的`ConfigCenter`中，开发者实现了自己的配置源后，可调用`config.RegisterConfigCenter()`完成配置中心的注册。

## 组件注册

组件的接口定义在`config/interface.go`的`Component`中，开发者实现了自己的组件后，可调用`config.RegisterComponent()`完成组件的注册