package main

import (
	poker "Go-Simple-Web-Server"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Let's play poker!")
	fmt.Println("Type {Name} wins to record a win")

	store, err := poker.NewBoltPlayerStore("prod.db")
	if err != nil {
		log.Fatal(err)
	}

	game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.StdOutAlerter), store)

	poker.NewCLI(os.Stdin, os.Stdout, game).PlayPoker()
}
