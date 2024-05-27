package jwtmiddleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	authservice "github.com/ew-kislov/go-sample-microservice/internal/service/auth_service"
	authservicemocks "github.com/ew-kislov/go-sample-microservice/internal/service/auth_service/mocks"
	"github.com/ew-kislov/go-sample-microservice/pkg/api"
	"github.com/ew-kislov/go-sample-microservice/pkg/cfg"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestJwtMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Fails if no Authorization header has provided", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest(http.MethodGet, "/", http.NoBody)

		mockAuthService := authservicemocks.NewMockAuthService(ctrl)
		mockConfig := cfg.Config{Env: cfg.Production}

		middleware := JwtMiddleware(mockAuthService, &mockConfig)
		middleware(ctx)

		var response map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, w.Code, http.StatusUnauthorized)
		assert.Equal(t, response["error"], AuthorizationHeaderNotProvided)
	})

	t.Run("Fails if Authorization header does not have format `Bearer <token>`", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest(http.MethodGet, "/", http.NoBody)
		ctx.Request.Header.Set("Authorization", "oops")

		mockAuthService := authservicemocks.NewMockAuthService(ctrl)
		mockConfig := cfg.Config{Env: cfg.Production}

		middleware := JwtMiddleware(mockAuthService, &mockConfig)
		middleware(ctx)

		var response map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, w.Code, http.StatusUnauthorized)
		assert.Equal(t, response["error"], WrongTokenFormat)
	})

	t.Run("Fails if AuthService call fails", func(t *testing.T) {
		token := "<Token>"
		err := "Some error"

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest(http.MethodGet, "/", http.NoBody)
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		mockAuthService := authservicemocks.NewMockAuthService(ctrl)
		mockAuthService.EXPECT().Authenticate(ctx, token).Return(
			nil, api.Error{Code: http.StatusUnauthorized, Message: err},
		)

		mockConfig := cfg.Config{Env: cfg.Production}

		middleware := JwtMiddleware(mockAuthService, &mockConfig)
		middleware(ctx)

		var response map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, w.Code, http.StatusUnauthorized)
		assert.Equal(t, response["error"], err)
	})

	t.Run("Successfully puts user into context", func(t *testing.T) {
		token := "<Token>"
		user := &authservice.User{Id: 1, Username: "username", DisplayName: "Display Name", Email: "email@domain.com"}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest(http.MethodGet, "/", http.NoBody)
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		mockAuthService := authservicemocks.NewMockAuthService(ctrl)
		mockAuthService.EXPECT().Authenticate(ctx, token).Return(user, nil)

		mockConfig := cfg.Config{Env: cfg.Production}

		middleware := JwtMiddleware(mockAuthService, &mockConfig)
		middleware(ctx)

		contextUser, _ := ctx.Get("user")

		assert.NotEqual(t, w.Code, http.StatusUnauthorized)
		assert.Equal(t, contextUser, user)
	})
}
