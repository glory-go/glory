# 配置中心

## 使用

应用的main函数中，引入各个配置中心(如：nacos)后，并确保各个模块完成了注册(一般在init方法中，或需要手动Init)，调用`InitConfigCenter(configCenterName ...string) error`方法进行各个模块的初始化操作

## 配置管理

### 配置文件路径

`config/config_center_{GLORY_ENV}.yaml`，`GLORY_ENV`为环境变量，配置中心的配置会基于该环境变量进行读取

### 自定义配置内容

配置文件中，对于每一个配置中心来说，内容为：

```
{config_center_name}:
    config_source: env
    {custom_config}
```

`config_center_name`为配置中心的名称，注册配置中心时的名字与该值需保持一致；

`config_source`为配置的解析来源，目前只支持`env`，配置时，所有自定义配置的value从环境变量中读取并替换；

`{custom_config}`为自定义的配置内容，**key与value均需要为string类型**

### 配置中心注册

提供配置中心服务的一方，需调用`RegisterConfigBuilder(configCenterName string, builder ConfigCenterBuilder) error`方法。`builder`的入参为解析完成的配置，它返回的`ConfigCenter`即为**初始化完成的配置中心服务**
