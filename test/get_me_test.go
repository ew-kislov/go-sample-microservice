package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	authservice "github.com/ew-kislov/go-sample-microservice/internal/service/auth_service"
	"github.com/ew-kislov/go-sample-microservice/pkg/jwt"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestGetMe(t *testing.T) {
	endpoint := fmt.Sprintf("http://localhost:%d/auth/me", Config.ServerPort)

	t.Run("Throws error 401 if token was not provided", func(t *testing.T) {
		resp, _ := http.Get(endpoint)
		bodyRay, _ := io.ReadAll(resp.Body)

		var body map[string]interface{}
		_ = json.Unmarshal(bodyRay, &body)

		assert.Equal(t, resp.StatusCode, http.StatusUnauthorized)
		assert.Equal(t, body["success"], false)
	})

	t.Run("Returns user", func(t *testing.T) {
		username := fmt.Sprintf("user-%s", uuid.New().String())
		email := fmt.Sprintf("%s@test.com", uuid.New().String())
		salt := "salt"
		hash := "hash"

		result, _ := Db.Query(
			context.TODO(),
			"INSERT INTO users(email, username, display_name, salt, hash) VALUES($1, $2, $3, $4, $5) RETURNING id",
			email, username, username, salt, hash,
		)
		id := result[0]["id"].(int64)

		token, _ := jwt.CreateJwt(map[string]interface{}{"id": id}, Config.JwtSecret)

		client := &http.Client{}
		req, _ := http.NewRequest("GET", endpoint, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp, _ := client.Do(req)

		bodyRay, _ := io.ReadAll(resp.Body)

		var body map[string]interface{}
		_ = json.Unmarshal(bodyRay, &body)

		var actualData authservice.User
		_ = mapstructure.Decode(body["data"], &actualData)

		expectedData := authservice.User{
			Id:          int(id),
			Email:       email,
			Username:    username,
			DisplayName: username,
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expectedData, actualData)
	})
}
