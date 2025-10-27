package telegram

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

type TemplateEngine struct {
	templates *template.Template
}

func NewTemplateEngine() (*TemplateEngine, error) {
	templates, err := template.ParseFS(templateFS, "templates/*.tmpl")
	if err != nil {
		return nil, err
	}

	return &TemplateEngine{
		templates: templates,
	}, nil
}

func (e *TemplateEngine) Render(templateName string, data interface{}) (string, error) {
	var buf bytes.Buffer

	err := e.templates.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		return "", fmt.Errorf("failed to render template %s: %w", templateName, err)
	}

	return buf.String(), nil
}
