package statuscontroller

import (
	"net/http"

	"github.com/ew-kislov/go-sample-microservice/pkg"
	"github.com/gin-gonic/gin"
)

func (controller *statusController) GetStatus(ctx *gin.Context) {
	response := StatusResponse{
		Version:    pkg.Version,
		Commit:     pkg.Commit,
		BbuildDate: pkg.BuildDate,
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}
