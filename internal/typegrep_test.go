// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 03/2020

package internal

import (
	"errors"
	"reflect"
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestTypeGrep_Grep_1(t *testing.T) {
	load := func(c *packages.Config, p ...string) ([]*packages.Package, error) {
		return nil, errors.New("error")
	}
	grepper := TypeGrep{
		LoaderFunc: load,
	}
	_, err := grepper.Grep( /* empty */ )
	if err != nil {
		t.Error("unexpected")
	}
}

func TestTypeGrep_Grep_2(t *testing.T) {
	load := func(c *packages.Config, p ...string) ([]*packages.Package, error) {
		return nil, errors.New("error")
	}
	grepper := TypeGrep{
		LoaderFunc: load,
	}
	_, err := grepper.Grep(".")
	if err == nil {
		t.Error("unexpected")
	}
}

func TestTypeGrep_Grep_3(t *testing.T) {
	errs := make([]packages.Error, 0)
	errs = append(errs, packages.Error{})

	pkgs := make([]*packages.Package, 0)
	pkgs = append(pkgs, &packages.Package{Errors: errs})

	load := func(c *packages.Config, p ...string) ([]*packages.Package, error) {
		return pkgs, nil
	}
	grepper := TypeGrep{
		LoaderFunc: load,
	}
	_, err := grepper.Grep(".")
	if err == nil {
		t.Error("unexpected")
	}
}

type TestType struct {
	private [10]chan bool
	Public  map[int]int
}

func TestTypeGrep_Grep_4(t *testing.T) {
	pkgPath := reflect.TypeOf(TestType{}).PkgPath()

	grepper := TypeGrep{
		FilterFunc:   CreateFilterFunc([]string{"TestType", "invalid[regex"}),
		IncludeTests: true,
		IncludeUnexp: true,
	}

	types, err := grepper.Grep(pkgPath)
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
