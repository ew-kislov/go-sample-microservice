package statuscontroller

import (
	authservice "github.com/ew-kislov/go-sample-microservice/internal/service/auth_service"

	"github.com/gin-gonic/gin"
)

type StatusResponse struct {
	Version    string `json:"version"`
	Commit     string `json:"commit"`
	BbuildDate string `json:"buildDate"`
}

type statusController struct {
	AuthService authservice.AuthService
}

type StatusController interface {
	GetStatus(ctx *gin.Context)
}

func NewStatusController() StatusController {
	return &statusController{}
}
