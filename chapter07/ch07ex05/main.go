// Ch07ex05 implements a LimitReader.
package main

import (
	"fmt"
	"io"
	"strings"
)

// MyLimitReader implements an io.Reader that returns EOF after a given number
// of bytes have been read.
type MyLimitReader struct {
	inner      io.Reader
	limit, idx int64
}

func (r *MyLimitReader) Read(p []byte) (int, error) {
	var (
		read int
		err  error
	)

	if r.idx+int64(len(p)) < r.limit {
		read, err = r.inner.Read(p)
		r.idx += int64(read)
		return read, err
	}

	read, err = r.inner.Read(p[:r.limit-r.idx])
	r.idx += int64(read)
	if err == nil {
		return read, io.EOF
	}
	return read, err
}

// LimitReader wraps r into a reader that reads from r and reports EOF after n
// bytes have been read.
func LimitReader(r io.Reader, n int64) io.Reader {
	return &MyLimitReader{
		inner: r,
		limit: n,
		idx:   0,
	}
}

func main() {
	reader := strings.NewReader("1234567890")
	lReader := LimitReader(reader, 5)
	p := make([]byte, 3)
	fmt.Println("Limit reader with 5 bytes")
	n, err := lReader.Read(p)
	fmt.Printf("Reading three bytes [123, 3, nil]:\t'%v'\t%v\t%v\n", string(p), n, err)
	p = make([]byte, 3)
	n, err = lReader.Read(p)
	fmt.Printf("Reading three more bytes [45, 2, io.EOF]:\t'%v'\t%v\t%v\n", string(p), n, err)
	p = make([]byte, 3)
	n, err = lReader.Read(p)
	fmt.Printf("Reading three more bytes ['', 0, io.EOF]:\t'%v'\t%v\t%v\n", string(p), n, err)

	reader = strings.NewReader("12345")
	lReader = LimitReader(reader, 6)
	p = make([]byte, 6)
	fmt.Println("Limit reader with 6 bytes, read from 5 byte string")
	n, err = lReader.Read(p)
	fmt.Printf("Reading six bytes [12345, 5, io.EOF]:\t'%v'\t%v\t%v\n", string(p), n, err)
}
