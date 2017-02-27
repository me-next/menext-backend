package queue

import (
	"fmt"
)

// BasicQueue manages a list of songs
type BasicQueue struct {
	songList map[SongUUID]Song

	// tracks the order songs were added in
	orderCounter int
}

// NewBasicQueue with an empty song list
func NewBasicQueue() Queue {
	q := BasicQueue{
		songList:     make(map[SongUUID]Song),
		orderCounter: 0,
	}
	return &q
}

// Add Song to queue
func (q *BasicQueue) Add(song Song) error {

	if _, found := q.songList[song.UUID]; found {
		return fmt.Errorf("song already in queue")
	}

	q.orderCounter++
	song.Priority = q.orderCounter
	q.songList[song.UUID] = song

	return nil
}

// Remove by UUID from teh queue
func (q *BasicQueue) Remove(songID SongUUID) error {
	if _, found := q.songList[songID]; !found {
		return fmt.Errorf("song %s not in queue", songID)
	}

	delete(q.songList, songID)

	return nil
}

// Pop the first in song
func (q *BasicQueue) Pop() (Song, error) {
	if len(q.songList) == 0 {
		return Song{}, fmt.Errorf("no songs in queue")
	}

	var oldest Song
	oldest.Priority = q.orderCounter + 1
	for _, song := range q.songList {
		if song.Priority < oldest.Priority {
			oldest = song
		}
	}

	q.Remove(oldest.UUID)
	return oldest, nil
}

// type cast checks that the BasicQueue implements Queue interface
var _ Queue = (*BasicQueue)(nil)
