package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(player string) int {
	return s.scores[player]
}

func (s *StubPlayerStore) RecordWin(player string) {
	s.winCalls = append(s.winCalls, player)
}

func TestGETPlayers(t *testing.T) {
	store := &StubPlayerStore{
		map[string]int{
			"Linda": 20,
			"Steve": 4,
		},
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
	store := &StubPlayerStore{}
	server := NewPlayerServer(store)

	t.Run("it returns 200 on /league", func(t *testing.T) {
		response := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
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
