package gass

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func BuildString(str string) (css string, err error) {
	return "", nil
}

func BuildFile(file string) (err error) {
	return nil
}

func BuildFolder(folder string) (err error) {
	return nil
}

func FormatString(str string) (gass string, err error) {
	return "", nil
}

func FormatFile(file string) (gass string, err error) {
	return "", nil
}

func FormatFolder(folder string) (err error) {
	return nil
}

// https://www.reddit.com/r/golang/comments/2gkofb/yosssigcss_pure_go_css_preprocessor/

func ParseString(s string) (string, error) {
	p := newParser()
	return p.parseString(s)
}

func CompileFile(name string) {
	src, err := ioutil.ReadFile(name)

	if err != nil {
		fmt.Println(err)
		return
	}

	css, _ := ParseString(string(src))
	nfn := strings.TrimSuffix(name, filepath.Ext(name)) + ".css"

	ioutil.WriteFile(nfn, []byte(css), 0)
}
