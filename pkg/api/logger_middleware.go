package api

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LoggerMiddleware struct {
	Logger *logrus.Logger
}

type copyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (cw copyWriter) Write(b []byte) (int, error) {
	cw.body.Write(b)
	return cw.ResponseWriter.Write(b)
}

func (middleware *LoggerMiddleware) Log(ctx *gin.Context) {
	requestBody, _ := io.ReadAll(ctx.Request.Body)
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	requestFields := logrus.Fields{
		"method":  ctx.Request.Method,
		"path":    ctx.Request.URL.Path,
		"headers": ctx.Request.Header,
		"params":  ctx.Request.URL.Query(),
		"body":    string(requestBody),
	}
	middleware.Logger.WithContext(ctx).WithFields(requestFields).Infof("--> %s %s", ctx.Request.Method, ctx.Request.URL.Path)

	cw := &copyWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
	ctx.Writer = cw

	ctx.Next()

	responseFields := logrus.Fields{
		"headers": ctx.Writer.Header(),
		"status":  ctx.Writer.Status(),
		"body":    cw.body.String(),
	}
	middleware.Logger.WithContext(ctx).WithFields(responseFields).Infof("<-- %s %s", ctx.Request.Method, ctx.Request.URL.Path)
}
