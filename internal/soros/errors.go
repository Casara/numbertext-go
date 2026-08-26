package soros

import "errors"

// errUnterminated is wrapped (via fmt.Errorf's %w) into the more
// specific, context-carrying errors the template parser returns for an
// unclosed "[" or "$(", so a caller can match on it with errors.Is
// instead of the message text.
var errUnterminated = errors.New("unterminated")
