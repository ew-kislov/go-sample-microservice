package getmecontroller

import (
	"net/http"

	authservice "github.com/ew-kislov/go-sample-microservice/internal/service/auth_service"
	"github.com/gin-gonic/gin"
)

// GetMe godoc
//
// @Summary			  Get user info
// @Accept			  application/json
// @Produce			  application/json
// @Tags			    Auth
// @Security			BearerAuth
// @Success			  200 {object} UserResponse
// @Failure			  401 {object} api.ErrorResponse
// @Failure			  401 {object} api.ErrorResponse
// @Failure			  404 {object} api.ErrorResponse
// @Router			  /auth/me [get]
func (controller *getMeController) GetMe(ctx *gin.Context) {
	userUntyped, _ := ctx.Get("user")
	user, _ := userUntyped.(*authservice.User)

	userResponse := UserResponse{
		Id:          user.Id,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Username:    user.Username,
	}

	ctx.JSON(http.StatusOK, userResponse)
}
