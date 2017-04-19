package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/me-next/menext-backend/party"
	"net/http"
	"strconv"
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

// CreatePartyWithName allows a user to create a party with a custom name.
// If the event name is taken, suggests an alternate name.
// Path is: /createWithName/{uid}/{uname}/{pid}
func (s *Server) CreatePartyWithName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uidStr, ufound := vars["uid"]
	pidStr, pfound := vars["pid"]
	uname, nfound := vars["uname"]

	if !ufound || !pfound || !nfound {
		urlerror(w)
		return
	}

	pid, alt, err := s.pm.CreatePartyWithName(party.UserUUID(uidStr), uname, pidStr)
	if err != nil {
		// need to return alt

		data := map[string]interface{}{
			"error":       err.Error(),
			"alternative": alt,
		}

		// try to mashal
		raw, err := json.Marshal(data)
		if err != nil {
			errMsg := jsonError("created alternative but failed to marshal %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(errMsg)
			return
		}

		// just write and return
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(raw)

		return
	} // else created the party

	// write the json
	data := map[string]string{
		"pid": string(pid),
	}

	raw, err := json.Marshal(data)
	if err != nil {

		// cleanup party
		rmvErr := s.pm.Remove(pid)
		if rmvErr != nil {
			errMsg := jsonError("failed to marshal (%s)  and failed to rmv (%s)", err.Error(), rmvErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(errMsg)

			return
		}

		errMsg := jsonError("created party but failed to marshal %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	w.Write(raw)
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(data)

		return
	}

	// write the json
	data := map[string]string{
		"pid": string(pid),
	}

	raw, err := json.Marshal(data)
	if err != nil {

		// cleanup party
		rmvErr := s.pm.Remove(pid)
		if rmvErr != nil {
			errMsg := jsonError("failed to marshal (%s)  and failed to rmv (%s)", err.Error(), rmvErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(errMsg)

			return
		}

		errMsg := jsonError("created party but failed to marshal %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	w.Write(raw)
}

// JoinParty with owner uuid ownerName and party uuid
// url is /{pid}/joinParty/{uuid}/{uname}
func (s *Server) JoinParty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uidStr, ufound := vars["uid"]
	pidStr, pfound := vars["pid"]
	uname, nfound := vars["uname"]

	if !ufound || !pfound || !nfound {
		urlerror(w)
		return
	}

	pid := PartyUUID(pidStr)
	p, err := s.pm.Party(pid)
	if err != nil {
		errMsg := jsonError("no such party %s", pid)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// verify that we can join the party
	if err = p.AddUser(party.UserUUID(uidStr), uname); err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// just exit with OK status code
}

// LeaveParty removes the user from the party
// url is /leaveParty/{pid}/{uid}
func (s *Server) LeaveParty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uidStr, ufound := vars["uid"]
	pidStr, pfound := vars["pid"]

	if !ufound || !pfound {
		urlerror(w)
		return
	}

	pid := PartyUUID(pidStr)
	p, err := s.pm.Party(pid)
	if err != nil {
		errMsg := jsonError("no such party %s", pid)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// try to leave the party
	if err = p.RemoveUser(party.UserUUID(uidStr)); err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// just exit with OK status code
}

// RemoveParty with owner uuid and party uuid
// url is /{uid}/{pid}/removeParty
func (s *Server) RemoveParty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uidStr, ufound := vars["uid"]
	pidStr, pfound := vars["pid"]

	if !ufound || !pfound {
		urlerror(w)
		return
	}

	// need to parse the pid to a uuid
	pid := PartyUUID(pidStr)

	p, err := s.pm.Party(pid)
	if err != nil {
		errMsg := jsonError("no such party %s", pid)
		// when you write header after writting msg the header doesn't get written
		// TODO: figure out why this happens
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// verify that we can finish the party
	// NOTE: should we have some locking after this point
	if canEnd := p.CanUserEndParty(party.UserUUID(uidStr)); !canEnd {
		errMsg := jsonError("user can not end party")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// assume we can end if we got here
	err = s.pm.Remove(pid)
	if err != nil {
		// this is super duper hokey b/c we said we could end but now we can't
		// suposedly no one else can hop in to finish this
		// TODO: migth want to think about putting a diff status code / logging / panicing here
		errMsg := jsonError("failed to remove party we should be allowed to end... very very bad")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// just exit with OK status code
}

// Pull all of the data for the client if there is a recent change. This is the most frequent getter.
// URL is /{pid}/{uid}/pull/{cid}
func (s *Server) Pull(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uidStr, ufound := vars["uid"]
	pidStr, pfound := vars["pid"]
	cidStr, cfound := vars["cid"]

	if !ufound || !pfound || !cfound {
		urlerror(w)
		return
	}

	pid := PartyUUID(pidStr)

	p, err := s.pm.Party(pid)
	if err != nil {
		errMsg := jsonError("no such party %s", pid)
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
	data, err := p.Pull(party.UserUUID(uidStr), cid)
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

// Permissions returns a map of permission keys to descriptions
// Path is /permissions/
func (s *Server) Permissions(w http.ResponseWriter, r *http.Request) {
	// get the permissions map
	data := party.GetPermissionDescriptions()

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

// SetPermissions for all users in a party.
// path is /setPermission/{pid}/{uid}/{permKey}/{val}.
// val == "true" when trying to set to true, otherwise "false"
func (s *Server) SetPermissions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pidStr, pfound := vars["pid"]
	uidStr, ufound := vars["uid"]
	permStr, permFound := vars["perm"]
	valStr, valFound := vars["val"]

	if !ufound || !pfound || !permFound || !valFound {
		urlerror(w)
		return
	}

	// get the party
	pid := PartyUUID(pidStr)

	p, err := s.pm.Party(pid)
	if err != nil {
		errMsg := jsonError("no such party %s", pid)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// try to set the permission
	err = p.SetPermission(permStr, valStr == "true", party.UserUUID(uidStr))
	if err != nil {
		errMsg := jsonError(fmt.Sprintf("error setting permission: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// exit with OK status
}

// GetAPI provides the server router. This is broken off from Start to make testing easier.
func (s *Server) GetAPI() http.Handler {
	router := mux.NewRouter()

	// debugging endpoint
	router.Path("/hello").HandlerFunc(s.sayHello).Methods("GET")

	// general party management
	router.Path("/createParty/{uid}/{uname}").HandlerFunc(s.CreateParty).Methods("GET")
	router.Path("/createPartyWithName/{uid}/{uname}/{pid}").HandlerFunc(s.CreatePartyWithName).Methods("GET")
	router.Path("/removeParty/{uid}/{pid}").HandlerFunc(s.RemoveParty).Methods("GET")
	router.Path("/pull/{uid}/{pid}/{cid}").HandlerFunc(s.Pull).Methods("GET")
	router.Path("/joinParty/{pid}/{uid}/{uname}").HandlerFunc(s.JoinParty).Methods("GET")
	router.Path("/leaveParty/{pid}/{uid}").HandlerFunc(s.LeaveParty).Methods("GET")

	// permissions
	router.Path("/permissions").HandlerFunc(s.Permissions).Methods("GET")
	router.Path("/setPermission/{pid}/{uid}/{perm}/{val}").HandlerFunc(s.Permissions).Methods("GET")

	// nowPlaying
	router.Path("/seek/{pid}/{uid}/{pos}").HandlerFunc(s.Seek).Methods("GET")
	router.Path("/songFinished/{pid}/{uid}/{sid}").HandlerFunc(s.SongFinished).Methods("GET")
	router.Path("/setVolume/{pid}/{uid}/{volume}").HandlerFunc(s.SetVolume).Methods("GET")
	router.Path("/play/{pid}/{uid}").HandlerFunc(s.Play).Methods("GET")
	router.Path("/pause/{pid}/{uid}/{pos}").HandlerFunc(s.Pause).Methods("GET")

	router.Path("/skip/{pid}/{uid}/{sid}").HandlerFunc(s.Skip).Methods("GET")
	router.Path("/previous/{pid}/{uid}/{sid}").HandlerFunc(s.Previous).Methods("GET")
	router.Path("/playNow/{pid}/{uid}/{sid}").HandlerFunc(s.PlayNow).Methods("GET")

	// queues
	router.Path("/suggest/{pid}/{uid}/{sid}").HandlerFunc(s.Suggest).Methods("GET")
	router.Path("/suggestDown/{pid}/{uid}/{sid}").HandlerFunc(s.SuggestionDownvote).Methods("GET")
	router.Path("/suggestUp/{pid}/{uid}/{sid}").HandlerFunc(s.SuggestionUpvote).Methods("GET")
	router.Path("/suggestClearvote/{pid}/{uid}/{sid}").HandlerFunc(s.SuggestionClearvote).Methods("GET")

	router.Path("/addPlayNext/{pid}/{uid}/{sid}").HandlerFunc(s.AddPlayNext).Methods("GET")
	router.Path("/addTopPlayNext/{pid}/{uid}/{sid}").HandlerFunc(s.AddTopPlayNext).Methods("GET")

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
