// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 06/2020

package p1

import (
	"time"

	"github.com/dtgorski/typex/internal/testdata/p2"
	"github.com/dtgorski/typex/internal/testdata/p2/p3"
)

type (
	A time.Duration
	B []time.Duration

	G struct {
		D map[string]time.Duration
		E map[string]B
		U
		Y <-chan chan<- D
		Z z
	}

	D struct {
		e bool
		F bool `json:"-,opt1,opt2"`
		G `embedded:"" json:"-"`
		H []map[int]struct {
			I []D
			J [][]*p2.T
			K struct {
				L []struct {
					m map[string]map[p2.I]**T
					N map[p3.U][10]string
					O map[int]*[]struct {
						P func() interface{}
					}
					Q map[int][]*func(interface{}) (<-chan chan<- bool, int)
				}
			}
		}
		R map[*int64]**W
		S bool `json:"other,omitempty"`
	}

	T struct {
		D
		W
		U **Y
		V *A
		X
		p3.Y
		Z U
	}

	U []struct {
		V *U
	}

	W map[int64]time.Time

	X interface {
		E() error
		e()
	}
	Y interface {
		X
		P() uintptr
	}
	z interface{}
)
