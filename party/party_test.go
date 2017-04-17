package party_test

import (
	//"bytes"
	//"encoding/gob"
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

func TestPartyPull(t *testing.T) {
	// NOTE: we need to get some actions that increase the change counter
	// before we can properly test the pull
	ownerUUID := party.UserUUID("1")
	p := party.New(ownerUUID, "bob")

	data, err := p.Pull(ownerUUID, 0)
	assert.Nil(t, data)
	assert.Nil(t, err)

	// bad change ID, too high
	data, err = p.Pull(ownerUUID, 2)
	assert.Nil(t, data)
	assert.NotNil(t, err)

	baduid := party.UserUUID("2")
	data, err = p.Pull(baduid, 0)
	assert.Nil(t, data)
	assert.Nil(t, err)

	data, err = p.Pull(baduid, 1)
	assert.Nil(t, data)
	assert.NotNil(t, err)
}

func TestPartySeek(t *testing.T) {
	ownerUUID := party.UserUUID("1")
	p := party.New(ownerUUID, "bob")

	// parses out the position
	getPos := func(p *party.Party, uid party.UserUUID, cid uint64) (uint32, error) {
		raw, err := p.Pull(uid, cid)
		if err != nil {
			return 0, err
		}

		pullData := raw.(map[string]interface{})
		changeData := pullData[party.PullPlayingKey].(map[string]interface{})

		pos := changeData[party.KSongPosition]
		return pos.(uint32), nil
	}

	p.Suggest(ownerUUID, "a")

	// try a valid seek
	var seekTo uint32 = 5
	err := p.Seek(ownerUUID, seekTo)
	assert.Nil(t, err)

	// fetch the pos, check that it matches
	pos, err := getPos(p, ownerUUID, 0)
	assert.Nil(t, err)
	assert.Equal(t, pos, seekTo)

	// we already know that seek works, test bad gets
	err = p.Seek(party.UserUUID("2"), 1)
	assert.NotNil(t, err)

	// check that the changeID didn't move
	data, err := p.Pull(ownerUUID, 2)
	assert.Nil(t, err)
	assert.Empty(t, data)
}

func TestPartyPermissions(t *testing.T) {
	ouid := party.UserUUID("1")
	p := party.New(ouid, "bob")

	otherUID := party.UserUUID("other")
	assert.Nil(t, p.AddUser(otherUID, "fred"))

	type testCase struct {
		permission string         // what to set
		value      bool           // attempted value of permission
		uid        party.UserUUID // user trying to set permission
		expectNil  bool           // if the case should have an error or not
	}

	cases := []testCase{
		{party.UserCanSeekPermission, false, ouid, true},        // change permission
		{party.UserCanSeekPermission, true, ouid, true},         // change it back
		{party.UserCanSeekPermission, false, otherUID, false},   // try to have a follower set permissions
		{"bad", true, ouid, false},                              // set bad permisison
		{party.UserCanPlayPausePermission, true, ouid, false},   // set permission without changing state
		{party.UserCanPlayPausePermission, false, "bad", false}, // do a good set but with a bad uid
	}

	// run the test cases
	for _, test := range cases {
		err := p.SetPermission(test.permission, test.value, test.uid)

		// switch on expected value of error
		if test.expectNil {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}
