package integration_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ew-kislov/go-sample-microservice/internal"
	"github.com/ew-kislov/go-sample-microservice/pkg/cfg"
	"github.com/ew-kislov/go-sample-microservice/pkg/logging"
	"github.com/ew-kislov/go-sample-microservice/pkg/sql"
)

var (
	Config *cfg.Config
	Db     sql.Database
)

func TestMain(m *testing.M) {
	go internal.StartApp("../.env.test")

	Config = cfg.ParseConfig("../.env.test")
	Db = sql.CreateDatabase(Config, logging.CreateLogger(Config))

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
			if serverReady() {
				return
			}
		}
	}
}

func serverReady() bool {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/v1/internal/status", Config.ServerPort))

	if err == nil {
		defer resp.Body.Close()
	}

	if err == nil && resp.StatusCode == http.StatusOK {
		return true
	}

	return false
}
