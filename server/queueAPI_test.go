package server_test

import (
	"fmt"
	"github.com/me-next/menext-backend/party"
	"github.com/me-next/menext-backend/server"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func (ts *testServer) suggestSong(pid server.PartyUUID, uid party.UserUUID, sid party.SongUID) error {

	resp := ts.getHTTPResponse(
		fmt.Sprintf("/suggest/%s/%s/%s", pid, uid, sid))

	// treat all the errors the same
	if resp.Code != http.StatusOK {
		return fmt.Errorf("error: %d %s", resp.Code, resp.Body.String())
	}

	return nil
}

func (ts *testServer) suggestDownvote(pid server.PartyUUID, uid party.UserUUID, sid party.SongUID) error {

	resp := ts.getHTTPResponse(
		fmt.Sprintf("/suggestDown/%s/%s/%s", pid, uid, sid))

	// treat all the errors the same
	if resp.Code != http.StatusOK {
		return fmt.Errorf("error: %d %s", resp.Code, resp.Body.String())
	}

	return nil
}

func (ts *testServer) joinEvent(pid server.PartyUUID, uid party.UserUUID, name string) error {

	resp := ts.getHTTPResponse(
		fmt.Sprintf("/%s/joinParty/%s/%s", pid, uid, name))

	// treat all the errors the same
	if resp.Code != http.StatusOK {
		return fmt.Errorf("error: %d %s", resp.Code, resp.Body.String())
	}

	return nil
}

// parses the suggestion queue info out of teh server
func parseSuggestionQueue(raw interface{}) []map[string]interface{} {
	data := raw.(map[string]interface{})

	songs := data["songs"].([]interface{})

	ret := make([]map[string]interface{}, len(songs))

	for i, val := range songs {
		ret[i] = val.(map[string]interface{})
	}

	return ret

}

func TestSuggestSongSimple(t *testing.T) {

	s := newTestServer()
	ouid := party.UserUUID("1")

	// create a party
	pid, err := s.createParty(ouid, "bob")
	assert.Nil(t, err)

	songs := []party.SongUID{"a", "b", "c"}

	for _, song := range songs {
		assert.Nil(t, s.suggestSong(pid, ouid, song))
	}

	assert.NotNil(t, s.suggestSong(pid, ouid, songs[1]))
	assert.NotNil(t, s.suggestSong(pid, ouid, songs[1]))
}

func TestSuggestUpvote(t *testing.T) {
	s := newTestServer()
	ouid := party.UserUUID("1")

	// create a party
	pid, err := s.createParty(ouid, "bob")
	assert.Nil(t, err)

	songs := []party.SongUID{"a", "b", "c"}

	// add the songs to teh queue
	for _, song := range songs {
		assert.Nil(t, s.suggestSong(pid, ouid, song))
	}

	assert.Nil(t, s.suggestDownvote(pid, ouid, songs[1]))
	fmt.Println(s.pull(ouid, pid, 1))
}
