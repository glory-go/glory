# 配置中心

配置中心提供了一组抽象的接口定义，从而支持开发者按需接入不同的配置来源

## 基础操作

### 配置文件默认位置及格式

配置文件默认的读取位置为`./glory/glory.${env}.yaml`，其中`env`为环境名称，将从`GLORY_ENV`环境变量中读取。用户可通过调用`config.ChangeDefaultConfigPath()`来指定新的文件位置，支持yaml后缀的文件，且`${env}`会被默认加入到读取配置文件的路径中

若启动时找不到带环境后缀的配置文件时，将默认读取不带后缀的配置文件。当所有配置文件均无法找到时，启动过程会被终止，用户可根据打印的panic信息进行定位和排查

配置文件由多个部分组成，各个部分的格式为：

```
config_key:
    key: value
```

其中，config_key为各个组件的名称，而key和value则是真正用于初始化的组件配置

### 配置文件中的变量

glory支持用户在配置文件中定义变量作为配置中心的输入，这些变量将在读取配置时，原封不动地传递给配置中心的实现，由配置中心的实现负责识别变量并加以使用

配置文件中定义的变量格式为`$config_center_name{VARS}`，其中`config_center`为该配置使用的配置中心，其中`file`、`env`和空的内容为保留类型，不允许外部配置中心注册使用这些名称；`VARS`为JSON格式的变量信息，变量仅支持数组类型。例如`$env{["key", true]}`则会依次将`"key"`和`true`传入`env`配置中心的读取参数中

当用户传入的参数过多】过少或者类型不匹配时，程序将在初始化时panic退出。目前仅支持用户将变量定义在字符串的字段中，且允许用户嵌套在结构体字段中使用，例如：

```(yaml)
key1: $center1{["param1", "param2"]}
key2:
    key2_1: $center2{["param1"]}
    key2_2:
        - $center3{["param1"]} # 此种写法不受支持
key3:
    - key3_1: $center1{["param1", "param2"]} # 此种写法不受支持
```

#### 使用环境变量替换配置

环境变量也是配置中心的一种，是由glory默认提供的配置中心，只接受一个字符串类型的参数，代表该配置读取的环境变量名称信息


## 启动流程

### 配置源加载

配置源本身也是一个组件，但由于其存储了用户配置，因此将其加载与其他组件的配置独立开来，并限制了其配置的来源方式。配置源加载遵循以下步骤：

|-------------|     |----------------------------|      |------------- |
|             |     |                            |      |              |
| 读取文件配置 | ->  | 基于环境变量更新初始配置内容 |  ->  | 初始化配置源 |
|             |     |                            |      |              |
|-------------|     |----------------------------|      |--------------|

### 用户组件加载

完成配置源加载后，就可以基于配置源初始化用户注册的组件了。组件的加载遵循以下步骤：

|------------------------|      |------------|
|                        |      |            |
| 基于配置中心更新配置内容 |  ->  | 初始化组件 |
|                        |      |            |
|------------------------|      |------------|

## 社区模块接入

### 配置源与组件的接入

配置源的接口信息定义在`config/interface.go`的`ConfigCenter`中，开发者实现了自己的配置源后，可调用`config.RegisterConfigCenter()`完成配置中心的注册

组件的接口定义在`config/interface.go`的`Component`中，开发者实现了自己的组件后，可调用`config.RegisterComponent()`完成组件的注册

### 组件的使用

glory不提供组件的获取方法，因此需要组件的提供方实现单例模式并透出组件实例。这是因为不同的组件拥有各自的能力，因此我们只帮助用户完成组件的初始化，具体的使用则需要调用具体的组件实例。