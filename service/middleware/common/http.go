package common

import (
	ghttp "github.com/glory-go/glory/http"
	"github.com/urfave/negroni"
)

var (
	gloryMWs   []ghttp.Filter
	negroniMWs []negroni.Handler
)

func RegisterGloryMWs(filters ...ghttp.Filter) {
	gloryMWs = append(gloryMWs, filters...)
}

func RegisterNegroniMWs(filters ...negroni.Handler) {
	negroniMWs = append(negroniMWs, filters...)
}

func GetGloryMWs() []ghttp.Filter {
	return gloryMWs
}

func GetNegroniMWs() []negroni.Handler {
	return negroniMWs
}
