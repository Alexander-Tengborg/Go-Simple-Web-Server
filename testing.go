package poker

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StubPlayerStore struct {
	Scores   map[string]int
	WinCalls []string
	League   []Player
}

func (s *StubPlayerStore) GetPlayerScore(player string) int {
	return s.Scores[player]
}

func (s *StubPlayerStore) RecordWin(player string) {
	s.WinCalls = append(s.WinCalls, player)
}

func (s *StubPlayerStore) GetLeague() []Player {
	return s.League
}

func getLeagueFromResponse(t *testing.T, body io.Reader) []Player {
	t.Helper()

	var got []Player

	err := json.NewDecoder(body).Decode(&got)
	assert.Nilf(t, err, "Unable to parse response %q into []Player", body)

	return got
}
