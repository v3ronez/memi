package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/v3ronez/memi/handler"
	"github.com/v3ronez/memi/handler/middleware"
	"github.com/v3ronez/memi/internal/data"
)

//go:embed public
var FS embed.FS

const version = "1.0"

type dbConfig struct {
	host         string
	dbname       string
	port         int64
	user         string
	password     string
	sslmode      string
	maxOpenConns int
	maxIdleConns int
	maxIndleTime string
}

type config struct {
	servPort int64
	envMode  string
	db       dbConfig
	smtp     struct {
		host     string
		port     int
		username string
		password string
	}
}
type application struct {
	config  config
	models  data.Models
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	wg sync.WaitGroup
}

func main() {
	if err := initEnv(); err != nil {
		log.Fatalf("Error to init server: %v", err)
	}
	cfg, err := initConfig()
	if err != nil {
		slog.Error("Error to init configurtions", "err", err)
	}
	app, err := newApp(*cfg)
	if err != nil {
		slog.Error("Error to init application", "err", err)
	}
	dbConn, err := initDatabase(*cfg)
	if err != nil {
		slog.Error("Error to init database connection", "err", err)
	}

	app.models = data.NewModel(dbConn)
	router := chi.NewMux()
	router.Use(chiMiddleware.Logger)
	router.Use(middleware.WithUser)

	renderTempl := handler.HandleTempl
	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))
	router.Get("/", renderTempl(handler.HandleHomeIndex))
	router.Get("/login", renderTempl(handler.HandleLoginIndex))

	slog.Info("Running server on port", "port", port)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("error to init server on port: %s", port)
	}
}

func initConfig() (*config, error) {
	port, err := strconv.ParseInt(os.Getenv("DB_PORT"), 10, 32)
	if err != nil {
		return nil, err
	}
	servPort, err := strconv.ParseInt(os.Getenv("HTTP_PORT"), 10, 32)
	if err != nil {
		return nil, err
	}
	return &config{
		servPort: servPort,
		envMode:  os.Getenv("ENV_MODE"),
		db: dbConfig{
			host:         os.Getenv("DB_HOST"),
			port:         port,
			dbname:       os.Getenv("DB_DATABASE"),
			user:         os.Getenv("DB_USER"),
			password:     os.Getenv("DB_PASSWORD"),
			sslmode:      os.Getenv("DB_SSL_MODE"),
			maxOpenConns: 10,
			maxIdleConns: 10,
			maxIndleTime: "5s",
		},
	}, nil
}

func newApp(cfg config) (*application, error) {
	return &application{
		config: cfg,
	}, nil
}

func initDatabase(cfg config) (*sql.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		cfg.db.user,
		cfg.db.password,
		cfg.db.host,
		cfg.db.port,
		cfg.db.dbname,
		cfg.db.sslmode)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	conn.SetConnMaxIdleTime(time.Duration(cfg.db.maxIdleConns))
	conn.SetMaxOpenConns(cfg.db.maxOpenConns)
	duration, err := time.ParseDuration(cfg.db.maxIndleTime)

	if err != nil {
		return nil, err
	}
	conn.SetConnMaxLifetime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := conn.PingContext(ctx); err != nil {
		return nil, err
	}
	fmt.Println("db connection success!")
	return conn, nil
}

func initEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return nil
}
