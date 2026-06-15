package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/woragis/management/backend/server/internal/httpserver"
	"github.com/woragis/management/backend/server/internal/middleware"
	"github.com/woragis/management/backend/server/internal/migrate"
	"github.com/woragis/management/backend/server/internal/platform/postgres"
)

func main() {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	adminKey := strings.TrimSpace(os.Getenv("ADMIN_API_KEY"))
	if adminKey == "" {
		log.Fatal("ADMIN_API_KEY is required")
	}

	db, err := postgres.Open(dsn)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	if skip := strings.TrimSpace(os.Getenv("SKIP_SQL_MIGRATIONS")); skip != "1" && !strings.EqualFold(skip, "true") {
		dir := strings.TrimSpace(os.Getenv("MIGRATIONS_DIR"))
		if dir == "" {
			dir = migrate.ResolveDir()
		}
		if dir != "" {
			sqlDB, err := db.DB()
			if err != nil {
				log.Fatalf("sql db: %v", err)
			}
			mctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			err = migrate.Up(mctx, sqlDB, dir)
			cancel()
			if err != nil {
				log.Fatalf("sql migrate: %v", err)
			}
			log.Printf("sql migrations applied from %s", dir)
		} else {
			log.Print("warning: SQL migrations skipped (set MIGRATIONS_DIR or run from a tree that contains migrations/)")
		}
	}

	mediaDir := strings.TrimSpace(os.Getenv("MEDIA_STORAGE_DIR"))
	if mediaDir == "" {
		mediaDir = "./data/media"
	}
	mediaBaseURL := strings.TrimSpace(os.Getenv("MEDIA_PUBLIC_BASE_URL"))
	if mediaBaseURL == "" {
		mediaBaseURL = "http://127.0.0.1:8080/v1/public/media"
	}

	app := &httpserver.App{
		DB:           db,
		AdminAPIKey:  adminKey,
		MediaDir:     mediaDir,
		MediaBaseURL: mediaBaseURL,
		SecretsKey:   loadSecretsKey(),
	}

	cfg := middleware.LoadConfigFromEnv()
	handler := httpserver.NewHandler(app, cfg)

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func loadSecretsKey() []byte {
	raw := strings.TrimSpace(os.Getenv("SECRETS_ENCRYPTION_KEY"))
	if raw == "" {
		return nil
	}
	return []byte(raw)
}
