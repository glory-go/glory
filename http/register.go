package http

import (
	"errors"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/glory-go/glory/log"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var rspImpPackage RspPackage

func init() {
	// 框架默认使用default空回包
	rspImpPackage = &DefaultRspPackage{}
}

// 添加默认的 header
func writeDefaultHeader(rsp http.ResponseWriter, req *http.Request) {
	// 防止 xss 攻击
	rsp.Header().Add("X-Content-Type-Options", "nosniff")
	// 设置 Content-Type
	ct := rsp.Header().Get("Content-Type")
	if ct == "" {
		ct = req.Header.Get("Content-Type")
		if req.Method == "GET" || ct == "" {
			ct = "application/json"
		}
		rsp.Header().Add("Content-Type", ct)
	}
}

// 处理函数
func getGloryHttpHandler(handler func(*GRegisterController) error, req, rsp interface{}, filters []Filter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		retPkg := rspImpPackage

		// recovery
		defer func() {
			if e := recover(); e != nil {
				buf := make([]byte, 1024)
				buf = buf[:runtime.Stack(buf, false)]
				log.Panic(e, buf)
				retPkg.SetErrorPkg(w, errors.New("server panic"), DefaultHttpErrorCode)
			}
		}()

		writeDefaultHeader(w, r)

		// 创建针对当前接口的过滤器chain
		chain := Chain{}
		chain.AddFilter(filters) // 注册过滤器

		tRegisterController := GRegisterController{
			Ctx:         r.Context(),
			R:           r,
			W:           w,
			RspCode:     UnsetHttpCode,
			IfNeedWrite: true,
		}

		// 处理 req
		if req != nil {
			requestType := reflect.TypeOf(req).Elem()
			tRegisterController.Req = reflect.New(requestType).Interface()
			if err := tRegisterController.GetReqData(r); err != nil {
				retPkg.SetErrorPkg(w, err, tRegisterController.RspCode) // go
				return
			}
		}

		// 处理 rsp
		if rsp != nil {
			rspType := reflect.TypeOf(rsp).Elem()
			tRegisterController.Rsp = reflect.New(rspType).Interface()
		}

		// 执行业务函数
		if err := chain.Handle(&tRegisterController, handler); err != nil {
			retPkg.SetErrorPkg(w, err, tRegisterController.RspCode)
			return
		}

		// 用户如果自己处理了
		if !tRegisterController.IfNeedWrite {
			return
		}

		// 最终回包
		retPkg.SetSuccessPkg(w, tRegisterController.Rsp, tRegisterController.RspCode)

		return
	}
}

func getGloryWSHandler(handler func(*GRegisterWSController)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		retPkg := rspImpPackage

		// recovery
		defer func() {
			if e := recover(); e != nil {
				buf := make([]byte, 1024)
				buf = buf[:runtime.Stack(buf, false)]
				log.Panic(e, buf)
				retPkg.SetErrorPkg(w, errors.New("server panic"), DefaultHttpErrorCode)
			}
		}()
		// 升级接口到websocket
		conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if err != nil {
			return
		}

		tRegisterController := GRegisterWSController{
			WSConn: conn,
			R:      r,
		}
		handler(&tRegisterController)
	}
}

// 自定义回包函数
func RegisterRspPackage(rspUserImpPackage RspPackage) {
	rspImpPackage = rspUserImpPackage
}

func checkMethod(method string) (string, bool) {
	if method == "GET" || method == "POST" || method == "DELETE" ||
		method == "PATCH" || method == "PUT" {
		return method, true
	}
	if method == "get" || method == "post" || method == "delete" ||
		method == "patch" || method == "put" {
		return strings.ToUpper(method), true
	}
	return "", false
}

// 入口函数
func RegisterRouter(path string, r *mux.Router, handler func(*GRegisterController) error, req, rsp interface{}, method string, filters []Filter) {
	gloryHttpHandler := getGloryHttpHandler(handler, req, rsp, filters)
	afterCheckedMethod, ok := checkMethod(method)
	if !ok {
		log.Panic("RegisterRouter: method unsupported")
		return
	}
	r.HandleFunc(path, gloryHttpHandler).Methods(afterCheckedMethod)
}

func RegisterWSRouter(path string, r *mux.Router, handler func(*GRegisterWSController)) {
	trpcWSHandler := getGloryWSHandler(handler)
	r.HandleFunc(path, trpcWSHandler)
}

func NewHttpRegister() {

}
