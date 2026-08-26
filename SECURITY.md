# Security Policy

## Supported Versions

Before `v1.0.0`, only the latest released version is supported with
security fixes. Once `v1.0.0` ships, this section will list which
major versions still receive patches.

## Threat model

`numbertext-go` has no network, filesystem, or auth surface of its
own - the two things worth being deliberate about are:

- **`RegisterLocale` compiles caller-supplied `.sor` text into regular
  expressions, via a hybrid engine (see `internal/soros/regex.go` and
  AGENTS.md's "Regex engine" section): Go's `regexp`/RE2 first, which
  guarantees linear-time matching, falling back to the backtracking
  `regexp2` engine only for the rare pattern RE2 can't compile at all
  (backreferences/lookaround).** That fallback is timeout-bounded
  (`fallbackMatchTimeout`) so a pathological pattern degrades to "no
  match" rather than hanging, but a timeout is a mitigation, not
  immunity - a `.sor` source is still arbitrary input if it comes from
  outside your own codebase, and this is the specific reason to treat
  it like any other configuration you'd review before loading, not
  like untrusted end-user input.
- **`Cardinal`/`Ordinal`/`Year`/`Currency`/`Money`/`Convert` recurse**
  (see CLAUDE.md/AGENTS.md's description of the Soros evaluator). The
  bundled `data/*.sor` files are bounded in practice; a hostile custom
  locale registered via `RegisterLocale` could in principle be written
  to recurse very deeply on a crafted input. The same caution as above
  applies: only register locale sources you trust.

## Reporting a Vulnerability

Please **do not** open a public issue for a security vulnerability.

Use GitHub's private vulnerability reporting instead: go to the
[Security tab](https://github.com/casara/numbertext-go/security) of
this repository and select "Report a vulnerability." This opens a
private advisory visible only to the maintainers until a fix is
ready.

Include, if possible:

- The affected version/commit.
- A minimal reproduction (the exact call and input that triggers it).
- The impact you believe it has.

We'll acknowledge the report and work with you on a disclosure
timeline before any public advisory is published.
