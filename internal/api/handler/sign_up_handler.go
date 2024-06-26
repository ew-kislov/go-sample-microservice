package handler

import (
	"net/http"

	authservice "github.com/ew-kislov/go-sample-microservice/internal/service/auth_service"
	"github.com/ew-kislov/go-sample-microservice/pkg/api"

	"github.com/gin-gonic/gin"
)

type SignUpRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Username    string `json:"username" binding:"required"`
	DisplayName string `json:"displayName" binding:"required"`
	Password    string `json:"password" binding:"required"`
} // @name SignUpRequest

type SignUpResponse struct {
	UserId int64  `json:"userId"`
	Token  string `json:"token"`
} // @name SignUpResponse

// SignUp godoc
//
// @Summary     Sign up user
// @Tags        Auth
// @Accept      application/json
// @Produce     application/json
// @Param       account body SignUpRequest true "Sign up request"
// @Success     201 {object} SignUpResponse
// @Failure     400 {object} ErrorResponse
// @Failure     409 {object} ErrorResponse
// @Failure     422 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /auth/sign-up [post]
func NewSignUpHandler(service authservice.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body SignUpRequest

		if err := ctx.ShouldBindJSON(&body); err != nil {
			_ = ctx.Error(api.Error{Code: http.StatusUnprocessableEntity, Message: "Could not parse body"})
			return
		}

		response, err := service.SignUp(ctx, authservice.SignUpParams(body))

		if err != nil {
			_ = ctx.Error(err)
		} else {
			ctx.JSON(http.StatusCreated, SignUpResponse(*response))
		}
	}
}
