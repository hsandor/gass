package gass

import (
	"testing"
)

func testFile(t *testing.T, file string) {
	if _, err := CompileFile(file, true); err != nil {
		t.Error(err)
	}
}

func TestParser(t *testing.T) {
	testFile("test/parser.gass", "test/parser.css")
	testFile("test/variables.gass", "test/variables.css")
	testFile("test/functions.gass", "test/functions.css")
}

/*
func processExpect(path string, info os.FileInfo, err error) error {
	if filepath.Ext(path) == ".gass" && !strings.HasPrefix(path, "_") {
		in, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		out := CompileString(string(in))
		exp, err := ioutil.ReadFile(path[:len(path)-5] + ".expect")
		if err != nil {
			return err
		}
		ioutil.WriteFile(path[:len(path)-5]+".log", []byte(out), 0)
		if out != string(exp) {
			return errors.New("mismatch:" + path)
		}
	}
	return nil
}

func TestExpect(t *testing.T) {
	if err := filepath.Walk("test/expect", processExpect); err != nil {
		t.Error(err)
	}
}
*/
