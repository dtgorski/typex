// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 06/2020

package ts

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"testing"

	typex "github.com/dtgorski/typex/internal"
	"github.com/dtgorski/typex/internal/testdata/p1"
)

// typex -l="ts-type" -r=".*/testdata:" ./internal/testdata/p1 > ./internal/testdata/fixture-ts-type.txt
func TestNewModuleLayout_Type(t *testing.T) {
	fix := "../testdata/fixture-ts-type.txt"

	pac := typex.Packagist{}
	pkg := reflect.TypeOf(p1.D{}).PkgPath()

	types, err := pac.Inspect(pkg)
	if err != nil {
		t.Error("unexpected")
	}

	re := typex.CreatePathReplaceFunc([]string{".*/testdata:", "[:", "", ":"})
	tr := TypeRender{PathReplaceFunc: re}

	buf := &bytes.Buffer{}
	tw := typex.TreeWalk{Layout: NewModuleLayout(buf)}

	if err := tw.Walk(tr.Render(types, false)); err != nil { // <- false
		t.Error("unexpected")
	}
	diff(t, buf, fix)
}

// typex -l="ts-class" -r=".*/testdata:" ./internal/testdata/p1 > ./internal/testdata/fixture-ts-class.txt
func TestNewModuleLayout_Class(t *testing.T) {
	fix := "../testdata/fixture-ts-class.txt"

	pac := typex.Packagist{}
	pkg := reflect.TypeOf(p1.D{}).PkgPath()

	types, err := pac.Inspect(pkg)
	if err != nil {
		t.Error("unexpected")
	}

	re := typex.CreatePathReplaceFunc([]string{".*/testdata:", "[:", "", ":"})
	tr := TypeRender{PathReplaceFunc: re}

	buf := &bytes.Buffer{}
	tw := typex.TreeWalk{Layout: NewModuleLayout(buf)}

	if err := tw.Walk(tr.Render(types, true)); err != nil { // <- true
		t.Error("unexpected")
	}
	diff(t, buf, fix)
}

func diff(t *testing.T, buf *bytes.Buffer, fixture string) {
	fix, err := ioutil.ReadFile(fixture)
	if err != nil {
		t.Error("unexpected")
	}
	if !bytes.EqualFold(fix, buf.Bytes()) {
		tmp, err := ioutil.TempFile("", "")
		if err != nil {
			t.Error(err)
		}
		defer func() { _ = os.Remove(tmp.Name()) }()

		err = ioutil.WriteFile(tmp.Name(), buf.Bytes(), 0644)
		if err != nil {
			t.Error(err)
		}
		cmd := exec.Command("diff", "-u", fixture, tmp.Name())
		out, err := cmd.Output()
		if err != nil {
			t.Error(err)
		}
		t.Error(string(out))
	}
}
