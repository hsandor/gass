package gass

import (
	"os"
	"path/filepath"
	"testing"
)

func processFile(path string, info os.FileInfo, err error) error {
	if filepath.Ext(path) == ".gass" {
		CompileFile(path)
	}
	return nil
}

func TestExperiment(t *testing.T) {
	filepath.Walk("test", processFile)
	//CompileFile("test/deep.gass")
}
