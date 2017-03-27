package party

import (
	"fmt"
	"sort"
)

// SongUID uniquely identifies a song
type SongUID string

// Song represents the actual song data
type Song struct {
	id SongUID
}

// VotableQueue defines a queue that can be voted on
type VotableQueue struct {
	songs map[SongUID]VotableSongElement

	addCounter uint64
}

// NewVotableQueue returns a queue can can be voted on
func NewVotableQueue() VotableQueue {
	return VotableQueue{
		songs:      make(map[SongUID]VotableSongElement),
		addCounter: 0,
	}
}

// AddSong to the queue.
// Upvotes the song by default.
func (q *VotableQueue) AddSong(uid UserUUID, sid SongUID) error {
	if _, has := q.songs[sid]; has {
		return fmt.Errorf("song already in queue")
	}

	// add song to queue and move the song counter
	vse := NewVotableSongElement(q.addCounter, sid)
	vse.Upvote(uid)

	// incr counter
	q.addCounter++

	q.songs[sid] = vse

	return nil
}

// Pull the data from the queue. Use the uid to find which
// songs the user voted on. Sorts the songs.
// ret is:
// {"q":[{"id":<song ID>, "vote":<{1, 0, -1}>}]}
// where vote is 1 if the user upvoted, 0 if no vote, -1 if downvote
func (q *VotableQueue) Pull(uid UserUUID) interface{} {
	// order the songs
	arr := make([]VotableSongElement, len(q.songs))
	i := 0

	// move everything to an array
	for _, vse := range q.songs {
		arr[i] = vse
		i++
	}

	// sort the array by total votes. Break ties with order added.
	sort.SliceIsSorted(arr, func(i, j int) bool {
		a := arr[i]
		b := arr[j]

		if a.Sum() == b.Sum() {
			return a.posAdded > b.posAdded
		}

		return a.Sum() < b.Sum()
	})

	// now make an array of the data
	dataArr := make([]interface{}, len(q.songs))
	for i, vse := range arr {
		// pull only the info for this user's request
		dataArr[i] = vse.Pull(uid)
	}

	data := make(map[string]interface{})
	data["songs"] = dataArr

	return data
}

// RemoveSong from the queue.
func (q *VotableQueue) RemoveSong(sid SongUID) error {
	if _, has := q.songs[sid]; !has {
		return fmt.Errorf("song not in queue")
	}

	delete(q.songs, sid)
	return nil
}

// Upvote song by one
func (q *VotableQueue) Upvote(uid UserUUID, sid SongUID) error {
	vse, has := q.songs[sid]
	if !has {
		return fmt.Errorf("song not in queue")
	}

	// check that the up / down votes are cleaned up properly
	vse.ClearUserVotes(uid)
	vse.Upvote(uid)
	return nil
}

// Downvote song by one
func (q *VotableQueue) Downvote(uid UserUUID, sid SongUID) error {
	vse, has := q.songs[sid]
	if !has {
		return fmt.Errorf("song not in queue")
	}

	// check that the up / down votes are cleaned up properly
	vse.ClearUserVotes(uid)
	vse.Downvote(uid)
	return nil
}

// ClearVotes for a user from the song
func (q *VotableQueue) ClearVotes(uid UserUUID, sid SongUID) error {
	vse, has := q.songs[sid]
	if !has {
		return fmt.Errorf("song not in queue")
	}

	// check that the up / down votes are cleaned up properly
	vse.ClearUserVotes(uid)
	return nil
}

// Pop the top song off of the queue.
func (q *VotableQueue) Pop() (SongUID, error) {

	// check that there are songs in the queue
	if len(q.songs) == 0 {
		return "", fmt.Errorf("no songs in queue")
	}

	// find the "top" song
	topScore := 0
	var topSong VotableSongElement

	for _, song := range q.songs {

		// if == just check position
		if song.Sum() == topScore {
			if song.posAdded < topSong.posAdded {
				topSong = song
			}
		}

		if song.Sum() > topScore {
			topScore = song.Sum()
			topSong = song
		}
	}

	return topSong.songID, nil
}

// values for the song element voting
const (
	UpVoteValue   = 1
	DownVoteValue = -1
	NoVoteValue   = 0
)

// VotableSongElement defines a song that can be voted on
type VotableSongElement struct {
	votes map[UserUUID]int

	songID   SongUID
	posAdded uint64
}

// Pull the song data.
func (vse VotableSongElement) Pull(uid UserUUID) interface{} {
	data := make(map[string]interface{})

	data["id"] = vse.songID

	// check the song data
	if val, has := vse.votes[uid]; has {
		data["vote"] = val
	} else {
		data["vote"] = 0
	}

	return data
}

// NewVotableSongElement returns a song element that can be voted on
func NewVotableSongElement(pos uint64, songID SongUID) VotableSongElement {
	return VotableSongElement{
		votes:    make(map[UserUUID]int),
		songID:   songID,
		posAdded: pos,
	}
}

// ClearUserVotes from both up and down.
func (vse *VotableSongElement) ClearUserVotes(uid UserUUID) {
	if _, has := vse.votes[uid]; has {
		delete(vse.votes, uid)
	}
}

// Upvote this song.
func (vse *VotableSongElement) Upvote(uid UserUUID) {
	vse.votes[uid] = UpVoteValue
}

// Downvote this song.
func (vse *VotableSongElement) Downvote(uid UserUUID) {
	vse.votes[uid] = DownVoteValue
}

// Sum the down votes and the up votes
func (vse VotableSongElement) Sum() int {
	sum := 0
	for _, vote := range vse.votes {
		sum += vote
	}

	return sum
}
