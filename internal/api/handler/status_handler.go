package handler

import (
	"net/http"

	"github.com/ew-kislov/go-sample-microservice/pkg/version"
	"github.com/gin-gonic/gin"
)

type StatusResponse struct {
	Version    string `json:"version"`
	Commit     string `json:"commit"`
	BbuildDate string `json:"buildDate"`
}

// GetStatus godoc
//
// @Summary			  Get status
// @Accept			  application/json
// @Produce			  application/json
// @Tags			    Internal
// @Success			  200 {object} StatusResponse
// @Failure       500 {object} api.ErrorResponse
// @Router			  /internal/status [get]
func GetStatusHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response := StatusResponse{
			Version:    version.Version,
			Commit:     version.Commit,
			BbuildDate: version.BuildDate,
		}

		ctx.JSON(http.StatusOK, response)
	}
}
