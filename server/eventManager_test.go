package server_test

import (
	"github.com/me-next/menext-backend/server"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManagerSimple(t *testing.T) {
	em := server.NewEventManager()
	assert.NotNil(t, em)

	eid, err := em.CreateEvent("1", "ted")
	assert.Nil(t, err)
	assert.NotEqual(t, eid, "")

	// lookup by event
	e, err := em.Event(eid)
	assert.NotNil(t, e)
	assert.Nil(t, err)

	// make sure we can update event
	err = e.AddUser("2", "bob")
	assert.Nil(t, err)

	echanged, err := em.Event(eid)
	assert.Nil(t, err)
	assert.Equal(t, echanged, e)

	err = echanged.RemoveUser("2")
	assert.Nil(t, err)

	// bad lookup
	baduuuid := server.EventUID("1")
	assert.NotEqual(t, baduuuid, eid)

	bade, err := em.Event(baduuuid)
	assert.Nil(t, bade)
	assert.NotNil(t, err)

	err = em.Remove(baduuuid)
	assert.NotNil(t, err)

	// remove a event
	err = em.Remove(eid)
	assert.Nil(t, err)

	// check that it actually got removed
	err = em.Remove(eid)
	assert.NotNil(t, err)

	e, err = em.Event(eid)
	assert.Nil(t, e)
	assert.NotNil(t, err)
}
