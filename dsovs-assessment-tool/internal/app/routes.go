package app

import (
"context"
"log/slog"
"net/http"
"sort"
"strings"
"time"

"dsovs-assessment-tool/internal/dsovs"
)

func (s *Server) registerRoutes(mux *http.ServeMux) {
mux.Handle("/static/", http.StripPrefix("/static/", webStaticHandler()))
mux.HandleFunc("GET /", s.handleDashboard)
mux.HandleFunc("GET /dashboard", s.handleDashboard)
mux.HandleFunc("POST /catalogue/sync", s.handleCatalogueSync)
mux.HandleFunc("GET /projects", s.handleProjects)
mux.HandleFunc("GET /project", s.handleProject)
mux.HandleFunc("GET /assessment", s.handleAssessment)
mux.HandleFunc("GET /results", s.handleResults)
mux.HandleFunc("GET /report", s.handleReport)
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
catalogues, _ := s.store.ListCatalogueVersions()
sort.Slice(catalogues, func(i, j int) bool { return catalogues[i].FetchedAt.After(catalogues[j].FetchedAt) })

data := map[string]any{
"Title":             "Dashboard",
"CatalogueCount":    len(catalogues),
"HasCatalogue":      len(catalogues) > 0,
"LastCatalogue":     firstOrNil(catalogues),
"ProjectCount":      len(s.store.ListProjects()),
"AssessmentCount":   len(s.store.ListAssessments()),
"ImprovementCount":  len(s.store.ListImprovements()),
"ControlCount":      s.store.CurrentControlCount(),
"ActiveNav":         "dashboard",
"SyncSuccessMessage": r.URL.Query().Get("synced"),
}
_ = s.renderer.Render(w, "dashboard", data)
}

func (s *Server) handleCatalogueSync(w http.ResponseWriter, r *http.Request) {
ctx, cancel := context.WithTimeout(r.Context(), s.cfg.SyncTimout)
defer cancel()

result, err := dsovs.Sync(ctx, s.client, s.store)
if err != nil {
http.Error(w, "sync failed", http.StatusBadGateway)
slog.Error("catalogue sync failed", "error", err)
return
}

slog.Info("catalogue synced", "version", result.Version, "changed", result.Changed)
http.Redirect(w, r, "/dashboard?synced="+urlQueryEscape(syncMessage(result.Changed)), http.StatusSeeOther)
}

func (s *Server) handleProjects(w http.ResponseWriter, r *http.Request) {
projects := s.store.ListProjects()
sort.Slice(projects, func(i, j int) bool { return projects[i].UpdatedAt.After(projects[j].UpdatedAt) })
_ = s.renderer.Render(w, "projects", map[string]any{"Title": "Projects", "Projects": projects, "ActiveNav": "projects"})
}

func (s *Server) handleProject(w http.ResponseWriter, r *http.Request) {
_ = s.renderer.Render(w, "project", map[string]any{"Title": "Project", "ActiveNav": "projects"})
}

func (s *Server) handleAssessment(w http.ResponseWriter, r *http.Request) {
_ = s.renderer.Render(w, "assessment", map[string]any{"Title": "Assessment", "ActiveNav": "assessment"})
}

func (s *Server) handleResults(w http.ResponseWriter, r *http.Request) {
_ = s.renderer.Render(w, "results", map[string]any{"Title": "Results", "ActiveNav": "results"})
}

func (s *Server) handleReport(w http.ResponseWriter, r *http.Request) {
_ = s.renderer.Render(w, "report", map[string]any{"Title": "Report", "ActiveNav": "report"})
}

func firstOrNil[T any](items []T) any {
if len(items) == 0 {
return nil
}
return items[0]
}

func syncMessage(changed bool) string {
if changed {
return "Catalogue synced"
}
return "Catalogue already up to date"
}

func urlQueryEscape(v string) string {
return strings.ReplaceAll(v, " ", "+")
}

func webStaticHandler() http.Handler {
return http.FileServer(http.FS(embeddedStaticFS))
}

var _ = time.Second
