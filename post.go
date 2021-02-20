package home

import (
	"io/fs"
	"strings"
	"time"
)

type PostMetadata struct {
	Path string `yaml:"-"`
	Slug string `yaml:"-"`

	Title       string    `yaml:"title"`
	Tags        []string  `yaml:"tags"`
	ShowUpdated bool      `yaml:"show_updated"`
	Draft       bool      `yaml:"draft"`
	Date        time.Time `yaml:"date"`
	Update      time.Time `yaml:"updated"`
}

func ReadPostMetadata(targetFS fs.FS, targetPath string) (*PostMetadata, error) {
	targetPath = strings.TrimSuffix(targetPath, ".gmi")

	return &PostMetadata{}, nil
}
