package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(player string) int {
	return s.scores[player]
}

func (s *StubPlayerStore) RecordWin(player string) {
	s.winCalls = append(s.winCalls, player)
}

func (s *StubPlayerStore) GetLeague() []Player {
	return s.league
}

func TestGETPlayers(t *testing.T) {
	store := &StubPlayerStore{
		map[string]int{
			"Linda": 20,
			"Steve": 4,
		},
		nil,
		nil,
	}
	server := NewPlayerServer(store)

	t.Run("returns Linda's score", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetScoreRequest("Linda")

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "20", response.Body.String())
	})

	t.Run("returns Steve's score", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetScoreRequest("Steve")

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "4", response.Body.String())
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetScoreRequest("Drogba")

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	store := &StubPlayerStore{
		map[string]int{},
		[]string{},
		nil,
	}
	server := NewPlayerServer(store)

	t.Run("it record wins when POST", func(t *testing.T) {
		player := "Linda"
		response := httptest.NewRecorder()
		request := newPostWinRequest(player)

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusAccepted, response.Code)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to recordWin, want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], player)
		}
	})
}

func TestLeague(t *testing.T) {
	t.Run("it returns 200 on /league", func(t *testing.T) {
		want := []Player{
			{"Linda", 28},
			{"George", 38},
			{"Benedict", 2},
		}

		store := &StubPlayerStore{nil, nil, want}
		server := NewPlayerServer(store)

		response := httptest.NewRecorder()
		request := newLeagueRequest()

		server.ServeHTTP(response, request)

		gotContentType := response.Result().Header.Get("content-type")
		got := getLeagueFromResponse(t, response.Body)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "application/json", gotContentType)
		assert.ElementsMatch(t, want, got)
	})
}

func newGetScoreRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", player), nil)
	return request
}

func newPostWinRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", player), nil)
	return request
}

func newLeagueRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return request
}

func getLeagueFromResponse(t *testing.T, body io.Reader) []Player {
	t.Helper()

	var got []Player

	err := json.NewDecoder(body).Decode(&got)
	assert.Nilf(t, err, "Unable to parse response %q into []Player", body)

	return got
}
