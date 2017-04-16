package party_test

import (
	"github.com/me-next/menext-backend/party"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSongInfo(t *testing.T) {
	np := &party.NowPlaying{}

	np.ChangeSong("1")

	getRaw := func(np *party.NowPlaying) map[string]interface{} {
		raw := np.Data()
		data := raw.(map[string]interface{})
		return data
	}

	// times come out in MS, but go only likes unix seconds or NS
	// need to convert from MS to time.Time
	toUnix := func(raw interface{}) time.Time {
		ms := raw.(int64)

		ns := int64(time.Millisecond) * ms / int64(time.Nanosecond)

		return time.Unix(0, ns)
	}

	// get the start time
	data := getRaw(np)
	startT := toUnix(data[party.KSongStartTimeMs])

	// wait a little
	time.Sleep(100 * time.Millisecond)

	data = getRaw(np)
	unchangedStartT := toUnix(data[party.KSongStartTimeMs])
	currentT := toUnix(data[party.KCurrentTimeMs])

	// no seek
	assert.EqualValues(t, startT, unchangedStartT)

	// sleep duration close
	assert.WithinDuration(t,
		startT.Add(100*time.Millisecond),
		currentT,
		time.Millisecond*1,
	)

	// now try a seek
	np.Seek(5)
	data = getRaw(np)
	startT = toUnix(data[party.KSongStartTimeMs])

	// sleep duration close
	assert.WithinDuration(t,
		time.Now(),
		currentT,
		time.Millisecond*1,
	)

	time.Sleep(100 * time.Millisecond)

	data = getRaw(np)
	unchangedStartT = toUnix(data[party.KSongStartTimeMs])
	currentT = toUnix(data[party.KCurrentTimeMs])

	// no seek
	assert.EqualValues(t, startT, unchangedStartT)

	// sleep duration close
	assert.WithinDuration(t,
		startT.Add(100*time.Millisecond),
		currentT,
		time.Millisecond*3,
		"sometimes the time is a little bit off on this one, thanks GC",
	)
}
