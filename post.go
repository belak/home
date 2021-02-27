package home

import (
	"bufio"
	"bytes"
	"html/template"
	"io"
	"io/fs"
	"path"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

func (s *Server) lookupPostFS(pm *PostMetadata) (fs.FS, error) {
	isDir, _ := checkIsDir(s.Content, pm.Path)
	if !isDir {
		return nopFS, nil
	}

	return fs.Sub(s.Content, pm.Path)
}

func (s *Server) readPost(targetPath string) (*PostContext, error) {
	isDir, _ := checkIsDir(s.Content, targetPath)

	var err error
	var f fs.File

	if isDir {
		f, err = s.Content.Open(path.Join(targetPath, "index.gmi"))
	} else {
		f, err = s.Content.Open(strings.TrimSuffix(targetPath, ".gmi") + ".gmi")
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
		Slug: strings.TrimSuffix(path.Base(targetPath), ".gmi"),
		Path: targetPath,
	}

	// If the first line is ---, we need to scan through the file until another
	// --- is reached and parse what's between that as the frontmatter.
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

		if meta.Lastmod.IsZero() {
			info, err := f.Stat()
			if err != nil {
				return nil, err
			}
			meta.Lastmod = info.ModTime()
		}

		if meta.Date.IsZero() {
			meta.Date = meta.Lastmod
		}
	}

	// Skip any leading spaces at the start of the file
	for {
		b, err := buf.ReadByte()
		if err != nil {
			return nil, err
		}

		if !unicode.IsSpace(rune(b)) {
			_ = buf.UnreadByte()
			break
		}
	}

	data, err := io.ReadAll(buf)
	if err != nil {
		return nil, err
	}

	geminiContent := string(data)
	htmlContent := gemtextToHtml(geminiContent)

	return &PostContext{
		PostMetadata:  meta,
		GeminiContent: geminiContent,
		HtmlContent:   template.HTML(htmlContent),
	}, nil
}
