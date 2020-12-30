package core

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func Test_AllRuleFile(t *testing.T) {
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

		ast, err := constructNodeFromString(nil,data)
		if err != nil {
			t.Fatal("file: " + fn + "\nerror:" + err.Error())
		} else {
			t.Log(ast)
		}
	}
}
