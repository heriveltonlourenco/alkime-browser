// cmd/fetch proves alkime-browser can reach the web over HTTP and
// HTTPS using Go's standard net/http client, and shows the result in
// the native UI. An address bar lets the user type any URL. HTML
// responses are run through pkg/htmltext to show readable content
// with basic visual structure (bigger headings, underlined links) and
// inline style="color"/"background-color" — the first slice of real
// CSS. Still no selectors, no cascade, no real layout engine, no
// clickable links.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"alkime-browser/pkg/htmltext"
	"alkime-browser/pkg/reactive"
	"alkime-browser/pkg/ui"
)

const (
	defaultURL = "https://example.com/"

	// maxRawBytes caps how much of the response body we read/parse.
	// maxBlocks caps how many blocks we render (there's no scrolling,
	// so anything past the fixed-size window would be wasted work).
	// wrapWidth is how many characters fit on one line of the window.
	maxRawBytes = 512 * 1024
	maxBlocks   = 20
	wrapWidth   = 100
	blockGap    = 8

	addressY = 20
	addressH = 32
	statusY  = 64
	blocksY  = 92

	windowWidth  = 700
	windowHeight = 560
)

var client = &http.Client{Timeout: 10 * time.Second}

func main() {
	addressBar := reactive.NewSignal(defaultURL)
	statusLine := reactive.NewSignal("Press Enter to fetch the URL above.")

	address := &ui.Node{
		Kind:    ui.KindTextInput,
		Text:    addressBar.Get,
		Focused: true,
		OnChar: func(r rune) {
			addressBar.Set(addressBar.Get() + string(r))
		},
		OnBackspace: func() {
			v := addressBar.Get()
			if len(v) > 0 {
				addressBar.Set(v[:len(v)-1])
			}
		},
		X: 20, Y: addressY, W: windowWidth - 40, H: addressH,
	}

	status := &ui.Node{
		Kind: ui.KindText,
		Text: statusLine.Get,
		X:    20, Y: statusY,
	}

	app := ui.NewApp([]*ui.Node{address, status})
	app.Width, app.Height = windowWidth, windowHeight

	address.OnSubmit = func() {
		url := normalizeURL(addressBar.Get())
		line, blocks := doFetch(url)
		statusLine.Set(line)

		nodes := []*ui.Node{address, status}
		nodes = append(nodes, buildBlockNodes(blocks, blocksY)...)
		app.Nodes = nodes
	}

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("alkime-browser — address bar demo")

	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}

// buildBlockNodes converts extracted content blocks into ui.Nodes,
// stacking them vertically: each node's Y is the previous block's
// bottom edge plus a fixed gap. This is the entirety of this
// project's "layout engine" — a single top-to-bottom pass, no
// wrapping around other elements, no inline flow.
func buildBlockNodes(blocks []htmltext.Block, startY int) []*ui.Node {
	if len(blocks) > maxBlocks {
		blocks = blocks[:maxBlocks]
	}

	nodes := make([]*ui.Node, 0, len(blocks))
	y := startY
	for _, blk := range blocks {
		kind := ui.KindText
		switch blk.Kind {
		case htmltext.BlockHeading:
			kind = ui.KindHeading
		case htmltext.BlockLink:
			kind = ui.KindLink
		}

		text := wrap(blk.Text, wrapWidth)
		nodes = append(nodes, &ui.Node{
			Kind:       kind,
			Text:       func() string { return text },
			Color:      blk.Color,
			Background: blk.BackgroundColor,
			X:          20, Y: y,
		})
		y += ui.LineHeight(kind, text) + blockGap
	}
	return nodes
}

// normalizeURL adds an https:// scheme when the user didn't type one,
// mirroring how real browsers treat a bare "example.com" address.
func normalizeURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}
	return "https://" + raw
}

// doFetch performs a blocking HTTP(S) GET request and returns a
// status line plus the content as a list of blocks. HTML responses
// are parsed into paragraph/heading/link blocks; anything else (or
// any error along the way) comes back as a single paragraph block, so
// callers only ever need to deal with one shape of result. It runs
// synchronously on ebiten's Update goroutine (triggered from
// OnSubmit), so there's no risk of a data race on the Signals — the
// tradeoff is the window is unresponsive for the duration of the
// request, which is acceptable for this MVP.
func doFetch(url string) (statusLine string, blocks []htmltext.Block) {
	resp, err := client.Get(url)
	if err != nil {
		return "", errorBlock(fmt.Sprintf("Fetching %s\n\nError: %v", url, err))
	}
	defer resp.Body.Close()

	statusLine = fmt.Sprintf("%s — %s", url, resp.Status)
	limited := io.LimitReader(resp.Body, maxRawBytes)

	if isHTML(resp.Header.Get("Content-Type")) {
		bs, err := htmltext.ExtractBlocks(limited)
		if err != nil {
			return statusLine, errorBlock(fmt.Sprintf("Error parsing HTML: %v", err))
		}
		return statusLine, bs
	}

	body, err := io.ReadAll(limited)
	if err != nil {
		return statusLine, errorBlock(fmt.Sprintf("Error reading body: %v", err))
	}
	return statusLine, []htmltext.Block{{Kind: htmltext.BlockParagraph, Text: string(body)}}
}

func errorBlock(msg string) []htmltext.Block {
	return []htmltext.Block{{Kind: htmltext.BlockParagraph, Text: msg}}
}

// isHTML reports whether a response's Content-Type header indicates
// an HTML document, as opposed to plain text, JSON, images, etc.
func isHTML(contentType string) bool {
	return strings.Contains(contentType, "html")
}

// wrap breaks s into lines of at most width characters, preserving
// existing newlines, so a block fits inside the fixed-size native
// window instead of running off screen.
func wrap(s string, width int) string {
	var out strings.Builder
	for _, line := range strings.Split(s, "\n") {
		for len(line) > width {
			out.WriteString(line[:width])
			out.WriteByte('\n')
			line = line[width:]
		}
		out.WriteString(line)
		out.WriteByte('\n')
	}
	return strings.TrimRight(out.String(), "\n")
}
