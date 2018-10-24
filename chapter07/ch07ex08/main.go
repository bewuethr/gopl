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

var tracks = []*statefulsort.Track{
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

func printTracks(tracks []*statefulsort.Track) {
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
	statefulTracks := statefulsort.NewStatefulTracks(tracks)
	fmt.Println("Initial order:")
	printTracks(statefulTracks.Tracks())
	sort.Sort(statefulTracks)
	fmt.Println("\nAfter sort:")
	printTracks(statefulTracks.Tracks())

	fmt.Println("\nSet primary to Year")
	statefulTracks.SetPrimary("Year")
	sort.Sort(statefulTracks)
	printTracks(statefulTracks.Tracks())

	fmt.Println("\nSet primary to Title")
	statefulTracks.SetPrimary("Title")
	sort.Sort(statefulTracks)
	printTracks(statefulTracks.Tracks())

	fmt.Println("\nSet primary to Artist, then Title")
	statefulTracks.SetPrimary("Artist")
	statefulTracks.SetPrimary("Title")
	sort.Sort(statefulTracks)
	printTracks(statefulTracks.Tracks())

	fmt.Println("\nSet primary to Title again, reverse sort order")
	statefulTracks.SetPrimary("Title")
	sort.Sort(statefulTracks)
	printTracks(statefulTracks.Tracks())
}
