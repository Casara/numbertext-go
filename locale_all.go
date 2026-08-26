//go:build !numbertext_select

package numbertext

import (
	"embed"
	"fmt"
	"path"
	"strings"
)

// Default build: every language under data/*.sor is embedded and
// registered. Build with -tags numbertext_select (plus one
// -tags numbertext_lang_<code> per language you actually want) to embed
// only a subset instead — see locale_<code>.go and "Selecting which
// languages get embedded" in README.md.
//
//go:embed data/*.sor
var allLocaleData embed.FS

func init() {
	entries, err := allLocaleData.ReadDir("data")
	if err != nil {
		panic(fmt.Sprintf("numbertext: embedded data directory is missing: %v", err))
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".sor") {
			continue
		}
		content, err := allLocaleData.ReadFile(path.Join("data", name))
		if err != nil {
			panic(fmt.Sprintf("numbertext: reading embedded %s: %v", name, err))
		}
		reg.register(strings.TrimSuffix(name, ".sor"), string(content))
	}
}
