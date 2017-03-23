package party

import (
	"github.com/me-next/menext-backend/queue"
	"time"
)

// NowPlaying has info about the currently playing song.
// The client can look at changes in start time to see if a seek occurred when they pull state.
// The client can determine the current song position with the song start time,
// song pos, and current server time.
// The time diff gives a relative "time since seek'd to songPos"
type NowPlaying struct {
	nowPlaying queue.SongUUID

	// when we started the song
	startTime time.Time
	songPos   uint32
}

// consts for song info map
const (
	KSongStartTimeMs = "SongStartTimeMs"
	KCurrentTimeMs   = "CurrentTimeMs"
	KSongPosition    = "SongPos"
)

// ChangeSong changes the currently playing song.
// Always start at 0:00 when changing.
func (np *NowPlaying) ChangeSong(song queue.SongUUID) {
	np.nowPlaying = song
	np.songPos = 0
	np.startTime = time.Now()
}

// Seek to a position in the song.
// Client needs to make sure that this makes sense.
func (np *NowPlaying) Seek(pos uint32) {
	np.startTime = time.Now()
	np.songPos = pos
}

// Data returns {songStartTime, pos, currTime}.
func (np NowPlaying) Data() interface{} {
	data := make(map[string]interface{})

	toMs := func(t time.Time) int64 {
		return int64(time.Nanosecond) * t.UnixNano() / int64(time.Millisecond)
	}

	data[KSongStartTimeMs] = toMs(np.startTime)
	data[KCurrentTimeMs] = toMs(time.Now())
	data[KSongPosition] = np.songPos

	return data
}
