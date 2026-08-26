# Go Style Guide

*[Leia em português](coding-style.pt-BR.md)*

## Purpose

This document defines additional conventions adopted by the project beyond the rules already enforced automatically
by tools such as `gofmt`, `goimports`, and `golangci-lint`.

Whenever possible, decisions should prioritize:

* readability;
* consistency;
* maintainability;
* diff quality.

---

## Formatting

All code must be formatted with:

* gofmt
* goimports

Manual changes that go against the formatting produced by these tools must not be made.

---

## Readability over conciseness

Prefer explicit, easy-to-understand code over excessively compact versions.

Preferred:

```go
if err != nil {
    return err
}
```

Avoid:

```go
if err != nil { return err }
```

---

## Function calls

Simple calls should stay on a single line.

```go
prog, err := reg.program(lang)

groups := match.Groups()
```

Calls with multiple arguments, options, or nested structures should use the vertical format.

```go
return nil, fmt.Errorf(
    "%w %q (see Languages for the supported base codes; "+
        "a regional variant like %q must build on one of them, e.g. \"en-GB\" on \"en\")",
    ErrUnknownLanguage, lang, lang,
)
```

---

## Diff quality

Whenever a construct has a natural tendency to grow, prefer the vertical format.

Example:

```go
declare -a targets=(
    "README.md|the prose around Usage/Regional variants/Gender"
    "CHANGELOG.md|every public change belongs here"
)
```

This format reduces merge conflicts and produces smaller diffs when new elements are added.

---

## Comments

Comments should explain:

* purpose;
* behavior;
* limitations;
* design decisions.

Comments that merely repeat the function name should be avoided.

Bad:

```go
// Compile compiles a program.
```

Better:

```go
// Compile parses Soros source code into a runnable Program for the
// given language code. It only fails outright if the source cannot be
// tokenized at all; individual rules it cannot compile are recorded in
// Program.Skipped instead of failing the whole locale.
```

---

## Error Handling (wrapcheck)

Errors returned by external dependencies or lower layers should receive additional context before being
propagated.

The goal is to make the origin of the failure evident to a caller matching on it with `errors.Is`, and to ease
diagnosis when a `.sor` file fails to compile.

### Rule

When returning an error received from another function, add context using `fmt.Errorf` and `%w`.

Correct:

```go
prog, err := soros.Compile(source, lang)
if err != nil {
    return nil, fmt.Errorf("numbertext: compiling locale %q: %w", lang, err)
}
```

Avoid:

```go
prog, err := soros.Compile(source, lang)
if err != nil {
    return nil, err
}
```

---

### Error message

The message should describe the operation that failed, not repeat the original error's text.

Correct:

```go
return fmt.Errorf("numbertext: RegisterLocale(%q): %w", code, err)
```

Avoid:

```go
return fmt.Errorf("error: %w", err)
```

```go
return fmt.Errorf("failed: %w", err)
```

```go
return fmt.Errorf("unexpected error: %w", err)
```

---

### When wrapping is not necessary

Do not wrap when:

* creating a new error;
* returning a sentinel error (e.g. `ErrUnknownLanguage`, `ErrEmptyLocaleCode` from `errors.go`);
* the error already contains sufficient context and the current layer adds no relevant information.

Examples:

```go
return ErrEmptyLocaleCode
```

```go
return errUnterminated
```

---

### Guidance for AIs

When fixing `wrapcheck` violations:

1. Preserve the error chain using `%w`.
2. Describe the operation that failed (e.g. "compiling locale %q", not "error compiling").
3. Do not use generic messages such as:

   * "error"
   * "failed"
   * "unexpected error"
4. Do not repeat context already present in lower layers.
5. Prefer short, lowercase messages, prefixed with the package name where the existing code already does so
   (`"numbertext: ..."`).
6. Include relevant identifiers when they add value:

```go
return fmt.Errorf("numbertext: compiling locale %q: %w", lang, err)
```

---

## Dependencies

Prefer explicit dependencies via function parameters and constructor-style values over global variables.

Preferred:

```go
func compilePattern(wrapped string) (matcher, error)
```

Avoid:

```go
var globalCompiledPattern matcher
```

Unless the nature of the component clearly justifies a shared singleton, as `reg` (the package-level locale
registry in `locales.go`) does: it exists once per process, is safe for concurrent use, and is exactly the
handle `RegisterLocale`/`Languages`/`Convert` all need to share.
