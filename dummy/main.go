package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Request structure for the POST request
type GreetingRequest struct {
	Username string `json:"username"`
}

// Handler for the "Hello, World!" route
func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

// Handler for the personalized greeting route
func greetUser(w http.ResponseWriter, r *http.Request) {
	var req GreetingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	response := fmt.Sprintf("Hello, %s!", req.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": response})
}

func main() {

	// Setup router
	http.HandleFunc("GET /hello", helloWorld)
	http.HandleFunc("POST /hello", greetUser)

	// Setup server
	server := &http.Server{
		Addr:    "localhost:3000",
		Handler: nil, // Use the default router
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
		fmt.Println("failed to shutdown server")
		// slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
}
