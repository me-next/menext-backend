package server

import (
	"fmt"
	"github.com/me-next/menext-backend/party"
	"math/rand"
	"sync"
	"time"
)

// PartyUUID uniquely identifies a party
type PartyUUID string

// PartyManager manages parties
type PartyManager struct {
	parties map[PartyUUID]*party.Party
	mux     *sync.RWMutex
}

// NewPartyManager from nothing.
func NewPartyManager() *PartyManager {
	pm := &PartyManager{
		parties: make(map[PartyUUID]*party.Party),
		mux:     &sync.RWMutex{},
	}

	// spin up the cleanup chron job
	go func(pm *PartyManager) {
		ticker := time.NewTicker(chronPeriodHours * time.Hour)
		for _ = range ticker.C {
			pm.Cleanup(partyExpirationTime * time.Hour)
		}
	}(pm)

	return pm
}

// CreateParty with a unique identifier
// TODO: should this have a check to see if the owner is in another party?
func (pm *PartyManager) CreateParty(owner party.UserUUID, ownerName string) (PartyUUID, error) {

	// create a new party
	p := party.New(owner, ownerName)

	// need to lock / unlock
	pm.mux.Lock()
	defer pm.mux.Unlock()

	// uuid
	pid := pm.generateUUID()

	if _, found := pm.parties[pid]; found {
		return "", fmt.Errorf("oh nose, failed to create unique pid")
	}

	pm.parties[pid] = p

	return pid, nil
}

// Party by uuid.
// NOTE: this is unsafe, the party may be editted / messed up while we are working on it
func (pm *PartyManager) Party(pid PartyUUID) (*party.Party, error) {
	pm.mux.RLock()
	defer pm.mux.RUnlock()

	// party my not be found because:
	// 1) bad key
	// 2) a party was disbanded but the client doesn't know about that yet
	// either way the user should behave the same way
	p, found := pm.parties[pid]
	if !found {
		return nil, fmt.Errorf("could not find party %s", pid)
	}

	return p, nil
}

// Remove a party from the manager by uuid
func (pm *PartyManager) Remove(pid PartyUUID) error {
	pm.mux.Lock()
	defer pm.mux.Unlock()

	if _, found := pm.parties[pid]; !found {
		return fmt.Errorf("could not find party %s", pid)
	}

	// NOTE: disbanding a party is the same as it not existing
	delete(pm.parties, pid)
	return nil
}

// consts for party cleanup
const (
	chronPeriodHours    = 6
	partyExpirationTime = 48
)

// Cleanup removes all events older than expirationTime.
// It is called periodically by a chron job.
func (pm *PartyManager) Cleanup(expirationTime time.Duration) {
	// Maps in go don't support concurrent access, and will even throw exceptions
	// if they suspect there is concurrent access.
	// Go doesn't support threadsafe deep copies of maps without iterating though the
	// whole thing. We take the read-lock to find expired events, release the read-lock
	// then delete expired events one-by-one to share the lock.

	// block all write activities while looping through the map
	pm.mux.RLock()
	expiredParties := make(map[PartyUUID]struct{})

	// find all expired parties and store the keys
	for key, event := range pm.parties {
		if event.TimeSinceLastChange().Minutes() > expirationTime.Minutes() {
			expiredParties[key] = struct{}{}
		}
	}

	// now let other people write
	pm.mux.RUnlock()

	// try to let other people through by giving up the WLock
	for key := range expiredParties {
		pm.mux.Lock()

		delete(pm.parties, key)

		pm.mux.Unlock()
	}
}

const (
	partyUUIDSizeConst       = 6
	partyUUIDCreateLoopLimit = 50
)

// generateUUID with 6 letters / numbers
// panics if can't create a uuid within partUUIDCreateLoopLimit tries
func (pm PartyManager) generateUUID() PartyUUID {
	letterBytes := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// try to generate a uuid
	// cap so if things are weird we can get out
	for j := 0; j < partyUUIDCreateLoopLimit; j++ {
		b := make([]byte, partyUUIDSizeConst)
		for i := range b {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		}

		// check for membership
		if _, found := pm.parties[PartyUUID(b)]; !found {
			return PartyUUID(b)
		}
	}

	panic("oh nose! couldn't generate a uuid")
}
