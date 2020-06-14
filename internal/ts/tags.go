// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 06/2020
// Partly based on reflect/type.go · Copyright 2009 The Go Authors
// Partly based on encoding/json/tags.go · Copyright 2011 The Go Authors

package ts

import (
	"strconv"
	"strings"
)

type (
	// A StructTag is the tag string in a struct field.
	//
	// By convention, tag strings are a concatenation of optionally
	// space-separated key:"value" pairs. Each key is a non-empty
	// string consisting of non-control characters other than space
	// (U+0020 ' '), quote (U+0022 '"'), and colon (U+003A ':').
	// Each value is quoted using U+0022 '"' characters and Go string
	// literal syntax.
	StructTag string

	// TagOptions is the string following a comma in a struct field's
	// tag, or the empty string. It does not include the leading comma.
	TagOptions string
)

// Get returns the values associated with key and options in the tag string.
func (t StructTag) Get(key string) (string, TagOptions) {
	tag := t
	for tag != "" {
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		value := string(tag[:i+1])
		tag = tag[i+1:]

		if key == name {
			value, err := strconv.Unquote(value)
			if err != nil {
				break
			}
			if i := strings.Index(value, ","); i != -1 {
				return value[:i], TagOptions(value[i+1:])
			}
			return value, ""
		}
	}
	return "", ""
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded
// by a string boundary or commas.
func (o TagOptions) Contains(name string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == name {
			return true
		}
		s = next
	}
	return false
}
