package infrastructure

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RequestIdMiddleware struct{}

func (middleware *RequestIdMiddleware) AddRequestId(ctx *gin.Context) {
	requestId := ctx.Request.Header.Get("x-request-id")

	if requestId != "" {
		ctx.Set("requestId", requestId)
	} else {
		ctx.Set("requestId", uuid.New().String())
	}

	ctx.Next()
}
