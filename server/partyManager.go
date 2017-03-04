package server

import (
	"fmt"
	"github.com/me-next/menext-backend/party"
	"math/rand"
	"sync"
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
	return &PartyManager{
		parties: make(map[PartyUUID]*party.Party),
		mux:     &sync.RWMutex{},
	}
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

// generateUUID with 6 letters / numbers
// panics if can't create a uuid
func (pm PartyManager) generateUUID() PartyUUID {
	letterBytes := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// try to generate a uuid
	for j := 0; j < 10; j++ {
		b := make([]byte, 6)
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
