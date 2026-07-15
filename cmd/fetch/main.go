// cmd/fetch proves alkime-browser can reach the web over HTTP and
// HTTPS using Go's standard net/http client, and shows the raw
// response in the native UI. An address bar lets the user type any
// URL — no HTML parsing yet, the response body is treated as opaque
// plain text.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"alkime-browser/pkg/reactive"
	"alkime-browser/pkg/ui"
)

const (
	defaultURL = "https://example.com/"

	// how much of the response body to show on screen, and how many
	// characters fit on one line of the fixed-size window below.
	maxBodyChars = 600
	wrapWidth    = 100

	windowWidth  = 700
	windowHeight = 480
)

var client = &http.Client{Timeout: 10 * time.Second}

func main() {
	addressBar := reactive.NewSignal(defaultURL)
	result := reactive.NewSignal("Press Enter to fetch the URL above.")

	navigate := func() {
		result.Set(doFetch(normalizeURL(addressBar.Get())))
	}

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
		OnSubmit: navigate,
		X:        20, Y: 20, W: windowWidth - 40, H: 32,
	}

	status := &ui.Node{
		Kind: ui.KindText,
		Text: result.Get,
		X:    20, Y: 70,
	}

	app := ui.NewApp([]*ui.Node{address, status})
	app.Width, app.Height = windowWidth, windowHeight

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("alkime-browser — address bar demo")

	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
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

// doFetch performs a blocking HTTP(S) GET request and formats the
// result as plain text. It runs synchronously on ebiten's Update
// goroutine (triggered from OnSubmit), so there's no risk of a data
// race on the Signal — the tradeoff is the window is unresponsive for
// the duration of the request, which is acceptable for this MVP.
func doFetch(url string) string {
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Sprintf("Fetching %s\n\nError: %v", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyChars))
	if err != nil {
		return fmt.Sprintf("Fetching %s\n\n%s\n\nError reading body: %v", url, resp.Status, err)
	}

	return fmt.Sprintf("%s\n%s\n\n%s", url, resp.Status, wrap(string(body), wrapWidth))
}

// wrap breaks s into lines of at most width characters, preserving
// existing newlines, so the response fits inside the fixed-size
// native window instead of running off screen.
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
