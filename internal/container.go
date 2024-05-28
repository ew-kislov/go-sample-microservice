package internal

import (
	"github.com/ew-kislov/go-sample-microservice/internal/api/handler"
	"github.com/ew-kislov/go-sample-microservice/internal/api/middleware"
	userrepository "github.com/ew-kislov/go-sample-microservice/internal/repository/user_repository"
	authservice "github.com/ew-kislov/go-sample-microservice/internal/service/auth_service"
	"github.com/ew-kislov/go-sample-microservice/pkg/api"
	"github.com/ew-kislov/go-sample-microservice/pkg/cfg"
	"github.com/ew-kislov/go-sample-microservice/pkg/logging"
	"github.com/ew-kislov/go-sample-microservice/pkg/sql"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Infrastructure struct {
	Config *cfg.Config
	Logger *logrus.Logger
	Db     sql.Database
}

type Repository struct {
	UserRepository userrepository.UserRepository
}

type Service struct {
	AuthService authservice.AuthService
}

type Middleware struct {
	Jwt          gin.HandlerFunc
	Logger       gin.HandlerFunc
	RequestId    gin.HandlerFunc
	ErrorHandler gin.HandlerFunc
}

type Handler struct {
	GetStatus gin.HandlerFunc
	GetMe     gin.HandlerFunc
	SignUp    gin.HandlerFunc
}

type Api struct {
	Middleware
	Handler
}

type Container struct {
	Infrastructure
	Repository
	Service
	Api
}

func BuildContainer(configPath string) *Container {
	config := cfg.ParseConfig(configPath)
	logger := logging.CreateLogger(config)
	db := sql.CreateDatabase(config, logger)

	infrastructure := Infrastructure{
		Config: config,
		Logger: logger,
		Db:     db,
	}

	repository := Repository{
		UserRepository: userrepository.NewUserRepository(infrastructure.Db),
	}

	service := Service{
		AuthService: authservice.NewAuthService(infrastructure.Config, repository.UserRepository),
	}

	//nolint: gocritic, revive
	api := Api{
		Middleware{
			ErrorHandler: api.ErrorHandlerMiddleware(infrastructure.Config),
			Logger:       api.LoggerMiddleware(infrastructure.Logger),
			RequestId:    api.RequestIdMiddleware(),

			Jwt: middleware.JwtMiddleware(service.AuthService, infrastructure.Config),
		},
		Handler{
			GetStatus: handler.GetStatusHandler(),
			SignUp:    handler.NewSignUpHandler(service.AuthService),
			GetMe:     handler.NewGetMeHandler(),
		},
	}

	return &Container{
		Infrastructure: infrastructure,
		Repository:     repository,
		Service:        service,
		Api:            api,
	}
}

func DisposeContainer(container *Container) {
	container.Infrastructure.Db.Close()
}
