package soros

import (
	"errors"
	"testing"
)

// TestParseTemplateUnterminatedBracket checks that a "[" with no
// matching "]" is reported as errUnterminated, propagating correctly
// out of parseNodes/parseBracket/parseTemplate.
func TestParseTemplateUnterminatedBracket(t *testing.T) {
	_, err := parseTemplate("foo[bar")
	if !errors.Is(err, errUnterminated) {
		t.Errorf("parseTemplate(%q) error = %v, want errUnterminated", "foo[bar", err)
	}
}

// TestParseTemplateUnterminatedCall checks that a "$(" with no matching
// ")" is reported as errUnterminated, propagating out of
// scanParens/parseCall/parseNodes/parseTemplate.
func TestParseTemplateUnterminatedCall(t *testing.T) {
	_, err := parseTemplate("$(foo")
	if !errors.Is(err, errUnterminated) {
		t.Errorf("parseTemplate(%q) error = %v, want errUnterminated", "$(foo", err)
	}
}

// TestParseTemplateUnterminatedBracketInsideCall checks that an
// unterminated "[" nested inside a call's argument propagates all the
// way out through parseCallArg and parseCall too, not just parseBracket
// itself.
func TestParseTemplateUnterminatedBracketInsideCall(t *testing.T) {
	_, err := parseTemplate("$([foo)")
	if !errors.Is(err, errUnterminated) {
		t.Errorf("parseTemplate(%q) error = %v, want errUnterminated", "$([foo)", err)
	}
}

// TestParseTemplateUnterminatedCallInsideBracket checks the opposite
// nesting: an unterminated "$(" inside a bracket propagates out through
// parseBracket, not just parseCall/scanParens.
func TestParseTemplateUnterminatedCallInsideBracket(t *testing.T) {
	_, err := parseTemplate("[$(foo]")
	if !errors.Is(err, errUnterminated) {
		t.Errorf("parseTemplate(%q) error = %v, want errUnterminated", "[$(foo]", err)
	}
}

// TestParseTemplateNestedCall checks that scanParens correctly tracks
// depth through a literal '(' inside an outer call's argument (here, a
// nested "$(...)" call), requiring two matching ')' before the outer
// call closes - not just the first one found.
func TestParseTemplateNestedCall(t *testing.T) {
	nodes, err := parseTemplate("$(a$(b)c)")
	if err != nil {
		t.Fatalf("parseTemplate: %v", err)
	}
	if len(nodes) != 1 {
		t.Fatalf("parseTemplate(%q) = %d nodes, want 1", "$(a$(b)c)", len(nodes))
	}
	outer, ok := nodes[0].(callNode)
	if !ok {
		t.Fatalf("nodes[0] is %T, want callNode", nodes[0])
	}
	if len(outer.arg) != 3 {
		t.Fatalf("outer.arg = %d nodes, want 3 (literal \"a\", inner call, literal \"c\")", len(outer.arg))
	}
	if _, ok := outer.arg[1].(callNode); !ok {
		t.Fatalf("outer.arg[1] is %T, want callNode (the nested $(b) call)", outer.arg[1])
	}
}

// TestAttachPipesTrailingPipe checks that a '|' immediately after a call
// with nothing following it (the call is the last node) still sets
// pipeAfter, the mirror image of TestPipeBoundaryModifier's leading-pipe
// case in soros_test.go.
func TestAttachPipesTrailingPipe(t *testing.T) {
	nodes, err := parseTemplate("$1|")
	if err != nil {
		t.Fatalf("parseTemplate: %v", err)
	}
	if len(nodes) != 1 {
		t.Fatalf("parseTemplate(%q) = %d nodes, want 1", "$1|", len(nodes))
	}
	call, ok := nodes[0].(callNode)
	if !ok {
		t.Fatalf("nodes[0] is %T, want callNode", nodes[0])
	}
	if !call.pipeAfter {
		t.Error("call.pipeAfter = false, want true for a trailing '|'")
	}
}
