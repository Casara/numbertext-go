---
name: Bug report
about: Something in numbertext-go isn't behaving as documented
title: ''
labels: bug
assignees: ''
---

**Version**
`numbertext-go` version/commit, `go version`.

**What happened**
A clear description of the behavior you saw (the exact output you got).

**What you expected**
What you expected instead. If it's about a specific language's wording,
say whether you checked it against the actual `data/<code>.sor` file or
against [numbertext.github.io](https://numbertext.github.io/) - locale
files encode real, sometimes surprising stylistic choices (see AGENTS.md's
"Verify locale output" note), so a wording difference isn't necessarily a
bug in this port.

**Minimal reproduction**
The exact call that reproduces it.

```go
// e.g. numbertext.Cardinal("en", "1001")
```

**Additional context**
Logs, stack trace, anything else relevant.
