package numbertext

// The embedded ".sor" sources themselves are not declared in this file:
// see locale_all.go (default: every bundled language) and locale_<code>.go
// (one per language, opt-in via build tags — see "Selecting which
// languages get embedded" in README.md). Both register into reg below
// through the same registry.register call, from an init function, so
// which set of files actually gets embedded is decided entirely by the
// build tags passed to `go build`/`go test`, with no other code change.

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/casara/numbertext-go/internal/soros"
)

// registry holds every known locale program (embedded plus any added at
// runtime via RegisterLocale), keyed by its language code exactly as
// used in the corresponding "*.sor" file name (e.g. "en", "hu_Hung").
type registry struct {
	mu       sync.RWMutex
	sources  map[string]string
	compiled map[string]*soros.Program
}

// reg starts out empty; it is populated before main() runs by the init
// functions in locale_all.go or locale_<code>.go (see the package doc
// comment above), and can be extended further at any time via
// RegisterLocale.
var reg = newRegistry()

func newRegistry() *registry {
	return &registry{
		sources:  make(map[string]string),
		compiled: make(map[string]*soros.Program),
	}
}

// program returns the compiled Soros program for lang, compiling and
// caching it on first use.
//
// A regional variant such as "en-GB" or "pt-BR" has no ".sor" file of
// its own: it is a set of "[:lang-code:]"-tagged lines inside the base
// language's file (spec 2.7.1), e.g. en.sor carries en-AU/en-GB/en-IE/
// en-NZ variants and pt.sor carries pt-BR. So when lang isn't a known
// source by itself, this falls back to the part before the first '-'
// (its base language) for the source text, while still compiling with
// the full requested code so those tagged lines take effect.
func (r *registry) program(lang string) (*soros.Program, error) {
	r.mu.RLock()
	prog, ok := r.compiled[lang]
	r.mu.RUnlock()

	if ok {
		return prog, nil
	}

	r.mu.RLock()
	source, ok := r.sources[lang]
	if !ok {
		if base, _, found := strings.Cut(lang, "-"); found {
			source, ok = r.sources[base]
		}
	}
	r.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf(
			"%w %q (see Languages for the supported base codes; "+
				"a regional variant like %q must build on one of them, e.g. \"en-GB\" on \"en\")",
			ErrUnknownLanguage, lang, lang,
		)
	}

	prog, err := soros.Compile(source, lang)
	if err != nil {
		return nil, fmt.Errorf("numbertext: compiling locale %q: %w", lang, err)
	}

	r.mu.Lock()
	r.compiled[lang] = prog
	r.mu.Unlock()

	return prog, nil
}

// register adds or replaces a locale's Soros source, invalidating any
// previously compiled program for that code.
func (r *registry) register(code, source string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sources[code] = source
	delete(r.compiled, code)
}

// languages returns every known language code, sorted.
func (r *registry) languages() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	codes := make([]string, 0, len(r.sources))
	for code := range r.sources {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	return codes
}

// Languages returns the language codes currently available, i.e. the
// bundled "data/*.sor" locales plus any added with RegisterLocale. Each
// code matches a ".sor" file's base name, for example "en", "pt", or the
// Hungarian long-scale variant "hu_Hung".
func Languages() []string {
	return reg.languages()
}

// RegisterLocale adds a new language at runtime from raw Soros source
// (the contents of a ".sor" file), without requiring a code change or a
// rebuild: this is the intended way to add a language Numbertext.org
// does not ship yet. code becomes the value passed as lang to every
// other function in this package (e.g. Cardinal, Ordinal).
//
// Registering a code that already exists replaces it.
func RegisterLocale(code, sorSource string) error {
	if code == "" {
		return ErrEmptyLocaleCode
	}

	_, err := soros.Compile(sorSource, code)
	if err != nil {
		return fmt.Errorf("numbertext: RegisterLocale(%q): %w", code, err)
	}

	reg.register(code, sorSource)

	return nil
}
