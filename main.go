package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	compile "github.com/theMitocondria/compiler/internal/http/handlers/compile"
)

func main() {

	// Setup router
	http.HandleFunc("POST /api/v1/compile", compile.CompileCode())

	// Setup server
	server := &http.Server{
		Addr:    ":3000",
		Handler: nil,
	}

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
	fmt.Print("shutting down the server gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("faled to shutdown server")
	}

	fmt.Print("shutting down the server gracefull finally")
}
