// Package htmltext extracts the human-readable text content of an
// HTML document, discarding markup, scripts and styles. It's the
// first step towards rendering fetched pages instead of just
// dumping raw markup as plain text — no layout, no structure, just
// the words a reader would actually see.
package htmltext

import (
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// skipTags never contribute visible text (their content is markup,
// code or metadata, not page content).
var skipTags = map[string]bool{
	"script":   true,
	"style":    true,
	"head":     true,
	"noscript": true,
	"template": true,
}

// blockTags get a blank line after their content, so paragraphs,
// list items, headings, etc. don't run together into a wall of text.
var blockTags = map[string]bool{
	"p": true, "div": true, "li": true, "br": true, "hr": true,
	"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
	"header": true, "footer": true, "section": true, "article": true,
	"tr": true, "table": true, "ul": true, "ol": true,
	"blockquote": true, "pre": true, "nav": true, "main": true, "aside": true, "form": true,
}

// Extract parses the HTML document read from r and returns its
// visible text content.
func Extract(r io.Reader) (string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	walk(doc, &b)
	return collapse(b.String()), nil
}

func walk(n *html.Node, b *strings.Builder) {
	if n.Type == html.ElementNode && skipTags[n.Data] {
		return
	}
	if n.Type == html.TextNode {
		if text := strings.TrimSpace(n.Data); text != "" {
			b.WriteString(text)
			b.WriteString(" ")
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walk(c, b)
	}
	if n.Type == html.ElementNode && blockTags[n.Data] {
		b.WriteString("\n\n")
	}
}

var (
	runSpaces  = regexp.MustCompile(`[ \t]+`)
	blankLines = regexp.MustCompile(`\n{3,}`)
)

// collapse tidies up the whitespace produced by walk: runs of spaces
// become one, and more than one blank line between blocks collapses
// down to a single blank line.
func collapse(s string) string {
	s = runSpaces.ReplaceAllString(s, " ")
	s = strings.ReplaceAll(s, " \n", "\n")
	s = blankLines.ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s)
}
