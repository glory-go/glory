package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/glory-go/glory/v2/config"
	"github.com/glory-go/glory/v2/sub"
	mockprovider "github.com/glory-go/glory/v2/sub/mock"
	"github.com/golang/mock/gomock"
)

func Test_Sub(t *testing.T) {
	setUp()
	ctrl := gomock.NewController(t)

	/** 注册Provider **/
	mockProvider := mockprovider.NewMockSubProvider(ctrl)
	mockProvider.EXPECT().Name().Return("mock_sub1").AnyTimes()
	mockProvider.EXPECT().Init(gomock.Any()).Return(nil).Times(1)
	mockProvider.EXPECT().Run().Return(nil).Times(1)
	sub.GetSub().RegisterSubProvider(mockProvider)

	// 配置初始化
	config.Init()
	// 启动
	sub.GetSub().Run()
}

func setUp() {
	path := fmt.Sprintf("/tmp/test_sub_%d.yaml", time.Now().Unix())
	config.ChangeDefaultConfigPath(path)
	content := `
service:
  sub:
    mock_sub1:
      test: 1
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
