# 配置管理

## 使用

**初始化前，应先初始化`config_manager`包，且调用config包的`Init()`方法完成手动初始化**

应用的main函数中，引入各个模块(mysql、redis等)后，并确保各个模块完成了注册(一般在init方法中，或需要手动Init)，调用`InitModules(modules ...string) error`方法进行各个模块的初始化操作

## 配置管理

### 配置文件路径

`config/config_{GLORY_ENV}.v2.yaml`，`GLORY_ENV`为环境变量，配置中心的配置会基于该环境变量进行读取

### 自定义配置内容

配置文件中，对于每一个配置中心来说，内容为：

```
{config_name}:
    {sub_name}:
        config_source: env
        {custom_config}
```

`config_name`为配置的名称，注册配置时的名字与该值需保持一致；

`sub_name`为子配置名，例如：mysql中一般是主从集群，可以为不同的集群设置不同的配置；

`config_source`为配置的解析来源，名称需要与`config_manager`包中注册的配置保持一致。需要被替换的配置需要以`group$$key`方式填写内容，此时配置将会从配置中心加载并替换。配置仅替换string类型的内容，支持递归解析；

`{custom_config}`为自定义的配置内容，key支持string类型，value可以为任意类型。

### 配置中心注册

需要解析配置的服务，需调用`RegisterConfig(name string, builder ComponentBuilder) error`方法。`builder`的第一个参数为`sub_name`，第二个参数为解析完成的配置
