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
	server := &PlayerServer{store}

	t.Run("returns Linda's score", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetScoreRequest("Linda")

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, response.Body.String(), "20")
	})

	t.Run("returns Steve's score", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetScoreRequest("Steve")

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, response.Body.String(), "4")
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
	server := &PlayerServer{store}

	t.Run("it record wins when POST", func(t *testing.T) {
		player := "Linda"
		response := httptest.NewRecorder()
		request := newPostWinRequest(player)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to recordWin, want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], player)
		}
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