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

	// try to parse the float
	seekPosition, err := strconv.ParseFloat(seekPositionStr, 32)
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
	err = p.Seek(party.UserUUID(uidStr), float32(seekPosition))
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

	// try to finish the song
	// error could be from invalid user or bad song
	err = p.SongFinished(party.UserUUID(uidStr), party.SongUID(sidStr))
	if err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// exit with OK status code
}

// Skip the currently playing song.
// The path is /skip/{pid}/{uid}/{sid}
func (s *Server) Skip(w http.ResponseWriter, r *http.Request) {
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

	// try to skip
	// error could be from invalid user or bad seek
	err = p.Skip(party.UserUUID(uidStr), party.SongUID(sidStr))
	if err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// exit with OK status code
}

// Previous plays the previous song.
// The currently playing song is put on the top of the play-next queue
// The path is /previous/{pid}/{uid}/{sid}
func (s *Server) Previous(w http.ResponseWriter, r *http.Request) {
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

	// try to play the previous song
	// error could be from invalid user or bad seek
	err = p.Previous(party.UserUUID(uidStr), party.SongUID(sidStr))
	if err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// exit with OK status code
}

// SetVolume changes the current volume
// The path is /setVolume/{pid}/{uid}/{volume}
func (s *Server) SetVolume(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uidStr, ufound := vars["uid"]
	pidStr, pfound := vars["pid"]
	volStr, vfound := vars["volume"]

	if !ufound || !pfound || !vfound {
		urlerror(w)
		return
	}

	// try to parse the int
	// this should just be a 32
	volume, err := strconv.ParseUint(volStr, 10, 32)
	if err != nil {
		errMsg := jsonError("failed to parse volume")
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

	// try to change the volume
	// err could be bad uid or bad volume
	err = p.SetVolume(party.UserUUID(uidStr), uint32(volume))
	if err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// exit with OK status code
}

// Play song if a song is playing. This is just a play/pause control.
// path is /play/{pid}/{uid}
func (s *Server) Play(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uidStr, ufound := vars["uid"]
	pidStr, pfound := vars["pid"]

	if !ufound || !pfound {
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

	// try to play
	// error could be from bad user, nothing to play
	err = p.Play(party.UserUUID(uidStr))
	if err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// exit with OK status code
}

// Pause the song at a certain position.
// path is /pause/{pid}/{uid}/pos
func (s *Server) Pause(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uidStr, ufound := vars["uid"]
	pidStr, pfound := vars["pid"]
	posStr, posfound := vars["pos"]

	if !ufound || !pfound || !posfound {
		urlerror(w)
		return
	}

	// try to parse the float
	pos, err := strconv.ParseFloat(posStr, 32)
	if err != nil {
		errMsg := jsonError("failed to parse position")
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

	// try to pause
	// err could be bad uid or song state not changing
	err = p.Pause(party.UserUUID(uidStr), float32(pos))
	if err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)

		return
	}

	// exit with OK status code
}

// PlayNow plays a song right now.
func (s *Server) PlayNow(w http.ResponseWriter, r *http.Request) {
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

	// try to play the song
	err = p.PlayNow(party.UserUUID(uidStr), party.SongUID(sidStr))
	if err != nil {
		errMsg := jsonError("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)
	}

	// exit with OK status code
}
