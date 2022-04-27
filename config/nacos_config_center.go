package config

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/pkg/errors"

	"gopkg.in/yaml.v2"
)

type NacosConfigCenter struct {
	client config_client.IConfigClient
	config *ConfigCenterConfig
	env    string
}

func newNacosConfigCenter(env string) *NacosConfigCenter {
	return &NacosConfigCenter{
		env: env,
	}
}

func (ncc *NacosConfigCenter) Conn(conf *ConfigCenterConfig) error {
	targetEnvNamespaceID, ok := conf.EnvNamespaceIDMap[ncc.env]
	if !ok {
		return errors.Errorf("can't find env read from GLORY_ENV that match any key in config_center config's env_map field, pls check your config!")
	}

	cc := constant.ClientConfig{
		Endpoint:       conf.EndPoint + ":8080",
		NamespaceId:    targetEnvNamespaceID,
		AccessKey:      conf.AccessKeyID,
		SecretKey:      conf.AccessSecret,
		TimeoutMs:      5 * 1000,
		ListenInterval: 30 * 1000,
	}

	client, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig": cc,
	})
	if err != nil {
		return err
	}
	ncc.config = conf
	ncc.client = client
	return nil
}

func (ncc *NacosConfigCenter) LoadConfig() (*ServerConfig, error) {
	content, err := ncc.client.GetConfig(vo.ConfigParam{
		DataId: ncc.config.ServerName,
		Group:  ncc.config.OrgName,
		// todo onchange restart
		OnChange: nil,
	})
	if err != nil {
		return nil, err
	}
	serverConfig := ServerConfig{}
	err = yaml.Unmarshal([]byte(content), &serverConfig)
	return &serverConfig, err
}
