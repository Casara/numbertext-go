package numbertext

import "errors"

// Sentinel errors this package returns, wrapped (via fmt.Errorf's %w)
// with the context that produced them, so a caller can match on them
// with errors.Is instead of the message text.
var (
	// ErrUnknownLanguage is returned when a language code (and, for a
	// regional variant, its base language too) has neither a bundled
	// nor a RegisterLocale-added source.
	ErrUnknownLanguage = errors.New("numbertext: unknown language")

	// ErrEmptyLocaleCode is returned by RegisterLocale when code is "".
	ErrEmptyLocaleCode = errors.New("numbertext: locale code must not be empty")
)
