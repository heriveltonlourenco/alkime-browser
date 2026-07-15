// Minimal demo: a native window with a counter and a button,
// rendered 100% in Go (no HTML, no JS, no V8).
//
// This proves the project's core thesis: it's possible to have a
// reactive UI running natively, without depending on a traditional browser.
package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"alkime-browser/pkg/reactive"
	"alkime-browser/pkg/ui"
)

func main() {
	// 1. Reactive state: a signal holding the counter value.
	counter := reactive.NewSignal(0)

	// 2. UI tree: a text node showing the counter and a button that
	// increments it. Note that Text is a function — it always reads
	// the current signal value at Draw time.
	label := &ui.Node{
		Kind: ui.KindText,
		Text: func() string {
			return fmt.Sprintf("Counter: %d", counter.Get())
		},
		X: 150, Y: 100,
	}

	button := &ui.Node{
		Kind: ui.KindButton,
		Text: func() string { return "Increment" },
		OnClick: func() {
			counter.Set(counter.Get() + 1)
		},
		X: 130, Y: 140, W: 140, H: 40,
	}

	app := ui.NewApp([]*ui.Node{label, button})

	ebiten.SetWindowSize(400, 300)
	ebiten.SetWindowTitle("MVP — Native UI in Go")

	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}
