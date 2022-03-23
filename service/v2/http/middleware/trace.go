package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/glory-go/glory/service/middleware/jaeger"
	"github.com/opentracing-contrib/go-gin/ginhttp"
)

type MTraceMW struct {
}

func NewMTraceMW() *MTraceMW {
	mw := &MTraceMW{}

	return mw
}

func (mw *MTraceMW) HandlerFunc(c *gin.Context) {
	ginhttp.Middleware(jaeger.GetTracer())(c)
}
