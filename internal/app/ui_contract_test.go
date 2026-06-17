package app

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"dsovs-assessment-tool/internal/storage"
)

func TestRouteSmokePages(t *testing.T) {
	s := newTestServer(t)

	for _, path := range []string{"/", "/catalogue", "/projects", "/assessments", "/reports"} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		s.http.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET %s: got %d want %d", path, rec.Code, http.StatusOK)
		}
	}
}

func TestSharedLayoutContract(t *testing.T) {
	s := newTestServer(t)
	html := getPageBody(t, s, "/projects")

	for _, want := range []string{
		`class="wcag-skip-link"`,
		"<header",
		"<nav",
		"<main",
		"<footer",
		`id="main-content"`,
		`href="/static/wcag-lite.css"`,
		"Secure Delivery Compass",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("shared layout missing %q", want)
		}
	}
}

func TestAccessibilityPatterns(t *testing.T) {
	s := newTestServer(t)

	projectNew := getPageBody(t, s, "/projects/new")
	for _, want := range []string{"<label", "hint-text"} {
		if !strings.Contains(projectNew, want) {
			t.Fatalf("/projects/new missing %q", want)
		}
	}

	projectID := seedProject(t, s)
	assessmentID := seedAssessment(t, s, projectID)
	assessmentPage := getPageBody(t, s, "/assessments/"+assessmentID)
	for _, want := range []string{"<fieldset", "<legend", "<label"} {
		if !strings.Contains(assessmentPage, want) {
			t.Fatalf("/assessments/{id} missing %q", want)
		}
	}

	errorPage := getPageBody(t, s, "/projects/new?error=Project+name+is+required")
	if !strings.Contains(errorPage, `tabindex="-1"`) {
		t.Fatal("error summary markup is missing tabindex=\"-1\"")
	}
}

func TestReportPageContract(t *testing.T) {
	s := newTestServer(t)
	projectID := seedProject(t, s)
	assessmentID := seedAssessment(t, s, projectID)

	report := getPageBody(t, s, "/assessments/"+assessmentID+"/report")
	want := "This report may contain sensitive information about project security gaps. Review before sharing."
	if !strings.Contains(report, want) {
		t.Fatalf("report page missing sensitive info warning")
	}

	css := getPageBody(t, s, "/static/wcag-lite.css")
	if !strings.Contains(css, "@media print") {
		t.Fatal("wcag-lite.css missing @media print styles")
	}
}

func TestStaticAssetsContract(t *testing.T) {
	s := newTestServer(t)

	for _, path := range []string{"/static/wcag-lite.css", "/static/wcag-lite.js"} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		s.http.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET %s: got %d want %d", path, rec.Code, http.StatusOK)
		}
	}

	css := getPageBody(t, s, "/static/wcag-lite.css")
	for _, want := range []string{":focus-visible", "@media print"} {
		if !strings.Contains(css, want) {
			t.Fatalf("wcag-lite.css missing %q", want)
		}
	}
}

func TestNoExternalDependencyGuard(t *testing.T) {
	repoRoot := chdirRepoRoot(t)

	forbiddenFiles := map[string]struct{}{
		"package.json": {},
	}
	forbiddenTerms := []string{
		"cdn.",
		"unpkg",
		"jsdelivr",
		"bootstrap",
		"tailwind",
		"react",
		"govuk-frontend",
		"fonts.googleapis.com",
		"fonts.gstatic.com",
	}

	err := filepath.WalkDir(repoRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if strings.Contains(path, string(filepath.Separator)+".git"+string(filepath.Separator)) {
			return nil
		}
		if d.IsDir() {
			if d.Name() == "node_modules" {
				return fs.ErrPermission
			}
			return nil
		}

		rel, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return err
		}
		if _, bad := forbiddenFiles[d.Name()]; bad {
			t.Fatalf("forbidden frontend dependency file found: %s", rel)
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".html" && ext != ".css" && ext != ".js" && ext != ".md" {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content := strings.ToLower(string(b))
		for _, term := range forbiddenTerms {
			if strings.Contains(content, term) {
				t.Fatalf("forbidden frontend dependency marker %q found in %s", term, rel)
			}
		}
		return nil
	})
	if err != nil && err != fs.ErrPermission {
		t.Fatal(err)
	}
}

func newTestServer(t *testing.T) *Server {
	t.Helper()
	chdirRepoRoot(t)

	s, err := NewServer(Config{
		Address:  "127.0.0.1:0",
		DataDir:  t.TempDir(),
		DSOVSURL: "https://example.invalid/dsovs.json",
	})
	if err != nil {
		t.Fatalf("NewServer: %v", err)
	}
	return s
}

func chdirRepoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("Chdir(%s): %v", root, err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })
	return root
}

func getPageBody(t *testing.T, s *Server, path string) string {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	s.http.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET %s: got %d want %d", path, rec.Code, http.StatusOK)
	}
	return rec.Body.String()
}

func seedProject(t *testing.T, s *Server) string {
	t.Helper()
	now := time.Now().UTC()
	p := storage.Project{
		ID:        "project-1",
		Name:      "Project One",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.SaveProject(p); err != nil {
		t.Fatalf("SaveProject: %v", err)
	}
	return p.ID
}

func seedAssessment(t *testing.T, s *Server, projectID string) string {
	t.Helper()
	now := time.Now().UTC()
	catalogue := storage.CatalogueRecord{
		Version:   "test-catalogue",
		SHA256:    "abc123",
		FetchedAt: now,
		Body: map[string]any{
			"phases": []any{
				map[string]any{
					"id":   "phase-1",
					"name": "Plan",
					"controls": []any{
						map[string]any{
							"id":      "C-1",
							"title":   "Control One",
							"summary": "Control summary",
							"level_0": "not started",
							"level_1": "basic",
							"level_2": "managed",
							"level_3": "optimized",
						},
					},
				},
			},
		},
	}
	if err := s.store.SaveCatalogue(catalogue); err != nil {
		t.Fatalf("SaveCatalogue: %v", err)
	}

	a := storage.Assessment{
		ID:             "assessment-1",
		ProjectID:      projectID,
		Name:           "Assessment One",
		AssessmentDate: "2026-01-01",
		CreatedAt:      now,
		UpdatedAt:      now,
		Scores: []storage.ScoreEntry{
			{
				ControlID:   "C-1",
				Current:     1,
				Target:      2,
				Evidence:    "Evidence",
				ActionNotes: "Actions",
			},
		},
	}
	if err := s.store.SaveAssessment(a); err != nil {
		t.Fatalf("SaveAssessment: %v", err)
	}
	return a.ID
}
