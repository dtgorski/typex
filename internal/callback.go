// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 06/2020

package internal

import (
	"regexp"
	"strings"
)

// CreatePathFilterFunc gets a list of text patterns which are converted
// to regular expressions. The result of this function is a PathFilterFunc
// which can be used as a filter for matching patterns against strings.
func CreatePathFilterFunc(include, exclude []string) PathFilterFunc {
	included := make([]*regexp.Regexp, 0)
	excluded := make([]*regexp.Regexp, 0)

	for _, expr := range include {
		re0, err := regexp.Compile(expr)
		if err != nil {
			esc := regexp.QuoteMeta(expr)
			re1, _ := regexp.Compile(esc)
			re0 = re1
		}
		included = append(included, re0)
	}

	for _, expr := range exclude {
		re0, err := regexp.Compile(expr)
		if err != nil {
			esc := regexp.QuoteMeta(expr)
			re1, _ := regexp.Compile(esc)
			re0 = re1
		}
		excluded = append(excluded, re0)
	}

	return func(name string) bool {
		incl, excl := false, false

		for _, f := range included {
			if incl = f.MatchString(name); incl {
				break
			}
		}
		if incl {
			for _, f := range excluded {
				if excl = f.MatchString(name); excl {
					break
				}
			}
		}
		return incl && !excl
	}
}

// CreatePathReplaceFunc returns the default path and name
// relocation function. The result is basically sanitized
// but bad input will lead to bad output.
func CreatePathReplaceFunc(list []string) PathReplaceFunc {
	type replace struct {
		old *regexp.Regexp
		new string
	}
	pairs := make([]replace, 0)

	for _, expr := range list {
		parts := strings.Split(expr, ":")
		if len(parts) != 2 {
			continue
		}
		p0 := strings.TrimSpace(parts[0])
		p1 := strings.TrimSpace(parts[1])
		if p0 == "" {
			continue
		}
		re0, err := regexp.Compile(p0)
		if err != nil {
			esc := regexp.QuoteMeta(p0)
			re1, _ := regexp.Compile(esc)
			re0 = re1
		}
		pairs = append(pairs, replace{re0, p1})
	}

	re1 := regexp.MustCompile(`\./+`)
	re2 := regexp.MustCompile(`/\.+`)
	re3 := regexp.MustCompile(`/+`)
	re4 := regexp.MustCompile(`\.+`)

	return func(s string) string {
		for i := 0; i < len(pairs); i++ {
			s = pairs[i].old.ReplaceAllString(s, pairs[i].new)
		}
		// sanitize
		s = re1.ReplaceAllString(s, ".")
		s = re2.ReplaceAllString(s, "/")
		s = re3.ReplaceAllString(s, "/")
		s = re4.ReplaceAllString(s, ".")

		s = strings.Trim(s, "/.")
		s = strings.ReplaceAll(s, "-", "_")

		return s
	}
}
