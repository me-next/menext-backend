package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/me-next/menext-backend/event"
	"net/http"
	"strconv"
)

// Server for the backend.
// format of requests is <stuff to id command>/<command>/<params>.
// ie to add a song to a event queue: /eventid/userid/addsong/songid
type Server struct {
	em *EventManager
}

// New server
func New() *Server {
	return &Server{
		em: NewEventManager(),
	}
}

// just for testing, no error checking or anything
func (s *Server) sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

// CreateEvent with uname and owner uuid.
// URL is /createEvent/{uuid}/{uname}
func (s *Server) CreateEvent(w http.ResponseWriter, r *http.Request) {
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

	eid, err := s.em.CreateEvent(event.UserUID(uid), uname)
	if err != nil {
		data := jsonError("failed to create event")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(data)

		return
	}

	// write the json
	data := map[string]string{
		"eid": string(eid),
	}

	raw, err := json.Marshal(data)
	if err != nil {

		// cleanup event
		rmvErr := s.em.Remove(eid)
		if rmvErr != nil {
			errMsg := jsonError("failed to marshal (%s)  and failed to rmv (%s)", err.Error(), rmvErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(errMsg)

			return
		}

		errMsg := jsonError("created event but failed to marshal %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	w.Write(raw)
}

// JoinEvent with owner uuid ownerName and event uuid
// url is /{eid}/joinEvent/{uuid}/{uname}
func (s *Server) JoinEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uidStr, ufound := vars["uid"]
	eidStr, efound := vars["eid"]
	uname, nfound := vars["uname"]

	if !ufound || !efound || !nfound {
		urlerror(w)
		return
	}

	eid := EventUID(eidStr)
	e, err := s.em.Event(eid)
	if err != nil {
		errMsg := jsonError("no such event %s", eid)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// verify that we can join the event
	if err = e.AddUser(event.UserUID(uidStr), uname); err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// just exit with OK status code
}

// RemoveEvent with owner uuid and event uuid
// url is /{uid}/{eid}/removeEvent
func (s *Server) RemoveEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uidStr, ufound := vars["uid"]
	eidStr, efound := vars["eid"]

	if !ufound || !efound {
		urlerror(w)
		return
	}

	// need to parse the eid to a uuid
	eid := EventUID(eidStr)

	e, err := s.em.Event(eid)
	if err != nil {
		errMsg := jsonError("no such event %s", eid)
		// when you write header after writting msg the header doesn't get written
		// TODO: figure out why this happens
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// verify that we can finish the event
	// NOTE: should we have some locking after this point
	if canEnd := e.CanUserEndEvent(event.UserUID(uidStr)); !canEnd {
		errMsg := jsonError("user can not end event")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// assume we can end if we got here
	err = s.em.Remove(eid)
	if err != nil {
		// this is super duper hokey b/c we said we could end but now we can't
		// suposedly no one else can hop in to finish this
		// TODO: migth want to think about putting a diff status code / logging / panicing here
		errMsg := jsonError("failed to remove event we should be allowed to end... very very bad")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// just exit with OK status code
}

// Pull all of the data for the client if there is a recent change. This is the most frequent getter.
// URL is /{eid}/{uid}/pull/{cid}
func (s *Server) Pull(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uidStr, ufound := vars["uid"]
	eidStr, efound := vars["eid"]
	cidStr, cfound := vars["cid"]

	if !ufound || !efound || !cfound {
		urlerror(w)
		return
	}

	eid := EventUID(eidStr)

	e, err := s.em.Event(eid)
	if err != nil {
		errMsg := jsonError("no such event %s", eid)
		// when you write header after writting msg the header doesn't get written
		// TODO: figure out why this happens
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// base 10, want a u64
	cid, err := strconv.ParseUint(cidStr, 10, 64)
	if err != nil {
		errMsg := jsonError("failed to parse changeID")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// need to get specifics for the user
	data, err := e.Pull(event.UserUID(uidStr), cid)
	if err != nil {
		errMsg := jsonError("err pulling from event:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// nil if nothing changed
	if data == nil {
		return
	}

	raw, err := json.Marshal(data)
	if err != nil {
		errMsg := jsonError("failed to serialize")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// write and exit
	w.Write(raw)
}

// GetAPI provides the server router. This is broken off from Start to make testing easier.
func (s *Server) GetAPI() http.Handler {
	router := mux.NewRouter()
	router.Path("/hello").HandlerFunc(s.sayHello).Methods("GET")
	router.Path("/createEvent/{uid}/{uname}").HandlerFunc(s.CreateEvent).Methods("GET")
	router.Path("/{uid}/{eid}/removeEvent").HandlerFunc(s.RemoveEvent).Methods("GET")
	router.Path("/{uid}/{eid}/pull/{cid}").HandlerFunc(s.Pull).Methods("GET")
	router.Path("/{eid}/joinEvent/{uid}/{uname}").HandlerFunc(s.JoinEvent).Methods("GET")

	return router
}

// Start the server.
func (s *Server) Start(port string) error {

	router := s.GetAPI()
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
