// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 06/2020

package g0

import (
	"fmt"
	"io"
	"strings"

	typex "github.com/dtgorski/typex/internal"
)

type (
	treeLayout struct {
		writer  io.Writer
		edgeMap map[int]bool
		indent  int
	}
)

// NewTreeLayout implements the Layout interface.
func NewTreeLayout(w io.Writer) typex.Layout {
	return &treeLayout{
		writer:  w,
		edgeMap: make(map[int]bool),
	}
}

func (t *treeLayout) Enter(path string, last bool) {
	i := strings.LastIndex(path, "/")
	t.write(t.writer, "%s%s\n", t.edges(last), path[i+1:])
	t.indent++
}

func (t *treeLayout) Print(line string, first, last bool) {
	conn := []rune(t.edges(last))
	edge := string(conn[:len(conn)-4])
	switch {
	case last && first:
		edge = string(conn)
	case last && !first:
		edge += "    "
	case !last && first:
		edge += "├── "
	case !last && !first:
		edge += "│   "
	}
	t.write(t.writer, "%s%s\n", edge, line)
}

func (t *treeLayout) Leave(_ string, _ bool) {
	t.indent--
}

func (t *treeLayout) edges(last bool) string {
	edge := ""
	for j, k := 0, t.indent; j < k; j++ {
		if t.edgeMap[j] {
			edge += "│   "
		} else {
			edge += "    "
		}
	}
	if t.edgeMap[t.indent] = !last; t.edgeMap[t.indent] {
		edge += "├── "
	} else {
		edge += "└── "
	}
	return edge
}

func (treeLayout) write(w io.Writer, f string, a ...interface{}) {
	_, _ = fmt.Fprintf(w, f, a...)
}
