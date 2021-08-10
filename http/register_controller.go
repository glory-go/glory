package http

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/glory-go/glory/log"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

var v *validator.Validate
var defaultSchemaDecoder *schema.Decoder

func init() {
	defaultSchemaDecoder = schema.NewDecoder()
	defaultSchemaDecoder.IgnoreUnknownKeys(true)
	v = validator.New()
}

type GRegisterController struct {
	// Req 存放请求的数据
	Req interface{}
	// Rsp 用来设置回包
	Rsp interface{}
	// R 暴露 http.Request
	R *http.Request
	// W 暴露 http.ResponseWriter
	W http.ResponseWriter
	// Ctx 传递出来的 context 供 handle 使用
	Ctx context.Context
	// VarsMap url 变量
	VarsMap map[string]string
	// RspCode 返回http状态码,不设则使用默认成功、失败状态码
	RspCode httpCode
	// IfNeedWrite 为false时，框架将不会向 http.ResponseWriter中写入值，用户需自主完成返回值的写入
	// 默认为true
	IfNeedWrite bool
}

type GRegisterWSController struct {
	WSConn *websocket.Conn
	R      *http.Request
}

// 获取请求参数并进行参数校验
func (trc *GRegisterController) GetReqData(r *http.Request) error {
	var err error

	// 解析qury。body数据
	if err = r.ParseForm(); err != nil {
		return err
	}
	// 解析 POSt 请求的 header 不是 x-www-form-urlencoded 的情况
	if r.Header == nil {
		return errors.New("r.Header nil ptr error")
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("read body err:", err)
		return err
	}
	if len(data) != 0 {
		if err := json.Unmarshal(data, trc.Req); err != nil {
			log.Error("request unmarshal err:", err)
			return err
		}
	}

	// 解构参数
	if err := defaultSchemaDecoder.Decode(trc.Req, r.Form); err != nil {
		log.Error("decode r.form err:", err)
		return err
	}

	// 参数校验，具体语法参考：https://godoc.org/github.com/go-playground/validator
	if err = v.Struct(trc.Req); err != nil {
		log.Error("validator check failed:", err)
		return err
	}

	return nil
}

// Key is made of $(path)_$(method)
func (trc *GRegisterController) Key() string {
	return trc.R.URL.Path + "_" + trc.R.Method
}
