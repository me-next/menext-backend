package party

import (
	"container/list"
	"fmt"
)

// PreviousStack holds all of the previous songs
type PreviousStack struct {
	// head is at the bottom, so the most recently added
	// song is on the back
	songs *list.List
}

// NewPreviousStack creates a new empty previous stack
func NewPreviousStack() PreviousStack {
	return PreviousStack{
		songs: list.New(),
	}
}

// Push a song onto the stack
func (s *PreviousStack) Push(song SongUID) {
	s.songs.PushBack(song)
}

// Pop a song off of the stack
func (s *PreviousStack) Pop() (SongUID, error) {
	if s.songs.Len() == 0 {
		return "", fmt.Errorf("no previous songs")
	}

	val := s.songs.Remove(s.songs.Back())

	return val.(SongUID), nil
}
