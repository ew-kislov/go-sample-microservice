package integration_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ew-kislov/go-sample-microservice/internal"
	"github.com/ew-kislov/go-sample-microservice/pkg"
	"github.com/jmoiron/sqlx"
)

var (
	Config pkg.Config
	Db     *sqlx.DB
)

func TestMain(m *testing.M) {
	go internal.StartApp("../../.env.test")

	Config = pkg.ParseConfig("../../.env.test")
	Db = pkg.CreateDatabase(Config)

	waitServer()

	defer Db.Close()

	m.Run()
}

func waitServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			panic("Server timed out")
		case <-ticker.C:
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/internal/status", Config.ServerPort))
			if err == nil && resp.StatusCode == http.StatusOK {
				return
			}
		}
	}
}
