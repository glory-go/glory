package middleware

import "github.com/gin-gonic/gin"

type GLoggerMW struct{}

func (mw GLoggerMW) HandlerFunc(ctx *gin.Context) {}
