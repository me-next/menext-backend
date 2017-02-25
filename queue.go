package queue

import (
	"fmt"
)

type Song string

type Queue struct {
	 songList map[Song]int
}

func NewQueue() *Queue {
	q := Queue{
		songList: make(map[Song]int),
	}
	return &q
}

func (q *Queue) AddSong(songID Song) error {
	 song := q.songList[songID]
	if song == 0 {
		q.songList[songID] = len(q.songList) + 1
		return nil
	}
	
	return	fmt.Errorf("song %s already in queue", songID)
}

func (q *Queue) RemoveSong(songID Song) error {
	if q.songList[songID] != 0{
		delete(q.songList, songID)
		return nil
	}	
	return fmt.Errorf("%s not in queue", songID)
}