package filter_impl

import (
	"strconv"
	"time"
)

import (
	ghttp "github.com/glory-go/glory/http"
	"github.com/glory-go/glory/metrics"
)

const (
	HTTPRequestTimeGaugeName = "http_request_time_ms_"
	HTTPQueryCountName       = "http_query_count_"
	HttpRetCode              = "http_ret_code_"
)

// BasicHttpStatusMiddleware push count and retCode and pending time of this http to default metrics service
func BasicHttpStatusMiddleware(c *ghttp.GRegisterController, f ghttp.HandleFunc) (err error) {
	currTime := time.Now()
	err = f(c)
	expireTime := time.Since(currTime)
	metrics.GaugeSet(HTTPRequestTimeGaugeName+c.Key(), float64(expireTime.Milliseconds()))
	metrics.CounterInsc(HttpRetCode + c.Key() + "_" + strconv.Itoa(int(c.RspCode)))
	metrics.CounterInsc(HTTPQueryCountName + c.Key())
	return err
}
