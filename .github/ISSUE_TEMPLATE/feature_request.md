---
name: Feature request
about: Propose something new for numbertext-go
title: ''
labels: enhancement
assignees: ''
---

**Problem**
What are you trying to do that numbertext-go doesn't support today?
Prefer a concrete scenario over an abstract one.

**Proposed solution**
What you'd like numbertext-go to do. If it's a new language or regional
variant, note whether it can be added purely as a `data/<code>.sor` file
(see `data/SOURCE.md#adding-a-new-language`) rather than a code change -
that's almost always true.

**Alternatives considered**
Anything else you tried, including `RegisterLocale` with your own
`.sor` source as a workaround.

**Soros spec, if relevant**
If this touches the interpreter itself (`internal/soros`), point at the
relevant part of the
[Soros language specification](https://github.com/Numbertext/libnumbertext/blob/master/doc/sorosspec.pdf).
See AGENTS.md for how this port's known, deliberate limitations
(backreferences in a matching pattern, one `bg.sor` bracket/call edge
case) are already scoped.
