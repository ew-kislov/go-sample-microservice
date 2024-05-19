package jwtmiddleware

import (
	authservice "github.com/ew-kislov/go-sample-microservice/internal/service/auth_service"

	"github.com/gin-gonic/gin"
)

type JwtMiddleware interface {
	CheckJwt(ctx *gin.Context)
}

type jwtMiddleware struct {
	authService authservice.AuthService
}

const (
	TokenNotProvided = "Token was not provided"
	WrongTokenFormat = "Authorization header must have format: Bearer <Token>"
)
