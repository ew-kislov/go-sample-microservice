package internal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ew-kislov/go-sample-microservice/docs"
	"github.com/ew-kislov/go-sample-microservice/pkg/cfg"
	"github.com/ew-kislov/go-sample-microservice/pkg/version"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

// @title           Sample microservice API
// @description     Sample microservice.

// @BasePath  /api/v1

// @securityDefinitions.apiKey  BearerAuth
// @in header
// @name Authorization
func BuildServer(container *Container) *http.Server {
	if container.Config.Env == cfg.Production {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	setMiddlewares(router, container)
	setEndpoints(router, container)
	setSwagger(router, container)

	return &http.Server{
		Addr:        fmt.Sprintf(":%d", container.Infrastructure.Config.ServerPort),
		Handler:     router.Handler(),
		ReadTimeout: time.Minute,
	}
}

func setMiddlewares(router *gin.Engine, container *Container) {
	router.Use(gin.Recovery())

	router.Use(container.Api.Middleware.RequestId)
	router.Use(container.Api.Middleware.Logger)
	router.Use(container.Api.ErrorHandler)
}

func setEndpoints(router *gin.Engine, container *Container) {
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/sign-up", container.Api.Handler.SignUp)
			auth.GET("/me", container.Api.Middleware.Jwt, container.Api.Handler.GetMe)
		}

		internal := v1.Group("/internal")
		{
			internal.GET("/status", container.Api.Handler.GetStatus)
		}
	}
}

func setSwagger(router *gin.Engine, container *Container) {
	docs.SwaggerInfo.Version = version.Version
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", container.Infrastructure.Config.ServerPort)

	router.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
}
