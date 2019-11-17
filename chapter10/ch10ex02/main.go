// Ch10ex02 prints the files and their sizes from archive files given as command
// line parameters. The archive type has to be registered by using a blank
// import of the corresponding package.
package main

import (
	"fmt"
	"os"

	"github.com/bewuethr/gopl/chapter10/ch10ex02/archive"
	_ "github.com/bewuethr/gopl/chapter10/ch10ex02/archive/tar"
	_ "github.com/bewuethr/gopl/chapter10/ch10ex02/archive/zip"
)

func main() {
	for _, name := range os.Args[1:] {
		fdata, kind, err := archive.Read(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ch10ex02: could not read archive %v: %v\n", name, err)
			continue
		}
		fmt.Printf("Files in archive %v (type %v):\n", name, kind)
		for _, d := range fdata {
			fmt.Printf("%v: %dB\n", d.Name, d.Size)
		}
	}
}
