package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/glory-go/glory/v2/config"
	"github.com/glory-go/glory/v2/service"
	mockservice "github.com/glory-go/glory/v2/service/mock"
	"github.com/golang/mock/gomock"
)

func Test_Service(t *testing.T) {
	setUp()
	ctrl := gomock.NewController(t)

	/** 注册服务 **/
	mockSrv := mockservice.NewMockService(ctrl)
	mockSrv.EXPECT().Name().Return("mock_srv1").AnyTimes()
	mockSrv.EXPECT().Init(gomock.Any()).Return(nil).Times(1)
	mockSrv.EXPECT().Run().Return(nil).Times(1)
	service.GetService().RegisterService(mockSrv)

	// 配置初始化
	config.Init()
	// 启动
	service.GetService().Run()
}

func setUp() {
	path := fmt.Sprintf("/tmp/test_service_%d.yaml", time.Now().Unix())
	config.ChangeDefaultConfigPath(path)
	content := `
service:
  mock_srv1:
    a: 1
`
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(content)
	if err != nil {
		panic(err)
	}
}
