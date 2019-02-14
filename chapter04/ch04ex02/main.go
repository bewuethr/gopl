// Ch04ex02 prints the SHA256 hash of its standard input, but supports command
// line flags for SHA384 or SHA512 hashes as well.
package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	bitPtr := flag.Int("bits", 256, "256, 384 or 512 bit checksum")

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	flag.Parse()
	switch *bitPtr {
	case 256:
		fmt.Printf("%x\n", sha256.Sum256(input))
	case 384:
		fmt.Printf("%x\n", sha512.Sum384(input))
	case 512:
		fmt.Printf("%x\n", sha512.Sum512(input))
	default:
		fmt.Fprintln(os.Stderr, "Usage: ch04ex02 [-bits=256|284|512]")
	}
}
