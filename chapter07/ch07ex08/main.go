// Ch07ex08 tests the statefulsort implementation for music tracks.
package main

import (
	"fmt"
	"sort"

	"github.com/bewuethr/gopl/chapter07/ch07ex08/track"
)

var tracks = []*track.Track{
	{Title: "Go", Artist: "Moby", Album: "Moby", Year: 1992, Length: track.Length("3m37s")},
	{Title: "Ready 2 Go", Artist: "Martin Solveig", Album: "Smash", Year: 2011, Length: track.Length("4m24s")},
	{Title: "Go Ahead", Artist: "Alicia Keys", Album: "As I Am", Year: 2007, Length: track.Length("4m36s")},
	{Title: "Go", Artist: "Delilah", Album: "From the Roots Up", Year: 2012, Length: track.Length("3m38s")},
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
