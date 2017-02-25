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

func getHTTPResponse(url string, s *server.Server) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)
	s.GetAPI().ServeHTTP(recorder, req)

	return recorder
}

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

// TODO: check coverage

func TestServerMultiUser(t *testing.T) {
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

	// add a user
	fid := "2"
	resp = getHTTPResponse(fmt.Sprintf("/%s/joinParty/%s/%s", pid, fid, "fred"), s)
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
