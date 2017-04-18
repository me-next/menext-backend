package party

import (
	"fmt"
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

	// range [0, 100]
	volume uint32

	playing bool
}

// CurrentlyHasSong checks if there is a song currently playing
func (np *NowPlaying) CurrentlyHasSong() bool {
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

	// make sure that the song doesn't get paused
	np.playing = true
}

// Seek to a position in the song.
// Client needs to make sure that this makes sense.
func (np *NowPlaying) Seek(pos uint32) {
	np.startTime = time.Now()
	np.songPos = pos
}

// SetVolume to level in range [0, 100]
func (np *NowPlaying) SetVolume(level uint32) error {
	if level > 100 {
		return fmt.Errorf("bad volume")
	}

	np.volume = level
	return nil
}

// SetPaused pauses if not already paused.
// Need to provide the position of the pause
func (np *NowPlaying) SetPaused(pos uint32) error {
	if !np.playing {
		return fmt.Errorf("already paused")
	}

	np.songPos = pos
	np.playing = false
	return nil
}

// SetPlaying plays the song if not already playing
func (np *NowPlaying) SetPlaying() error {
	if np.playing {
		return fmt.Errorf("already playing")
	}

	// need to update time
	np.startTime = time.Now()

	np.playing = true
	return nil
}

// consts for song info map
const (
	KSongStartTimeMs = "SongStartTimeMs"
	KCurrentTimeMs   = "CurrentTimeMs"
	KSongPosition    = "SongPos"
	KCurrentSongID   = "CurrentSongId"
	KHasSong         = "HasSong"
	KVolume          = "Volume"
	KPlaying         = "Playing"
)

// Data returns {songStartTime, pos, currTime}.
func (np NowPlaying) Data() interface{} {
	data := make(map[string]interface{})

	toMs := func(t time.Time) int64 {
		return int64(time.Nanosecond) * t.UnixNano() / int64(time.Millisecond)
	}

	if np.CurrentlyHasSong() {
		data[KSongStartTimeMs] = toMs(np.startTime)
		data[KCurrentTimeMs] = toMs(time.Now())
		data[KSongPosition] = np.songPos
		data[KCurrentSongID] = np.nowPlaying
		data[KHasSong] = true
		data[KVolume] = np.volume
		data[KPlaying] = np.playing
	} else {
		data[KHasSong] = false
	}

	return data
}
