// Package ui contains the minimal node tree and the renderer that
// draws these nodes on screen using ebiten. This is the simplest
// version possible: no CSS cascade, no real layout engine — nodes
// are assembled and positioned directly in Go code.
package ui

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// Kind identifies the type of a node.
type Kind int

const (
	KindText Kind = iota
	KindButton
	KindTextInput
	KindHeading // like KindText, but rendered larger — no other styling
	KindLink    // like KindText, but rendered with an underline — not clickable yet
)

// contentFace is the font used for page content (KindText/KindHeading
// /KindLink), chosen specifically because it supports per-call color
// — buttons and the address bar keep using ebitenutil.DebugPrintAt's
// fixed white-only font, since UI chrome doesn't need CSS styling.
var contentFace = basicfont.Face7x13
var contentAscent = contentFace.Metrics().Ascent.Round()

// Approximate metrics used to size and stack nodes without any real
// layout library: debugCharWidth/debugLineHeight match contentFace's
// actual advance/line spacing (with a little breathing room on the
// line height).
const (
	debugCharWidth  = 7
	debugLineHeight = 16
	headingScale    = 1.6
)

// Node is the basic unit of our UI tree. Each node knows how to draw
// itself and, depending on its Kind, how to react to clicks or
// keyboard input.
type Node struct {
	Kind    Kind
	Text    func() string // a function instead of a fixed string, so it always reads the current signal value
	OnClick func()
	X, Y    int
	W, H    int

	// KindText/KindHeading/KindLink only: optional inline-style
	// overrides (see pkg/htmltext's Block.Color/BackgroundColor).
	// nil means "use the kind's default" (white text; underline-blue
	// for links; no background).
	Color      *color.RGBA
	Background *color.RGBA

	// KindTextInput only: the node doesn't own its text — like
	// KindText, the value comes from Text(). These callbacks let the
	// caller update its own reactive.Signal in response to keystrokes,
	// keeping "who owns the state" consistent across node kinds.
	Focused     bool
	OnChar      func(r rune)
	OnBackspace func()
	OnSubmit    func()
}

// LineHeight returns the approximate rendered height, in pixels, of
// a text-bearing node's content (KindText, KindHeading, KindLink) —
// useful for stacking blocks vertically without a real layout engine.
func LineHeight(kind Kind, text string) int {
	lines := strings.Count(text, "\n") + 1
	h := lines * debugLineHeight
	if kind == KindHeading {
		h = int(float64(h) * headingScale)
	}
	return h
}

// App is the minimal "engine": it holds the node tree and implements
// the ebiten.Game interface (Update/Draw/Layout), which is the main
// loop of any ebiten application.
type App struct {
	Nodes       []*Node
	NeedsRedraw bool // kept for conceptual clarity; ebiten already redraws every frame
	Width       int  // logical resolution, in pixels
	Height      int
}

// NewApp creates a new application with the given nodes. Width and
// Height default to the original MVP's fixed 400x300 canvas; set them
// on the returned App before RunGame if a demo needs more room (e.g.
// to show a fetched HTTP response).
func NewApp(nodes []*Node) *App {
	return &App{Nodes: nodes, Width: 400, Height: 300}
}

// Update is called by ebiten every frame (typically 60x/second): it
// dispatches mouse clicks to buttons and keyboard input to whichever
// text input node is focused.
func (a *App) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		for _, n := range a.Nodes {
			if n.Kind == KindButton && n.OnClick != nil {
				if mx >= n.X && mx <= n.X+n.W && my >= n.Y && my <= n.Y+n.H {
					n.OnClick()
				}
			}
		}
	}

	for _, n := range a.Nodes {
		if n.Kind != KindTextInput || !n.Focused {
			continue
		}
		for _, r := range ebiten.AppendInputChars(nil) {
			if n.OnChar != nil {
				n.OnChar(r)
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && n.OnBackspace != nil {
			n.OnBackspace()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && n.OnSubmit != nil {
			n.OnSubmit()
		}
	}

	return nil
}

// Draw is called by ebiten every frame to render the current state
// of the screen. In this MVP we redraw the entire tree every time —
// there's no diffing. At this scale (a handful of nodes), that's
// irrelevant for performance.
func (a *App) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 24, G: 24, B: 27, A: 255}) // dark background

	for _, n := range a.Nodes {
		switch n.Kind {
		case KindButton:
			drawButton(screen, n)
		case KindTextInput:
			drawTextInput(screen, n)
		case KindHeading:
			drawHeading(screen, n)
		case KindLink:
			drawLink(screen, n)
		case KindText:
			drawContentText(screen, n)
		}
	}
}

func drawButton(screen *ebiten.Image, n *Node) {
	buttonColor := color.RGBA{R: 63, G: 63, B: 70, A: 255}
	rect := ebiten.NewImage(n.W, n.H)
	rect.Fill(buttonColor)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(n.X), float64(n.Y))
	screen.DrawImage(rect, opts)

	ebitenutil.DebugPrintAt(screen, n.Text(), n.X+8, n.Y+n.H/2-4)
}

func drawTextInput(screen *ebiten.Image, n *Node) {
	fieldColor := color.RGBA{R: 40, G: 40, B: 46, A: 255}
	if n.Focused {
		fieldColor = color.RGBA{R: 58, G: 58, B: 70, A: 255} // lighter when active, as a simple focus ring substitute
	}
	rect := ebiten.NewImage(n.W, n.H)
	rect.Fill(fieldColor)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(n.X), float64(n.Y))
	screen.DrawImage(rect, opts)

	str := n.Text()
	if n.Focused {
		str += "_" // static cursor: good enough to signal "you can type here" without a blink timer
	}
	ebitenutil.DebugPrintAt(screen, str, n.X+8, n.Y+n.H/2-4)
}

// textBounds estimates the pixel size of str using the approximate
// content-font metrics, for sizing background fills and the offscreen
// buffer drawHeading scales up.
func textBounds(s string) (w, h int) {
	longest := 0
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		if len(line) > longest {
			longest = len(line)
		}
	}
	return longest*debugCharWidth + 4, len(lines)*debugLineHeight + 4
}

// drawTextBlock draws str onto dst starting at (x, y): an optional
// background fill sized to the text, followed by each line in fg
// using contentFace — the one shared primitive behind
// drawContentText, drawHeading (via an offscreen buffer) and drawLink.
func drawTextBlock(dst *ebiten.Image, str string, x, y int, fg color.Color, bg *color.RGBA) {
	if bg != nil {
		w, h := textBounds(str)
		rect := ebiten.NewImage(w, h)
		rect.Fill(*bg)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x), float64(y))
		dst.DrawImage(rect, opts)
	}
	for i, line := range strings.Split(str, "\n") {
		text.Draw(dst, line, contentFace, x, y+i*debugLineHeight+contentAscent, fg)
	}
}

func drawContentText(screen *ebiten.Image, n *Node) {
	fg := color.Color(color.White)
	if n.Color != nil {
		fg = *n.Color
	}
	drawTextBlock(screen, n.Text(), n.X, n.Y, fg, n.Background)
}

// drawHeading renders text larger than normal by drawing it onto a
// small offscreen image and scaling that image up — there's no
// variable-size font available, so this is the cheapest way to make
// a heading look bigger without a new font dependency.
func drawHeading(screen *ebiten.Image, n *Node) {
	str := n.Text()
	w, h := textBounds(str)
	if w <= 0 || h <= 0 {
		return
	}

	fg := color.Color(color.White)
	if n.Color != nil {
		fg = *n.Color
	}
	buf := ebiten.NewImage(w, h)
	drawTextBlock(buf, str, 0, 0, fg, n.Background)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(headingScale, headingScale)
	opts.GeoM.Translate(float64(n.X), float64(n.Y))
	screen.DrawImage(buf, opts)
}

// drawLink renders text at normal size with a solid underline drawn
// beneath each wrapped line, in the same color as the text — the
// closest approximation of "looks like a hyperlink" available without
// a real font or per-character styling. Defaults to a classic link
// blue when no inline color was set.
func drawLink(screen *ebiten.Image, n *Node) {
	fg := color.Color(color.RGBA{R: 96, G: 165, B: 250, A: 255})
	if n.Color != nil {
		fg = *n.Color
	}
	str := n.Text()
	drawTextBlock(screen, str, n.X, n.Y, fg, n.Background)

	for i, line := range strings.Split(str, "\n") {
		width := len(line) * debugCharWidth
		if width <= 0 {
			continue
		}
		rect := ebiten.NewImage(width, 1)
		rect.Fill(fg)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(n.X), float64(n.Y+i*debugLineHeight+14))
		screen.DrawImage(rect, opts)
	}
}

// Layout defines the logical resolution of the window. Defaults to
// 400x300 (see NewApp) but callers can override App.Width/Height.
func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return a.Width, a.Height
}
