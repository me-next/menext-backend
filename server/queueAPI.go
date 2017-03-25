package server

// this file contains the API for running queue operations

import (
	"github.com/gorilla/mux"
	"github.com/me-next/menext-backend/party"
	"net/http"
)

// Suggest a song to a party's suggestion queue.
// Path is /suggest/{pid}/{uid}/{sid}
// The client must verify that the song id is good.
func (s *Server) Suggest(w http.ResponseWriter, r *http.Request) {
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

	// try to suggest teh song
	err = p.Suggest(party.UserUUID(uidStr), party.SongUID(sidStr))
	if err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)
	}

	// exit with OK status code
}

// SuggestionUpvote a song in an event's suggestion queue.
// Path is /suggest/{pid}/{uid}/{sid}
// The client must verify that the song id is good.
func (s *Server) SuggestionUpvote(w http.ResponseWriter, r *http.Request) {
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

	// try to downvote
	err = p.SuggestionUpvote(party.UserUUID(uidStr), party.SongUID(sidStr))
	if err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)
	}

	// exit with OK status code
}

// SuggestionDownvote a song in an even't suggestion queue.
// Path is /suggestDown/{pid}/{uid}/{sid}
// The client must verify that the song id is good.
func (s *Server) SuggestionDownvote(w http.ResponseWriter, r *http.Request) {
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

	// try to upvote
	err = p.SuggestionDownvote(party.UserUUID(uidStr), party.SongUID(sidStr))
	if err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)
	}

	// exit with OK status code
}

// SuggestionClearvote clears a users votes for a song
// Path is /suggestDown/{pid}/{uid}/{sid}
// The client must verify that the song id is good.
func (s *Server) SuggestionClearvote(w http.ResponseWriter, r *http.Request) {
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

	// try to upvote
	err = p.SuggestionClearvote(party.UserUUID(uidStr), party.SongUID(sidStr))
	if err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)
	}

	// exit with OK status code
}
