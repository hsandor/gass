package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hsandor/gass"
)

func fileExists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	return false
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: gass <build|fmt> [inputfile]")
		os.Exit(1)
	}

	cmd := os.Args[1]

	if cmd == "fmt" {

	} else {
		if len(os.Args) > 2 {
			if os.Args[2] == "-a" {
				// all *.gass in current directory
			} else if fileExists(os.Args[2]) {
				// given file
			} else {
				fmt.Println("can't open file: ", os.Args[2])
			}
		} else {
			if src, err := ioutil.ReadAll(os.Stdin); err != nil {
				fmt.Println(err)
			} else {
				if css, err := gass.ParseString(string(src)); err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(css)
				}
			}
		}
	}
}
