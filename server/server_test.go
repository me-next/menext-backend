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

// parses the PID out of a string
func parsePID(str string) (string, error) {
	data := make(map[string]string)

	err := json.Unmarshal([]byte(str), &data)
	pid, found := data["pid"]
	if err != nil || !found {
		return "", err
	}

	return pid, nil
}

func TestServerSingleUser(t *testing.T) {
	s := server.New()

	ownerName := "bob"
	ouid := "1"

	// create a party
	resp := getHTTPResponse(fmt.Sprintf("/createParty/%s/%s", ouid, ownerName), s)
	assert.Equal(t, http.StatusOK, resp.Code)
	pidJson := resp.Body.String()

	pid, err := parsePID(pidJson)
	assert.Nil(t, err)
	assert.NotEqual(t, "", pid)

	// do a bad remove
	resp = getHTTPResponse(fmt.Sprintf("/%s/%s/removeParty", "2", pid), s)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.NotEqual(t, "", resp.Body.String())

	// remove the party
	resp = getHTTPResponse(fmt.Sprintf("/%s/%s/removeParty", ouid, pid), s)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "", resp.Body.String())

	// check that we can't still remove party
	resp = getHTTPResponse(fmt.Sprintf("/%s/%s/removeParty", ouid, pid), s)
	assert.Equal(t, 500, resp.Code)
	assert.NotEqual(t, "", resp.Body.String())
}

// helper function for creating a party on the server
// needs the owner name and id, plus the server and testing objects
// parses the pid
func createParty(ouid, oname string, s *server.Server, t *testing.T) string {
	resp := getHTTPResponse(fmt.Sprintf("/createParty/%s/%s", ouid, oname), s)
	assert.Equal(t, http.StatusOK, resp.Code)
	pidJson := resp.Body.String()

	pid, err := parsePID(pidJson)
	assert.Nil(t, err)
	assert.NotEqual(t, "", pid)

	return pid
}

func TestServerMultiUser(t *testing.T) {
	s := server.New()

	ownerName := "bob"
	ouid := "1"
	pid := createParty(ouid, ownerName, s, t)

	// create a party

	// add a user
	fid := "2"
	resp := getHTTPResponse(fmt.Sprintf("/%s/joinParty/%s/%s", pid, fid, "fred"), s)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "", resp.Body.String())

	// double add the user
	resp = getHTTPResponse(fmt.Sprintf("/%s/joinParty/%s/%s", pid, fid, "fred"), s)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.NotEqual(t, "", resp.Body.String())

	// have the user violate permissions
	resp = getHTTPResponse(fmt.Sprintf("/%s/%s/removeParty", fid, pid), s)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.NotEqual(t, "", resp.Body.String())

	// remove the party
	resp = getHTTPResponse(fmt.Sprintf("/%s/%s/removeParty", ouid, pid), s)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "", resp.Body.String())
}

// helper function to pull from the client
// need to provide the uid, event id, change id plus the server and testing objects
// expects that the pulls are good, and will test as such
func pull(ouid, pid string, change uint64, s *server.Server, t *testing.T) map[string]interface{} {
	resp := getHTTPResponse(fmt.Sprintf("/%s/%s/pull/%d", ouid, pid, change), s)

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
	pid := createParty(ouid, oname, s, t)

	// should be empty
	data := pull(ouid, pid, 0, s, t)
	assert.Empty(t, data)

	// check a bad pull
	resp := getHTTPResponse(fmt.Sprintf("/%s/%s/pull/%d", ouid, pid, 1), s)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	assert.Contains(t, resp.Body.String(), "bad pull id")
}
