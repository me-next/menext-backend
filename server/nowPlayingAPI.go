package server

// contains the API functions for the playback features

import (
	"github.com/gorilla/mux"
	"github.com/me-next/menext-backend/party"
	"net/http"
)

// Seek to a position in the song.
// The path is /{pid}/{uid}/seek/{where}.
// The client must validate that the seek is correct.
func (s *Server) Seek(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uidStr, ufound := vars["uid"]
	pidStr, pfound := vars["pid"]
	seekPosition, sfound := vars["seekPos"]

	if !ufound || !pfound || !sfound {
		urlerror(w)
		return
	}

	p, err := s.pm.Party(PartyUUID(pidStr))
	if err != nil {
		errMsg := jsonError("no such party", pid)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// try to seek

}
