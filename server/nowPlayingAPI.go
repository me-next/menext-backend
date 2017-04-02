package server

// contains the API functions for the playback features

import (
	"github.com/gorilla/mux"
	"github.com/me-next/menext-backend/party"
	"net/http"
	"strconv"
)

// Seek to a position in the song.
// The path is /{pid}/{uid}/seek/{pos}.
// The client must validate that the seek is correct.
// The party checks that the seek is valid
func (s *Server) Seek(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uidStr, ufound := vars["uid"]
	pidStr, pfound := vars["pid"]
	seekPositionStr, sfound := vars["pos"]

	if !ufound || !pfound || !sfound {
		urlerror(w)
		return
	}

	// try to parse the int
	// this should just be a 32
	seekPosition, err := strconv.ParseUint(seekPositionStr, 10, 32)
	if err != nil {
		errMsg := jsonError("failed to parse seek position")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	p, err := s.pm.Party(PartyUUID(pidStr))
	if err != nil {
		errMsg := jsonError("no such party")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// try to seek
	// error could be from invalid user or bad seek
	err = p.Seek(party.UserUUID(uidStr), uint32(seekPosition))
	if err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// exit with OK status code
}

// SongFinished notifies the server to play the next song
// The path is /songFinished/{pid}/{uid}/{sid}
func (s *Server) SongFinished(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uidStr, ufound := vars["uid"]
	pidStr, pfound := vars["pid"]
	sidStr, sfound := vars["sid"]

	if !ufound || !pfound || !sfound {
		urlerror(w)
		return
	}

	p, err := s.pm.Party(PartyUUID(pidStr))
	if err != nil {
		errMsg := jsonError("no such party")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// try to seek
	// error could be from invalid user or bad seek
	err = p.SongFinished(party.UserUUID(uidStr), party.SongUID(sidStr))
	if err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// exit with OK status code
}
