# Glory MQ 抽象层

## 如何使用

1. 找到业务上使用到的mq类型的抽象实现，例如基于阿里云的rocketmq的实现：https://github.com/glory-go/mq

2. 在具体的mq实现中，找到：
   - 用来注册的init函数，该init函数将mq的初始化方法注册到了mq的类型映射关系中
   - 配置文件的定义

3. 将目标的mq实现中，init函数所在包import到你的业务代码中，例如：`import _ "github.com/glory-go/mq/rocketmq/aliyun"`

4. 根据目标mq实现中的mq配置，按照以下的方式编写你的配置文件

```(yaml)
mq:
    mq_name: # 自定义的名称
        type: aliyun_rocketmq # 目标mq实现中指定的类型
        config_source: env # 可选，包含该配置时，将从指定来源中解析mq配置，目前仅支持环境变量
        config: # 以下内容为用户自定义mq所消费
            key: value
```

5. 调用`mq.Init()`方法，初始化mq配置

完成以上步骤后，可在代码中通过`mq.GetMQInstance("mq_name")`，获取你在配置中定义的mq实例

## 支持能力

`MQService`的能力如下：

```(golang)
type MQService interface {
	Connect() error
	Send(topic string, msg []byte) (msgID string, err error)
	DelaySend(topic string, msg []byte, handleTime time.Time) (msgID string, err error)
	RegisterHandler(topic string, handler MQMsgHandler)
}
```

需要注意的是，你需要仔细了解引入的mq具体实现，从而了解对方是否支持全部的能力，例如：未加插件的rabbitmq是不支持发送延迟消息的，而基于阿里云的rocketmq则天然支持延迟消息
