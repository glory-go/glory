# logrus

[logrus](https://github.com/sirupsen/logrus)是go语言下一款十分热门的日志组件。glory在logrus提供的hook的基础上进行了简单封装，帮助用户通过配置文件更好地扩展logrus的日志打印

## 定义并使用Hook

用户可以在代码中注册日志hooks，并在配置文件中定义hooks配置，从而支持输出到文件、error告警等功能。glory集成了一些常用的hook及其配置定义，一个可供参考的配置如下：

```
logrus: # logrus组件名称，不可修改
  hook1: # hook名称，由用户自定义
    type: file # 注册的hook类型
    # 后面则是hook自定义的配置内容
    level_path:
      "0": "/var/log/panic.log"
```
