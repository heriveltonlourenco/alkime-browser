# Contributing to alkime-browser

Thanks for considering a contribution — this project is intentionally
early-stage, and early contributors have an outsized impact on its
direction.

## Philosophy

This is not meant to be any single person's project. The goal is a
native, reactive UI runtime that doesn't depend on a browser or on
V8/JavaScript. If that idea interests you, you're welcome here
regardless of experience level.

## Before you start

- Check open [Issues](../../issues) first — especially ones labeled
  `good first issue`, which are scoped to be approachable without
  needing to understand the whole codebase.
- For anything larger than a small fix (a new subsystem, a breaking
  API change, a new rendering backend), please open an issue or
  discussion first so we can align on direction before you invest
  time in a PR.
- No idea is too rough for a Discussion post. Early architecture
  debates are welcome and expected.

## Development setup

```bash
git clone <this-repo>
cd alkime-browser
go mod tidy
go run ./cmd/demo
```

Requires Go 1.22+.

## Code style

- Run `go fmt ./...` and `go vet ./...` before submitting.
- Keep comments explaining *why*, not just *what* — this project
  values being approachable to newcomers more than being terse.
- Prefer small, focused PRs over large ones. A PR that does one thing
  well is easier to review and merge.

## Project structure
/cmd/demo        → example application(s)
/pkg/reactive    → the signal-based reactivity system
/pkg/ui          → the node tree and renderer

As the project grows, expect this to expand into `/pkg/layout`,
`/pkg/parser`, etc. — see open issues and discussions for current
architectural plans.

## Submitting a PR

1. Fork the repo and create a branch from `main`.
2. Make your change, with tests if applicable.
3. Make sure `go build ./...` and `go vet ./...` pass.
4. Open a PR describing what changed and why. Link the issue it
   addresses, if any.

## Code of conduct

Be respectful, assume good intent, and remember that this is a
learning project for many contributors — questions are always
welcome, no matter how basic they seem.
