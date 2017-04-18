package party_test

// tests specially for playing
import (
	"github.com/me-next/menext-backend/party"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getCurrentlyPlaying(p *party.Party, ouid party.UserUUID) (
	party.SongUID, error) {
	raw, err := p.Pull(ouid, 0)
	if err != nil {
		return "", err
	}

	// convert to map
	data := raw.(map[string]interface{})

	nowPlayingRaw := data[party.PullPlayingKey]
	nowPlayingData := nowPlayingRaw.(map[string]interface{})

	currentSong := nowPlayingData[party.KCurrentSongID].(party.SongUID)
	return currentSong, nil
}

func TestPlaySongOrderCorrect(t *testing.T) {
	// behavior is to pull from play-next first then try suggestion
	// check that it acts this way
	ouid := party.UserUUID("1")

	p := party.New(ouid, "bob")

	assert.Nil(t, p.PlayNext(ouid, "a"))
	assert.Nil(t, p.PlayNext(ouid, "b"))
	assert.Nil(t, p.Suggest(ouid, "c")) // should get played last
	assert.Nil(t, p.PlayNext(ouid, "d"))
	assert.Nil(t, p.Suggest(ouid, "e"))

	expecteds := []party.SongUID{"a", "b", "d", "c"}

	for _, expected := range expecteds {
		actual, err := getCurrentlyPlaying(p, ouid)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)

		// now skip
		assert.Nil(t, p.SongFinished(ouid, actual))
	}

	actual, err := getCurrentlyPlaying(p, ouid)
	assert.Nil(t, err)
	assert.Equal(t, party.SongUID("e"), actual)
}
