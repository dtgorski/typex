// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 06/2020

package p2

import (
	"github.com/dtgorski/typex/internal/testdata/p2/p3"
)

type (
	I int
	F func(*I, ...*p3.U) (T, error)
	S struct{ Fn **F }

	T struct {
		ArrayType     [10]string
		BoolType      bool
		IntType       I
		Int8Type      int8
		Int16Type     int16
		Int32Type     int32
		Int64Type     int64
		UintType      p3.U
		Uint8Type     uint8
		Uint16Type    uint16
		Uint32Type    uint32
		Uint64Type    uint64
		ByteType      byte
		RuneType      rune
		UintPtrType   uintptr
		Float32Type   float32
		Float64Type   float64
		InterfaceType interface {
			Foo() (int, error)
		}
		FuncType       **func(x, y int, z ...int) (int, error)
		FuncType_      func(x, y, z int) chan *struct{}
		ChanType       <-chan *bool
		Complex64Type  complex64
		Complex128Type complex128
		MapType        map[*int]string
		MapType_       map[string]chan *struct{}
		StringType     string
		StructType     struct{}
		SliceType      []string

		FuncStruct chan<- *S `json:"funcStruct"`
		Types      *T        `tags:"types"`
	}
)
