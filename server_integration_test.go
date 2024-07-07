package main

import (
	"Go-Simple-Web-Server/stores"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	DB_NAME := "dev.db"

	store, err := stores.NewBoltPlayerStore(DB_NAME)
	assert.Nil(t, err)
	defer store.Close()

	err = store.ResetBucket()
	assert.Nil(t, err)

	server := PlayerServer{store}
	player := "Linda"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "3", response.Body.String())
}
