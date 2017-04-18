package party

import (
	"container/list"
	"fmt"
)

// PlayNextQueue is a FIFO queue implemented as a list.
type PlayNextQueue struct {
	// head of the list is the next song to be pulled off
	songs *list.List
}

// AddSong to the back of the PlayNext queue.
// Error if the song is already in the queue.
func (pnq *PlayNextQueue) AddSong(sid SongUID) error {
	if element := pnq.getSong(sid); element != nil {
		return fmt.Errorf("song already in pnq")
	}

	pnq.songs.PushBack(sid)

	return nil
}

// NewPlayNextQueue with empty list
func NewPlayNextQueue() PlayNextQueue {
	return PlayNextQueue{
		songs: list.New(),
	}
}

// Pop the top item off the queue. Error if nothing is in the queue
func (pnq *PlayNextQueue) Pop() (SongUID, error) {
	if pnq.songs.Len() == 0 {
		return "", fmt.Errorf("play next queue is empty")
	}

	// get the head, remove it, then return value
	front := pnq.songs.Front()
	val := pnq.songs.Remove(front)
	return val.(SongUID), nil
}

// SetTop song in the playnext. Error if the song is already in the queue
func (pnq *PlayNextQueue) SetTop(sid SongUID) error {
	// check if the song is already in the queue
	existingElem := pnq.getSong(sid)
	if existingElem != nil {
		pnq.songs.MoveToFront(existingElem)
		return nil
	}

	// put the new song on the top
	pnq.songs.PushFront(sid)
	return nil
}

// Pull the values in the PlayNextQueue.
// Returns the next items in play order
func (pnq PlayNextQueue) Pull() interface{} {
	if pnq.songs.Len() == 0 {
		return nil
	}

	// loop over the list, adding each elem to song
	var vals []SongUID
	elem := pnq.songs.Front()
	for elem != nil {
		vals = append(vals, elem.Value.(SongUID))

		elem = elem.Next()
	}

	return vals
}

// return element with song id, nil if no such element found
func (pnq PlayNextQueue) getSong(sid SongUID) *list.Element {

	elem := pnq.songs.Front()
	for elem != nil {
		if elem.Value.(SongUID) == sid {
			return elem
		}

		elem = elem.Next()
	}

	return nil
}
