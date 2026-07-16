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
- An HTTP/HTTPS network layer with a real address bar (`cmd/fetch`),
  using Go's standard `net/http` client: type any URL, press Enter,
  and see the page's readable text in the native UI
- HTML text extraction (`pkg/htmltext`), using `golang.org/x/net/html`
  to strip tags, scripts and styles from a fetched page, leaving only
  the words a reader would actually see — no visual structure yet
  (headings, lists and links all look like plain text)

## How to run

\`\`\`bash
go mod tidy
go run ./cmd/demo   # native reactive UI counter demo
go run ./cmd/fetch  # HTTP/HTTPS network fetch demo
\`\`\`

## What doesn't exist yet (on purpose)

- Visual HTML structure: headings, lists and links are extracted as
  plain text (`pkg/htmltext`), not rendered with any distinct look,
  and links aren't clickable yet
- A real layout engine: `cmd/fetch`'s node positions are still
  hardcoded; the only "layout" is `pkg/htmltext` inserting blank
  lines between blocks
- Reusable components
- Support for multiple platforms/screens

This project is in its early stages and **open to contributions**.

## License

MIT — use it, modify it, contribute to it.
