// Package statefulsort provides a sort interface for music tracks with a
// mutable order in which track fields are taken into account for sorting.
package statefulsort

import (
	"fmt"
	"time"
)

// Track describes a music track on an album.
type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

// These functions behave like C's strcmp
func cmpTitle(x, y *Track) int {
	if x.Title < y.Title {
		return -1
	} else if x.Title == y.Title {
		return 0
	} else {
		return 1
	}
}

func cmpArtist(x, y *Track) int {
	if x.Artist < y.Artist {
		return -1
	} else if x.Artist == y.Artist {
		return 0
	} else {
		return 1
	}
}

func cmpAlbum(x, y *Track) int {
	if x.Album < y.Album {
		return -1
	} else if x.Album == y.Album {
		return 0
	} else {
		return 1
	}
}

func cmpYear(x, y *Track) int {
	if x.Year < y.Year {
		return -1
	} else if x.Year == y.Year {
		return 0
	} else {
		return 1
	}
}

func cmpLength(x, y *Track) int {
	if x.Length < y.Length {
		return -1
	} else if x.Length == y.Length {
		return 0
	} else {
		return 1
	}
}

type sortStruct struct {
	name     string
	sortFunc func(x, y *Track) int
}

var sortStructs = []*sortStruct{
	{"Title", cmpTitle},
	{"Artist", cmpArtist},
	{"Album", cmpAlbum},
	{"Year", cmpYear},
	{"Length", cmpLength},
}

// StatefulSort implements sort.Interface and keeps track of the order in which
// the fields are looked at when comparing two tracks.
type StatefulSort struct {
	t           []*Track
	sortStructs []*sortStruct
	reverse     bool
}

// NewStatefulTracks returns a statful sort for the slice of tracks t.
func NewStatefulTracks(t []*Track) StatefulSort {
	return StatefulSort{t, sortStructs, false}
}

// Tracks returns the tracks of a stateful sort.
func (s StatefulSort) Tracks() []*Track {
	return s.t
}

// Len returns the number of tracks in s.
func (s StatefulSort) Len() int { return len(s.t) }

// Swap swaps the tracks at indices i and j.
func (s StatefulSort) Swap(i, j int) { s.t[i], s.t[j] = s.t[j], s.t[i] }

// Less compares two tracks based on the current ordering of the comparison
// functions.
func (s StatefulSort) Less(i, j int) bool {
	for _, sStr := range s.sortStructs {
		switch sStr.sortFunc(s.t[i], s.t[j]) {
		case -1:
			return !s.reverse
		case 1:
			return s.reverse
		}
	}
	return s.reverse
}

// SetPrimary moves the sort function corresponding to name n to the front of
// the sort function slice. If it already is at the front, the sort order is
// reversed by flipping reverse.
func (s *StatefulSort) SetPrimary(n string) error {
	idx := -1
	for i, v := range s.sortStructs {
		if v.name == n {
			idx = i
			break
		}
	}
	if idx == -1 {
		// Name not found
		return fmt.Errorf("sortStruct with name %v not found", n)
	}

	if idx == 0 {
		// Reverse sort order
		s.reverse = !s.reverse
		return nil
	}

	// Rearrange to move new primary to front
	newSortStructs := []*sortStruct{s.sortStructs[idx]}
	newSortStructs = append(newSortStructs, s.sortStructs[:idx]...)
	newSortStructs = append(newSortStructs, s.sortStructs[idx+1:]...)
	s.sortStructs = newSortStructs
	return nil
}
