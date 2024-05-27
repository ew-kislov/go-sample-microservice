package jwtmiddleware

import (
	"net/http"
	"strings"

	authservice "github.com/ew-kislov/go-sample-microservice/internal/service/auth_service"
	"github.com/ew-kislov/go-sample-microservice/pkg/api"
	"github.com/ew-kislov/go-sample-microservice/pkg/cfg"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeaderNotProvided = "Authorization header was not provided"
	WrongTokenFormat               = "Authorization header must have format: Bearer <Token>"
)

func JwtMiddleware(authService authservice.AuthService, config *cfg.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")

		if header == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse{Error: AuthorizationHeaderNotProvided})
			return
		}

		splitted := strings.Split(header, " ")

		if len(splitted) != 2 || splitted[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse{Error: WrongTokenFormat})
			return
		}

		token := splitted[1]

		user, err := authService.Authenticate(ctx, token)

		if err != nil {
			code, response := api.CreateErrorResponse(err, config)
			ctx.AbortWithStatusJSON(code, response)

			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}
