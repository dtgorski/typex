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

func TestNewModuleLayout(t *testing.T) {

	// typex -l="ts" -r=".*/testdata:" ./internal/testdata/p1 > ./internal/testdata/fixture-ts.txt
	fixture := "../testdata/fixture-ts.txt"

	pac := typex.Packagist{}
	pkg := reflect.TypeOf(p1.D{}).PkgPath()

	types, err := pac.Inspect(pkg)
	if err != nil {
		t.Error("unexpected")
	}

	re := typex.CreatePathReplaceFunc([]string{".*/testdata:", "[:", "", ":"})
	tr := TypeRender{PathReplaceFunc: re}

	buf := &bytes.Buffer{}
	tw := typex.TreeWalker{Layout: NewModuleLayout(buf)}

	err = tw.Walk(tr.Render(types))
	if err != nil {
		t.Error("unexpected")
	}

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
