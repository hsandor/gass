package gass

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestRegression(t *testing.T) {
	files, err := filepath.Glob("test/regression/*.gass")
	if err != nil {
		t.Error(err)
	} else {
		for _, file := range files {
			src, _ := ioutil.ReadFile(file)
			dst, err := CompileString(string(src))
			if err != nil {
				t.Error(err)
			}
			css := strings.TrimSuffix(file, filepath.Ext(file)) + ".css"
			chk, _ := ioutil.ReadFile(css)
			if strings.TrimSpace(dst) != strings.TrimSpace(string(chk)) {
				out := strings.TrimSuffix(file, filepath.Ext(file)) + ".err"
				ioutil.WriteFile(out, []byte(dst), 0)
				t.Error(errors.New("result/expected mismatch:" + out))
			}
		}
	}
}
