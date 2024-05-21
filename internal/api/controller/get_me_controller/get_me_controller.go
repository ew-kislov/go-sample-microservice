package getmecontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	userservice "github.com/ew-kislov/go-sample-microservice/internal/service/user_service"
)

func (controller *getMeController) GetMe(ctx *gin.Context) {
	userUntyped, _ := ctx.Get("user")
	user, _ := userUntyped.(*userservice.User)

	userResponse := UserResponse{
		Id:          user.Id,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Username:    user.Username,
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": userResponse})
}