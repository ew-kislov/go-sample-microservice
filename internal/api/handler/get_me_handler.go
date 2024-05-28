package handler

import (
	"errors"
	"net/http"

	authservice "github.com/ew-kislov/go-sample-microservice/internal/service/auth_service"
	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
} // @name UserResponse

// GetMe godoc
//
// @Summary			  Get user info
// @Accept			  application/json
// @Produce			  application/json
// @Tags			    Auth
// @Security			BearerAuth
// @Success			  200 {object} UserResponse
// @Failure			  401 {object} ErrorResponse
// @Failure			  401 {object} ErrorResponse
// @Failure			  404 {object} ErrorResponse
// @Router			  /auth/me [get]
func NewGetMeHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
}
