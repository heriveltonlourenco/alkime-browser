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
- HTML content extraction (`pkg/htmltext`), using
  `golang.org/x/net/html` to strip tags, scripts and styles from a
  fetched page and split it into headings, paragraphs and links
- Basic visual structure (`cmd/fetch` renders headings larger and
  links underlined) plus a first real slice of CSS: inline
  `style="color: ...; background-color: ..."` on an element is parsed
  and actually rendered — still no selectors, `<style>` blocks or
  cascade, just one element's own declarations

## How to run

\`\`\`bash
go mod tidy
go run ./cmd/demo   # native reactive UI counter demo
go run ./cmd/fetch  # HTTP/HTTPS network fetch demo
\`\`\`

## What doesn't exist yet (on purpose)

- Most of CSS: no selectors, no `<style>`/external stylesheets, no
  cascade, no properties beyond `color`/`background-color` — only an
  element's own inline `style="..."` is read
- A real layout engine: blocks are stacked top-to-bottom in a single
  pass (`cmd/fetch`'s `buildBlockNodes`); no inline flow (a link
  inside a paragraph gets pulled onto its own line), no reflow, no
  scrolling for pages longer than the window
- Clickable links
- Reusable components
- Support for multiple platforms/screens

This project is in its early stages and **open to contributions**.

## License

MIT — use it, modify it, contribute to it.
