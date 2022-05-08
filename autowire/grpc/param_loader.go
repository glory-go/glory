package grpc

import (
	"errors"
)

import (
	"google.golang.org/grpc"
)

import (
	"github.com/glory-go/glory/autowire"
	"github.com/glory-go/glory/config"
)

type paramLoader struct {
}

/*
Load support load grpc field:
```go
ResourceServiceClient resources.ResourceServiceClient `grpc:"resource-service"`
```go

from:

```yaml
autowire:
  grpc:
    resource-service:
      address: 127.0.0.1:8080
```

Make Dial and generate *grpc.ClientConn as param
*/
func (p *paramLoader) Load(_ *autowire.StructDescriber, fi *autowire.FieldInfo) (interface{}, error) {
	if fi == nil {
		return nil, errors.New("not supported")
	}
	grpcConfig := &Config{}
	if err := config.LoadConfigByPrefix("autowire.grpc."+fi.TagValue, grpcConfig); err != nil {
		return nil, err
	}
	return grpc.Dial(grpcConfig.Address, grpc.WithInsecure())
}

type Config struct {
	Address string
}
