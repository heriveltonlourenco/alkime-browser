# Draft "good first issue" tickets

Copy each section below into a separate GitHub Issue. Suggested labels
are noted at the top of each one.

---

## 1. Support multiple counters / buttons on screen

**Labels:** `good first issue`, `enhancement`

Right now `cmd/demo/main.go` hardcodes exactly one label and one
button. Extend the demo (or add a new one under `cmd/`) that renders
several independent counters, each with its own button and signal.

This is a good way to get familiar with `pkg/ui.Node` and
`pkg/reactive.Signal` without touching the renderer internals.

**Acceptance criteria:**
- At least 3 counters on screen, each incrementing independently
- No changes needed to `pkg/ui` or `pkg/reactive` (if you find you
  need one, mention it in the PR — that's useful signal too)

---

## 2. Make button colors configurable

**Labels:** `good first issue`, `enhancement`

`drawButton` in `pkg/ui/node.go` currently hardcodes the button color.
Add a `Color color.RGBA` field to `Node` (with a sensible default) so
callers can customize button appearance.

**Acceptance criteria:**
- `Node` has a new optional color field
- Existing demo still works with the default color if the field is
  left unset
- No hardcoded colors remain in `drawButton`

---

## 3. Add a decrement button and prevent negative counters

**Labels:** `good first issue`

Add a second button to the demo that decrements the counter, and
make sure the counter doesn't go below zero.

**Acceptance criteria:**
- Two buttons: increment and decrement
- Counter is clamped at 0 (decrementing at 0 does nothing)
- A short comment explains the clamping logic

---

## 4. Write unit tests for `reactive.Signal`

**Labels:** `good first issue`, `testing`

There are currently no tests. Add a `signal_test.go` file covering:
- `Get` returns the initial value before any `Set`
- `Set` updates the value returned by subsequent `Get`
- `Subscribe`d listeners are called on `Set`
- Multiple listeners are all called, in registration order

**Acceptance criteria:**
- `go test ./...` passes
- Tests use Go's standard `testing` package (no new dependencies)

---

## 5. Add a hover state to buttons

**Labels:** `good first issue`, `enhancement`

Buttons currently look the same whether the mouse is over them or
not. Use `ebiten.CursorPosition()` in `Update` (or `Draw`) to detect
hover and change the button's fill color slightly when hovered.

**Acceptance criteria:**
- Visually distinct hover state
- No dependency added beyond what's already imported

---

## 6. Document the reactivity model in `pkg/reactive/README.md`

**Labels:** `good first issue`, `documentation`

Write a short README explaining, in plain language, how `Signal`
works and why this pattern is used instead of a virtual DOM diffing
approach. Aimed at contributors who are new to reactive UI patterns.

**Acceptance criteria:**
- New `pkg/reactive/README.md`
- Explains `Get`/`Set`/`Subscribe` with a short example
- Links back to the root `README.md`
