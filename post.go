package home

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"path"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type PostMetadata struct {
	Slug string `yaml:"-"`
	Path string `yaml:"-"`

	Title       string    `yaml:"title"`
	Tags        []string  `yaml:"tags"`
	ShowUpdated bool      `yaml:"show_updated"`
	Draft       bool      `yaml:"draft"`
	Date        time.Time `yaml:"date"`
	Updated     time.Time `yaml:"updated"`
}

func (pm *PostMetadata) DisplayTitle() string {
	if pm.Title != "" {
		return pm.Title
	}

	formattedDate := pm.Date.Format("2006-01-02")

	slug := pm.Slug

	// Trim off common slug prefixes (like the date) because we
	// already display this.
	slug = strings.TrimPrefix(slug, formattedDate)
	slug = strings.TrimPrefix(slug, " - ")
	slug = strings.TrimPrefix(slug, "-")

	return slug
}

type cachedPost struct {
	Meta *PostMetadata
	Body string
}

func (s *Server) lookupPosts(includeDrafts bool) []*PostMetadata {
	s.cachedPostsLock.RLock()
	defer s.cachedPostsLock.RUnlock()

	var posts []*PostMetadata
	for _, post := range s.cachedPosts {
		// Skip all draft posts
		if !includeDrafts && post.Meta.Draft {
			continue
		}

		posts = append(posts, post.Meta)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	return posts
}

func (s *Server) lookupCachedPost(name string) *cachedPost {
	s.cachedPostsLock.RLock()
	defer s.cachedPostsLock.RUnlock()

	return s.cachedPosts[name]
}

func (s *Server) lookupPostFS(pm *PostMetadata) (fs.FS, error) {
	isDir, err := checkIsDir(s.Content, pm.Path)
	if err != nil {
		return nil, err
	}

	if !isDir {
		return nopFS, nil
	}

	return fs.Sub(s.Content, pm.Path)
}

func (s *Server) updateCachedPosts() error {
	s.cachedPostsLock.Lock()
	defer s.cachedPostsLock.Unlock()

	entries, err := fs.ReadDir(s.Content, "posts")
	if err != nil {
		return err
	}

	var errs []error
	newPosts := make(map[string]*cachedPost)

	for _, entry := range entries {
		name := entry.Name()

		post, err := s.readCachedPost(name)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
			continue
		}

		if _, ok := newPosts[post.Meta.Slug]; ok {
			errs = append(errs, fmt.Errorf("%s: slug collision with %s, overwriting", name, post.Meta.Path))
		}

		newPosts[post.Meta.Slug] = post
	}

	s.cachedPosts = newPosts

	if len(errs) > 0 {
		return &compositeError{
			Text:   "failed to update cached files",
			Errors: errs,
		}
	}

	return nil
}

func (s *Server) readCachedPost(name string) (*cachedPost, error) {
	targetPath := path.Join("posts", name)

	isDir, err := checkIsDir(s.Content, targetPath)
	if err != nil {
		return nil, err
	}

	var f fs.File

	if isDir {
		f, err = s.Content.Open(path.Join(targetPath, "index.gmi"))
	} else {
		f, err = s.Content.Open(targetPath)
	}

	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := bufio.NewReader(f)

	line, err := buf.Peek(4)
	if err != nil {
		return nil, err
	}

	meta := PostMetadata{
		Slug: strings.TrimSuffix(name, ".gmi"),
		Path: targetPath,
	}

	var frontmatterLines []string
	if bytes.Equal(line, []byte("---\n")) {
		_, _ = buf.Discard(4)

		for {
			line, err := buf.ReadString('\n')
			if err != nil {
				return nil, err
			}
			line = strings.TrimSuffix(line, "\n")

			if line == "---" {
				break
			}

			frontmatterLines = append(frontmatterLines, line)
		}

		err = yaml.Unmarshal([]byte(strings.Join(frontmatterLines, "\n")), &meta)
		if err != nil {
			return nil, err
		}

		if meta.Updated.IsZero() {
			info, err := f.Stat()
			if err != nil {
				return nil, err
			}

			meta.Updated = info.ModTime()
		}

		if meta.Date.IsZero() {
			meta.Date = meta.Updated
		}
	}

	// Skip any newlines at the start of the file
	for {
		b, err := buf.ReadByte()
		if err != nil {
			return nil, err
		}

		if b != '\n' {
			_ = buf.UnreadByte()
			break
		}
	}

	data, err := io.ReadAll(buf)
	if err != nil {
		return nil, err
	}

	return &cachedPost{
		Meta: &meta,
		Body: string(data),
	}, nil
}
