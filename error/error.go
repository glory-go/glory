package err

import (
	"fmt"
)

// golang doesn't support struct as constant, which is a drawback!
// I just use var to replace const, depressing.
var GloryFrameworkErrorServerUnaryMethodNotFound = &GloryFrameworkError{Code: -201, Msg: "server doesn't have target unary method"}
var GloryFrameworkErrorServerStreamMethodNotFound = &GloryFrameworkError{Code: -202, Msg: "server doesn't have target stream method"}
var GloryFrameworkErrorServerStreamReceiveUniqueKeyNotFound = &GloryFrameworkError{Code: -203, Msg: "server doesn't have target receive unique key"}
var GloryFrameworkErrorServerStreamReceiveChanOffsetNotFound = &GloryFrameworkError{Code: -204, Msg: "server doesn't have target receive channel offset"}

var GloryFrameworkErrorTargetInvokerNotFound = &GloryFrameworkError{Code: -301, Msg: "registry: can't found target provider"}

type GloryFrameworkError struct {
	Code int32
	Msg  string
}

func (g *GloryFrameworkError) Error() string {
	return fmt.Sprintf("Code: %d, Msg: %s", g.Code, g.Msg)
}
