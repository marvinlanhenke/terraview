package tree

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// matcher holds a compiled search query. It supports two modes: plain
// case-insensitive substring matching and regex matching (when the query is
// wrapped in forward slashes, e.g. /pattern/).
type matcher struct {
	raw     string
	literal string
	re      *regexp.Regexp
}

// newMatcher parses query and returns a ready-to-use matcher. If the query is
// a valid regex literal (surrounded by '/'), the pattern is compiled
// case-insensitively; invalid regexes fall back to plain substring matching.
func newMatcher(query string) matcher {
	raw := strings.TrimSpace(query)

	m := matcher{
		raw:     raw,
		literal: strings.ToLower(raw),
	}

	if !isRegexQuery(raw) {
		return m
	}

	re, err := regexp.Compile("(?i)" + unwrapRegex(raw))
	if err == nil {
		m.re = re
	}

	return m
}

// active reports whether the matcher has a non-empty query.
func (m matcher) active() bool {
	return m.raw != ""
}

// matchNode reports whether the node's id, label, action, or pre-computed
// search payload satisfies the matcher.
func (m matcher) matchNode(n *Node) bool {
	if n == nil {
		return false
	}

	if !m.active() {
		return true
	}

	return m.matchString(n.Id) ||
		m.matchString(n.Label) ||
		m.matchString(string(n.Action)) ||
		m.matchString(n.searchPayload)
}

// matchString reports whether v satisfies the matcher. When no query is active
// every string is considered a match.
func (m matcher) matchString(v string) bool {
	if !m.active() {
		return true
	}

	if m.re != nil {
		return m.re.MatchString(v)
	}

	return strings.Contains(strings.ToLower(v), m.literal)
}

// isRegexQuery reports whether query is a regex literal, i.e. a non-empty
// string of the form /pattern/. Single-character inputs like "/" and "//"
// are treated as plain text.
func isRegexQuery(query string) bool {
	if len(query) < 3 {
		// Require at least /x/, so "/" and "//" stay plain text.
		return false
	}

	if query[0] != '/' || query[len(query)-1] != '/' {
		return false
	}

	return strings.TrimSpace(unwrapRegex(query)) != ""
}

// unwrapRegex strips the surrounding '/' delimiters from a regex literal and
// returns the inner pattern. It returns an empty string for inputs shorter
// than three characters.
func unwrapRegex(query string) string {
	if len(query) < 3 {
		return ""
	}
	return query[1 : len(query)-1]
}

// prepareSearchPayloads walks the subtree rooted at n and populates the
// searchPayload field on every node. This is called once when a new root is
// set so that subsequent searches can match against the pre-computed strings
// without re-serialising on every keystroke.
func prepareSearchPayloads(n *Node) {
	if n == nil {
		return
	}

	n.searchPayload = searchablePayload(n.Payload)

	for _, child := range n.Children {
		prepareSearchPayloads(child)
	}
}

// searchablePayload converts an arbitrary node payload into a searchable
// string. Strings are returned as-is; all other types are JSON-marshalled so
// that nested keys and values become searchable text.
func searchablePayload(v any) string {
	if v == nil {
		return "null"
	}

	switch t := v.(type) {
	case string:
		return t
	default:
		b, err := json.Marshal(t)
		if err != nil {
			return fmt.Sprintf("%v", t)
		}
		return "\n" + string(b)
	}
}
