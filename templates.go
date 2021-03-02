package home

import (
	"html/template"
	"net/http"
	"time"
)

func (s *Server) siteContext() *SiteContext {
	return &SiteContext{
		Title: "Coded by Kaleb",
	}
}

func (s *Server) templateContext(path, section string, page interface{}) *TemplateContext {
	return &TemplateContext{
		Site:      s.siteContext(),
		Permalink: template.URL(path),
		Section:   section,
		Page:      page,
	}
}

func (s *Server) httpExecuteTemplate(w http.ResponseWriter, r *http.Request, tmplName, section string, page interface{}) {
	tmplCtx := s.templateContext(r.URL.Path, section, page)

	tmpl, err := template.ParseFS(s.Templates, "base.html", "partials/*.html", tmplName)
	if err != nil {
		panic(err.Error())
	}

	err = tmpl.ExecuteTemplate(w, tmplName, tmplCtx)
	if err != nil {
		panic(err.Error())
	}
}

type TemplateContext struct {
	Site      *SiteContext
	Feeds     []*FeedContext
	Permalink template.URL
	Section   string
	Page      interface{}
}

type SiteContext struct {
	Title string
	Menus []*MenuItemContext
}

type FeedContext struct {
	Title     string
	Rel       string
	MediaType string
	URL       template.URL
}

type MenuItemContext struct {
	Name    string
	Section string
	URL     template.URL
}

type TagContext struct {
	Name  string
	Count int
	URL   template.URL
}

type ArticleListContext struct {
	Posts []*ArticleMetadata
}

type ArticleMetadata struct {
	Title       string    `yaml:"title"`
	Slug        string    `yaml:"-"`
	Path        string    `yaml:"-"`
	Draft       bool      `yaml:"draft"`
	Date        time.Time `yaml:"date"`
	Lastmod     time.Time `yaml:"lastmod"`
	ShowUpdated bool      `yaml:"show_updated"`
	Tags        []string  `yaml:"tags"`
}

type ArticleContext struct {
	Meta          *ArticleMetadata
	GeminiContent string
	HtmlContent   template.HTML
}

type NotFoundContext struct {
	Path string
}
