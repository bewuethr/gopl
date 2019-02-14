// Ch07ex12 has the list endpoint rewritten to use html/template.
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var itemTable = template.Must(template.New("table").Parse(`<table>
<tr>
  <th>Item</th>
  <th>Price</th>
</tr>
{{- range $item, $price := .}}
<tr>
  <td>{{$item}}</td>
  <td>{{$price}}</td>
</tr>
{{- end}}
</table>`))

func main() {
	db := database{"shoes": 50, "socks": 5}
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	if err := itemTable.Execute(w, db); err != nil {
		panic(err)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}
