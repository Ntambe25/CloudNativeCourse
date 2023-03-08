package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	db := database{"shoes": 50, "socks": 5}
	mux := http.NewServeMux()
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)
	mux.HandleFunc("/update", db.update)
	mux.HandleFunc("/create", db.create)
	mux.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe(":8000", mux))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
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

func (db database) update(w http.ResponseWriter, req *http.Request) {

	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")
	priceFloat, _ := strconv.ParseFloat(price, 32)

	if _, ok := db[item]; ok {
		delete(db, item)
		db[item] = dollars(priceFloat)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}

}

func (db database) create(w http.ResponseWriter, req *http.Request) {

	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")
	priceFloat, _ := strconv.ParseFloat(price, 32)

	db[item] = dollars(priceFloat)

}

func (db database) delete(w http.ResponseWriter, req *http.Request) {

	item := req.URL.Query().Get("item")

	//delete(db, item)

	if _, ok := db[item]; ok {
		delete(db, item)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}

}
