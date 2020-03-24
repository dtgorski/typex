// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 03/2020

package p1

import (
	"time"

	"github.com/dtgorski/typex/internal/testdata/p2"
	"github.com/dtgorski/typex/internal/testdata/p2/p3"
)

type (
	T struct {
		notExported bool
		X           `embedded:"",json:"-"`
		UnExported  bool `json:"-"`
		AnyTagName  bool `json:"otherTagName"`

		Ä struct {
			B []map[int]struct {
				C   *T
				DDD [][]*p2.T
				E   struct {
					F []struct {
						G   map[string]map[p2.I]**S
						HHH map[p3.U][10]string
					}
				}
			}
		}
	}

	S struct {
		T
		x
		X
		XX *X
		I
	}

	X struct {
		Y map[*int64]**x
		Z map[*X]time.Duration
	}

	x map[int64]time.Time

	i interface {
		g()
	}
	I interface {
		i
		F()
		f()
	}

	j struct {
		*J
		uintptr
	}
	J struct {
		j
	}
)
