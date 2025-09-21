package main

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/v3ronez/memi/handler"
	"github.com/v3ronez/memi/handler/middleware"
)

//go:embed public
var FS embed.FS

func main() {
	if err := initServer(); err != nil {
		log.Fatalf("Error to init server: %v", err)
	}
	router := chi.NewMux()
	router.Use(chiMiddleware.Logger)
	router.Use(middleware.WithUser)

	renderTempl := handler.HandleTempl
	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))
	router.Get("/", renderTempl(handler.HomeHandleIndex))

	port := os.Getenv("HTTP_ADDR")
	addr := fmt.Sprintf(":%s", port)
	slog.Info("Running server on port", "port", port)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("error to init server on port: %s", port)
	}
}

func initServer() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return nil
}
