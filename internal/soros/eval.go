package soros

import "strings"

// Run evaluates the program's implicit "$" function against input,
// treating input as the whole (global) call: it is eligible for both
// leading- and trailing-boundary-anchored rules (spec 2.3.2).
func (p *Program) Run(input string) string {
	return p.run(input, true, true)
}

// run is the recursive evaluator behind Run. begin/end record whether
// this call is still attached to the leading/trailing boundary of the
// original, outermost input (spec 2.3.2, 2.4.4): a rule whose pattern had
// an explicit '^' only applies when begin is true, and likewise '$'/end.
func (p *Program) run(input string, begin, end bool) string {
	for _, r := range p.rules {
		if (r.begin && !begin) || (r.end && !end) {
			continue
		}
		groups := r.re.FindStringSubmatch(input)
		if groups == nil {
			continue
		}

		return p.render(r.tmpl, groups, begin, end)
	}

	return ""
}

// render evaluates a parsed replacement template against the capture
// groups of the rule that matched, recursing into the program for every
// "$(...)"/"$N" call it contains.
func (p *Program) render(tmpl []node, groups []string, begin, end bool) string {
	var b strings.Builder
	for i, n := range tmpl {
		atStart := i == 0
		atEnd := i == len(tmpl)-1
		b.WriteString(p.renderNode(n, groups, begin, end, atStart, atEnd))
	}

	return b.String()
}

// renderNode evaluates a single template node. atStart/atEnd tell a call
// node whether it sits at the very start/end of the template it belongs
// to, which (absent an explicit '|' boundary modifier) is how it inherits
// the enclosing begin/end boundary status (spec 2.4.4).
func (p *Program) renderNode(n node, groups []string, begin, end bool, atStart, atEnd bool) string {
	switch elem := n.(type) {
	case literalNode:
		return elem.text
	case backrefNode:
		return group(groups, elem.group)
	case callNode:
		arg := p.render(elem.arg, groups, begin, end)
		b := elem.pipeBefore || (atStart && begin)
		e := elem.pipeAfter || (atEnd && end)

		return p.run(arg, b, e)
	case optionalNode:
		return p.renderOptional(elem, groups, begin, end, atStart, atEnd)
	default:
		return ""
	}
}

// renderOptional implements conditional text (spec 2.4.5): the rendered
// content of a "[...]" block is kept only if every direct-child call it
// contains evaluated to a non-empty string. outerAtStart/outerAtEnd carry
// through whether the bracket itself sits at the start/end of the
// template it is embedded in, since brackets are inline and do not start
// a new template of their own.
func (p *Program) renderOptional(
	v optionalNode, groups []string, begin, end bool, outerAtStart, outerAtEnd bool,
) string {
	var b strings.Builder
	nonEmpty := true
	for i, child := range v.children {
		atStart := i == 0 && outerAtStart
		atEnd := i == len(v.children)-1 && outerAtEnd
		out := p.renderNode(child, groups, begin, end, atStart, atEnd)
		if _, isCall := child.(callNode); isCall && out == "" {
			nonEmpty = false
		}
		b.WriteString(out)
	}
	if !nonEmpty {
		return ""
	}

	return b.String()
}

// group returns capture group i of groups (0 is the whole match), or ""
// if i is out of range.
func group(groups []string, i int) string {
	if i < 0 || i >= len(groups) {
		return ""
	}

	return groups[i]
}
