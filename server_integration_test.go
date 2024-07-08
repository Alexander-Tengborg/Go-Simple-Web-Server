package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	DB_NAME := "dev.db"

	store, err := NewBoltPlayerStore(DB_NAME)
	assert.Nil(t, err)
	defer store.Close()

	err = store.ResetBucket()
	assert.Nil(t, err)

	server := NewPlayerServer(store)
	player := "Linda"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	player2 := "George"
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player2))

	player3 := "Ted"
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player3))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player3))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player3))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player3))

	t.Run("get score of one player", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "3", response.Body.String())
	})

	t.Run("get league from highest to lowest amount of wins", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())

		gotContentType := response.Result().Header.Get("content-type")
		got := getLeagueFromResponse(t, response.Body)
		want := []Player{
			{"Ted", 4},
			{"Linda", 3},
			{"George", 1},
		}

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "application/json", gotContentType)
		assert.Equal(t, want, got)
	})
}
