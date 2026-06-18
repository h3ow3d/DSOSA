package web

import (
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

// Renderer renders named page templates composed with the shared layout.
type Renderer struct {
	baseDir string
}

// NewRenderer returns a Renderer rooted at the current working directory.
func NewRenderer() (*Renderer, error) {
	return &Renderer{baseDir: "."}, nil
}

// Render writes the named page template wrapped in layout.html to w.
// Each page template must define a {{define "content"}}...{{end}} block.
func (r *Renderer) Render(w http.ResponseWriter, name string, data any) error {
	layoutPath := filepath.Join(r.baseDir, "web", "templates", "layout.html")
	pagePath := filepath.Join(r.baseDir, "web", "templates", name+".html")

	t, err := template.New("layout.html").Funcs(templateFuncs).ParseFiles(layoutPath, pagePath)
	if err != nil {
		slog.Error("template parse error", "name", name, "error", err)
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.ExecuteTemplate(w, "layout.html", data); err != nil {
		slog.Error("template render error", "name", name, "error", err)
		return err
	}
	return nil
}

var templateFuncs = template.FuncMap{
	"fmtTime": func(t time.Time) string {
		if t.IsZero() {
			return "—"
		}
		return t.Format("2 Jan 2006 15:04 UTC")
	},
	"fmtDate": func(t time.Time) string {
		if t.IsZero() {
			return "—"
		}
		return t.Format("2 Jan 2006")
	},
	"pct": func(n, total int) int {
		if total == 0 {
			return 0
		}
		return n * 100 / total
	},
	"eqIntPtr": func(v *int, n int) bool {
		return v != nil && *v == n
	},
	"intPtr": func(v *int) string {
		if v == nil {
			return ""
		}
		return strconv.Itoa(*v)
	},
}
