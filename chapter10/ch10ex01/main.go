// Ch10ex01 converts an image file read from standard input to the format
// specified with the -to flag (default: png) and writes it to standard output.
// Supported formats are png, jpeg and gif.
package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"
)

func main() {
	var toFmt string
	flag.StringVar(&toFmt, "to", "png", "output format (png, jpeg, gif)")
	flag.Parse()

	if err := toOutputFmt(os.Stdin, os.Stdout, toFmt); err != nil {
		fmt.Fprintf(os.Stderr, "ch10ex01: %s\n", err)
		os.Exit(1)
	}
}

func toOutputFmt(in io.Reader, out io.Writer, toFmt string) error {
	img, kind, err := image.Decode(in)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Input format =", kind)

	switch strings.ToLower(toFmt) {
	case "jpg", "jpeg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
	case "png":
		return png.Encode(out, img)
	case "gif":
		return gif.Encode(out, img, nil)
	default:
		return fmt.Errorf("unknown output format %v", toFmt)
	}
}
