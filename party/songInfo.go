package party

import (
	"time"
)

// NowPlaying has info about the currently playing song.
// The client can look at changes in start time to see if a seek occurred when they pull state.
// The client can determine the current song position with the song start time,
// song pos, and current server time.
// The time diff gives a relative "time since seek'd to songPos"
type NowPlaying struct {
	nowPlaying SongUID

	// when we started the song
	startTime time.Time
	songPos   uint32
}

// CurrentlyPlaying checks if there is a song currently playing
func (np *NowPlaying) CurrentlyPlaying() bool {
	return np.nowPlaying != ""
}

// SetNonePlaying indicates that no song is playing.
func (np *NowPlaying) SetNonePlaying() {
	np.nowPlaying = ""
}

// ChangeSong changes the currently playing song.
// Always start at 0:00 when changing.
func (np *NowPlaying) ChangeSong(song SongUID) {
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

// consts for song info map
const (
	KSongStartTimeMs = "SongStartTimeMs"
	KCurrentTimeMs   = "CurrentTimeMs"
	KSongPosition    = "SongPos"
	KCurrentSongID   = "CurrentSongId"
	KHasSong         = "HasSong"
)

// Data returns {songStartTime, pos, currTime}.
func (np NowPlaying) Data() interface{} {
	data := make(map[string]interface{})

	toMs := func(t time.Time) int64 {
		return int64(time.Nanosecond) * t.UnixNano() / int64(time.Millisecond)
	}

	if np.CurrentlyPlaying() {
		data[KSongStartTimeMs] = toMs(np.startTime)
		data[KCurrentTimeMs] = toMs(time.Now())
		data[KSongPosition] = np.songPos
		data[KCurrentSongID] = np.nowPlaying
		data[KHasSong] = true
	} else {
		data[KHasSong] = false
	}

	return data
}
