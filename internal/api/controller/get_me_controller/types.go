package getmecontroller

import "github.com/gin-gonic/gin"

type UserResponse struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type getMeController struct {
}

type GetMeController interface {
	GetMe(ctx *gin.Context)
}

func NewGetMeController() GetMeController {
	return &getMeController{}
}
