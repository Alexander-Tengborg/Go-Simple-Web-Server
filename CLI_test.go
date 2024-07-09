package poker_test

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	poker "Go-Simple-Web-Server"

	"github.com/stretchr/testify/assert"
)

type SpyBlindAlerter struct {
	alerts []ScheduledAlert
}

type ScheduledAlert struct {
	At     time.Duration
	Amount int
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	s.alerts = append(s.alerts, ScheduledAlert{duration, amount})
}

type GameSpy struct {
	StartCalled  bool
	StartedWith  int

	FinishCalled bool
	FinishedWith string

	BlindAlert []byte
}

func (g *GameSpy) Start(numberOfPlayers int, alertsDestination io.Writer) {
	g.StartedWith = numberOfPlayers
	g.StartCalled = true
	alertsDestination.Write(g.BlindAlert)
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
	g.FinishCalled = true
}

var dummySpyAlerter = &SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyGame = &GameSpy{}

func TestCLI(t *testing.T) {
	t.Run("Start a game with 2 players and record 'Linda' as the winner", func(t *testing.T) {
		in := userSends("2", "Linda wins")
		stdout := &bytes.Buffer{}

		game := &GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		gotPrompt := stdout.String()
		wantPrompt := poker.PlayerPrompt

		assert.Equal(t, wantPrompt, gotPrompt)
		assert.Equal(t, 2, game.StartedWith)
		assert.Equal(t, "Linda", game.FinishedWith)
	})

	t.Run("Start a game with 10 players and record 'Ted' as the winner", func(t *testing.T) {
		in := userSends("10", "Ted wins")
		stdout := &bytes.Buffer{}

		game := &GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		gotPrompt := stdout.String()
		wantPrompt := poker.PlayerPrompt

		assert.Equal(t, wantPrompt, gotPrompt)
		assert.Equal(t, 10, game.StartedWith)
		assert.Equal(t, "Ted", game.FinishedWith)
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		in := userSends("Hehe")
		stdout := &bytes.Buffer{}

		game := &GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assert.False(t, game.StartCalled)

		want := poker.PlayerPrompt + poker.BadPlayerInputErrMsg
		got := stdout.String()

		assert.Equal(t, want, got)
	})

	t.Run("after starting the game, if '{name} wins' is not written, the game ends with no winner", func(t *testing.T) {
		in := userSends("5", "Dada wins hehe")
		stdout := &bytes.Buffer{}

		game := &GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assert.True(t, game.StartCalled)
		assert.False(t, game.FinishCalled)
	})
}

func userSends(messages ...string) io.Reader {
	builder := strings.Builder{}
	for _, message := range messages {
		builder.WriteString(message + "\n")
	}

	return strings.NewReader(builder.String())
}
