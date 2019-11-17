package archive

import "errors"

// FileData describes a file in an archive by name and size (uncompressed).
type FileData struct {
	Name string
	Size uint64
}

// A format holds an archive type's name and read function; the match function
// checks if a file specified by its name is a valid archive of that format.
type format struct {
	name  string
	read  func(string) ([]FileData, error)
	match func(string) bool
}

// ErrFormat is the error returned when archive is asked to get the files from
// an unkown archive format.
var ErrFormat = errors.New("archive: unknown format")

// format is the list of registered formats.
var formats []format

// RegisterFormat registers an archive format by mapping its name to a read
// function and a function that checks if a file is of a specific archive type.
// It is typically called from an init function in the archive package
// registering itself.
func RegisterFormat(name string, read func(string) ([]FileData, error), match func(string) bool) {
	formats = append(formats, format{
		name:  name,
		read:  read,
		match: match,
	})
}

// Read returns a slice of file data structs for the archive in name and the
// type of archive, or ErrFormat if the archive type isn't registered.
func Read(name string) ([]FileData, string, error) {
	for _, f := range formats {
		if f.match(name) {
			fdata, err := f.read(name)
			return fdata, f.name, err
		}
	}
	return nil, "", ErrFormat
}
