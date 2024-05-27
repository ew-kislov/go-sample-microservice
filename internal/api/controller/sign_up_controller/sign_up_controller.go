package signupcontroller

import (
	"net/http"

	authservice "github.com/ew-kislov/go-sample-microservice/internal/service/auth_service"
	"github.com/ew-kislov/go-sample-microservice/pkg/api"

	"github.com/gin-gonic/gin"
)

// SignUp godoc
//
// @Summary     Sign up user
// @Tags        Auth
// @Accept      application/json
// @Produce     application/json
// @Param       account body SignUpRequest true "Sign up request"
// @Success     201 {object} SignUpResponse
// @Failure     400 {object} api.ErrorResponse
// @Failure     409 {object} api.ErrorResponse
// @Failure     422 {object} api.ErrorResponse
// @Failure     500 {object} api.ErrorResponse
// @Router      /auth/sign-up [post]
func (controller *signUpController) SignUp(ctx *gin.Context) {
	var body SignUpRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		_ = ctx.Error(api.Error{Code: http.StatusUnprocessableEntity, Message: "Could not parse body"})
		return
	}

	response, err := controller.authService.SignUp(ctx, authservice.SignUpParams(body))

	if err != nil {
		_ = ctx.Error(err)
	} else {
		ctx.JSON(http.StatusCreated, SignUpResponse(*response))
	}
}
