package signupcontroller

import (
	"net/http"

	userservice "github.com/ew-kislov/go-sample-microservice/internal/service/user_service"

	"github.com/gin-gonic/gin"
)

func (controller *signUpController) SignUp(ctx *gin.Context) {
	var body SignUpRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response, err := controller.AuthService.SignUp(ctx, userservice.CreateUserParams(body))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusCreated, gin.H{"success": true, "data": response})
	}
}
