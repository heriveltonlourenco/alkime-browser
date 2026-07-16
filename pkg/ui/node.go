// Package ui contains the minimal node tree and the renderer that
// draws these nodes on screen using ebiten. This is the simplest
// version possible: no CSS, no real layout engine — nodes are
// assembled and positioned directly in Go code.
package ui

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

// Approximate metrics of ebitenutil's built-in debug font, used to
// size and stack nodes without any real font/layout library.
const (
	debugCharWidth  = 6
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
			ebitenutil.DebugPrintAt(screen, n.Text(), n.X, n.Y)
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

	text := n.Text()
	if n.Focused {
		text += "_" // static cursor: good enough to signal "you can type here" without a blink timer
	}
	ebitenutil.DebugPrintAt(screen, text, n.X+8, n.Y+n.H/2-4)
}

// drawHeading renders text larger than normal by drawing the regular
// debug font onto a small offscreen image and scaling that image up —
// there's no variable-size font available, so this is the cheapest
// way to make a heading look bigger without a new font dependency.
func drawHeading(screen *ebiten.Image, n *Node) {
	text := n.Text()
	lines := strings.Split(text, "\n")
	longest := 0
	for _, line := range lines {
		if len(line) > longest {
			longest = len(line)
		}
	}
	w := longest*debugCharWidth + 10
	h := len(lines)*debugLineHeight + 4
	if w <= 0 || h <= 0 {
		return
	}

	buf := ebiten.NewImage(w, h)
	ebitenutil.DebugPrintAt(buf, text, 0, 0)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(headingScale, headingScale)
	opts.GeoM.Translate(float64(n.X), float64(n.Y))
	screen.DrawImage(buf, opts)
}

// drawLink renders text at normal size with a solid underline drawn
// beneath each wrapped line — the closest approximation of "looks
// like a hyperlink" available without per-character text coloring.
func drawLink(screen *ebiten.Image, n *Node) {
	text := n.Text()
	ebitenutil.DebugPrintAt(screen, text, n.X, n.Y)

	underlineColor := color.RGBA{R: 96, G: 165, B: 250, A: 255}
	for i, line := range strings.Split(text, "\n") {
		width := len(line) * debugCharWidth
		if width <= 0 {
			continue
		}
		rect := ebiten.NewImage(width, 1)
		rect.Fill(underlineColor)

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
