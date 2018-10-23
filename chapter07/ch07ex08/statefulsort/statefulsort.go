// Package statefulsort provides a sort interface for music tracks with a
// mutable order in which track fields are taken into account for sorting.
package statefulsort

import "time"

// Track describes a music track on an album.
type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

// These functions behave like C's strcmp
func lessTitle(x, y *Track) int {
	if x.Title < y.Title {
		return -1
	} else if x.Title == y.Title {
		return 0
	} else {
		return 1
	}
}

func lessArtist(x, y *Track) int {
	if x.Artist < y.Artist {
		return -1
	} else if x.Artist == y.Artist {
		return 0
	} else {
		return 1
	}
}

func lessAlbum(x, y *Track) int {
	if x.Album < y.Album {
		return -1
	} else if x.Album == y.Album {
		return 0
	} else {
		return 1
	}
}

func lessYear(x, y *Track) int {
	if x.Year < y.Year {
		return -1
	} else if x.Year == y.Year {
		return 0
	} else {
		return 1
	}
}

func lessLength(x, y *Track) int {
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
	{"Title", lessTitle},
	{"Artist", lessArtist},
	{"Album", lessAlbum},
	{"Year", lessYear},
	{"Length", lessLength},
}

// StatefulSort implements sort.Interface and keeps track of the order in which
// the fields are looked at when comparing two tracks.
type StatefulSort struct {
	t           []*Track
	sortStructs []*sortStruct
}

// NewStatefulTracks returns a statful sort for the slice of tracks t.
func NewStatefulTracks(t []*Track) StatefulSort {
	return StatefulSort{t, sortStructs}
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
			return true
		case 1:
			return false
		}
	}
	return false
}
