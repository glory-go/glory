package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RspPackage interface {
	SetSuccessPkg(w http.ResponseWriter, msg interface{}, retCode httpCode) // 成功回包
	SetErrorPkg(w http.ResponseWriter, err error, retCode httpCode)         // 错误回包
}

// DefaultFomattedRspPackage
// 框架提供的默认格式化回包，包含三个字段,与下面的 DefaultRspPackage 选择使用
type DefaultFomattedRspPackage struct {
	Retcode int32       `json:"retcode"`
	Retmsg  string      `json:"retmsg"`
	Result  interface{} `json:"result"`
}

func (rpkg *DefaultFomattedRspPackage) SetSuccessPkg(w http.ResponseWriter, result interface{}, retCode httpCode) {
	rpkg.Retmsg = "ok"
	rpkg.Retcode = 0
	rpkg.Result = result
	if retCode == UnsetHttpCode {
		retCode = DefaultHttpSuccessCode
	}
	w.WriteHeader(int(retCode))
	var err error
	rpkgBody := make([]byte, 0)
	if rpkgBody, err = json.Marshal(*rpkg); err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf(`{"retcode":%d, "retmsg":"fatel err: marshal rspPackage failed", "result": null}`, -1)))
		return
	}
	_, _ = w.Write(rpkgBody)
}

func (rpkg *DefaultFomattedRspPackage) SetErrorPkg(w http.ResponseWriter, err error, retCode httpCode) {
	rpkg.Retmsg = err.Error()
	rpkg.Retcode = -1
	rpkg.Result = nil
	if retCode == UnsetHttpCode {
		retCode = DefaultHttpSuccessCode
	}
	w.WriteHeader(int(retCode))
	rpkgBody := make([]byte, 0)
	if rpkgBody, err = json.Marshal(*rpkg); err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf(`{"retcode":%d, "retmsg":"fatel err: marshal rspPackage failed", "result": null}`, -1)))
		return
	}
	_, _ = w.Write(rpkgBody)
}

//

// DefaultRspPackage
// 框架提供的默认空回包，包含一个字段
type DefaultRspPackage struct {
	Result interface{} `json:"result"`
	OK     bool        `json:"ok"`
}

func (rpkg *DefaultRspPackage) SetSuccessPkg(w http.ResponseWriter, result interface{}, retCode httpCode) {
	rpkg.Result = result
	rpkg.OK = true
	if retCode == UnsetHttpCode {
		retCode = DefaultHttpSuccessCode
	}
	w.WriteHeader(int(retCode))
	var err error
	rpkgBody := make([]byte, 0)
	if rpkgBody, err = json.Marshal(*rpkg); err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf(`{"retcode":%d, "retmsg":"fatel err: marshal rspPackage failed", "result": null}`, -1)))
		return
	}
	_, _ = w.Write(rpkgBody)
}

func (rpkg *DefaultRspPackage) SetErrorPkg(w http.ResponseWriter, err error, retCode httpCode) {
	rpkg.Result = err.Error()
	rpkg.OK = false
	if retCode == UnsetHttpCode {
		retCode = DefaultHttpErrorCode
	}
	w.WriteHeader(int(retCode))
	rpkgBody := make([]byte, 0)
	if rpkgBody, err = json.Marshal(*rpkg); err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf(`{"retcode":%d, "retmsg":"fatel err: marshal rspPackage failed", "result": null}`, -1)))
		return
	}
	_, _ = w.Write(rpkgBody)
}
