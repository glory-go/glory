package config

import (
	"fmt"
	"reflect"
)

// ServerConfig 整个服务的配置
type ServerConfig struct {
	// aliCloud common
	AliCloudCommonConfig *AliCloudCommonConfig `yaml:"alicloud"`
	// ServiceConfigs 服务端单个Service配置
	ServiceConfigs map[string]*ServiceConfig `yaml:"provider"`
	// LogConfigs 日志配置
	LogConfigs map[string]*LogConfig `yaml:"log"`
	// MetricsConfigs 数据上报配置
	MetricsConfigs []*MetricsConfig `yaml:"metrics"`
	// ClientConfig 客户端配置 目前goonline服务端的主调只涉及到grpcclient、jsonrpcclient和gloryclient，其中grpc为使用原生client
	ClientConfig map[string]*ClientConfig `yaml:"consumer"`
	// RegistryConfig 注册中心配置
	RegistryConfig map[string]*RegistryConfig `yaml:"registry"`
	// MysqlConfigs Mysql 数据库配置
	MysqlConfigs map[string]*MysqlConfig `yaml:"mysql"`
	// RedisConfig redis 配置
	RedisConfig map[string]*RedisConfig `yaml:"redis"`
	// K8SConfig K8s配置
	K8SConfig map[string]*K8SConfig `yaml:"k8s"`
	// MQConfig MQ配置
	MQConfig map[string]*MQConfig `yaml:"mq"`
	// OssConfigs 持久云存储配置
	OssConfigs map[string]*OssConfig `yaml:"oss"`
	// MongoDBConfig Mongo 数据库配置
	MongoDBConfig map[string]*MongoDBConfig `yaml:"mongodb"`
	// rpc Filter 配置
	FilterConfigMap map[string]*FilterConfig `yaml:"filter"`
	// rpc Filter 配置
	UserConfig map[interface{}]interface{} `yaml:"config"`

	OrgName    string `yaml:"org_name"` // 可选 classroom|ide|children|goonline: goonline为公共服务，比如前端数据上报
	ServerName string `yaml:"server_name"`
}

func (s *ServerConfig) Print() {
	for _, v := range s.ServiceConfigs {
		fmt.Print(v.Protocol, v.Port)
	}
}

func (s *ServerConfig) GetAppKey() string {
	return s.OrgName + "_" + s.ServerName
}

// checkAndFixConfigs 所有子Config的检查和配置填充
// 目前只支持:
// map[string]*Config map[string]map[string]*Config 和 []*Config 类型和*Struct类型 的配置填充
func (s *ServerConfig) checkAndFix() {
	val := reflect.ValueOf(s).Elem()
	typ := reflect.TypeOf(s).Elem()
	num := val.NumField()
	for i := 0; i < num; i++ {
		if typ.Field(i).Type.Kind() == reflect.Ptr {
			if v, ok := val.Field(i).Interface().(config); ok && v != nil {
				val := reflect.ValueOf(v).Elem()
				// nil config continue
				if val.Kind() == reflect.Invalid {
					continue
				}
				v.checkAndFix()
			}
			fmt.Println("error: the val is no config")
		}
		if typ.Field(i).Type.Kind() == reflect.Map {
			//log.Println("field ", i, "is map with name ", typ.Field(i).Name)
			iter := val.Field(i).MapRange()
			for iter.Next() {
				v1, ok1 := iter.Value().Interface().(config)
				if ok1 {
					v1.checkAndFix()
				}
				v2, ok2 := iter.Value().Interface().(map[string]config)
				if ok2 {
					for k := range v2 {
						v2[k].checkAndFix()
					}
				}
				if !ok1 && !ok2 {
					fmt.Println("error,the val is not config")
				}

			}
		}
		if typ.Field(i).Type.Kind() == reflect.Slice {
			//log.Println("field ", i, "is slice with name ", typ.Field(i).Name)
			configItemNum := val.Field(i).Len()
			for idx := 0; idx < configItemNum; idx++ {
				configItem, ok := val.Field(i).Index(idx).Interface().(config)
				if !ok {
					fmt.Println("error,the val is not config")
					continue
				}
				configItem.checkAndFix()
			}
		}
	}
	if s.OrgName == "" {
		panic("please add your service org_name in config file from: classroom|ide|children|goonline")
	}

	if s.ServerName == "" {
		panic("please add your server_name in config file!")
	}
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		ServiceConfigs: make(map[string]*ServiceConfig),
		LogConfigs:     make(map[string]*LogConfig),
		MetricsConfigs: make([]*MetricsConfig, 0),
		MysqlConfigs:   make(map[string]*MysqlConfig),
		K8SConfig:      make(map[string]*K8SConfig),
		RedisConfig:    make(map[string]*RedisConfig),
		MQConfig:       map[string]*MQConfig{},
		OssConfigs:     make(map[string]*OssConfig),
		ClientConfig:   make(map[string]*ClientConfig),
		RegistryConfig: make(map[string]*RegistryConfig),
	}
}
