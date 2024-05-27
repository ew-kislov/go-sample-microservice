package getmecontroller

import (
	"errors"
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
func (*getMeController) GetMe(ctx *gin.Context) {
	userUntyped, _ := ctx.Get("user")
	user, ok := userUntyped.(*authservice.User)

	if !ok {
		_ = ctx.Error(errors.New("Could not cast ctx.user to User type"))
		return
	}

	userResponse := UserResponse{
		Id:          user.Id,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Username:    user.Username,
	}

	ctx.JSON(http.StatusOK, userResponse)
}
