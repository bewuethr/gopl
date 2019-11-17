package zip

import (
	"archive/zip"

	"github.com/bewuethr/gopl/chapter10/ch10ex02/archive"
)

func init() {
	archive.RegisterFormat("zip", Read, Match)
}

// Read reads the zip archive file designated by name and returns a slice of
// file data.
func Read(name string) ([]archive.FileData, error) {
	r, err := zip.OpenReader(name)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var fdata []archive.FileData
	for _, f := range r.File {
		fdata = append(fdata, archive.FileData{
			Name: f.Name,
			Size: f.UncompressedSize64,
		})
	}

	return fdata, nil
}

// Match checks if name is a file containing a zip archive.
func Match(name string) bool {
	r, err := zip.OpenReader(name)
	if err != nil {
		return false
	}
	r.Close()
	return true
}
