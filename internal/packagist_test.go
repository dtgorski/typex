// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 06/2020

package internal

import (
	"errors"
	"reflect"
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestInspector_Inspect_1(t *testing.T) {
	loader := func(c *packages.Config, p ...string) ([]*packages.Package, error) {
		return nil, errors.New("error")
	}
	pac := Packagist{
		PackageLoaderFunc: loader,
	}
	_, err := pac.Inspect( /* empty */ )
	if err != nil {
		t.Error("unexpected")
	}
}

func TestInspector_Inspect_2(t *testing.T) {
	loader := func(c *packages.Config, p ...string) ([]*packages.Package, error) {
		return nil, errors.New("error")
	}
	pac := Packagist{
		PackageLoaderFunc: loader,
	}
	_, err := pac.Inspect(".")
	if err == nil {
		t.Error("unexpected")
	}
}

func TestInspector_Inspect_3(t *testing.T) {
	errs := make([]packages.Error, 0)
	errs = append(errs, packages.Error{})

	pkgs := make([]*packages.Package, 0)
	pkgs = append(pkgs, &packages.Package{Errors: errs})

	loader := func(c *packages.Config, p ...string) ([]*packages.Package, error) {
		return pkgs, nil
	}
	pac := Packagist{
		PackageLoaderFunc: loader,
	}
	_, err := pac.Inspect(".")
	if err == nil {
		t.Error("unexpected")
	}
}

type TestType struct {
	private [10]chan bool
	Public  map[int]int
}

func TestInspector_Inspect_4(t *testing.T) {
	pkgPath := reflect.TypeOf(TestType{}).PkgPath()

	pac := Packagist{
		PathFilterFunc:    CreatePathFilterFunc([]string{"TestType", "invalid[regex"}),
		IncludeTestFiles:  true,
		IncludeUnexported: true,
	}

	types, err := pac.Inspect(pkgPath)
	if err != nil {
		t.Error("unexpected")
		return
	}

	typeName := pkgPath + ".TestType"
	if types[typeName] == nil {
		t.Error("unexpected")
		return
	}
}
