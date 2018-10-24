// ch07ex08 tests the statefulsort implementation for music tracks.
package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"gopl/chapter07/ch07ex08/statefulsort"
)

// Track describes a music track on an album.
type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

var tracks = []*Track{
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

// These functions behave like C's strcmp
func cmpTitle(x, y *Track) int {
	switch {
	case x.Title < y.Title:
		return -1
	case x.Title == y.Title:
		return 0
	default:
		return 1
	}
}

func cmpArtist(x, y *Track) int {
	switch {
	case x.Artist < y.Artist:
		return -1
	case x.Artist == y.Artist:
		return 0
	default:
		return 1
	}
}

func cmpAlbum(x, y *Track) int {
	switch {
	case x.Album < y.Album:
		return -1
	case x.Album == y.Album:
		return 0
	default:
		return 1
	}
}

func cmpYear(x, y *Track) int {
	switch {
	case x.Year < y.Year:
		return -1
	case x.Year == y.Year:
		return 0
	default:
		return 1
	}
}

func cmpLength(x, y *Track) int {
	switch {
	case x.Length < y.Length:
		return -1
	case x.Length == y.Length:
		return 0
	default:
		return 1
	}
}

var (
	cmpFuncs = []func(x, y interface{}) int{cmpTitle, cmpArtist, cmpAlbum, cmpYear, cmpLength}
	names    = []string{"Title", "Artist", "Album", "Year", "Length"}
)

func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush() // calculate column widths and print table
}

func main() {
	statefulTracks := statefulsort.NewStatefulSort(tracks, names, cmpFuncs)
	fmt.Println("Initial order:")
	printTracks(statefulTracks.Elements())
	sort.Sort(statefulTracks)
	fmt.Println("\nAfter sort:")
	printTracks(statefulTracks.Elements())

	fmt.Println("\nSet primary to Year")
	statefulTracks.SetPrimary("Year")
	sort.Sort(statefulTracks)
	printTracks(statefulTracks.Elements())

	fmt.Println("\nSet primary to Title")
	statefulTracks.SetPrimary("Title")
	sort.Sort(statefulTracks)
	printTracks(statefulTracks.Elements())

	fmt.Println("\nSet primary to Artist, then Title")
	statefulTracks.SetPrimary("Artist")
	statefulTracks.SetPrimary("Title")
	sort.Sort(statefulTracks)
	printTracks(statefulTracks.Elements())

	fmt.Println("\nSet primary to Title again, reverse sort order")
	statefulTracks.SetPrimary("Title")
	sort.Sort(statefulTracks)
	printTracks(statefulTracks.Elements())
}
