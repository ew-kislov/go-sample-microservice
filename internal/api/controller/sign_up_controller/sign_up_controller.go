package signupcontroller

import (
	"net/http"

	authservice "github.com/ew-kislov/go-sample-microservice/internal/service/auth_service"
	"github.com/ew-kislov/go-sample-microservice/pkg/api"

	"github.com/gin-gonic/gin"
)

func (controller *signUpController) SignUp(ctx *gin.Context) {
	var body SignUpRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Error(api.ApiError{Code: http.StatusUnprocessableEntity, Message: "Could not parse body"})
		return
	}

	response, err := controller.authService.SignUp(ctx, authservice.SignUpParams(body))

	if err != nil {
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusCreated, SignUpResponse(*response))
	}
}
