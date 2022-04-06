package tools

import "context"

func GetCtxWithLogID() context.Context {
	ctx := context.Background()
	// TODO: 添加logid信息
	return ctx
}
