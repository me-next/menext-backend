package server

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/me-next/menext-backend/party"
	"sync"
)

// PartyManager manages parties
type PartyManager struct {
	parties map[uuid.UUID]*party.Party
	mux     *sync.RWMutex
}

// NewPartyManager from nothing.
func NewPartyManager() *PartyManager {
	return &PartyManager{
		parties: make(map[uuid.UUID]*party.Party),
		mux:     &sync.RWMutex{},
	}
}

// CreateParty with a unique identifier
// TODO: should this have a check to see if the owner is in another party?
func (pm *PartyManager) CreateParty(owner party.UserUUID, ownerName string) (uuid.UUID, error) {

	// create a new party
	p := party.New(owner, ownerName)

	// need to lock / unlock
	pm.mux.Lock()
	defer pm.mux.Unlock()

	// uuid
	pid := uuid.New()

	if _, found := pm.parties[pid]; found {
		return uuid.Nil, fmt.Errorf("oh nose, failed to create unique pid")
	}

	pm.parties[pid] = p

	return pid, nil
}

// Party by uuid.
// NOTE: this is unsafe, the party may be editted / messed up while we are working on it
func (pm *PartyManager) Party(pid uuid.UUID) (*party.Party, error) {
	pm.mux.RLock()
	defer pm.mux.RUnlock()

	// party my not be found because:
	// 1) bad key
	// 2) a party was disbanded but the client doesn't know about that yet
	// either way the user should behave the same way
	p, found := pm.parties[pid]
	if !found {
		return nil, fmt.Errorf("could not find party %s", pid.String())
	}

	return p, nil
}

// Remove a party from the manager by uuid
func (pm *PartyManager) Remove(pid uuid.UUID) error {
	pm.mux.Lock()
	defer pm.mux.Unlock()

	if _, found := pm.parties[pid]; !found {
		return fmt.Errorf("could not find party %s", pid.String())
	}

	// NOTE: disbanding a party is the same as it not existing
	delete(pm.parties, pid)
	return nil
}
