package queue

// SongUID is a universally unique identifier for a song
type SongUID string

// Song data includes priority
type Song struct {
	Priority int
	UID     SongUID
}

// Queue is the basic interface for a queue.
// NOTE: the queue itself provides no threadsafety.
type Queue interface {
	Add(Song) error
	Remove(SongUID) error
	Pop() (Song, error)
}
