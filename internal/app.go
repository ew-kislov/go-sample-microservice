package internal

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
)

func StartApp(configPath string) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	container := BuildContainer(configPath)
	defer DisposeContainer(container)

	server := BuildServer(container)

	go func() {
		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("got error while listening to server: %w", err))
		}
	}()

	<-ctx.Done()
}
