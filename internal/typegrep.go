// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 03/2020

package internal

import (
	"errors"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

type (
	// TypeGrep ...
	TypeGrep struct {
		LoaderFunc LoaderFunc
		FilterFunc FilterFunc

		IncludeTests bool
		IncludeUnexp bool

		typeMap TypeMap
	}

	// LoaderFunc ...
	LoaderFunc func(*packages.Config, ...string) ([]*packages.Package, error)

	// FilterFunc ...
	FilterFunc func(string) bool

	// TypeMap ...
	TypeMap map[string]types.Type
)

// Grep ...
func (t *TypeGrep) Grep(pkgPaths ...string) (TypeMap, error) {
	if t.postInit(); len(pkgPaths) == 0 {
		return t.typeMap, nil
	}
	pkgs, err := t.load(pkgPaths...)
	if err != nil {
		return nil, err
	}
	for _, pkg := range pkgs {
		t.filter(pkg)
	}
	return t.typeMap, nil
}

func (t *TypeGrep) postInit() {
	if t.LoaderFunc == nil {
		t.LoaderFunc = packages.Load
	}
	if t.FilterFunc == nil {
		t.FilterFunc = func(string) bool { return true }
	}
	t.typeMap = make(TypeMap)
}

func (t *TypeGrep) load(patterns ...string) ([]*packages.Package, error) {
	mode := packages.NeedTypes
	mode |= packages.NeedTypesInfo
	conf := &packages.Config{Mode: mode, Tests: t.IncludeTests}

	pkgs, err := t.LoaderFunc(conf, patterns...)
	if err != nil {
		return nil, err
	}
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			err := errors.New(pkg.Errors[0].Msg)
			return nil, err
		}
	}
	return pkgs, nil
}

func (t *TypeGrep) filter(pkg *packages.Package) {
	scope := pkg.Types.Scope()

	for _, name := range scope.Names() {
		if !t.exportable(name) {
			continue
		}
		obj := scope.Lookup(name)
		path := obj.Pkg().Path() + "." + name

		if t.FilterFunc(path) {
			t.visit(obj.Type())
		}
	}
}

func (t *TypeGrep) visit(typ types.Type) {
	switch tt := typ.(type) {

	case *types.Array:
		t.visit(tt.Elem())

	case *types.Chan:
		t.visit(tt.Elem())

	case *types.Interface:
		for i, n := 0, tt.NumEmbeddeds(); i < n; i++ {
			n := tt.EmbeddedType(i).(*types.Named)
			if t.exportable(n.String()) {
				t.visit(tt.EmbeddedType(i))
			}
		}
		for i, n := 0, tt.NumMethods(); i < n; i++ {
			if t.exportable(tt.Method(i).Name()) {
				t.visit(tt.Method(i).Type())
			}
		}

	case *types.Map:
		t.visit(tt.Key())
		t.visit(tt.Elem())

	case *types.Named:
		p := tt.String()
		if _, ok := t.typeMap[p]; !ok {
			t.typeMap[p] = tt
			t.visit(tt.Underlying())
		}

	case *types.Pointer:
		t.visit(tt.Elem())

	case *types.Signature:
		params := tt.Params()
		t.visit(params)
		t.visit(tt.Results())

	case *types.Slice:
		t.visit(tt.Elem())

	case *types.Struct:
		for i := 0; i < tt.NumFields(); i++ {
			if t.exportable(tt.Field(i).Name()) {
				t.visit(tt.Field(i).Type())
			}
		}
	case *types.Tuple:
		for i, n := 0, tt.Len(); i < n; i++ {
			t.visit(tt.At(i).Type())
		}
	}
}

func (t *TypeGrep) exportable(s string) bool {
	if t.IncludeUnexp {
		return true
	}
	path := strings.ReplaceAll(s, ".", "/")
	part := strings.Split(path, "/")
	return token.IsExported(part[len(part)-1])
}
