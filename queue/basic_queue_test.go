package queue_test

import (
	"github.com/me-next/menext-backend/queue"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicQueue(t *testing.T) {
	// out of order in case the queue alphabetizes them
	names := []string{
		"b",
		"a",
		"d",
		"c",
	}

	q := queue.NewBasicQueue()

	for _, name := range names {
		song := queue.Song{UUID: queue.SongUUID(name)}
		assert.Nil(t, q.Add(song))
	}

	for _, name := range names {
		song, err := q.Pop()
		assert.Nil(t, err)
		assert.Equal(t, queue.SongUUID(name), song.UUID)
	}

	_, err := q.Pop()
	assert.NotNil(t, err)
}

func TestBasicQueueAdd(t *testing.T) {
	q := queue.NewBasicQueue()
	song := queue.Song{UUID: queue.SongUUID("a")}

	assert.Nil(t, q.Add(song))
	assert.NotNil(t, q.Add(song))
}

func TestBasicQueueRemove(t *testing.T) {
	// out of order in case the queue alphabetizes them
	names := []string{
		"b",
		"a",
		"d",
		"c",
	}

	oddOut := "b"

	q := queue.NewBasicQueue()

	for _, name := range names {
		song := queue.Song{UUID: queue.SongUUID(name)}
		assert.Nil(t, q.Add(song))
	}

	assert.Nil(t, q.Remove(queue.SongUUID(oddOut)))
	// check double remove
	assert.NotNil(t, q.Remove(queue.SongUUID(oddOut)))

	for _, name := range names {
		if name == oddOut {
			continue
		}
		song, err := q.Pop()
		assert.Nil(t, err)
		assert.Equal(t, queue.SongUUID(name), song.UUID)
	}

	_, err := q.Pop()
	assert.NotNil(t, err)
}
