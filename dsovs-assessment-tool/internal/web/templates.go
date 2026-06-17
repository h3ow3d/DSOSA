package web

import (
"embed"
"html/template"
"net/http"
)

//go:embed ../../web/templates/*.html
var templatesFS embed.FS

type Renderer struct {
templates *template.Template
}

func NewRenderer() (*Renderer, error) {
t, err := template.ParseFS(templatesFS, "../../web/templates/*.html")
if err != nil {
return nil, err
}
return &Renderer{templates: t}, nil
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data any) error {
w.Header().Set("Content-Type", "text/html; charset=utf-8")
return r.templates.ExecuteTemplate(w, name+".html", data)
}
