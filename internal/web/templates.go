package web

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type Renderer struct {
	baseDir string
}

func NewRenderer() (*Renderer, error) {
	return &Renderer{baseDir: "."}, nil
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data any) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	path := filepath.Join(r.baseDir, "web", "templates", name+".html")
	t, err := template.ParseFiles(path)
	if err != nil {
		return err
	}
	return t.Execute(w, data)
}
