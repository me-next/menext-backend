package party

import (
)

type Queue struct {
	 songList []string
}

func NewQueue() *Queue {
	q := Queue{
		songList:make([]string, 10),
	}
	return &q
}

func (q *Queue) AddSong(songID string) error {
	//TODO: check if song already added
	q.songList = append(q.songList, songID)
	return nil
}

func (q *Queue) RemoveSong(songID string) error {
	filteredList := q.songList[:0]
	for _, s := range q.songList {
		if s != songID{
				filteredList = append(filteredList, s)
		}	
	}
	q.songList = filteredList
	return nil
}
