package templatemanager

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type TemplateManager struct {
    templates *template.Template
}

func NewTemplateManager(templatesDir string) *TemplateManager {
    pattern := filepath.Join(templatesDir, "*.html")
    templates, err := template.ParseGlob(pattern)
    if err != nil {
        panic(err)
    }
    return &TemplateManager{templates: templates}
}

func (tm *TemplateManager) ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) error{
    return tm.templates.ExecuteTemplate(w, name, data)
}
