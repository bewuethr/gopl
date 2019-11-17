package tar

import (
	"archive/tar"
	"io"
	"os"

	"github.com/bewuethr/gopl/chapter10/ch10ex02/archive"
)

func init() {
	archive.RegisterFormat("tar", Read, Match)
}

// Read reads the tar archive file designated by name and returns a slice of
// file data.
func Read(name string) ([]archive.FileData, error) {
	tr, err := getTarReader(name)
	if err != nil {
		return nil, err
	}

	var fdata []archive.FileData
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		fdata = append(fdata, archive.FileData{
			Name: hdr.Name,
			Size: uint64(hdr.Size),
		})
	}
	return fdata, nil
}

// Match checks if name is a file containing a tar archive.
func Match(name string) bool {
	tr, err := getTarReader(name)
	if err != nil {
		return false
	}
	hdr, err := tr.Next()
	if err != nil {
		return false
	}
	return hdr.Format != tar.FormatUnknown
}

func getTarReader(name string) (*tar.Reader, error) {
	r, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return tar.NewReader(r), nil
}
