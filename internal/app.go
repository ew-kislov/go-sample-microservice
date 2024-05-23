package internal

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/ew-kislov/go-sample-microservice/pkg/api"
	"github.com/ew-kislov/go-sample-microservice/pkg/cfg"
	"github.com/ew-kislov/go-sample-microservice/pkg/db"
	"github.com/ew-kislov/go-sample-microservice/pkg/logger"
	"github.com/gin-gonic/gin"

	jwtmiddleware "github.com/ew-kislov/go-sample-microservice/internal/api/middleware/jwt_middleware"
	userrepository "github.com/ew-kislov/go-sample-microservice/internal/repository/user_repository"
	authservice "github.com/ew-kislov/go-sample-microservice/internal/service/auth_service"

	getmecontroller "github.com/ew-kislov/go-sample-microservice/internal/api/controller/get_me_controller"
	signupcontroller "github.com/ew-kislov/go-sample-microservice/internal/api/controller/sign_up_controller"
	statuscontroller "github.com/ew-kislov/go-sample-microservice/internal/api/controller/status_controller"
)

func StartApp(configPath string) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	config := cfg.ParseConfig(configPath)
	logger := logger.CreateLogger(config)
	db := db.CreateDatabase(&config, logger)

	defer db.Close()

	userRepository := userrepository.NewUserRepository(db)

	authService := authservice.NewAuthService(config, userRepository)

	jwtMiddleware := jwtmiddleware.NewJwtMiddleware(authService)
	loggerMiddleware := api.LoggerMiddleware{Logger: logger}
	requestIdMiddleware := api.RequestIdMiddleware{}

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
	router.Use(requestIdMiddleware.AddRequestId)
	router.Use(loggerMiddleware.Log)

	public := router.Group("/")
	{
		public.POST("/auth/sign-up", signUpController.SignUp)
		public.GET("/auth/me", jwtMiddleware.CheckJwt, getMeController.GetMe)
	}

	internal := router.Group("/internal")
	{
		internal.GET("/status", statusController.GetStatus)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.ServerPort),
		Handler: router.Handler(),
	}

	go func() {
		err := srv.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("got error while listening to server: %w", err))
		}
	}()

	<-ctx.Done()
}
