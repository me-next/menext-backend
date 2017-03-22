package event_test

import (
	"github.com/me-next/menext-backend/event"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEventUserAdd(t *testing.T) {
	ownerUID := event.UserUID("1")
	e := event.New(ownerUID, "fred")

	user := event.UserUID("2")
	assert.Nil(t, e.AddUser(user, "bob"))

	// check that we can't remove someone who isn't there
	assert.NotNil(t, e.RemoveUser("4"))

	// can't double add users
	assert.NotNil(t, e.AddUser(user, "bob"))

	// remove users should be ok
	assert.Nil(t, e.RemoveUser(user))

	// can't double-remove
	assert.NotNil(t, e.RemoveUser(user))

	// can remove then add back in
	assert.Nil(t, e.AddUser(user, "bob"))

	// check that owner was inserted
	assert.NotNil(t, e.AddUser(ownerUID, "fred"))

	// check that we can't remove owners
	assert.NotNil(t, e.RemoveUser(ownerUID))
}

func TestEventCanRemove(t *testing.T) {
	ownerUID := event.UserUID("1")
	e := event.New(ownerUID, "bob")

	user := event.UserUID("2")

	assert.False(t, e.CanUserEndEvent(user))
	assert.True(t, e.CanUserEndEvent(ownerUID))
}

func TestEventPull(t *testing.T) {
	// NOTE: we need to get some actions that increase the change counter
	// before we can properly test the pull
	ownerUID := event.UserUID("1")
	e := event.New(ownerUID, "bob")

	data, err := e.Pull(ownerUID, 0)
	assert.Nil(t, data)
	assert.Nil(t, err)

	// bad change ID, too high
	data, err = e.Pull(ownerUID, 2)
	assert.Nil(t, data)
	assert.NotNil(t, err)

	baduid := event.UserUID("2")
	data, err = e.Pull(baduid, 0)
	assert.Nil(t, data)
	assert.Nil(t, err)

	data, err = e.Pull(baduid, 1)
	assert.Nil(t, data)
	assert.NotNil(t, err)
}
