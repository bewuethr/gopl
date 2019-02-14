// Ch07ex09 implements a web server with a sortable table that keeps track of
// state.
package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"

	"github.com/bewuethr/gopl/chapter07/ch07ex09/track"
)

var statefulTracks = track.NewStatefulTracks([]*track.Track{
	{Title: "Go", Artist: "Moby", Album: "Moby", Year: 1992, Length: track.Length("3m37s")},
	{Title: "Ready 2 Go", Artist: "Martin Solveig", Album: "Smash", Year: 2011, Length: track.Length("4m24s")},
	{Title: "Go Ahead", Artist: "Alicia Keys", Album: "As I Am", Year: 2007, Length: track.Length("4m36s")},
	{Title: "Go", Artist: "Delilah", Album: "From the Roots Up", Year: 2012, Length: track.Length("3m38s")},
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
