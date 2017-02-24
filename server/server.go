package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/me-next/menext-backend/party"
	"net/http"
)

// Server for the backend.
// format of requests is <stuff to id command>/<command>/<params>.
// ie to add a song to a party queue: /partyid/userid/addsong/songid
type Server struct {
	pm *PartyManager
}

// New server
func New() *Server {
	return &Server{
		pm: NewPartyManager(),
	}
}

// just for testing, no error checking or anything
func (s *Server) sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

// CreateParty with uname and owner uuid.
// URL is /createParty/{uuid}/{uname}
func (s *Server) CreateParty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, found := vars["uid"]
	if !found {
		urlerror(w)
		return
	}

	uname, found := vars["uname"]
	if !found {
		urlerror(w)
		return
	}

	pid, err := s.pm.CreateParty(party.UserUUID(uid), uname)
	if err != nil {
		data := jsonError("failed to create party")
		w.Write(data)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// write the json
	data := map[string]string{
		"pid": pid.String(),
	}

	raw, err := json.Marshal(data)
	if err != nil {
		errMsg := jsonError("created party but failed to marshal %s", err.Error())
		w.Write(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(raw)
}

// Start the server.
func (s *Server) Start(port string) error {
	router := mux.NewRouter()
	router.Path("/hello").HandlerFunc(s.sayHello).Methods("GET")
	router.Path("/createParty/{uid}/{uname}").HandlerFunc(s.CreateParty).Methods("GET")

	// shouldn't ever return
	return http.ListenAndServe(port, router)
}

// write a urlerror to the header.
// writes the status.
func urlerror(w http.ResponseWriter) {
	data := jsonError("error, bad url")
	w.Write(data)
	w.WriteHeader(http.StatusInternalServerError)
}

// helper converts a string to bytes for writing msgs
func asbytes(vars ...interface{}) []byte {
	str := fmt.Sprint(vars...)
	return []byte(str)
}

// jsonError converts an error message to a json format
func jsonError(fmtString string, vars ...interface{}) []byte {
	data := make(map[string]string)
	data["error"] = fmt.Sprintf(fmtString, vars...)

	raw, err := json.Marshal(data)
	if err != nil {
		return []byte(fmt.Sprintf(fmtString, vars...))
	}

	return raw
}

// helper converts a formatted string to bytes for writing msgs
func asbytesf(fmtString string, vars ...interface{}) []byte {
	str := fmt.Sprintf(fmtString, vars...)
	return []byte(str)
}
