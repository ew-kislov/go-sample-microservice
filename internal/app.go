package internal

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	getmecontroller "github.com/ew-kislov/go-sample-microservice/internal/api/controller/get_me_controller"
	signupcontroller "github.com/ew-kislov/go-sample-microservice/internal/api/controller/sign_up_controller"
	statuscontroller "github.com/ew-kislov/go-sample-microservice/internal/api/controller/status_controller"
	jwtmiddleware "github.com/ew-kislov/go-sample-microservice/internal/api/middleware/jwt_middleware"
	userrepository "github.com/ew-kislov/go-sample-microservice/internal/repository/user_repository"
	authservice "github.com/ew-kislov/go-sample-microservice/internal/service/auth_service"
	userservice "github.com/ew-kislov/go-sample-microservice/internal/service/user_service"
	"github.com/ew-kislov/go-sample-microservice/pkg"
)

func StartApp(configPath string) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	config := pkg.ParseConfig(configPath)
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
	statusController := statuscontroller.NewStatusController()

	if config.Env == pkg.Production {
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
