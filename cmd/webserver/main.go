package main

import (
	"log"
	"net/http"

	"Go-Simple-Web-Server"
)

func main() {
	store, err := poker.NewBoltPlayerStore("prod.db")
	if err != nil {
		log.Fatal(err)
	}

	server := poker.NewPlayerServer(store)

	err = http.ListenAndServe(":8080", server)
	if err != nil {
		log.Fatal(err)
	}
}
