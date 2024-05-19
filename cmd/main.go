package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	getmecontroller "github.com/ew-kislov/go-sample-microservice/internal/api/controller/get_me_controller"
	signupcontroller "github.com/ew-kislov/go-sample-microservice/internal/api/controller/sign_up_controller"
	jwtmiddleware "github.com/ew-kislov/go-sample-microservice/internal/api/middleware/jwt_middleware"
	userrepository "github.com/ew-kislov/go-sample-microservice/internal/repository/user_repository"
	authservice "github.com/ew-kislov/go-sample-microservice/internal/service/auth_service"
	userservice "github.com/ew-kislov/go-sample-microservice/internal/service/user_service"
	"github.com/ew-kislov/go-sample-microservice/pkg"
)

func main() {
	config := pkg.ParseConfig()
	logger := pkg.CreateLogger(config)
	db := pkg.CreateDatabase(config)

	defer db.Close()

	baseRepository := pkg.BaseRepository{Db: db, Logger: logger}
	encryptionProvider := pkg.EncryptionProvider{}

	userRepository := userrepository.NewUserRepository(db, baseRepository)

	userService := userservice.NewUserService(userRepository, encryptionProvider)
	authService := authservice.NewAuthService(config, encryptionProvider, userService)

	jwtMiddleware := jwtmiddleware.NewJwtMiddleware(authService)
	loggerMiddleware := pkg.LoggerMiddleware{Logger: logger}
	requestIdMiddleware := pkg.RequestIdMiddleware{}

	signUpController := signupcontroller.NewSignUpController(authService)
	getMeController := getmecontroller.NewGetMeController()

	if config.Env == pkg.Production {
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
