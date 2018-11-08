// ch07ex09 implements a web server with a sortable table that keeps track of
// state.
package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"

	"gopl/chapter07/ch07ex09/track"
)

var statefulTracks = track.NewStatefulTracks([]*track.Track{
	{"Go", "Moby", "Moby", 1992, track.Length("3m37s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, track.Length("4m24s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, track.Length("4m36s")},
	{"Go", "Delilah", "From the Roots Up", 2012, track.Length("3m38s")},
})

var pageTemplate = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Music tracks</title>
    <style type="text/css">
    table {
      margin: auto;
    }
    table, th, td {
      border: solid;
      border-collapse: collapse;
    }
    </style>
  </head>
  <body>
    {{.}}
  </body>
</html>`))

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Starting server on port 8000")
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	if sortKey := r.FormValue("orderby"); sortKey != "" {
		if err := statefulTracks.SetPrimary(sortKey); err != nil {
			panic(err)
		}
	}
	sort.Sort(statefulTracks)

	table := new(bytes.Buffer)
	track.PrintTrackTable(statefulTracks, table)
	if err := pageTemplate.Execute(w, template.HTML(table.String())); err != nil {
		panic(err)
	}
}
