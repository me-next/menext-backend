package queue

// SongUUID is a universally unique identifier for a song
type SongUUID string

// Song data includes priority
type Song struct {
	Priority int
	UUID     SongUUID
}

// Queue is the basic interface for a queue.
// NOTE: the queue its self provides no threadsafety
type Queue interface {
	Add(Song) error
	Remove(SongUUID) error
	Pop() (Song, error)
}
