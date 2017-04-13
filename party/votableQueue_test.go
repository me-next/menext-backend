package party_test

import (
	"github.com/me-next/menext-backend/party"
	"github.com/stretchr/testify/assert"
	"testing"
)

func parseSongsFromVQPull(raw interface{}) []map[string]interface{} {
	data := raw.(map[string]interface{})

	songs := data["songs"].([]interface{})

	ret := make([]map[string]interface{}, len(songs))

	for i, val := range songs {
		ret[i] = val.(map[string]interface{})
	}

	return ret
}

func TestVotableQueueSimple(t *testing.T) {
	q := party.NewVotableQueue()

	// try adding songs and pulling
	songs := []party.SongUID{"a", "b", "c"}

	for _, song := range songs {
		assert.Nil(t, q.AddSong("1", song))
	}

	// try re-adding
	assert.NotNil(t, q.AddSong("2", songs[1]))

	// check that the songs are all there
	rawPull := q.Pull("1")
	data := parseSongsFromVQPull(rawPull)
	assert.Len(t, data, len(songs))

	for i, song := range data {
		assert.Equal(t, songs[i], song["id"])
		assert.Equal(t, song["vote"], 1)
	}
}

func TestVotableQueueSort(t *testing.T) {
	// check that the queue seems to sort properly
	songs := []party.SongUID{"a", "b", "c", "d", "e"}

	q := party.NewVotableQueue()
	for _, song := range songs {
		assert.Nil(t, q.AddSong("1", song))
	}

	assert.Nil(t, q.Upvote("1", "d"))
	assert.Nil(t, q.Upvote("2", "d"))
	assert.Nil(t, q.Downvote("2", "b"))

	expectedOrder := []party.SongUID{"d", "a", "c", "e", "b"}
	for _, song := range expectedOrder {
		actual, err := q.Pop()
		assert.Nil(t, err)
		assert.Equal(t, song, actual)
	}
}
