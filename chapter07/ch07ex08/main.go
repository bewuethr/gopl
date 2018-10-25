// ch07ex08 tests the statefulsort implementation for music tracks.
package main

import (
	"fmt"
	"sort"

	"gopl/chapter07/ch07ex08/track"
)

var tracks = []*track.Track{
	{"Go", "Moby", "Moby", 1992, track.Length("3m37s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, track.Length("4m24s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, track.Length("4m36s")},
	{"Go", "Delilah", "From the Roots Up", 2012, track.Length("3m38s")},
}

func main() {
	statefulTracks := track.NewStatefulTracks(tracks)
	fmt.Println("Initial order:")
	track.PrintTracks(statefulTracks)
	sort.Sort(statefulTracks)
	fmt.Println("\nAfter sort:")
	track.PrintTracks(statefulTracks)

	fmt.Println("\nSet primary to Year")
	statefulTracks.SetPrimary("Year")
	sort.Sort(statefulTracks)
	track.PrintTracks(statefulTracks)

	fmt.Println("\nSet primary to Title")
	statefulTracks.SetPrimary("Title")
	sort.Sort(statefulTracks)
	track.PrintTracks(statefulTracks)

	fmt.Println("\nSet primary to Artist, then Title")
	statefulTracks.SetPrimary("Artist")
	statefulTracks.SetPrimary("Title")
	sort.Sort(statefulTracks)
	track.PrintTracks(statefulTracks)

	fmt.Println("\nSet primary to Title again, reverse sort order")
	statefulTracks.SetPrimary("Title")
	sort.Sort(statefulTracks)
	track.PrintTracks(statefulTracks)
}
