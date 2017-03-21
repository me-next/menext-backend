package event_test

import (
	"github.com/me-next/menext-backend/event"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPartyUserAdd(t *testing.T) {
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
	
	//check add play next
	assert.Nil(t, e.AddPlayNext(user, "Song 1"))
	
	//check add suggestion
	assert.Nil(t, e.AddsuggestionQ(user, "Song 1"))
	
	//can't double add songs
	assert.NotNil(t, e.AddsuggestionQ(user, "Song 1"))

	// check remove play next
	assert.Nil(t, e.RemovePlayNext(user, "Song 1"))
	
	// check remove suggestion
	assert.Nil(t, e.RemovesuggestionQ(user, "Song 1"))
	
	// check that we can't remove  song that isn't there
	assert.NotNil(t, e.RemovesuggestionQ(user, "Not Song"))

}
