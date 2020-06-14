// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 06/2020
// Partly based on types/typestring.go · Copyright 2013 The Go Authors

package g0

import (
	"bytes"
	"fmt"
	"go/token"
	"go/types"
	"io"
	"strings"

	typex "github.com/dtgorski/typex/internal"
)

type (
	// TypeRender ...
	TypeRender struct {
		PathReplaceFunc   typex.PathReplaceFunc
		IncludeUnexported bool

		indent int
	}

	context struct {
		writer   io.Writer
		seenType []types.Type
	}
)

// Render ...
func (r *TypeRender) Render(m typex.TypeMap) typex.PathMap {
	r.indent = 0
	pathMap := make(typex.PathMap)

	for p, t := range m {
		path, name := r.pathAndName(p)
		buf, typ := bytes.Buffer{}, t.Underlying()
		ctx := context{&buf, make([]types.Type, 0)}

		r.writeType(ctx, typ)
		pathMap[path] = name + " " + buf.String()
	}
	return pathMap
}

func (r *TypeRender) writeType(ctx context, t types.Type) {
	for _, tt := range ctx.seenType {
		if tt == t {
			return
		}
	}
	ctx.seenType = append(ctx.seenType, t)

	switch tt := t.(type) {
	case *types.Basic:
		r.write(ctx, tt.Name())

	case *types.Array:
		r.write(ctx, "[%d]", tt.Len())
		r.writeType(ctx, tt.Elem())

	case *types.Chan:
		r.writeChan(ctx, tt)

	case *types.Interface:
		r.writeInterface(ctx, tt)

	case *types.Map:
		r.writeMap(ctx, tt)

	case *types.Named:
		r.writeNamed(ctx, tt)

	case *types.Pointer:
		r.write(ctx, "*")
		r.writeType(ctx, tt.Elem())

	case *types.Signature:
		r.write(ctx, "func")
		r.writeSignature(ctx, tt)

	case *types.Slice:
		r.write(ctx, "[]")
		r.writeType(ctx, tt.Elem())

	case *types.Struct:
		r.writeStruct(ctx, tt)
	}
}

func (r *TypeRender) writeChan(ctx context, t *types.Chan) {
	s := ""
	switch t.Dir() {
	case types.SendRecv:
		s = "chan "
	case types.SendOnly:
		s = "chan<- "
	case types.RecvOnly:
		s = "<-chan "
	}
	r.write(ctx, s)
	r.writeType(ctx, t.Elem())
}

func (r *TypeRender) writeInterface(ctx context, t *types.Interface) {
	r.indent++
	r.write(ctx, "interface {")

	void := true
	for i, n := 0, t.NumEmbeddeds(); i < n; i++ {
		if !r.isExported(t.EmbeddedType(i).String()) {
			continue
		}
		r.write(ctx, "\n")
		r.writePadding(ctx)
		r.writeType(ctx, t.EmbeddedType(i))
		void = false
	}
	for i, n := 0, t.NumMethods(); i < n; i++ {
		if !r.isExported(t.Method(i).Name()) {
			continue
		}
		r.write(ctx, "\n")
		r.writePadding(ctx)
		r.write(ctx, t.Method(i).Name())
		r.writeSignature(ctx, t.Method(i).Type().(*types.Signature))
		void = false
	}
	if r.indent--; !void {
		r.write(ctx, "\n")
		r.writePadding(ctx)
	}
	r.write(ctx, "}")
}

func (r *TypeRender) writeMap(ctx context, t *types.Map) {
	r.write(ctx, "map[")
	r.writeType(ctx, t.Key())
	r.write(ctx, "]")
	r.writeType(ctx, t.Elem())
}

func (r *TypeRender) writeNamed(ctx context, t *types.Named) {
	p := ""
	if obj := t.Obj(); obj != nil {
		if obj.Pkg() != nil {
			p = obj.Pkg().Path() + "."
		}
		p += obj.Name()
		p = r.replacePath(p)
		i := strings.LastIndex(p, "/")
		p = p[i+1:]
	}
	r.write(ctx, p)
}

func (r *TypeRender) writeStruct(ctx context, t *types.Struct) {
	void := true
	r.indent++
	r.write(ctx, "struct {")

	for i, n := 0, t.NumFields(); i < n; i++ {
		if !r.isExported(t.Field(i).Name()) {
			continue
		}
		if tt, ok := t.Field(i).Type().(*types.Named); ok {
			if !r.isExported(tt.String()) {
				continue
			}
		}
		r.write(ctx, "\n")
		r.writePadding(ctx)

		if !t.Field(i).Embedded() {
			r.write(ctx, t.Field(i).Name())
			r.write(ctx, " ")
		}
		r.writeType(ctx, t.Field(i).Type())

		if tag := t.Tag(i); tag != "" {
			r.write(ctx, "\t\t`%s`", tag)
		}
		void = false
	}
	if r.indent--; !void {
		r.write(ctx, "\n")
		r.writePadding(ctx)
	}
	r.write(ctx, "}")
}

func (r *TypeRender) writeTuple(ctx context, t *types.Tuple, variadic bool) {
	if t == nil {
		r.write(ctx, "()")
		return
	}
	r.write(ctx, "(")
	for i, n := 0, t.Len(); i < n; i++ {
		if i > 0 {
			r.write(ctx, ", ")
		}
		tt := t.At(i)
		if tt.Name() != "" {
			r.write(ctx, tt.Name())
			r.write(ctx, " ")
		}
		typ := tt.Type()
		if variadic && i == n-1 {
			if s, ok := typ.(*types.Slice); ok {
				r.write(ctx, "...")
				typ = s.Elem()
			}
		}
		r.writeType(ctx, typ)
	}
	r.write(ctx, ")")
}

func (r *TypeRender) writeSignature(ctx context, t *types.Signature) {
	r.writeTuple(ctx, t.Params(), t.Variadic())
	n := t.Results().Len()
	if n == 0 {
		return
	}
	r.write(ctx, " ")
	if n == 1 {
		if t.Results().At(0).Name() == "" {
			r.writeType(ctx, t.Results().At(0).Type())
			return
		}
	}
	r.writeTuple(ctx, t.Results(), false)
}

func (r *TypeRender) writePadding(ctx context) {
	r.write(ctx, strings.Repeat("    ", r.indent))
}

func (TypeRender) write(ctx context, f string, a ...interface{}) {
	_, _ = fmt.Fprintf(ctx.writer, f, a...)
}

func (r *TypeRender) pathAndName(s string) (p, n string) {
	p = r.replacePath(s)
	n = p
	i := strings.LastIndex(p, ".")
	if i > -1 && i < len(p)-1 {
		n = p[i+1:]
		p = p[:i] + "/" + n
	}
	return p, n
}

func (r *TypeRender) replacePath(s string) string {
	if r.PathReplaceFunc != nil {
		return r.PathReplaceFunc(s)
	}
	return s
}

func (r *TypeRender) isExported(s string) bool {
	if r.IncludeUnexported {
		return true
	}
	n := strings.ReplaceAll(s, ".", "/")
	i := strings.LastIndex(n, "/")
	return token.IsExported(n[i+1:])
}
