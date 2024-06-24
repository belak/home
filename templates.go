package home

import (
	"html/template"
	"net/http"
)

func (s *Server) siteContext() *SiteContext {
	return &SiteContext{
		Title: "Coded by Kaleb",
	}
}

func (s *Server) templateContext(page interface{}) *TemplateContext {
	return &TemplateContext{
		Site: s.siteContext(),
		Page: page,
	}
}

func (s *Server) httpExecuteTemplate(w http.ResponseWriter, r *http.Request, tmplName, section string, page interface{}) {
	tmplCtx := s.templateContext(page)

	tmpl, err := template.ParseFS(templatesFS, "base.html", "partials/*.html", tmplName)
	if err != nil {
		panic(err.Error())
	}

	err = tmpl.ExecuteTemplate(w, tmplName, tmplCtx)
	if err != nil {
		panic(err.Error())
	}
}

type TemplateContext struct {
	Site *SiteContext
	Page interface{}
}

type SiteContext struct {
	Title string
	Menus []*MenuItemContext
}

type MenuItemContext struct {
	Name    string
	Section string
	URL     template.URL
}

type NotFoundContext struct {
	Path string
}
