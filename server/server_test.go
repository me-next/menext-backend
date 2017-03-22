package server_test

import (
	"encoding/json"
	"fmt"
	"github.com/me-next/menext-backend/server"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// fires a request using the response recorder
func getHTTPResponse(url string, s *server.Server) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)
	s.GetAPI().ServeHTTP(recorder, req)

	return recorder
}

// parses the EID out of a string
func parseEID(str string) (string, error) {
	data := make(map[string]string)

	err := json.Unmarshal([]byte(str), &data)
	eid, found := data["eid"]
	if err != nil || !found {
		return "", err
	}

	return eid, nil
}

func TestServerSingleUser(t *testing.T) {
	s := server.New()

	ownerName := "bob"
	ouid := "1"

	// create a event
	resp := getHTTPResponse(fmt.Sprintf("/createEvent/%s/%s", ouid, ownerName), s)
	assert.Equal(t, http.StatusOK, resp.Code)
	eidJson := resp.Body.String()

	eid, err := parseEID(eidJson)
	assert.Nil(t, err)
	assert.NotEqual(t, "", eid)

	// do a bad remove
	resp = getHTTPResponse(fmt.Sprintf("/%s/%s/removeEvent", "2", eid), s)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.NotEqual(t, "", resp.Body.String())

	// remove the event
	resp = getHTTPResponse(fmt.Sprintf("/%s/%s/removeEvent", ouid, eid), s)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "", resp.Body.String())

	// check that we can't still remove event
	resp = getHTTPResponse(fmt.Sprintf("/%s/%s/removeEvent", ouid, eid), s)
	assert.Equal(t, 500, resp.Code)
	assert.NotEqual(t, "", resp.Body.String())
}

// helper function for creating a event on the server
// needs the owner name and id, plus the server and testing objects
// parses the eid
func createEvent(ouid, oname string, s *server.Server, t *testing.T) string {
	resp := getHTTPResponse(fmt.Sprintf("/createEvent/%s/%s", ouid, oname), s)
	assert.Equal(t, http.StatusOK, resp.Code)
	eidJson := resp.Body.String()

	eid, err := parseEID(eidJson)
	assert.Nil(t, err)
	assert.NotEqual(t, "", eid)

	return eid
}

func TestServerMultiUser(t *testing.T) {
	s := server.New()

	ownerName := "bob"
	ouid := "1"
	eid := createEvent(ouid, ownerName, s, t)

	// create a event

	// add a user
	fid := "2"
	resp := getHTTPResponse(fmt.Sprintf("/%s/joinEvent/%s/%s", eid, fid, "fred"), s)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "", resp.Body.String())

	// double add the user
	resp = getHTTPResponse(fmt.Sprintf("/%s/joinEvent/%s/%s", eid, fid, "fred"), s)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.NotEqual(t, "", resp.Body.String())

	// have the user violate permissions
	resp = getHTTPResponse(fmt.Sprintf("/%s/%s/removeEvent", fid, eid), s)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.NotEqual(t, "", resp.Body.String())

	// remove the event
	resp = getHTTPResponse(fmt.Sprintf("/%s/%s/removeEvent", ouid, eid), s)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "", resp.Body.String())
}

// helper function to pull from the client
// need to provide the uid, event id, change id plus the server and testing objects
// expects that the pulls are good, and will test as such
func pull(ouid, eid string, change uint64, s *server.Server, t *testing.T) map[string]interface{} {
	resp := getHTTPResponse(fmt.Sprintf("/%s/%s/pull/%d", ouid, eid, change), s)

	// check response
	assert.Equal(t, http.StatusOK, resp.Code)
	pullJson := resp.Body.String()

	data := make(map[string]interface{})
	if resp.Body.Len() == 0 {
		// if the body is empty then return 0
		return data
	}

	// should always be able to unmarshal correctly
	err := json.Unmarshal([]byte(pullJson), &data)
	assert.Nil(t, err)

	return data
}

func TestServerPullUser(t *testing.T) {
	s := server.New()

	ouid := "1"
	oname := "bob"
	eid := createEvent(ouid, oname, s, t)

	// should be empty
	data := pull(ouid, eid, 0, s, t)
	assert.Empty(t, data)

	// check a bad pull
	resp := getHTTPResponse(fmt.Sprintf("/%s/%s/pull/%d", ouid, eid, 1), s)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	assert.Contains(t, resp.Body.String(), "bad pull id")
}
