package event

import (
	"fmt"
	"github.com/me-next/menext-backend/queue"
	"net/http"
)

// event contains a queue and manages users
type Event struct {
	users     map[UserUID]*User
	ownerUID UserUID
	playNextQueue *queue.Queue
	suggestionQueue *queue.Queue
}

// New event
func New(ownerUID UserUID, ownerName string) *Event {
	e := Event{
		users:     make(map[UserUID]*User),
		ownerUID: ownerUID,
		playNextQueue: queue.NewQueue(),
		suggestionQueue: queue.NewQueue(),
	}

	e.AddUser(ownerUID, ownerName)
	e.SetOwner(ownerUID)
	return &e
}

// AddUser to the event, applies default permissions
func (e *Event) AddUser(userUID UserUID, name string) error {
	user := NewUser(name)
	if _, has := e.getUser(userUID); has == nil {
		return fmt.Errorf("event already contains user %s", userUID)
	}

	e.setDefaultPermission(user)
	e.users[userUID] = user
	return nil
}

// RemoveUser from the event
func (e *Event) RemoveUser(userUID UserUID) error {
	if userUID == e.ownerUID {
		// TODO: should terminate instead...
		return fmt.Errorf("removing owner from event")
	}

	if _, has := e.getUser(userUID); has != nil {
		return fmt.Errorf("user %s not in the event", userUID)
	}

	delete(e.users, userUID)
	return nil
}

// canUserPerformAction id'd by string
func (e *Event) canUserPerformAction(userUID UserUID, action string) (bool, error) {
	user, err := e.getUser(userUID)

	if err != nil {
		return false, err
	}

	return user.CanPerform(action), nil
}

func (e *Event) setDefaultPermission(user *User) {
	// TODO: replace with real permissions
	user.SetPermission("default", true)
	user.SetPermission("bad", false)
}

func (e *Event) getUser(userUID UserUID) (*User, error) {
	user, has := e.users[userUID]
	if !has {
		return nil, fmt.Errorf("user %s not found", userUID)
	}

	return user, nil
}

// SetOwner of the event (there can be only one)
func (e *Event) SetOwner(userUID UserUID) error {
	if _, has := e.getUser(userUID); has == nil {
		return fmt.Errorf("user %s not found", userUID)
	}

	// TODO: set the permissions
	e.ownerUID = userUID

	return nil
}

//TODO: check if user allowe to add/remove
func (e *Event) AddPlayNext(user UserUID, songID queue.SongUID) error {
	return e.playNextQueue.AddSong(songID)
}

func (e *Event) RemovePlayNext(user UserUID, songID queue.SongUID) error {
	return e.playNextQueue.RemoveSong(songID)
}

func (e *Event) AddsuggestionQ(user UserUID, songID queue.SongUID) error {
	return e.suggestionQueue.AddSong(songID)
}

func (e *Event) RemovesuggestionQ(user UserUID, songID queue.SongUID) error {
	return e.suggestionQueue.RemoveSong(songID)
}

func (e *Event) AddToSuggestion(w http.ResponseWriter, r *http.Request) {
// TODO: Implement
}

func (e *Event) RemoveFromSuggestion(w http.ResponseWriter, r *http.Request) {
// TODO: Implement
}

func (e *Event) AddToPlayNext(w http.ResponseWriter, r *http.Request) {
// TODO: Implement
}

func (e *Event) RemoveFromPlayNext(w http.ResponseWriter, r *http.Request) {
// TODO: Implement
}

//move in suggestion queue
func (e *Event) ChangeProirity(w http.ResponseWriter, r *http.Request) {
// TODO: Implement
}

//increase priority in suggestion queue
func (e *Event) ThumbsUp(w http.ResponseWriter, r *http.Request) {
// TODO: Implement
}

//decrease priority in suggestion queue
func (e *Event) ThumbsDown(w http.ResponseWriter, r *http.Request) {
// TODO: Implement
}

func (e *Event) Next(w http.ResponseWriter, r *http.Request) {
// TODO: Implement
}

func (e *Event) Previous(w http.ResponseWriter, r *http.Request) {
// TODO: Implement
}

func (e *Event) ChangeUserPermission(w http.ResponseWriter, r *http.Request) {
// TODO: Implement
}

func (e *Event) ChangeEventPermission(w http.ResponseWriter, r *http.Request) {
// TODO: Implement
}