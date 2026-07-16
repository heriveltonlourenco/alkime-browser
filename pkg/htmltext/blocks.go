package htmltext

import (
	"image/color"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// BlockKind identifies the kind of a content Block, letting a caller
// render headings and links differently from regular paragraphs
// without needing to understand HTML or CSS.
type BlockKind int

const (
	BlockParagraph BlockKind = iota
	BlockHeading
	BlockLink
)

// Block is a single piece of visible page content, extracted from an
// HTML document and stripped of markup.
type Block struct {
	Kind BlockKind
	Text string
	Href string // set only for BlockLink

	// Color and BackgroundColor come from the element's own inline
	// style="..." attribute (see style.go) — nil when unset. This is
	// the first slice of CSS support: no selectors, no <style>
	// blocks, no cascade, just one element's own declarations.
	Color           *color.RGBA
	BackgroundColor *color.RGBA
}

var headingTags = map[string]bool{
	"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
}

// ExtractBlocks parses the HTML document read from r and returns its
// visible content as a sequence of blocks, distinguishing headings
// and links from regular paragraph text. Unlike Extract, this keeps
// enough structure for a caller to render each block differently
// (bigger heading, underlined link, inline style colors) — still no
// real layout, blocks are meant to be stacked vertically in document
// order.
//
// A link nested inside a paragraph is pulled out as its own Block,
// splitting the surrounding paragraph in two, rather than rendering
// inline within a line of text: this project has no inline text
// layout (mixed runs of styled and unstyled text sharing one line),
// only block stacking.
func ExtractBlocks(r io.Reader) ([]Block, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	c := &collector{}
	c.walk(doc)
	c.flushParagraph(nil, nil)
	return c.blocks, nil
}

type collector struct {
	blocks  []Block
	current strings.Builder
}

func (c *collector) flushParagraph(fg, bg *color.RGBA) {
	if text := collapse(c.current.String()); text != "" {
		c.blocks = append(c.blocks, Block{Kind: BlockParagraph, Text: text, Color: fg, BackgroundColor: bg})
	}
	c.current.Reset()
}

func (c *collector) walk(n *html.Node) {
	if n.Type == html.ElementNode && skipTags[n.Data] {
		return
	}

	if n.Type == html.ElementNode && headingTags[n.Data] {
		c.flushParagraph(nil, nil)
		if text := collapse(textOf(n)); text != "" {
			fg, bg := parseInlineStyle(attr(n, "style"))
			c.blocks = append(c.blocks, Block{Kind: BlockHeading, Text: text, Color: fg, BackgroundColor: bg})
		}
		return
	}

	if n.Type == html.ElementNode && n.Data == "a" {
		c.flushParagraph(nil, nil)
		if text := collapse(textOf(n)); text != "" {
			fg, bg := parseInlineStyle(attr(n, "style"))
			c.blocks = append(c.blocks, Block{
				Kind: BlockLink, Text: text, Href: attr(n, "href"),
				Color: fg, BackgroundColor: bg,
			})
		}
		return
	}

	if n.Type == html.TextNode {
		if text := strings.TrimSpace(n.Data); text != "" {
			c.current.WriteString(text)
			c.current.WriteString(" ")
		}
	}

	for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
		c.walk(ch)
	}

	if n.Type == html.ElementNode && blockTags[n.Data] {
		fg, bg := parseInlineStyle(attr(n, "style"))
		c.flushParagraph(fg, bg)
	}
}

// textOf collects all visible text under n, ignoring skipTags — used
// for headings and links, which are extracted as one block instead
// of being walked node-by-node into the running paragraph buffer.
func textOf(n *html.Node) string {
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
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
			walk(c)
		}
	}
	walk(n)
	return b.String()
}

func attr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}
