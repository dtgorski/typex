// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 03/2020

package internal

import (
	"regexp"
)

// CreateFilterFunc ...
func CreateFilterFunc(list []string) FilterFunc {
	filter := make([]*regexp.Regexp, 0)

	for _, expr := range list {
		re0, err := regexp.Compile(expr)
		if err != nil {
			esc := regexp.QuoteMeta(expr)
			re1, err := regexp.Compile(esc)
			if err != nil {
				continue
			}
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
