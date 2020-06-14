// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 06/2020

package p3

type (
	U uint

	Y struct {
		M interface{}
		N map[[10]U]interface{}
		O map[int]map[int]interface{}
		P map[U]func()
		Q [10]**interface{}
		R []*map[U]*interface{}
		S map[*Z]*Z
	}

	Z interface{}
)
