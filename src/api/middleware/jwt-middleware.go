package middleware

import (
	"go-sample-microservice/src/infrastructure"
	"go-sample-microservice/src/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type JwtMiddleware struct {
	AuthService service.AuthService
}

func (middleware *JwtMiddleware) CheckJwt(ctx *gin.Context) {
	token, err := extractBearerToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	user, err := middleware.AuthService.Authenticate(ctx, token)

	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.Set("user", user)
	ctx.Next()
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", infrastructure.ApiError{Code: http.StatusUnauthorized, Message: "Token was not provided"}
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 || jwtToken[0] != "Bearer" {
		return "", infrastructure.ApiError{Code: http.StatusUnauthorized, Message: "Authorization header must have format: Bearer <Token>"}
	}

	return jwtToken[1], nil
}
