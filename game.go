package poker

import (
	"fmt"
	"time"
)

type Game interface {
	Start(numberOfPlayers int)
	Finish(winner string)
}

type TexasHoldem struct {
	store   PlayerStore
	alerter BlindAlerter
}

func NewTexasHoldem(alerter BlindAlerter, store PlayerStore) *TexasHoldem {
	return &TexasHoldem{
		alerter: alerter,
		store:   store,
	}
}

func (g *TexasHoldem) Start(numberOfPlayers int) {
	g.scheduleBlindAlerts(numberOfPlayers)
}

func (g *TexasHoldem) scheduleBlindAlerts(numberOfPlayers int) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		g.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime += blindIncrement
	}
}

func (g *TexasHoldem) Finish(winner string) {
	fmt.Printf("WINNER: %s\n", winner)
	g.store.RecordWin(winner)
}
