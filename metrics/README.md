# metrics 服务端数据上报使用

### 1. 一般配置
```yaml
# ide/classroom/children/goonline 分别代表对应组织
# 默认为goonline
org_name: ide

# server_name 对应当前服务名
# 默认是default_server_name
server_name: log-demo-client
 
# 日志配置
log :
  "console-log":
    log_type: console
    level: debug

# 指标配置！
metrics:
# 可开启多个配置，这里有几个条目，对应开启几个metrics service
# 每个metrics service 都有自己的jobname，默认为$(org_name).$(server_name)
  - client_port: 8082
    metrics_type: prometheus_client
    client_path: /prometheus
    action_type: pull
  - gateway_host: localhost
    config_source: file # 可通过config_source: env 尝试从环境读取
    gateway_port: 9091
    action_type: push
    metrics_type: prometheus_client
```
上述配置中：$(org_name)_$(server_name) \
默认为metrics配置列表中每个service的默认jobname,像这样
不指定jobname，会将数据依次上报到每个service，上报后的jobname字段为默认的。

上报数据

```
metrics.CounterInsc("in main")
```
### 2. jobname指定服务上报

如果希望在代码中需要指定数据上报的service，就需要在配置中为特定service增加jobname字段

注意！jobname字段不可出现除`.` `-` `$` `_` 之外的特殊字符，并且这几个都会变成 `_`
因为prometheus只支持`-`
```yaml
metrics:
# 无jobname service
  - gateway_host: PROMETHEUS_PUSH_GATEWAY_IP
    config_source: env
    gateway_port: PROMETHEUS_PUSH_GATEWAY_PORT
    action_type: push
    metrics_type: prometheus_client

# 指定 job_name: IDE.frontEndUpload
  - gateway_host: PROMETHEUS_PUSH_GATEWAY_IP
    config_source: env
    gateway_port: PROMETHEUS_PUSH_GATEWAY_PORT
    action_type: push
    metrics_type: prometheus_client
    job_name: IDE.frontEndUpload

# 指定 job_name: IDE.frontEndUpload
  - gateway_host: PROMETHEUS_PUSH_GATEWAY_IP
    config_source: env
    gateway_port: PROMETHEUS_PUSH_GATEWAY_PORT
    action_type: push
    metrics_type: prometheus_client
    job_name: Classroom.frontEndUpload
```

```go
// 数据上报到jobname为IDE.frontEndUpload的service, 上报条目的jobname字段也为IDE.frontEndUpload
metrics.GaugeSet(req.MetricsName, req.GaugeValue, "IDE.frontEndUpload")

// 数据上报到所有默认jobname（无配置jobname的service), 上报条目的jobname字段为$(org_name).$(server_name)
metrics.CounterInsc(HttpRetCode + "_" + strconv.Itoa(int(c.RspCode)))
```

### 3. 数据上报http中间件

目前已经集成至http/filter_impl

```go
// BasicHttpStatusMiddleware push count and retCode and pending time of this http to default metrics service
func BasicHttpStatusMiddleware(c *ghttp.GRegisterController, f ghttp.HandleFunc) (err error) {
	currTime := time.Now()
	err = f(c)
	expireTime := time.Now().Sub(currTime)
	metrics.GaugeSet(HTTPRequestTimeGaugeName+c.Key(), float64(expireTime.Milliseconds()))
	metrics.CounterInsc(HttpRetCode + c.Key() + "_" + strconv.Itoa(int(c.RspCode)))
	metrics.CounterInsc(HTTPQueryCountName + c.Key())
	return err
}
```
可直接在开启http service 时引入，默认数据上报jobname=$(org_name)_$(server_name)