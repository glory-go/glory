package glory

import (
	"fmt"
	hessian "github.com/apache/dubbo-go-hessian2"
)

type ErrorCode int32

const (
	GloryErrorUserErrorCode = -1004
)

var (
	GloryErrorNoErr = NewError(0, "")

	// client error
	GloryErrorConnErr          = NewError(-1001, "conntion error")
	GloryErrorTimeoutErr       = NewError(-1002, "waiting for response time out")
	GloryErrorEmptyResponseErr = NewError(-1003, "get empty response")

	// server error
	GloryErrorServerInternalErr = NewError(-2001, "server internal error")

	// protocol error
	GloryErrorProtocol     = NewError(-3001, "protocol error")
	GloryErrorVersion      = NewError(-3002, "glory version error")
	GloryErrorPkgTypeError = NewError(-3003, "glory package type error")

	// registry error (now useless)
	GloryErrorRegistryProviderNotFound = NewError(-4001, "registry provider not found in register center")
	GloryErrorRegistruyServerNotFound  = NewError(-4002, "registry server not found")

	GloryErrorUnknown = NewError(-5000, "unknow error type")
	//GloryErrorUserError = NewError(-2, "")
)

func NewGloryErrorUserError(err error) *Error {
	return NewError(GloryErrorUserErrorCode, err.Error())
}

func init() {
	hessian.RegisterPOJO(&Error{})
}

// Error is a field of rsp pkg
type Error struct {
	Code int32  // Code shows glory error code
	Msg  string // Msg shows error message
}

func (e *Error) Error() string {
	return fmt.Sprintf("Code: %d, Msg: %s", e.Code, e.Msg)
}

func (Error) JavaClassName() string {
	return "GLORY_Error"
}

func NewError(code int32, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}
