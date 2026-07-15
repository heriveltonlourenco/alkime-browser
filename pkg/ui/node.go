// Package ui contains the minimal node tree and the renderer that
// draws these nodes on screen using ebiten. This is the simplest
// version possible: no markup parser, no layout engine — nodes are
// assembled directly in Go code with fixed position and size.
package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Kind identifies the type of a node. Only two exist in this MVP.
type Kind int

const (
	KindText Kind = iota
	KindButton
)

// Node is the basic unit of our UI tree. Each node knows how to draw
// itself and, if it's a button, how to react to clicks.
type Node struct {
	Kind    Kind
	Text    func() string // a function instead of a fixed string, so it always reads the current signal value
	OnClick func()
	X, Y    int
	W, H    int
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

// Update is called by ebiten every frame (typically 60x/second).
// Here we only check whether a mouse click happened and, if so,
// whether it landed inside the bounds of any button node.
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

// Layout defines the logical resolution of the window. Defaults to
// 400x300 (see NewApp) but callers can override App.Width/Height.
func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return a.Width, a.Height
}
