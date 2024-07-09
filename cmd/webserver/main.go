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

	game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.Alerter), store)

	server, err := poker.NewPlayerServer(store, game)
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(":8080", server)
	if err != nil {
		log.Fatal(err)
	}
}
