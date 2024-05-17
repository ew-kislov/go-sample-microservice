package main

import (
	"fmt"

	"go-sample-microservice/src/api/controller"
	"go-sample-microservice/src/api/middleware"
	"go-sample-microservice/src/infrastructure"
	"go-sample-microservice/src/repository"
	"go-sample-microservice/src/service"

	"github.com/gin-gonic/gin"
)

func main() {
	config := infrastructure.ParseConfig()
	logger := infrastructure.CreateLogger(config)
	db := infrastructure.CreateDatabase(config)

	defer db.Close()

	baseRepository := infrastructure.BaseRepository{Db: db, Logger: logger}

	userRepository := repository.UserRepository{Db: db, BaseRepository: baseRepository}

	userService := service.UserService{UserRepository: userRepository}
	authService := service.AuthService{UserService: userService, Config: config}

	jwtMiddleware := middleware.JwtMiddleware{AuthService: authService}
	loggerMiddleware := infrastructure.LoggerMiddleware{Logger: logger}
	requestIdMiddleware := infrastructure.RequestIdMiddleware{}

	signUpController := controller.SignUpController{AuthService: authService}
	getMeController := controller.GetMeController{}

	if config.Env == infrastructure.Production {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	app := gin.New()

	app.Use(gin.Recovery())
	app.Use(requestIdMiddleware.AddRequestId)
	app.Use(loggerMiddleware.Log)

	app.POST("/auth/sign-up", signUpController.SignUp)
	app.GET("/auth/me", jwtMiddleware.CheckJwt, getMeController.GetMe)

	app.Run(fmt.Sprintf(":%d", config.ServerPort))
}
