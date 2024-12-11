package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/theMitocondria/compiler/internal/config"
)

func main() {

	//load config
	cfg := config.MustLoad()

	//build images and spin the continers

	//setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /api/v1/compile", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("compiled successfully."))
	})

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("failed to start server")
	}

	slog.Info("Server Started")

	//setup server

}
