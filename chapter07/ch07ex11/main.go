// Ch07ex11 adds CRUD functionality to a simple e-commerce server.
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	db := database{"shoes": 50, "socks": 5}
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/read", db.read)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if item == "" {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprint(w, "request is missing item")
		return
	}
	if _, ok := db[item]; ok {
		w.WriteHeader(http.StatusConflict) // 409
		fmt.Fprintf(w, "item exists already: %s\n", item)
		return
	}

	priceStr := req.URL.Query().Get("price")
	if priceStr == "" {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprint(w, "request is missing price")
		return
	}

	price, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "price is not a number: %s\n", priceStr)
		return
	}

	db[item] = dollars(price)
	fmt.Fprintf(w, "inserted %s at price %v\n", item, db[item])
}

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) read(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if item == "" {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprint(w, "request is missing item")
		return
	}
	if _, ok := db[item]; !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}

	priceStr := req.URL.Query().Get("price")
	if priceStr == "" {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprint(w, "request is missing price")
		return
	}

	price, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "price is not a number: %s\n", priceStr)
		return
	}

	db[item] = dollars(price)
	fmt.Fprintf(w, "updated %s to price %v\n", item, db[item])
}

func (db database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if _, ok := db[item]; ok {
		delete(db, item)
		fmt.Fprintf(w, "deleted item %s\n", item)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}
