package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	getmecontroller "github.com/ew-kislov/go-sample-microservice/internal/api/controller/get_me_controller"
	"github.com/ew-kislov/go-sample-microservice/pkg/jwt"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestGetMe(t *testing.T) {
	endpoint := fmt.Sprintf("http://localhost:%d/api/v1/auth/me", Config.ServerPort)

	t.Run("Throws error 401 if token was not provided", func(t *testing.T) {
		resp, _ := http.Get(endpoint)
		bodyRay, _ := io.ReadAll(resp.Body)

		defer resp.Body.Close()

		var body map[string]any
		_ = json.Unmarshal(bodyRay, &body)

		assert.Equal(t, resp.StatusCode, http.StatusUnauthorized)
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

		token, _ := jwt.CreateJwt(map[string]any{"id": id}, Config.JwtSecret)

		client := &http.Client{}
		req, _ := http.NewRequest("GET", endpoint, http.NoBody)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp, _ := client.Do(req)
		bodyRaw, _ := io.ReadAll(resp.Body)

		defer resp.Body.Close()

		var body map[string]any
		_ = json.Unmarshal(bodyRaw, &body)

		var actualData getmecontroller.UserResponse
		_ = mapstructure.Decode(body, &actualData)

		expectedData := getmecontroller.UserResponse{
			Id:          int(id),
			Email:       email,
			Username:    username,
			DisplayName: username,
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expectedData, actualData)

		_, _ = Db.Exec(context.TODO(), "DELETE FROM users WHERE id = $1", id)
	})
}
