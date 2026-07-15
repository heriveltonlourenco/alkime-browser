# alkime-browser

A minimal proof of concept: a reactive UI rendered 100% in native Go,
with no HTML, no JavaScript, and no dependency on V8 or any
traditional browser engine.

## The thesis

Modern web frameworks (React, Vue, Svelte) still depend on the
browser — and therefore on V8 — to run. This project explores a
different path: a declarative, reactive UI that is compiled and
executed natively, without the JavaScript interpretation/JIT layer.

This repository is the first step: the smallest possible program that
proves the idea works.

## What this MVP proves

- A simple UI tree (`pkg/ui`)
- A signal-based reactivity system (`pkg/reactive`)
- Rendering via `ebiten`, with no HTML/CSS/JS involved
- Clicking a button updates state, which updates the screen — the
  full cycle of a reactive framework, in miniature

## How to run

\`\`\`bash
go mod tidy
go run ./cmd/demo
\`\`\`

## What doesn't exist yet (on purpose)

- Markup/template parser
- Layout engine (positions are hardcoded)
- Reusable components
- Support for multiple platforms/screens

This project is in its early stages and **open to contributions**.

## License

MIT — use it, modify it, contribute to it.
