package server

import (
	"fmt"
	"github.com/me-next/menext-backend/event"
	"math/rand"
	"sync"
)

// EventUID uniquely identifies a event
type EventUID string

// EventManager manages events
type EventManager struct {
	events map[EventUID]*event.Event
	mux     *sync.RWMutex
}

// NewEventManager from nothing.
func NewEventManager() *EventManager {
	return &EventManager{
		events: make(map[EventUID]*event.Event),
		mux:     &sync.RWMutex{},
	}
}

// CreateEvent with a unique identifier
// TODO: should this have a check to see if the owner is in another event?
func (em *EventManager) CreateEvent(owner event.UserUID, ownerName string) (EventUID, error) {

	// create a new event
	e := event.New(owner, ownerName)

	// need to lock / unlock
	em.mux.Lock()
	defer em.mux.Unlock()

	// uuid
	eid := em.generateUID()

	if _, found := em.events[eid]; found {
		return "", fmt.Errorf("oh nose, failed to create unique eid")
	}

	em.events[eid] = e

	return eid, nil
}

// Event by uuid.
// NOTE: this is unsafe, the event may be editted / messed up while we are working on it
func (em *EventManager) Event(eid EventUID) (*event.Event, error) {
	em.mux.RLock()
	defer em.mux.RUnlock()

	// event my not be found because:
	// 1) bad key
	// 2) a event was disbanded but the client doesn't know about that yet
	// either way the user should behave the same way
	e, found := em.events[eid]
	if !found {
		return nil, fmt.Errorf("could not find event %s", eid)
	}

	return e, nil
}

// Remove a event from the manager by uuid
func (em *EventManager) Remove(eid EventUID) error {
	em.mux.Lock()
	defer em.mux.Unlock()

	if _, found := em.events[eid]; !found {
		return fmt.Errorf("could not find event %s", eid)
	}

	// NOTE: disbanding a event is the same as it not existing
	delete(em.events, eid)
	return nil
}

const (
	eventUIDSizeConst       = 6
	eventUIDCreateLoopLimit = 50
)

// generateUID with 6 letters / numbers
// panics if can't create a uuid within partUIDCreateLoopLimit tries
func (em EventManager) generateUID() EventUID {
	letterBytes := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// try to generate a uuid
	// cap so if things are weird we can get out
	for j := 0; j < eventUIDCreateLoopLimit; j++ {
		b := make([]byte, eventUIDSizeConst)
		for i := range b {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		}

		// check for membership
		if _, found := em.events[EventUID(b)]; !found {
			return EventUID(b)
		}
	}

	panic("oh nose! couldn't generate a uuid")
}
