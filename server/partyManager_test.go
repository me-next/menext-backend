package server_test

import (
	"github.com/me-next/menext-backend/server"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManagerSimple(t *testing.T) {
	pm := server.NewPartyManager()
	assert.NotNil(t, pm)

	pid, err := pm.CreateParty("1", "ted")
	assert.Nil(t, err)
	assert.NotEqual(t, pid, "")

	// lookup by party
	p, err := pm.Party(pid)
	assert.NotNil(t, p)
	assert.Nil(t, err)

	// make sure we can update party
	err = p.AddUser("2", "bob")
	assert.Nil(t, err)

	pchanged, err := pm.Party(pid)
	assert.Nil(t, err)
	assert.Equal(t, pchanged, p)

	err = pchanged.RemoveUser("2")
	assert.Nil(t, err)

	// bad lookup
	baduuuid := server.PartyUUID("1")
	assert.NotEqual(t, baduuuid, pid)

	badp, err := pm.Party(baduuuid)
	assert.Nil(t, badp)
	assert.NotNil(t, err)

	err = pm.Remove(baduuuid)
	assert.NotNil(t, err)

	// remove a party
	err = pm.Remove(pid)
	assert.Nil(t, err)

	// check that it actually got removed
	err = pm.Remove(pid)
	assert.NotNil(t, err)

	p, err = pm.Party(pid)
	assert.Nil(t, p)
	assert.NotNil(t, err)
}
