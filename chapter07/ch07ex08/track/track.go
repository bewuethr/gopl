// Package track provides Track, an implementation of the statefulsort
// interface that represents a music track on an album.
package track

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/bewuethr/gopl/chapter07/ch07ex08/statefulsort"
)

// Track describes a music track on an album.
type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

// Length parses a string representing a duration. It can be used to simplify
// initialization of a new track.
func Length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

func toTrack(x, y interface{}) (*Track, *Track) {
	xx, ok1 := x.(*Track)
	yy, ok2 := y.(*Track)
	if !ok1 || !ok2 {
		log.Fatal("could not convert interface to *Track")
	}
	return xx, yy
}

// These functions behave like C's strcmp
func cmpTitle(x, y interface{}) int {
	xx, yy := toTrack(x, y)
	switch {
	case xx.Title < yy.Title:
		return -1
	case xx.Title == yy.Title:
		return 0
	default:
		return 1
	}
}

func cmpArtist(x, y interface{}) int {
	xx, yy := toTrack(x, y)
	switch {
	case xx.Artist < yy.Artist:
		return -1
	case xx.Artist == yy.Artist:
		return 0
	default:
		return 1
	}
}

func cmpAlbum(x, y interface{}) int {
	xx, yy := toTrack(x, y)
	switch {
	case xx.Album < yy.Album:
		return -1
	case xx.Album == yy.Album:
		return 0
	default:
		return 1
	}
}

func cmpYear(x, y interface{}) int {
	xx, yy := toTrack(x, y)
	switch {
	case xx.Year < yy.Year:
		return -1
	case xx.Year == yy.Year:
		return 0
	default:
		return 1
	}
}

func cmpLength(x, y interface{}) int {
	xx, yy := toTrack(x, y)
	switch {
	case xx.Length < yy.Length:
		return -1
	case xx.Length == yy.Length:
		return 0
	default:
		return 1
	}
}

var (
	cmpFuncs = []func(x, y interface{}) int{cmpTitle, cmpArtist, cmpAlbum, cmpYear, cmpLength}
	names    = []string{"Title", "Artist", "Album", "Year", "Length"}
)

// NewStatefulTracks takes a slice of track pointers as an argument and returns
// a statefulsort for these tracks.
func NewStatefulTracks(tracks []*Track) statefulsort.StatefulSort {
	return statefulsort.NewStatefulSort(toIface(tracks), names, cmpFuncs)
}

func toIface(tracks []*Track) []interface{} {
	ifaces := make([]interface{}, len(tracks))
	for i, t := range tracks {
		ifaces[i] = t
	}
	return ifaces
}

// PrintTracks prints the tracks as they are currently sorted as a table.
func PrintTracks(s statefulsort.StatefulSort) {
	tracks := GetTracks(s)
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush() // calculate column widths and print table
}

// GetTracks returns the track pointer slice for the statefulsort supplied as
// the argument.
func GetTracks(s statefulsort.StatefulSort) []*Track {
	elements := s.Elements()
	tracks := make([]*Track, len(elements))
	for i, e := range elements {
		t, ok := e.(*Track)
		if !ok {
			log.Fatal("could not convert interface to *Track")
		}
		tracks[i] = t
	}
	return tracks
}
