package home

import (
	"html/template"
	"time"
)

type TemplateContext struct {
	Site      SiteContext
	Permalink template.URL
	Section   string
	Page      interface{}
}

type SiteContext struct {
	Title string
	Feeds []FeedContext
	Menus []MenuItemContext
	Tags  []TagContext
}

type FeedContext struct {
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

type PostListContext struct {
	Posts []PostMetadata
}

type PostMetadata struct {
	Title       string    `yaml:"title"`
	Slug        string    `yaml:"-"`
	Path        string    `yaml:"-"`
	Draft       bool      `yaml:"draft"`
	Date        time.Time `yaml:"date"`
	Lastmod     time.Time `yaml:"lastmod"`
	ShowUpdated bool      `yaml:"show_updated"`
	Tags        []string  `yaml:"tags"`
}

type PostContext struct {
	PostMetadata
	GeminiContent string
	HtmlContent   template.HTML
}
