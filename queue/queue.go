package queue

import (
	"fmt"
)

type SongUUID string

type Song struct {
	Priority int
}

func NewSong(p int) Song{
	return Song{Priority : p}
}

type Queue struct {
	songList map[SongUUID]Song
}

func NewQueue() *Queue {
	q := Queue{
		songList: make(map[SongUUID]Song),
	}
	return &q
}

func (q *Queue) AddSong(songID SongUUID) error {
	_, found := q.songList[songID]
	if !found {
		for _, s := range q.songList{
			s.Priority ++
		}
		q.songList[songID] = NewSong(1)
		return nil
	}
	return fmt.Errorf("song %s already in queue", songID)
}

func (q *Queue) RemoveSong(songID SongUUID) error {
	_, found := q.songList[songID]
	if !found {
		return fmt.Errorf("song %s not in queue", songID)
	}
	q.decreaseAllAbove(songID)
	delete(q.songList, songID)
	return nil
	
}

func(s Song) increasePriority(){
	s.Priority++ 
}

func(s Song) decreasePriority(){
	s.Priority--
}

func (q *Queue) increaseAllAbove(s SongUUID){
	song := q.songList[s]
	p := song.Priority
	for _,i := range q.songList{
		if p < i.Priority {
			i.increasePriority()
		}
	}
}

func (q *Queue) decreaseAllAbove(s SongUUID){
	song := q.songList[s]
	p := song.Priority
	for _,i := range q.songList{
		if p < i.Priority {
			i.decreasePriority()
		}
	}
}