package api

import (
	"net/http"

	"github.com/ew-kislov/go-sample-microservice/pkg/cfg"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
} // @name ErrorResponse

func CreateErrorResponse(err error, config *cfg.Config) (int, ErrorResponse) {
	apiError, ok := err.(Error)

	if ok {
		return apiError.Code, ErrorResponse{Error: apiError.Message}
	}

	if config.Env == cfg.Production {
		return http.StatusInternalServerError, ErrorResponse{Error: "Internal error"}
	} else {
		return http.StatusInternalServerError, ErrorResponse{Error: "Internal error", Details: err.Error()}
	}
}
