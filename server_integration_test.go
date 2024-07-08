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

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "3", response.Body.String())
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())

		gotContentType := response.Result().Header.Get("content-type")
		got := getLeagueFromResponse(t, response.Body)
		want := []Player{
			{"Linda", 3},
		}

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "application/json", gotContentType)
		assert.ElementsMatch(t, want, got)
	})
}
