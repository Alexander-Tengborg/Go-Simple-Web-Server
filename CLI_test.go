package poker_test

import (
	"strings"
	"testing"

	poker "Go-Simple-Web-Server"

	"github.com/stretchr/testify/assert"
)

func TestCLI(t *testing.T) {
	t.Run("record Linda win from user input", func(t *testing.T) {
		in := strings.NewReader("Linda wins\n")
		store := &poker.StubPlayerStore{}

		cli := poker.NewCLI(store, in)
		cli.PlayPoker()

		assert.Equal(t, 1, len(store.WinCalls))

		got := store.WinCalls[0]
		want := "Linda"
		assert.Equal(t, want, got)
	})

	t.Run("record Ted win from user input", func(t *testing.T) {
		in := strings.NewReader("Ted wins\n")
		store := &poker.StubPlayerStore{}

		cli := poker.NewCLI(store, in)
		cli.PlayPoker()

		assert.Equal(t, 1, len(store.WinCalls))

		got := store.WinCalls[0]
		want := "Ted"
		assert.Equal(t, want, got)
	})
}
