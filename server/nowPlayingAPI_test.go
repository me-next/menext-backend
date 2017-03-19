package server_test

import (
	"encoding/json"
	"fmt"
	"github.com/me-next/menext-backend/party"
	"github.com/me-next/menext-backend/server"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// testServer provides helpful testing facilities for the server
type testServer struct {
	s *server.Server
}

func (ts *testServer) getHTTPResponse(url string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)
	ts.s.GetAPI().ServeHTTP(recorder, req)

	return recorder
}

func (ts *testServer) createParty(ouid party.UserUUID, oname string) (server.PartyUUID, error) {

	resp := ts.getHTTPResponse(fmt.Sprintf("/createParty/%s/%s",
		ouid, oname))

	// treat all the errors the same
	if resp.Code != http.StatusOK {
		return "", fmt.Errorf("error: %s", resp.Body.String())
	}

	str := resp.Body.String()
	data := make(map[string]string)

	err := json.Unmarshal([]byte(str), &data)
	pid, found := data["pid"]
	if err != nil || !found {
		return "", err
	}

	return server.PartyUUID(pid), nil
}

func (ts *testServer) pull(ouid party.UserUUID,
	pid server.PartyUUID, cid uint64) (map[string]interface{}, error) {
	resp := ts.getHTTPResponse(
		fmt.Sprintf("/%s/%s/pull/%d", ouid, pid, cid))

	// check response
	if resp.Code != http.StatusOK {
		return nil, fmt.Errorf("error: %s", resp.Body.String())
	}

	// parse the json
	pullJson := resp.Body.String()

	data := make(map[string]interface{})
	if resp.Body.Len() == 0 {
		// if the body is empty then return 0
		return data, nil
	}

	// should always be able to unmarshal correctly
	err := json.Unmarshal([]byte(pullJson), &data)

	return data, err
}

func (ts *testServer) seek(ouid party.UserUUID,
	pid server.PartyUUID, pos uint32) error {

	resp := ts.getHTTPResponse(
		fmt.Sprintf("/%s/%s/seek/%d", pid, ouid, pos))

	// check response
	if resp.Code != http.StatusOK {
		return fmt.Errorf("error: %s", resp.Body.String())
	}

	return nil
}

func newTestServer() *testServer {
	return &testServer{
		s: server.New(),
	}
}

// extracts song position out of pull data
func extractPos(raw interface{}) (uint32, error) {

	pullData := raw.(map[string]interface{})
	changeData := pullData[party.PullPlayingKey].(map[string]interface{})

	// json will make this f64
	posf := changeData[party.KSongPosition].(float64)
	pos := uint32(posf)

	return pos, nil
}

func TestSeek(t *testing.T) {
	s := newTestServer()
	ouid := party.UserUUID("1")

	pid, err := s.createParty(ouid, "bob")
	assert.Nil(t, err)
	assert.NotEmpty(t, pid)

	// check good seek for sanity
	var seekTo uint32 = 5
	err = s.seek(ouid, pid, seekTo)
	assert.Nil(t, err)

	data, err := s.pull(ouid, pid, 0)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.NotEmpty(t, data)

	// pull to check state
	pos, err := extractPos(data)
	assert.Nil(t, err)
	assert.EqualValues(t, seekTo, pos)

	// bad party
	err = s.seek(ouid, server.PartyUUID("bad"), 50)
	assert.NotNil(t, err)

	// bad user
	err = s.seek("bad", pid, 50)
	assert.NotNil(t, err)

	// check that bad requests didn't change anything
	data, err = s.pull(ouid, pid, 1)
	assert.Nil(t, err)
	assert.Empty(t, data)
}
