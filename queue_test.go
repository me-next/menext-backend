package queue_test

import (
	//"github.com/me-next/menext-backend/queue"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSongAdd(t *testing.T) {
	q := queue.NewQueue()
	
	//check add song
	assert.Nil(t, q.AddSong("Song 1"))

	// can't double add songs
	assert.NotNil(t, q.AddSong("Song 1"))

	// check remove songs
	assert.Nil(t, q.RemoveSong("Song 1"))
	
	// check that we can't remove  song that isn't there
	assert.NotNil(t, q.RemoveSong("Not Song"))
}