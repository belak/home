package home

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"unicode"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/gosimple/slug"
)

func gemtextToHtml(input string) string {
	w := &bytes.Buffer{}
	var preformatted bool
	var preformattedHint string
	var curBuf []string
	var inList bool

	for _, line := range strings.Split(input, "\n") {
		if preformatted {
			if line == "```" {
				var stripPrefix string

				// Determine what we will attempt to strip from preformatted
				// lines by using the first line as a template.
				if len(curBuf) > 0 {
					stripCount := len(curBuf[0]) - len(strings.TrimLeftFunc(curBuf[0], unicode.IsSpace))
					stripPrefix = curBuf[0][:stripCount]

					for i := range curBuf {
						curBuf[i] = strings.TrimPrefix(curBuf[i], stripPrefix)
					}
				}

				preformattedText := strings.Join(curBuf, "\n")

				lexer := lexers.Match(preformattedHint)
				if lexer == nil {
					lexer = lexers.Get(preformattedHint)
				}
				if lexer == nil {
					lexer = lexers.Analyse(preformattedText)
				}
				if lexer == nil {
					lexer = lexers.Fallback
				}
				lexer = chroma.Coalesce(lexer)

				formatter := html.New(
					html.WithClasses(true),
					html.PreventSurroundingPre(true),
				)

				iterator, err := lexer.Tokenise(nil, preformattedText)
				if err != nil {
					panic(err.Error())
				}

				// fmt.Println(lexer.Config().Name)

				fmt.Fprintf(w,
					"<pre class=\"chroma\"><code class=\"language-%s\">\n",
					strings.ToLower(lexer.Config().Name))

				err = formatter.Format(w, &chroma.Style{}, iterator)
				if err != nil {
					panic(err.Error())
				}

				fmt.Fprintf(w, "</code></pre>\n")

				curBuf = nil
				preformatted = !preformatted
				preformattedHint = ""
				continue
			}

			curBuf = append(curBuf, line)
			continue
		}

		// Lists are another special case because we need to detect when the
		// user starts and ends a list.
		if strings.HasPrefix(line, "*") {
			if !inList {
				fmt.Fprint(w, "<p><ul>\n")
			}
			inList = true

			fmt.Fprintf(w,
				"<li>%s</li>\n",
				template.HTMLEscapeString(strings.TrimSpace(line[1:])))

			continue
		} else {
			if inList {
				fmt.Fprint(w, "</ul></p>\n")
			}
			inList = false
		}

		if strings.HasPrefix(line, "```") {
			curBuf = nil
			preformattedHint = strings.TrimSpace(line[3:])
			preformatted = !preformatted
		} else if strings.HasPrefix(line, "###") {
			heading := strings.TrimSpace(line[3:])
			fmt.Fprintf(w,
				"<h3 id=\"%s\">%s</h3>\n",
				template.HTMLEscapeString(slug.Make(heading)),
				template.HTMLEscapeString(heading))
		} else if strings.HasPrefix(line, "##") {
			heading := strings.TrimSpace(line[2:])
			fmt.Fprintf(w,
				"<h2 id=\"%s\">%s</h2>\n",
				template.HTMLEscapeString(slug.Make(heading)),
				template.HTMLEscapeString(heading))
		} else if strings.HasPrefix(line, "#") {
			heading := strings.TrimSpace(line[1:])
			fmt.Fprintf(w,
				"<h1 id=\"%s\">%s</h1>\n",
				template.HTMLEscapeString(slug.Make(heading)),
				template.HTMLEscapeString(heading))
		} else if strings.HasPrefix(line, "=>") {
			split := strings.SplitN(strings.TrimSpace(line[2:]), " ", 2)
			fmt.Fprintf(w, "<p><a href=\"%s\">", template.HTMLEscapeString(split[0]))
			if len(split) == 2 {
				fmt.Fprintf(w, "%s</a></p>\n", template.HTMLEscapeString(split[1]))
			} else {
				fmt.Fprintf(w, "%s</a></p>\n", template.HTMLEscapeString(split[0]))
			}
		} else if line != "" {
			fmt.Fprintf(w, "<p>%s</p>\n", template.HTMLEscapeString(line))
		}
	}

	return w.String()
}
