// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 06/2020

package ts

import (
	"fmt"
	"io"
	"strings"

	typex "github.com/dtgorski/typex/internal"
)

type (
	moduleLayout struct {
		writer io.Writer
		indent int
	}
)

// NewModuleLayout ...
func NewModuleLayout(w io.Writer) typex.Layout {
	return &moduleLayout{writer: w}
}

func (m *moduleLayout) Enter(path string, _ bool) {
	i := strings.LastIndex(path, "/")
	p := strings.Repeat(" ", m.indent<<2)
	m.write(m.writer, "%sexport module %s {\n", p, path[i+1:])
	m.indent++
}

func (m *moduleLayout) Print(line string, _, _ bool) {
	p := strings.Repeat(" ", m.indent<<2)
	m.write(m.writer, "%s%s\n", p, line)
}

func (m *moduleLayout) Leave(_ string, _ bool) {
	m.indent--
	p := strings.Repeat(" ", m.indent<<2)
	m.write(m.writer, "%s}\n", p)
}

func (moduleLayout) write(w io.Writer, f string, a ...interface{}) {
	_, _ = fmt.Fprintf(w, f, a...)
}
