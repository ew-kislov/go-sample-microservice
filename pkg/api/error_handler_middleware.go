package api

import (
	"github.com/ew-kislov/go-sample-microservice/pkg/cfg"
	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware(config cfg.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Unwrap()
		code, response := CreateErrorResponse(err, config)

		c.JSON(code, response)
	}
}
