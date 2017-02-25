package party_test

import (
	"github.com/me-next/menext-backend/party"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPartyUserAdd(t *testing.T) {
	ownerUUID := party.UserUUID("1")
	p := party.New(ownerUUID, "fred")

	user := party.UserUUID("2")
	assert.Nil(t, p.AddUser(user, "bob"))

	// check that we can't remove someone who isn't there
	assert.NotNil(t, p.RemoveUser("4"))

	// can't double add users
	assert.NotNil(t, p.AddUser(user, "bob"))

	// remove users should be ok
	assert.Nil(t, p.RemoveUser(user))

	// can't double-remove
	assert.NotNil(t, p.RemoveUser(user))

	// can remove then add back in
	assert.Nil(t, p.AddUser(user, "bob"))

	// check that owner was inserted
	assert.NotNil(t, p.AddUser(ownerUUID, "fred"))

	// check that we can't remove owners
	assert.NotNil(t, p.RemoveUser(ownerUUID))
}

func TestPartyCanRemove(t *testing.T) {
	ownerUUID := party.UserUUID("1")
	p := party.New(ownerUUID, "bob")

	user := party.UserUUID("2")

	assert.False(t, p.CanUserEndParty(user))
	assert.True(t, p.CanUserEndParty(ownerUUID))
}
