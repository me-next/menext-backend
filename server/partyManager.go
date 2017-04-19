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

	// spin up the cleanup thread in the background
	go func(pm *PartyManager) {
		ticker := time.NewTicker(cleanupPeriodHours * time.Hour)
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

// CreatePartyWithName attempts to create a party with the custom ID.
// If the party is already used, return a map with suggested names.
// Suggested names format is: {"suggested": ["a", ...]}.
// All of the suggested names will be valid
func (pm *PartyManager) CreatePartyWithName(ouid party.UserUUID, oname string, pid string) (PartyUUID, PartyUUID, error) {
	// check that the party exists
	pm.mux.RLock()

	if _, found := pm.parties[PartyUUID(pid)]; !found {
		pm.mux.RUnlock()

		// get the write lock
		pm.mux.Lock()

		p := party.New(ouid, oname)

		// double check that our desired name is still available
		if _, found = pm.parties[PartyUUID(pid)]; !found {
			pm.parties[PartyUUID(pid)] = p
			pm.mux.Unlock()

			return PartyUUID(pid), "", nil
		}

		// someone ninja'd the name
		// release the write lock, get the read lock, and try to generate a new name
		pm.mux.Unlock()
		pm.mux.RLock()
	}

	// generate alternative name
	alternative, err := pm.attemptMutate(pid)

	if err != nil {
		defer pm.mux.RUnlock()

		return "", pm.generateUUID(), err
	}

	pm.mux.RUnlock()

	// TODO: should this ever return an error
	return "", alternative, fmt.Errorf("party name not available")
}

// tries to generate an alternative
func (pm PartyManager) attemptMutate(pid string) (PartyUUID, error) {

	// generate a slice that we'll reuse
	n := len(pid) + 1
	attempt := make([]byte, n)
	copy(attempt, pid)

	for i := 1; i < 10; i++ {
		attempt[n-1] = '0' + byte(i)
		if _, found := pm.parties[PartyUUID(attempt)]; !found {
			return PartyUUID(attempt), nil
		}
	}

	return "", fmt.Errorf("could not generate an alternative")
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
	cleanupPeriodHours  = 6
	partyExpirationTime = 48
)

// Cleanup removes all events older than expirationTime.
// It is called by a background thread every <cleanupPeriodHours>.
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
		if event.TimeSinceLastChange().Seconds() > expirationTime.Seconds() {
			expiredParties[key] = struct{}{}
		}
	}

	// now let other people write
	pm.mux.RUnlock()

	// try to let other people through by giving up the WLock
	for key := range expiredParties {
		pm.Remove(key)
	}
}

const (
	partyUUIDSizeConst       = 6
	partyUUIDCreateLoopLimit = 50
)

// generateUUID with 6 letters / numbers
// panics if can't create a uuid within partUUIDCreateLoopLimit tries
func (pm PartyManager) generateUUID() PartyUUID {
	letterBytes := "abcdefghijklmnopqrstuvwxyz"

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
