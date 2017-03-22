package event

import (
	"fmt"
	"sync"
)

// Event contains a queue and manages users
type Event struct {
	users     map[UserUID]*User
	ownerUID UserUID
	mux       *sync.Mutex
	changeID  uint64
}

// New event
func New(ownerUID UserUID, ownerName string) *Event {
	e := Event{
		users:     make(map[UserUID]*User),
		ownerUID: ownerUID,
		mux:       &sync.Mutex{},
	}

	e.AddUser(ownerUID, ownerName)
	e.SetOwner(ownerUID)
	return &e
}

// AddUser to the event, applies default permissions
func (e *Event) AddUser(userUID UserUID, name string) error {

	e.mux.Lock()
	defer e.mux.Unlock()

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

	e.mux.Lock()
	defer e.mux.Unlock()

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

// CanUserEndEvent id'ing the user by uuid
func (e *Event) CanUserEndEvent(userUID UserUID) bool {
	e.mux.Lock()
	defer e.mux.Unlock()

	return e.ownerUID == userUID
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

	e.mux.Lock()
	defer e.mux.Unlock()

	if _, has := e.getUser(userUID); has == nil {
		return fmt.Errorf("user %s not found", userUID)
	}

	// TODO: set the permissions
	e.ownerUID = userUID

	return nil
}

// updated increments the update tracker.
// Should call this whenever there's an update everyone should know about.
func (e *Event) setUpdated() {
	e.changeID++
}

// Pull returns the user data in a serializable format.
// NOTE: this checks for changes before checking uid.
func (e *Event) Pull(userUID UserUID, clientChangeID uint64) (interface{}, error) {
	e.mux.Lock()
	defer e.mux.Unlock()

	// if the client's change is larger than our current change
	if e.changeID < clientChangeID {
		return nil, fmt.Errorf("bad pull id")
	}

	// up to date
	if e.changeID == clientChangeID {
		return nil, nil
	}

	user, err := e.getUser(userUID)

	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data["permissions"] = user.Data()
	data["change"] = e.changeID

	return data, nil
}
