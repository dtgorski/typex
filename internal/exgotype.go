// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 03/2020
// Partly based on types/typestring.go · Copyright 2013 The Go Authors

package internal

import (
	"bytes"
	"fmt"
	"go/token"
	"go/types"
	"io"
	"strings"
)

type (
	// ExGoType ...
	ExGoType struct {
		Renderer     Renderer
		IncludeUnexp bool

		indent int
		qualif func(p *types.Package) string
	}
)

// Export ...
func (t *ExGoType) Export(w io.Writer, m TypeMap) error {
	if t.postInit(); len(m) == 0 {
		return nil
	}
	views := make(ViewMap)
	for path, typ := range m {
		name, i := path, strings.LastIndex(path, ".")
		if i > -1 && i < len(path)-1 {
			name = path[i+1:]
			path = path[:i] + "/" + name
		}
		body := t.stringify(typ.Underlying())
		views[path] = name + " " + body
	}
	return t.Renderer.Render(w, views)
}

func (t *ExGoType) postInit() {
	if t.Renderer == nil {
		t.Renderer = &TreeView{}
	}
	t.qualif = func(p *types.Package) string {
		parts := strings.Split(p.Path(), "/")
		return parts[len(parts)-1]
	}
	t.indent = 0
}

func (t *ExGoType) stringify(typ types.Type) string {
	buf := bytes.Buffer{}
	t.writeTyp(&buf, typ, make([]types.Type, 0, 8))
	return buf.String()
}

func (t *ExGoType) writeTyp(buf *bytes.Buffer, typ types.Type, seen []types.Type) {
	for _, tt := range seen {
		if tt == typ {
			t.writeStr(buf, "%T", typ)
			return
		}
	}
	seen = append(seen, typ)

	switch tt := typ.(type) {
	case nil:
		t.writeStr(buf, "<nil>")

	case *types.Basic:
		if tt.Kind() == types.UnsafePointer {
			t.writeStr(buf, "unsafe.")
		}
		t.writeStr(buf, tt.Name())

	case *types.Array:
		t.writeStr(buf, "[%d]", tt.Len())
		t.writeTyp(buf, tt.Elem(), seen)

	case *types.Slice:
		t.writeStr(buf, "[]")
		t.writeTyp(buf, tt.Elem(), seen)

	case *types.Struct:
		t.indent++
		t.writeStr(buf, "struct {")

		void := true
		for i, n := 0, tt.NumFields(); i < n; i++ {
			if !t.exportable(tt.Field(i).Name()) {
				continue
			}
			named, ok := tt.Field(i).Type().(*types.Named)
			if ok && !t.exportable(named.String()) {
				continue
			}

			t.writeStr(buf, "\n")
			t.writePad(buf)

			if !tt.Field(i).Embedded() {
				t.writeStr(buf, tt.Field(i).Name())
				t.writeStr(buf, " ")
			}
			t.writeTyp(buf, tt.Field(i).Type(), seen)

			if tag := tt.Tag(i); tag != "" {
				t.writeStr(buf, "\t\t`%s`", tag)
			}
			void = false
		}
		if t.indent--; !void {
			t.writeStr(buf, "\n")
			t.writePad(buf)
		}
		t.writeStr(buf, "}")

	case *types.Pointer:
		t.writeStr(buf, "*")
		t.writeTyp(buf, tt.Elem(), seen)

	case *types.Tuple:
		t.writeTup(buf, tt, seen, false)

	case *types.Signature:
		t.writeStr(buf, "func")
		t.writeSig(buf, tt, seen)

	case *types.Interface:
		t.indent++
		t.writeStr(buf, "interface {")

		void := true
		for i, n := 0, tt.NumEmbeddeds(); i < n; i++ {
			if !t.exportable(tt.EmbeddedType(i).String()) {
				continue
			}
			t.writeStr(buf, "\n")
			t.writePad(buf)
			t.writeTyp(buf, tt.EmbeddedType(i), seen)
			void = false
		}
		for i, n := 0, tt.NumMethods(); i < n; i++ {
			if !t.exportable(tt.Method(i).Name()) {
				continue
			}
			t.writeStr(buf, "\n")
			t.writePad(buf)
			t.writeStr(buf, tt.Method(i).Name())
			t.writeSig(buf, tt.Method(i).Type().(*types.Signature), seen)
			void = false
		}
		if t.indent--; !void {
			t.writeStr(buf, "\n")
			t.writePad(buf)
		}
		t.writeStr(buf, "}")

	case *types.Map:
		t.writeStr(buf, "map[")
		t.writeTyp(buf, tt.Key(), seen)
		t.writeStr(buf, "]")
		t.writeTyp(buf, tt.Elem(), seen)

	case *types.Chan:
		s, par := "", false
		switch tt.Dir() {
		case types.SendRecv:
			s = "chan "
			ch, _ := tt.Elem().(*types.Chan)
			par = ch != nil && ch.Dir() == types.RecvOnly
		case types.SendOnly:
			s = "chan<- "
		case types.RecvOnly:
			s = "<-chan "
		}
		t.writeStr(buf, s)
		if par {
			t.writeStr(buf, "(")
		}
		t.writeTyp(buf, tt.Elem(), seen)
		if par {
			t.writeStr(buf, ")")
		}

	case *types.Named:
		s := "<Named w/o object>"
		if obj := tt.Obj(); obj != nil {
			t.writePkg(buf, obj.Pkg())
			s = obj.Name()
		}
		t.writeStr(buf, s)

	default:
		t.writeStr(buf, tt.String())
	}
}

func (t *ExGoType) writeTup(buf *bytes.Buffer, tup *types.Tuple, seen []types.Type, variadic bool) {
	if tup == nil {
		t.writeStr(buf, "()")
		return
	}
	t.writeStr(buf, "(")
	for i, n := 0, tup.Len(); i < n; i++ {
		if i > 0 {
			t.writeStr(buf, ", ")
		}
		v := tup.At(i)
		if v.Name() != "" {
			t.writeStr(buf, v.Name())
			t.writeStr(buf, " ")
		}
		typ := v.Type()
		if variadic && i == n-1 {
			if s, ok := typ.(*types.Slice); ok {
				t.writeStr(buf, "...")
				typ = s.Elem()
			} else {
				t.writeTyp(buf, typ, seen)
				t.writeStr(buf, "...")
				continue
			}
		}
		t.writeTyp(buf, typ, seen)
	}
	t.writeStr(buf, ")")
}

func (t *ExGoType) writeSig(buf *bytes.Buffer, sig *types.Signature, seen []types.Type) {
	t.writeTup(buf, sig.Params(), seen, sig.Variadic())
	n := sig.Results().Len()
	if n == 0 {
		return
	}
	t.writeStr(buf, " ")
	if n == 1 {
		if sig.Results().At(0).Name() == "" {
			t.writeTyp(buf, sig.Results().At(0).Type(), seen)
			return
		}
	}
	t.writeTup(buf, sig.Results(), seen, false)
}

func (t *ExGoType) writePkg(buf *bytes.Buffer, pkg *types.Package) {
	if pkg == nil {
		return
	}
	s := pkg.Path()
	if t.qualif != nil {
		s = t.qualif(pkg)
	}
	if s != "" {
		t.writeStr(buf, s)
		t.writeStr(buf, ".")
	}
}

func (t *ExGoType) writePad(buf *bytes.Buffer) {
	t.writeStr(buf, strings.Repeat("    ", t.indent))
}

func (*ExGoType) writeStr(buf *bytes.Buffer, f string, a ...interface{}) {
	_, _ = fmt.Fprintf(buf, f, a...)
}

func (t *ExGoType) exportable(s string) bool {
	if t.IncludeUnexp {
		return true
	}
	path := strings.ReplaceAll(s, ".", "/")
	part := strings.Split(path, "/")
	return token.IsExported(part[len(part)-1])
}
