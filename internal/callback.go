// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 06/2020

package internal

import (
	"regexp"
	"strings"
)

// CreatePathFilterFunc gets a list of text patterns which are converted
// to regular expressions. The result of this function is a PathFilterFunc
// which can be used as a filter for matching patterns against strings.
func CreatePathFilterFunc(list []string) PathFilterFunc {
	filter := make([]*regexp.Regexp, 0)

	for _, expr := range list {
		re0, err := regexp.Compile(expr)
		if err != nil {
			esc := regexp.QuoteMeta(expr)
			re1, _ := regexp.Compile(esc)
			re0 = re1
		}
		filter = append(filter, re0)
	}

	return func(name string) bool {
		for _, f := range filter {
			if f.MatchString(name) {
				return true
			}
		}
		return false
	}
}

// CreatePathReplaceFunc ...
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
