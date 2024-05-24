package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIdMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestId := ctx.Request.Header.Get("x-request-id")

		if requestId != "" {
			ctx.Set("requestId", requestId)
			ctx.Writer.Header().Set("x-request-id", requestId)
		} else {
			requestId := uuid.New().String()
			ctx.Set("requestId", requestId)
			ctx.Writer.Header().Set("x-request-id", requestId)
		}

		ctx.Next()
	}
}
