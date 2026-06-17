package main

import (
"context"
"log/slog"
"os"
"os/signal"
"syscall"
"time"

"dsovs-assessment-tool/internal/app"
)

func main() {
logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
slog.SetDefault(logger)

cfg := app.Config{
Address:    envOrDefault("APP_ADDR", ":8080"),
DataDir:    envOrDefault("DATA_DIR", "/data"),
DSOVSURL:   envOrDefault("DSOVS_URL", "https://owasp.org/www-project-devsecops-verification-standard/dist/dsovs.json"),
AutoSync:   envOrDefault("AUTO_SYNC_CATALOGUE", "false") == "true",
SyncTimout: 30 * time.Second,
}

server, err := app.NewServer(cfg)
if err != nil {
logger.Error("failed to initialize server", "error", err)
os.Exit(1)
}

ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
defer cancel()

if err := server.Run(ctx); err != nil {
logger.Error("server terminated", "error", err)
os.Exit(1)
}
}

func envOrDefault(key, fallback string) string {
if value := os.Getenv(key); value != "" {
return value
}
return fallback
}
