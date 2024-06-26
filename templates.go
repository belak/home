package home

import (
	"net/http"

	"github.com/belak/home/internal"
)

func (s *Server) siteContext() *SiteContext {
	return &SiteContext{
		Title: "Belak's Tools",
	}
}

func (s *Server) templateContext(page interface{}) *TemplateContext {
	return &TemplateContext{
		Site: s.siteContext(),
		Page: page,
	}
}

func (s *Server) httpExecuteTemplate(w http.ResponseWriter, r *http.Request, tmplName string, page interface{}) {
	ctx := r.Context()
	tmplCtx := s.templateContext(page)

	internal.RenderTemplate(ctx, w, tmplName, tmplCtx)
}

type TemplateContext struct {
	Site *SiteContext
	Page interface{}
}

type SiteContext struct {
	Title string
}

type NotFoundContext struct {
	Path string
}
