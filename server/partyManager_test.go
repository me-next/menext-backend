package server_test

import (
	"github.com/me-next/menext-backend/server"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

func TestCleanup(t *testing.T) {
	// check that we can clean up old parties properly
	t.Skip("this test takes a little while to run")

	pm := server.NewPartyManager()

	// insert a party
	pida, err := pm.CreateParty("1", "a")
	assert.Nil(t, err)

	// sleep for a second
	time.Sleep(1 * time.Second)

	// add 2nd event
	pidb, err := pm.CreateParty("2", "b")
	assert.Nil(t, err)

	time.Sleep(500 * time.Millisecond)

	// cleanup A
	pm.Cleanup(1 * time.Second)

	_, err = pm.Party(pida)
	assert.NotNil(t, err)

	// should still have b
	_, err = pm.Party(pidb)
	assert.Nil(t, err)

	// wait then clean up B
	time.Sleep(700 * time.Millisecond)

	pm.Cleanup(1 * time.Second)
	_, err = pm.Party(pidb)
	assert.NotNil(t, err)
}
