package internal

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ew-kislov/go-sample-microservice/docs"
	"github.com/ew-kislov/go-sample-microservice/pkg/api"
	"github.com/ew-kislov/go-sample-microservice/pkg/cfg"
	"github.com/ew-kislov/go-sample-microservice/pkg/logging"
	"github.com/ew-kislov/go-sample-microservice/pkg/sql"
	"github.com/ew-kislov/go-sample-microservice/pkg/version"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"

	jwtmiddleware "github.com/ew-kislov/go-sample-microservice/internal/api/middleware/jwt_middleware"
	userrepository "github.com/ew-kislov/go-sample-microservice/internal/repository/user_repository"
	authservice "github.com/ew-kislov/go-sample-microservice/internal/service/auth_service"

	getmecontroller "github.com/ew-kislov/go-sample-microservice/internal/api/controller/get_me_controller"
	signupcontroller "github.com/ew-kislov/go-sample-microservice/internal/api/controller/sign_up_controller"
	statuscontroller "github.com/ew-kislov/go-sample-microservice/internal/api/controller/status_controller"
)

// @title           Sample microservice API
// @description     Sample microservice.

// @BasePath  /api/v1

// @securityDefinitions.apiKey  BearerAuth
// @in header
// @name Authorization
func StartApp(configPath string) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	config := cfg.ParseConfig(configPath)
	logger := logging.CreateLogger(&config)
	db := sql.CreateDatabase(&config, logger)

	defer db.Close()

	userRepository := userrepository.NewUserRepository(db)

	authService := authservice.NewAuthService(&config, userRepository)

	jwtMiddleware := jwtmiddleware.JwtMiddleware(authService, &config)

	signUpController := signupcontroller.NewSignUpController(authService)
	getMeController := getmecontroller.NewGetMeController()
	statusController := statuscontroller.NewStatusController()

	if config.Env == cfg.Production {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(api.RequestIdMiddleware())
	router.Use(api.LoggerMiddleware(logger))
	router.Use(api.ErrorHandlerMiddleware(&config))

	v1 := router.Group("/api/v1")
	{
		public := v1.Group("/auth")
		{
			public.POST("/sign-up", signUpController.SignUp)
			public.GET("/me", jwtMiddleware, getMeController.GetMe)
		}

		internal := v1.Group("/internal")
		{
			internal.GET("/status", statusController.GetStatus)
		}
	}

	docs.SwaggerInfo.Version = version.Version
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", config.ServerPort)

	router.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	srv := &http.Server{
		Addr:        fmt.Sprintf(":%d", config.ServerPort),
		Handler:     router.Handler(),
		ReadTimeout: time.Minute,
	}

	go func() {
		err := srv.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("got error while listening to server: %w", err))
		}
	}()

	<-ctx.Done()
}
