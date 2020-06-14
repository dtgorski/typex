// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 06/2020

package internal

import (
	"errors"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

type (
	// Packagist ...
	Packagist struct {
		PackageLoaderFunc PackageLoaderFunc
		PathFilterFunc    PathFilterFunc

		IncludeUnexported bool
		IncludeTestFiles  bool

		typeMap TypeMap
	}

	// PackageLoaderFunc returns the Go packages named by the given patterns.
	PackageLoaderFunc func(*packages.Config, ...string) ([]*packages.Package, error)

	// PathFilterFunc is a simple string filtering/matching function.
	PathFilterFunc func(string) bool

	// PathReplaceFunc is a simple string replacement function.
	PathReplaceFunc func(string) string

	// QualifierFunc returns the type qualifier (path) of a package.
	QualifierFunc func(p types.Package) string

	// TypeMap is a set of types indexed by their package import paths.
	TypeMap map[string]types.Type
)

// Inspect finds matching types in the packages named by the given patterns.
func (p *Packagist) Inspect(patterns ...string) (TypeMap, error) {
	if p.PackageLoaderFunc == nil {
		p.PackageLoaderFunc = packages.Load
	}
	if p.PathFilterFunc == nil {
		p.PathFilterFunc = func(string) bool { return true }
	}
	if p.typeMap = make(TypeMap); len(patterns) == 0 {
		return p.typeMap, nil
	}
	pkgs, err := p.load(patterns...)
	if err != nil {
		return nil, err
	}
	for _, pkg := range pkgs {
		p.filter(pkg)
	}
	return p.typeMap, nil
}

func (p *Packagist) load(patterns ...string) ([]*packages.Package, error) {
	mode := packages.NeedTypes
	mode |= packages.NeedTypesInfo
	conf := &packages.Config{Mode: mode, Tests: p.IncludeTestFiles}

	pkgs, err := p.PackageLoaderFunc(conf, patterns...)
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

func (p *Packagist) filter(pkg *packages.Package) {
	scope := pkg.Types.Scope()

	for _, name := range scope.Names() {
		if !p.isExported(name) {
			continue
		}
		obj := scope.Lookup(name)
		path := obj.Pkg().Path() + "." + name
		if !p.PathFilterFunc(path) {
			continue
		}
		p.visit(obj.Type())
	}
}

func (p *Packagist) visit(t types.Type) {
	switch tt := t.(type) {
	case *types.Array:
		p.visit(tt.Elem())

	case *types.Chan:
		p.visit(tt.Elem())

	case *types.Interface:
		for i, n := 0, tt.NumEmbeddeds(); i < n; i++ {
			nn := tt.EmbeddedType(i).(*types.Named)
			if p.isExported(nn.String()) {
				p.visit(tt.EmbeddedType(i))
			}
		}
		for i, n := 0, tt.NumMethods(); i < n; i++ {
			if p.isExported(tt.Method(i).Name()) {
				p.visit(tt.Method(i).Type())
			}
		}
	case *types.Map:
		p.visit(tt.Key())
		p.visit(tt.Elem())

	case *types.Named:
		s := tt.String()
		if _, ok := p.typeMap[s]; !ok {
			p.typeMap[s] = tt
			p.visit(tt.Underlying())
		}
	case *types.Pointer:
		p.visit(tt.Elem())

	case *types.Signature:
		p.visit(tt.Params())
		p.visit(tt.Results())

	case *types.Slice:
		p.visit(tt.Elem())

	case *types.Struct:
		for i := 0; i < tt.NumFields(); i++ {
			if p.isExported(tt.Field(i).Name()) {
				p.visit(tt.Field(i).Type())
			}
		}
	case *types.Tuple:
		for i, n := 0, tt.Len(); i < n; i++ {
			p.visit(tt.At(i).Type())
		}
	}
}

func (p *Packagist) isExported(s string) bool {
	if p.IncludeUnexported {
		return true
	}
	n := strings.ReplaceAll(s, ".", "/")
	i := strings.LastIndex(n, "/")
	return token.IsExported(n[i+1:])
}
