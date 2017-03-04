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

func pull(ouid, pid string, change uint64, s *server.Server, t *testing.T) map[string]interface{} {
	resp := getHTTPResponse(fmt.Sprintf("/%s/%s/pull/%d", ouid, pid, change), s)
	assert.Equal(t, http.StatusOK, resp.Code)
	pullJson := resp.Body.String()

	assert.Equal(t, resp.Code, http.StatusOK)

	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(pullJson), &data)

	assert.Nil(t, err)

	return data
}

func TestServerPullUser(t *testing.T) {
	s := server.New()

	ouid := "1"
	oname := "bob"
	pid := createParty(ouid, oname, s, t)

	data := pull(ouid, pid, 1, s, t)
	assert.NotEmpty(t, data)

	permissions, found := data["permissions"]
	fmt.Println(data)

	assert.True(t, found)
	assert.Contains(t, permissions, "default")
	assert.Contains(t, permissions, "bad")
	assert.NotContains(t, permissions, "moo")

}
