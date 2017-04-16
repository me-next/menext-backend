package party_test

import (
	"github.com/me-next/menext-backend/party"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayNextQueue(t *testing.T) {
	q := party.NewPlayNextQueue()

	// good add
	songs := []party.SongUID{"1", "2", "3", "4"}
	for _, song := range songs {
		assert.Nil(t, q.AddSong(song))
	}

	// bad add
	assert.NotNil(t, q.AddSong("1"))

	// now pop songs off
	for _, expected := range songs {
		actual, err := q.Pop()
		assert.Equal(t, expected, actual)
		assert.Nil(t, err)
	}

	// empty pop
	_, err := q.Pop()
	assert.NotNil(t, err)

	// add to empty
	assert.Nil(t, q.AddSong("1"))
}
