package handlers

import (
	"net/http"

	"github.com/wootaiklee/git-resume/internal/api/dto"
	"github.com/wootaiklee/git-resume/internal/llm"
)

// TemplatesHandler handles template-related endpoints
type TemplatesHandler struct{}

// NewTemplatesHandler creates a new templates handler
func NewTemplatesHandler() *TemplatesHandler {
	return &TemplatesHandler{}
}

// List handles GET /api/templates
func (h *TemplatesHandler) List(w http.ResponseWriter, r *http.Request) {
	mgr := llm.NewTemplateManager()
	builtinTemplates := mgr.ListTemplates()

	templates := make([]dto.TemplateResponse, 0, len(builtinTemplates))
	for _, name := range builtinTemplates {
		mgr.SetTemplate(name)
		tmpl := mgr.GetTemplate()
		templates = append(templates, dto.TemplateResponse{
			Name:        tmpl.Name,
			Description: tmpl.Description,
			ToneStyle:   tmpl.ToneStyle,
			Focus:       tmpl.Focus,
		})
	}

	respondOK(w, dto.TemplatesListResponse{
		Templates: templates,
	})
}
