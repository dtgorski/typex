// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 03/2020

package internal

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"golang.org/x/tools/godoc/vfs/mapfs"
)

type (
	// TreeView ...
	TreeView struct{}

	treeView struct {
		writer  io.Writer
		viewMap ViewMap
		edgeMap map[int]bool
		indent  int
	}

	// ViewMap ...
	ViewMap map[string]string

	// Renderer ...
	Renderer interface {
		Render(w io.Writer, m ViewMap) error
	}
)

// Render ...
func (*TreeView) Render(w io.Writer, m ViewMap) error {
	t := &treeView{
		writer:  w,
		viewMap: m,
		edgeMap: make(map[int]bool),
	}
	return t.render()
}

func (t *treeView) render() error {
	mapFs := mapfs.New(t.viewMap)
	t.indent = -1

	var walk func(dir string) error
	walk = func(dir string) error {
		t.indent++
		defer func() { t.indent-- }()

		fi, err := mapFs.ReadDir(dir)
		if err != nil {
			return err
		}
		for i, n := 0, len(fi); i < n; i++ {
			conn := t.edges(i, n)

			if fi[i].IsDir() {
				t.printf("%s%s\n", conn, fi[i].Name())
				err := walk(dir + fi[i].Name() + "/")
				if err != nil {
					return err
				}
				continue
			}

			file, err := mapFs.Open(dir + fi[i].Name())
			if err != nil {
				return err
			}
			view, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}

			lines := strings.Split(string(view), "\n")
			for j, m := 0, len(lines); j < m; j++ {
				edge := conn[:len(conn)-(3*3+1)]
				switch {
				case i == n-1 && j == 0:
					edge = conn
				case i == n-1 && j != 0:
					edge += "    "
				case i != n-1 && j == 0:
					edge += "├── "
				case i != n-1 && j != 0:
					edge += "│   "
				}
				t.printf("%s%s\n", edge, lines[j])
			}
		}
		return nil
	}
	return walk("/")
}

func (t *treeView) edges(i, n int) string {
	edge := ""
	for j, m := 0, t.indent; j < m; j++ {
		if t.edgeMap[j] {
			edge += "│   "
		} else {
			edge += "    "
		}
	}
	if t.edgeMap[t.indent] = i != n-1; t.edgeMap[t.indent] {
		edge += "├── "
	} else {
		edge += "└── "
	}
	return edge
}

func (t *treeView) printf(f string, a ...interface{}) {
	_, _ = fmt.Fprintf(t.writer, f, a...)
}
