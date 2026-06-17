package app

import (
"context"
"errors"
"fmt"
"log/slog"
"net/http"
"os"
"os/signal"
"path/filepath"
"syscall"
"time"

"dsovs-assessment-tool/internal/dsovs"
"dsovs-assessment-tool/internal/storage"
"dsovs-assessment-tool/internal/web"
)

type Config struct {
Address    string
DataDir    string
DSOVSURL   string
AutoSync   bool
SyncTimout time.Duration
}

type Server struct {
cfg      Config
http     *http.Server
store    *storage.JSONStore
client   *dsovs.Client
renderer *web.Renderer
}

func NewServer(cfg Config) (*Server, error) {
if err := os.MkdirAll(cfg.DataDir, 0o755); err != nil {
return nil, fmt.Errorf("create data dir: %w", err)
}

store, err := storage.NewJSONStore(cfg.DataDir)
if err != nil {
return nil, err
}

renderer, err := web.NewRenderer()
if err != nil {
return nil, err
}

s := &Server{
cfg:      cfg,
store:    store,
renderer: renderer,
client:   dsovs.NewClient(cfg.DSOVSURL),
}

mux := http.NewServeMux()
s.registerRoutes(mux)
s.http = &http.Server{Addr: cfg.Address, Handler: mux, ReadHeaderTimeout: 10 * time.Second}

return s, nil
}

func (s *Server) Run(ctx context.Context) error {
if s.cfg.AutoSync {
syncCtx, cancel := context.WithTimeout(ctx, s.cfg.SyncTimout)
defer cancel()
if _, err := dsovs.Sync(syncCtx, s.client, s.store); err != nil {
slog.Warn("auto sync failed", "error", err)
}
}

errCh := make(chan error, 1)
go func() {
slog.Info("listening", "address", s.cfg.Address, "data_dir", filepath.Clean(s.cfg.DataDir))
errCh <- s.http.ListenAndServe()
}()

sigCtx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
defer cancel()

select {
case <-sigCtx.Done():
shutdownCtx, stop := context.WithTimeout(context.Background(), 10*time.Second)
defer stop()
return s.http.Shutdown(shutdownCtx)
case err := <-errCh:
if errors.Is(err, http.ErrServerClosed) {
return nil
}
return err
}
}
