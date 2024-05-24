package signupcontroller

import (
	authservice "github.com/ew-kislov/go-sample-microservice/internal/service/auth_service"

	"github.com/gin-gonic/gin"
)

type SignUpRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Username    string `json:"username" binding:"required"`
	DisplayName string `json:"displayName" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type SignUpResponse struct {
	UserId int64  `json:"userId"`
	Token  string `json:"token"`
}

type signUpController struct {
	authService authservice.AuthService
}

type SignUpController interface {
	SignUp(ctx *gin.Context)
}

func NewSignUpController(authService authservice.AuthService) SignUpController {
	return &signUpController{authService}
}
