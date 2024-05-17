package controller

import (
	"go-sample-microservice/src/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignUpRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Username    string `json:"username" binding:"required"`
	DisplayName string `json:"displayName" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type SignUpController struct {
	AuthService service.AuthService
}

func (controller *SignUpController) SignUp(ctx *gin.Context) {
	var body SignUpRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response, err := controller.AuthService.SignUp(ctx, service.CreateUserParams(body))

	if err != nil {
		ctx.JSON(400, gin.H{"success": false, "message": err.Error()})
	} else {
		ctx.JSON(200, gin.H{"success": true, "data": response})
	}
}
