package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	compile "github.com/theMitocondria/compiler/internal/http/handlers/compile"
	spincontainers "github.com/theMitocondria/compiler/internal/utils/spinContainers"
)

func main() {

	// Build images and spin the containers
	go func() {
		spincontainers.SpinContainer()
		// slog.Info("containers spined")
		fmt.Println("containres spined")
	}()

	// Setup router
	http.HandleFunc("POST /api/v1/compile", compile.CompileCode())

	// Setup server
	server := &http.Server{
		Addr:    "localhost:3000",
		Handler: nil,
	}

	// slog.Info("server started")
	fmt.Println("server started")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-done
	slog.Info("shutting down the server gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("faled to shutdown server")
		// slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
}
