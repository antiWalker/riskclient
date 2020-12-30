package core

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func Test_Eval(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	files := []string{
		"demo.json",
	}

	for _, fn := range files {
		absFileName := filepath.Join(dir, "rules", fn)

		fp, err := os.Open(absFileName)
		if err != nil {
			continue
		}
		data, err := ioutil.ReadAll(fp)
		_ = fp.Close()
		if err != nil {
			continue
		}

		var thisParams InParams

		thisParams.BusinessLine = "sales"

		checkResult, err := Eval(string(data), &thisParams)

		if checkResult == false {
			t.Log("Pass Eval")
		} else {
			t.Error("Failed Eval")
		}
	}
}
