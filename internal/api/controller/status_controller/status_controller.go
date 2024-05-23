package statuscontroller

import (
	"net/http"

	"github.com/ew-kislov/go-sample-microservice/pkg/version"
	"github.com/gin-gonic/gin"
)

func (controller *statusController) GetStatus(ctx *gin.Context) {
	response := StatusResponse{
		Version:    version.Version,
		Commit:     version.Commit,
		BbuildDate: version.BuildDate,
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}
