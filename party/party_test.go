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

func TestPartySkipPrev(t *testing.T) {
	ouid := party.UserUUID("1")
	p := party.New(ouid, "bob")

	// there was an issue with using skip and previous together
	// repro is:
	// 1) add song
	// 2) let a song finish
	// 3) skip a few songs
	// 4) hit previous
	// 5) nothing plays
	assert.Nil(t, p.Suggest(ouid, "a"))
	assert.Nil(t, p.Suggest(ouid, "b"))
	assert.Nil(t, p.Suggest(ouid, "c"))
	assert.Nil(t, p.Suggest(ouid, "d"))
	assert.Nil(t, p.Suggest(ouid, "e"))

	// finish a, b
	// should be changes 6, 7
	assert.Nil(t, p.Skip(ouid, "a"))
	assert.Nil(t, p.Skip(ouid, "b"))
	t.Log(p.Pull(ouid, 6))

	// now try previous, current should be c
	// change 8
	assert.Nil(t, p.Previous(ouid, "c"))
	t.Log(p.Pull(ouid, 6))

	// now try pulling and check output

}

func TestPartyRemoveSuggest(t *testing.T) {
	// check that adding a song to play-next removes from suggest
	// however, it shouldn't keep it out of the suggest in the future

	ouid := party.UserUUID("1")
	p := party.New(ouid, "bob")

	assert.Nil(t, p.Suggest(ouid, "a"))
	assert.Nil(t, p.Suggest(ouid, "c"))

	assert.Nil(t, p.Suggest(ouid, "b"))
	assert.NotNil(t, p.Suggest(ouid, "b"))

	// now try adding b to playnext
	assert.Nil(t, p.PlayNext(ouid, "b"))

	// now try adding to the suggest
	assert.Nil(t, p.Suggest(ouid, "b"))
	assert.NotNil(t, p.Suggest(ouid, "b"))

	// skip to wipe playnext
	assert.Nil(t, p.Skip(ouid, "a"))
	assert.NotNil(t, p.Suggest(ouid, "b"))

	// now try with addTop
	assert.Nil(t, p.AddTopPlayNext(ouid, "b"))

	// check we can suggest something in playnext
	assert.Nil(t, p.Suggest(ouid, "b"))
	assert.NotNil(t, p.Suggest(ouid, "b"))
}
