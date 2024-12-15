package main

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	spincontainers "github.com/theMitocondria/compiler/internal/utils/spinContainers"
)

func main() {

	//build images and spin the continers
	go func() {
		spincontainers.SpinContainer()
		slog.Info("containers spined")
	}()
	//setup router
	// router := compile.InitRouter()

	// //setup server

	// server := http.Server{
	// 	Addr:    "localhost:3000",
	// 	Handler: router,
	// }

	// slog.Info("server stared")

	// done := make(chan os.Signal, 1)

	// signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// go func() {
	// 	err := server.ListenAndServe()

	// 	if err != nil {
	// 		log.Fatal("failed to start server")
	// 	}
	// }()

	// <-done

	// slog.Info("shutting down the server gracefully")

	// ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	// defer cancel()

	// if err := server.Shutdown(ctx); err != nil {
	// 	slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	// }

	// slog.Info("server shutdown successfully")

	app := fiber.New()

	app.Use(logger.New())

	group := app.Group("/api")

	app.Listen(":3000")
}
