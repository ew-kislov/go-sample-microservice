package controller

import (
	"go-sample-microservice/src/service"

	"github.com/gin-gonic/gin"
)

type GetMeController struct {
}

type UserResponse struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

func (controller *GetMeController) GetMe(ctx *gin.Context) {
	userUntyped, _ := ctx.Get("user")
	user, _ := userUntyped.(*service.User)

	userResponse := UserResponse{
		Id:          user.Id,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Username:    user.Username,
	}

	ctx.JSON(200, gin.H{"success": true, "data": userResponse})
}
