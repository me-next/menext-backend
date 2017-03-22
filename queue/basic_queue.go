package queue

import (
	"fmt"
)

// BasicQueue manages a list of songs
type BasicQueue struct {
	songList map[SongUID]Song

	// tracks the order songs were added in
	orderCounter int
}

// NewBasicQueue with an empty song list
func NewBasicQueue() Queue {
	q := BasicQueue{
		songList:     make(map[SongUID]Song),
		orderCounter: 0,
	}

	return &q
}

// Add Song to queue
func (q *BasicQueue) Add(song Song) error {

	if _, found := q.songList[song.UID]; found {
		return fmt.Errorf("song already in queue")
	}

	q.orderCounter++
	song.Priority = q.orderCounter
	q.songList[song.UID] = song

	return nil
}

// Remove by UID from teh queue
func (q *BasicQueue) Remove(songID SongUID) error {
	if _, found := q.songList[songID]; !found {
		return fmt.Errorf("song not in queue")
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

	err := q.Remove(oldest.UID)

	return oldest, err
}

// type cast checks that the BasicQueue implements Queue interface
var _ Queue = (*BasicQueue)(nil)
