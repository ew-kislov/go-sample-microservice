package getmecontroller

import (
	"net/http"

	authservice "github.com/ew-kislov/go-sample-microservice/internal/service/auth_service"
	"github.com/gin-gonic/gin"
)

func (controller *getMeController) GetMe(ctx *gin.Context) {
	userUntyped, _ := ctx.Get("user")
	user, _ := userUntyped.(*authservice.User)

	userResponse := UserResponse{
		Id:          user.Id,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Username:    user.Username,
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": userResponse})
}
