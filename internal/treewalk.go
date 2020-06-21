// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 06/2020

package internal

import (
	"io/ioutil"
	"strings"

	"golang.org/x/tools/godoc/vfs/mapfs"
)

type (
	// TreeWalk knows about Layout.
	TreeWalk struct {
		Layout
		indent int
	}

	// The Layout interface shapes a type visitor.
	Layout interface {
		Enter(path string, last bool)
		Print(line string, first, last bool)
		Leave(path string, last bool)
	}

	// PathMap carries the mapping between the fully
	// qualified import paths and their type projections.
	PathMap map[string]string
)

// Walk traverses a PathMap according to the paths.
func (t *TreeWalk) Walk(m PathMap) error {
	if len(m) == 0 {
		return nil
	}
	mapFs := mapfs.New(m)
	t.indent = -1

	var walk func(dir string) error
	walk = func(dir string) error {
		t.indent++
		defer func() { t.indent-- }()

		finfo, err := mapFs.ReadDir(dir)
		if err != nil {
			return err
		}
		for i, n := 0, len(finfo); i < n; i++ {
			path := dir + finfo[i].Name()

			if finfo[i].IsDir() {
				t.Enter(path, i == n-1)
				err = walk(path + "/")
				if err != nil {
					return err
				}
				t.Leave(path, i == n-1)
				continue
			}

			file, err := mapFs.Open(path)
			if err != nil {
				return err
			}
			data, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}
			lines := strings.Split(string(data), "\n")
			for j, m := 0, len(lines); j < m; j++ {
				t.Print(lines[j], j == 0, i == n-1)
			}
		}
		return nil
	}
	return walk("/")
}
