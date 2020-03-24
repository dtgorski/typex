// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 03/2020

package internal

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func TestExGoType_Export(t *testing.T) {

	fixPath := "./testdata/fixture-01.txt"
	pkgPath := reflect.TypeOf(TestType{}).PkgPath()

	find := TypeGrep{}
	types, err := find.Grep(pkgPath + "/testdata/p1")
	if err != nil {
		t.Error("unexpected")
	}

	variant := ExGoType{}
	result := &bytes.Buffer{}

	err = variant.Export(result, types)
	if err != nil {
		t.Error("unexpected")
	}

	fix, err := ioutil.ReadFile(fixPath)
	if err != nil {
		t.Error("unexpected")
	}

	if !bytes.EqualFold(fix, result.Bytes()) {
		tmp, err := ioutil.TempFile("", "")
		if err != nil {
			t.Error(err)
		}
		defer func() { _ = os.Remove(tmp.Name()) }()

		err = ioutil.WriteFile(tmp.Name(), result.Bytes(), 0644)
		if err != nil {
			t.Error(err)
		}
		cmd := exec.Command("diff", "-u", fixPath, tmp.Name())
		out, err := cmd.Output()
		if err != nil {
			t.Error(err)
		}
		t.Error(string(out))
	}
}
