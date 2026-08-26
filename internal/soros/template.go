// Package soros implements an interpreter for the Soros programming
// language, the small regex-rewriting DSL used by the Numbertext.org
// project (https://numbertext.github.io/) to describe number-to-text
// conversion rules on a per-language basis.
package soros

import (
	"fmt"
	"strings"
)

// node is one element of a parsed replacement template.
type node any

// literalNode is plain text emitted as-is.
type literalNode struct {
	text string
}

// backrefNode is a raw "\N" reference: the captured text of group N is
// substituted verbatim, without being re-run through the interpreter.
type backrefNode struct {
	group int
}

// callNode is a recursive call, either the abbreviated "$N" form or the
// full "$(...)" form. arg is the (already parsed) template used to build
// the call's argument string before recursing.
type callNode struct {
	arg        []node
	pipeBefore bool
	pipeAfter  bool
}

// optionalNode is a "[...]" conditional block: its rendered text is kept
// only if every callNode directly inside it (not nested inside another
// optionalNode) evaluates to a non-empty string.
type optionalNode struct {
	children []node
}

// parseTemplate parses a rule's replacement text (already unescaped, see
// unescapeSpecials) into a sequence of nodes. parseNodes is called here
// with insideOptional false, so on success it always fully consumes s:
// its only early-return-with-leftover-text path requires insideOptional
// true (an unclosed '[').
func parseTemplate(s string) ([]node, error) {
	nodes, _, err := parseNodes(s, false)
	if err != nil {
		return nil, err
	}

	return attachPipes(nodes), nil
}

// parseNodes parses a run of nodes until either the input is exhausted or,
// when insideOptional is true, a closing ']' is found. It returns the
// parsed nodes and whatever input remains (starting with ']' if stopped
// early). One switch case per Soros token kind (spec 2.4): splitting the
// dispatch itself (as opposed to each case's body, already extracted into
// parseBracket/parseCall) would trade a single-page state machine for
// several that have to be read together to see the whole grammar.
//
//nolint:cyclop,funlen // see above
func parseNodes(s string, insideOptional bool) ([]node, string, error) {
	var nodes []node
	var lit strings.Builder

	flush := func() {
		if lit.Len() > 0 {
			nodes = append(nodes, literalNode{text: revealEscapes(lit.String())})
			lit.Reset()
		}
	}

	for len(s) > 0 {
		c := s[0]
		switch {
		case insideOptional && c == ']':
			flush()

			return nodes, s, nil
		case c == '[':
			flush()
			opt, rest, err := parseBracket(s[1:])
			if err != nil {
				return nil, "", err
			}
			nodes = append(nodes, opt)
			s = rest
		case c == '|':
			flush()
			nodes = append(nodes, pipeNode{})
			s = s[1:]
		case c == '\\' && len(s) > 1 && s[1] >= '0' && s[1] <= '9':
			flush()
			nodes = append(nodes, backrefNode{group: int(s[1] - '0')})
			s = s[2:]
		case c == '$' && len(s) > 1 && s[1] >= '1' && s[1] <= '9':
			flush()
			nodes = append(nodes, callNode{arg: []node{backrefNode{group: int(s[1] - '0')}}})
			s = s[2:]
		case c == '$' && len(s) > 1 && s[1] == '(':
			flush()
			call, rest, err := parseCall(s[2:])
			if err != nil {
				return nil, "", err
			}
			nodes = append(nodes, call)
			s = rest
		default:
			lit.WriteByte(c)
			s = s[1:]
		}
	}
	flush()

	return nodes, "", nil
}

// parseBracket parses the content of a "[...]" conditional block whose
// opening '[' was already consumed, and returns the resulting node plus
// whatever text follows its closing ']'.
func parseBracket(s string) (node, string, error) {
	children, rest, err := parseNodes(s, true)
	if err != nil {
		return nil, "", err
	}
	if !strings.HasPrefix(rest, "]") {
		return nil, "", fmt.Errorf("%w '[' in template", errUnterminated)
	}

	return optionalNode{children: children}, rest[1:], nil
}

// parseCall parses a full-form "$(...)" call whose leading "$(" was
// already consumed, and returns the resulting node plus whatever text
// follows its closing ')'.
func parseCall(s string) (node, string, error) {
	inner, rest, err := scanParens(s)
	if err != nil {
		return nil, "", err
	}
	argNodes, err := parseCallArg(inner)
	if err != nil {
		return nil, "", err
	}

	return callNode{arg: argNodes}, rest, nil
}

// parseCallArg parses the argument text of a "$(...)" call. It reuses
// parseNodes (insideOptional false, so it always fully consumes s on
// success - see parseTemplate), but a call argument may itself contain
// balanced, unescaped parentheses (from nested "$(...)" calls already
// scanned out by scanParens), so this is just a thin wrapper kept
// separate for clarity/documentation.
func parseCallArg(s string) ([]node, error) {
	nodes, _, err := parseNodes(s, false)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

// scanParens finds the matching ')' for a "$(" that was already
// consumed, honoring nested parentheses that belong to inner "$(...)"
// calls. It returns the text between the parentheses and whatever
// follows the closing ')'.
func scanParens(s string) (string, string, error) {
	depth := 1
	for i := range len(s) {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return s[:i], s[i+1:], nil
			}
		}
	}

	return "", "", fmt.Errorf("%w '$(' in template", errUnterminated)
}

// pipeNode is a boundary-modifier marker. It never produces output; it
// only affects the begin/end boundary flags of an adjacent callNode (see
// attachPipes and the evaluator).
type pipeNode struct{}

// attachPipes folds pipeNode markers into the pipeBefore/pipeAfter flags
// of neighboring callNodes and removes them from the node list, since a
// pipe by itself never renders anything. A single pipe between two calls
// applies to both (spec 2.4.4: "$1|$2 is equivalent form of $1||$2").
func attachPipes(nodes []node) []node {
	out := make([]node, 0, len(nodes))
	for i, n := range nodes {
		if _, ok := n.(pipeNode); ok {
			continue
		}
		call, ok := n.(callNode)
		if !ok {
			out = append(out, n)

			continue
		}
		if i > 0 {
			if _, ok := nodes[i-1].(pipeNode); ok {
				call.pipeBefore = true
			}
		}
		if i+1 < len(nodes) {
			if _, ok := nodes[i+1].(pipeNode); ok {
				call.pipeAfter = true
			}
		}
		call.arg = attachPipes(call.arg)
		out = append(out, call)
	}
	// Recurse into optional blocks too.
	for i, n := range out {
		if opt, ok := n.(optionalNode); ok {
			opt.children = attachPipes(opt.children)
			out[i] = opt
		}
	}

	return out
}
