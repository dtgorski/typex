// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 06/2020
// Partly based on types/typestring.go · Copyright 2013 The Go Authors

package ts

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
	// TypeRender renders TypeScript types or classes.
	TypeRender struct {
		PathReplaceFunc   typex.PathReplaceFunc
		IncludeUnexported bool

		indent int
	}

	context struct {
		writer   io.Writer
		seenType []types.Type
		needCtor bool
		expoObjs bool
	}
)

// Render converts a TypeMap to a PathMap.
func (r *TypeRender) Render(m typex.TypeMap, exportObj bool) typex.PathMap {
	r.indent = 0

	exClass := false
	pathMap := make(typex.PathMap)

	for p, t := range m {
		if !r.isTranslatable(t.Underlying()) {
			continue
		}
		path, name := r.pathAndName(p)
		buf, typ := bytes.Buffer{}, t.Underlying()

		if exClass = exportObj; exClass {
			_, exClass = typ.(*types.Struct)
		}
		ctx := context{
			writer:   &buf,
			seenType: make([]types.Type, 0),
			needCtor: exClass,
			expoObjs: exClass,
		}
		r.writeType(ctx, typ)

		if exClass {
			pathMap[path] = "export class " + name + " " + buf.String()
		} else {
			pathMap[path] = "export type " + name + " = " + buf.String()
		}
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
		r.writeBasic(ctx, tt)

	case *types.Chan:
		r.write(ctx, "any")

	case *types.Array:
		r.writeType(ctx, tt.Elem())
		r.write(ctx, "[/* %d */]", tt.Len())

	case *types.Interface:
		r.write(ctx, "any")

	case *types.Map:
		r.writeMap(ctx, tt)

	case *types.Named:
		r.writeNamed(ctx, tt)

	case *types.Pointer:
		r.writeType(ctx, tt.Elem())

	case *types.Signature:
		r.write(ctx, "any")

	case *types.Slice:
		r.writeType(ctx, tt.Elem())
		r.write(ctx, "[]")

	case *types.Struct:
		r.writeStruct(ctx, tt)
	}
}

func (r *TypeRender) writeBasic(ctx context, t *types.Basic) {
	switch k := t.Info(); true {
	case (k & types.IsBoolean) != 0:
		r.write(ctx, "boolean")
	case (k & types.IsNumeric) != 0:
		r.write(ctx, "number")
	case (k & types.IsString) != 0:
		r.write(ctx, "string")
	default:
		r.write(ctx, "any")
	}
}

func (r *TypeRender) writeNamed(ctx context, t *types.Named) {
	p := "any"
	if r.isTranslatable(t.Underlying()) {
		if obj := t.Obj(); obj != nil {
			p = ""
			if obj.Pkg() != nil {
				p += obj.Pkg().Path() + "."
			}
			p += obj.Name()
			p = r.replacePath(p)
			p = strings.ReplaceAll(p, "/", ".")
		}
	}
	r.write(ctx, p)
}

func (r *TypeRender) writeMap(ctx context, t *types.Map) {
	r.write(ctx, "Record<")
	if r.isValidMapKey(t.Key()) {
		r.writeType(ctx, t.Key())
	} else {
		r.write(ctx, "symbol")
	}
	r.write(ctx, ", ")
	r.writeType(ctx, t.Elem())
	r.write(ctx, ">")
}

func (r *TypeRender) writeStruct(ctx context, t *types.Struct) {
	void := true
	r.indent++
	r.write(ctx, "{")

	ctor := ctx.needCtor && ctx.expoObjs
	if ctor {
		ctx.needCtor = false
		r.write(ctx, "\n")
		r.writePadding(ctx)
		r.write(ctx, "constructor(")
		r.indent++
	}

	for i, n := 0, t.NumFields(); i < n; i++ {
		tag, opt := (StructTag)(t.Tag(i)).Get("json")
		if tag == "-" {
			continue
		}
		fld := t.Field(i)
		typ := fld.Type()
		name := fld.Name()

		if !r.isExported(name) {
			continue
		}
		if tt, ok := t.Field(i).Type().(*types.Named); ok {
			if !r.isExported(tt.String()) {
				continue
			}
		}
		r.write(ctx, "\n")
		r.writePadding(ctx)
		if ctx.expoObjs {
			r.write(ctx, "readonly ")
		}
		if tag != "" {
			r.write(ctx, tag)
		} else {
			r.write(ctx, name)
		}
		if opt.Contains("omitempty") {
			r.write(ctx, "?")
		}
		r.write(ctx, ": ")
		r.writeType(ctx, typ)
		r.write(ctx, ",")
		void = false
	}
	if r.indent--; !void {
		r.write(ctx, "\n")
		r.writePadding(ctx)
	}
	if ctor {
		r.indent--
		r.write(ctx, ") {}\n")
		r.writePadding(ctx)
	}
	r.write(ctx, "}")
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

func (r *TypeRender) isValidMapKey(t types.Type) bool {
	switch tt := t.(type) {
	case *types.Basic:
		k := tt.Info()
		return (k&types.IsNumeric) != 0 || (k&types.IsString) != 0
	case *types.Map:
		return r.isValidMapKey(tt.Key())
	case *types.Named:
		return r.isValidMapKey(tt.Underlying())
	}
	return false
}

func (r *TypeRender) isTranslatable(t types.Type) bool {
	switch t.(type) {
	case *types.Chan, *types.Interface, *types.Signature:
		return false
	case *types.Named:
		return r.isTranslatable(t.Underlying())
	}
	return true
}
