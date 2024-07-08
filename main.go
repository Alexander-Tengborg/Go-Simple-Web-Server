package main

import (
	"log"
	"net/http"
)

func main() {
	store, err := NewBoltPlayerStore("prod.db")
	if err != nil {
		log.Fatal(err)
	}

	server := NewPlayerServer(store)

	err = http.ListenAndServe(":8080", server)
	if err != nil {
		log.Fatal(err)
	}
}
