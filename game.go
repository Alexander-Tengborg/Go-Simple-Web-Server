package poker

import (
	"fmt"
	"io"
	"time"
)

type Game interface {
	Start(numberOfPlayers int, alertsDestination io.Writer)
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

func (g *TexasHoldem) Start(numberOfPlayers int, alertsDestination io.Writer) {
	g.scheduleBlindAlerts(numberOfPlayers, alertsDestination)
}

func (g *TexasHoldem) scheduleBlindAlerts(numberOfPlayers int, alertsDestination io.Writer) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Second

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		g.alerter.ScheduleAlertAt(blindTime, blind, alertsDestination)
		blindTime += blindIncrement
	}
}

func (g *TexasHoldem) Finish(winner string) {
	fmt.Printf("WINNER: %s\n", winner)
	g.store.RecordWin(winner)
}
