package ginhttp

import "github.com/glory-go/glory/v2/service"

func init() {
	service.GetService().RegisterService(GetGinHttpService())
}
