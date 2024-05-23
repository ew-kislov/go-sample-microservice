package jwtmiddleware

import (
	"net/http"
	"strings"

	authservice "github.com/ew-kislov/go-sample-microservice/internal/service/auth_service"
	"github.com/ew-kislov/go-sample-microservice/pkg/api"

	"github.com/gin-gonic/gin"
)

func NewJwtMiddleware(authService authservice.AuthService) JwtMiddleware {
	return &jwtMiddleware{authService}
}

func (middleware *jwtMiddleware) CheckJwt(ctx *gin.Context) {
	token, err := extractBearerToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
		return
	}

	user, err := middleware.authService.Authenticate(ctx, token)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.Set("user", user)
	ctx.Next()
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", api.ApiError{Code: http.StatusUnauthorized, Message: TokenNotProvided}
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 || jwtToken[0] != "Bearer" {
		return "", api.ApiError{Code: http.StatusUnauthorized, Message: WrongTokenFormat}
	}

	return jwtToken[1], nil
}
