package poker_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	poker "Go-Simple-Web-Server"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestGETPlayers(t *testing.T) {
	store := &poker.StubPlayerStore{
		map[string]int{
			"Linda": 20,
			"Steve": 4,
		},
		nil,
		nil,
	}
	server := mustMakePlayerServer(t, store, dummyGame)

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
	store := &poker.StubPlayerStore{
		map[string]int{},
		[]string{},
		nil,
	}
	server := mustMakePlayerServer(t, store, dummyGame)

	t.Run("it record wins when POST", func(t *testing.T) {
		player := "Linda"
		response := httptest.NewRecorder()
		request := newPostWinRequest(player)

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusAccepted, response.Code)

		if len(store.WinCalls) != 1 {
			t.Errorf("got %d calls to recordWin, want %d", len(store.WinCalls), 1)
		}

		if store.WinCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", store.WinCalls[0], player)
		}
	})
}

func TestLeague(t *testing.T) {
	t.Run("it returns 200 on /league", func(t *testing.T) {
		want := []poker.Player{
			{"Linda", 28},
			{"George", 38},
			{"Benedict", 2},
		}

		store := &poker.StubPlayerStore{nil, nil, want}
		server := mustMakePlayerServer(t, store, dummyGame)

		response := httptest.NewRecorder()
		request := newLeagueRequest()

		server.ServeHTTP(response, request)

		gotContentType := response.Result().Header.Get("content-type")
		got := poker.GetLeagueFromResponse(t, response.Body)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "application/json", gotContentType)
		assert.ElementsMatch(t, want, got)
	})
}

func TestGame(t *testing.T) {
	t.Run("it returns 200 on /game", func(t *testing.T) {
		store := &poker.StubPlayerStore{}
		server := mustMakePlayerServer(t, store, dummyGame)

		response := httptest.NewRecorder()
		request := newGameRequest()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Start a game with 2 players and record 'Linda' as the winner", func(t *testing.T) {
		wantedBlindAlert := "Blind is 100"
		winner := "Linda"

		game := &GameSpy{BlindAlert: []byte(wantedBlindAlert)}
		server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.Nil(t, err)
		defer ws.Close()

		err = ws.WriteMessage(websocket.TextMessage, []byte("3"))
		assert.Nil(t, err)

		err = ws.WriteMessage(websocket.TextMessage, []byte(winner))
		assert.Nil(t, err)

		time.Sleep(10 * time.Millisecond)

		assert.Equal(t, 3, game.StartedWith)
		assert.Equal(t, winner, game.FinishedWith)

		_, gotBlindAlert, _ := ws.ReadMessage()
		assert.Equal(t, wantedBlindAlert, string(gotBlindAlert))
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

func newGameRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return request
}

func mustMakePlayerServer(t *testing.T, store poker.PlayerStore, game poker.Game) *poker.PlayerServer {
	server, err := poker.NewPlayerServer(store, game)
	assert.Nil(t, err)
	return server
}
