package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hsandor/gass"
)

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage:")
		fmt.Println("  gass build                 - compile all gass file in current directory")
		fmt.Println("  gass [build] <input file>  - compile the given file ")
		os.Exit(1)
	}

	cmd := os.Args[1]

	if cmd == "fmt" {
		fmt.Println("not implemented - yet")
	} else if cmd == "build" {
		if len(os.Args) >= 3 {
			if css, err := gass.CompileFile(os.Args[2], false); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("generate:", css)
			}
		} else {
			if files, err := filepath.Glob("./*.gass"); err != nil {
				fmt.Println(err)
			} else {
				for _, file := range files {
					if !strings.HasPrefix(file, "_") {
						if css, err := gass.CompileFile(file, false); err != nil {
							fmt.Println(err)
						} else {
							fmt.Println("generate:", css)
						}
					}
				}
			}
		}
	} else {
		if fileExists(os.Args[1]) {
			if css, err := gass.CompileFile(os.Args[1], false); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("generate:", css)
			}
		} else {
			fmt.Println("argument error:", os.Args[1])
		}
	}
}
