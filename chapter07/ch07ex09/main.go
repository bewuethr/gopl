// ch07ex08 tests the statefulsort implementation for music tracks.
package main

import (
	"log"
	"net/http"

	"gopl/chapter07/ch07ex09/track"
)

var statefulTracks = track.NewStatefulTracks([]*track.Track{
	{"Go", "Moby", "Moby", 1992, track.Length("3m37s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, track.Length("4m24s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, track.Length("4m36s")},
	{"Go", "Delilah", "From the Roots Up", 2012, track.Length("3m38s")},
})

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	track.PrintTracks(statefulTracks, w)
}
