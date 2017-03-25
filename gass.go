package gass

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Compile compiles a GASS formatted source provided by r (io.Reader)
// and write the generated CSS into w (io.Writer).
// Any hard error during the compilation returned as usual.
func Compile(w io.Writer, r io.Reader) error {
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	p := newParser()
	dst, err := p.parseString(string(src))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(dst))
	if err != nil {
		return err
	}
	return nil
}

// CompileString compiled a GASS formatted string provided in src into a CSS
// and returned with a possible error.
func CompileString(src string) (string, error) {
	r := bytes.NewBufferString(src)
	w := bytes.NewBuffer([]byte{})
	if err := Compile(w, r); err != nil {
		return "", err
	}
	return w.String(), nil
}

// CompileFile compiles the given file(name) in src from a GASS format into a CSS
// file and returned the new file's path and name with any error code encountered
// during the process. The build argument instruments the function to recompile
// the source even if it's already up to date (by file modificaton time check).
func CompileFile(src string, build bool) (string, error) {
	dst := strings.TrimSuffix(src, filepath.Ext(src)) + ".css"
	sinfo, err := os.Stat(src)
	if err != nil {
		return "", err
	}
	dinfo, err := os.Stat(dst)
	if err == nil {
		if !sinfo.ModTime().After(dinfo.ModTime()) && !build {
			return dst, nil
		}
	}
	fsrc, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer fsrc.Close()

	fdst, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer fdst.Close()

	if err = Compile(fdst, fsrc); err != nil {
		return "", err
	}
	return dst, nil
}
