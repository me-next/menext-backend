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
// Path is /suggestUp/{pid}/{uid}/{sid}
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

	// try to upvote
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
// Path is /suggestClearvote /{pid}/{uid}/{sid}
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

	// try to clear the votes
	err = p.SuggestionClearvote(party.UserUUID(uidStr), party.SongUID(sidStr))
	if err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)
	}

	// exit with OK status code
}

// AddPlayNext adds a song to play next queue
// Path is /addPlayNext/{pid}/{uid}/{sid}
// The client must verify that the song id is good.
func (s *Server) AddPlayNext(w http.ResponseWriter, r *http.Request) {
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
	err = p.PlayNext(party.UserUUID(uidStr), party.SongUID(sidStr))
	if err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)
	}

	// exit with OK status code
}

// AddTopPlayNext adds a song to the top of the play next queue.
func (s *Server) AddTopPlayNext(w http.ResponseWriter, r *http.Request) {
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

	// try to add the song
	err = p.AddTopPlayNext(party.UserUUID(uidStr), party.SongUID(sidStr))
	if err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)
	}

	// exit with OK status code
}

// RemovePlayNext removes a song from the play next queue.
// Error if the song to remove isn't there.
func (s *Server) RemovePlayNext(w http.ResponseWriter, r *http.Request) {
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

	// try to add the song
	err = p.RemoveFromPlayNext(party.UserUUID(uidStr), party.SongUID(sidStr))
	if err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)
	}

	// exit with OK status code
}
