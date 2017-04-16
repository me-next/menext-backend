package party

import (
	"fmt"
	"sync"
	"time"
)

// Party contains a queue and manages users
type Party struct {
	users     map[UserUUID]*User
	ownerUUID UserUUID
	mux       *sync.Mutex
	changeID  uint64

	// queues
	suggestionQueue VotableQueue
	nowPlaying      NowPlaying

	lastChangeT time.Time

	// map of the permission to bools of if a user can use them
	permMap map[string]bool
}

// New party
func New(ownerUUID UserUUID, ownerName string) *Party {
	p := Party{
		users:           make(map[UserUUID]*User),
		ownerUUID:       ownerUUID,
		mux:             &sync.Mutex{},
		nowPlaying:      NowPlaying{},
		suggestionQueue: NewVotableQueue(),

		lastChangeT: time.Now(),

		permMap: make(map[string]bool),
	}

	// initially set true for all permissions
	for key := range PermissionDescriptionMap {
		p.permMap[key] = true
	}

	p.AddUser(ownerUUID, ownerName)
	p.SetOwner(ownerUUID)
	return &p
}

// AddUser to the party, applies default permissions
func (p *Party) AddUser(userUUID UserUUID, name string) error {

	p.mux.Lock()
	defer p.mux.Unlock()

	user := NewUser(name)
	if _, has := p.getUser(userUUID); has == nil {
		return fmt.Errorf("party already contains user %s", userUUID)
	}

	p.setDefaultPermission(user)
	p.users[userUUID] = user
	return nil
}

// RemoveUser from the party
func (p *Party) RemoveUser(userUUID UserUUID) error {

	p.mux.Lock()
	defer p.mux.Unlock()

	if userUUID == p.ownerUUID {
		// TODO: should terminate instead...
		return fmt.Errorf("removing owner from party")
	}

	if _, has := p.getUser(userUUID); has != nil {
		return fmt.Errorf("user %s not in the party", userUUID)
	}

	delete(p.users, userUUID)
	return nil
}

// CanUserEndParty id'ing the user by uuid
func (p *Party) CanUserEndParty(userUUID UserUUID) bool {
	p.mux.Lock()
	defer p.mux.Unlock()

	return p.ownerUUID == userUUID
}

// canUserPerformAction id'd by string
func (p *Party) canUserPerformAction(userUUID UserUUID, action string) (bool, error) {
	// check that the user exists
	if _, err := p.getUser(userUUID); err != nil {
		return false, err
	}

	// check the permission
	value, has := p.permMap[action]
	if !has {
		return false, fmt.Errorf("unknown permisison")
	}

	return value, nil
}

func (p *Party) setDefaultPermission(user *User) {
	// TODO: replace with real permissions
	user.SetPermission("default", true)
	user.SetPermission("bad", false)

	/*
		user.SetPermission(UserCanSeekPermission, true)
		user.SetPermission(UserCanSuggestSongPermission, true)
		user.SetPermission(UserCanVoteSuggestionPermission, true)
		user.SetPermission(UserCanChangeVolumePermission, true)
		user.SetPermission(UserCanPlayPausePermission, true)

	*/
}

// GetPermissions returns a map of permission values to their descriptions
func (p Party) GetPermissions() map[string]string {
	p.mux.Lock()
	defer p.mux.Unlock()

	return PermissionDescriptionMap
}

// SetPermission by key.
// uid of person trying to set the permissions.
func (p *Party) SetPermission(which string, value bool, uid UserUUID) error {

	// not a valid permission
	if _, has := PermissionDescriptionMap[which]; !has {
		return fmt.Errorf("not a valid permission")
	}

	// lock later b/c above should be threadsafe
	p.mux.Lock()
	defer p.mux.Unlock()

	// check that the owner is setting perms
	if uid != p.ownerUUID {
		return fmt.Errorf("only owner can set permissions")
	}

	// adding a permission
	if permValue, has := p.permMap[which]; !has {
		return fmt.Errorf("something very wrong! permission map is missing a valid permission")
	} else if permValue == value {
		return fmt.Errorf("not changing anything")
	}

	// else update permission
	p.permMap[which] = value
	p.setUpdated()

	return nil
}

func (p *Party) getUser(userUUID UserUUID) (*User, error) {
	user, has := p.users[userUUID]
	if !has {
		return nil, fmt.Errorf("user %s not found", userUUID)
	}

	return user, nil
}

// SetOwner of the party (there can be only one)
func (p *Party) SetOwner(userUUID UserUUID) error {

	p.mux.Lock()
	defer p.mux.Unlock()

	if _, has := p.getUser(userUUID); has == nil {
		return fmt.Errorf("user %s not found", userUUID)
	}

	// TODO: set the permissions
	p.ownerUUID = userUUID

	return nil
}

// SuggestionUpvote with user ID, song ID
func (p *Party) SuggestionUpvote(uid UserUUID, sid SongUID) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	if can, err := p.canUserPerformAction(uid, UserCanVoteSuggestionPermission); err != nil {
		return err
	} else if !can {
		return fmt.Errorf("user can't upvote")
	}

	err := p.suggestionQueue.Upvote(uid, sid)
	if err != nil {
		return err
	}

	p.setUpdated()
	return nil
}

// SuggestionDownvote with user ID, song ID
func (p *Party) SuggestionDownvote(uid UserUUID, sid SongUID) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	if can, err := p.canUserPerformAction(uid, UserCanVoteSuggestionPermission); err != nil {
		return err
	} else if !can {
		return fmt.Errorf("user can't downvote")
	}

	err := p.suggestionQueue.Downvote(uid, sid)
	if err != nil {
		return err
	}

	p.setUpdated()
	return nil
}

// SuggestionClearvote song to suggestion queue
func (p *Party) SuggestionClearvote(uid UserUUID, sid SongUID) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	if can, err := p.canUserPerformAction(uid, UserCanSuggestSongPermission); err != nil {
		return err
	} else if !can {
		return fmt.Errorf("user can't suggest")
	}

	err := p.suggestionQueue.ClearVotes(uid, sid)
	if err != nil {
		return err
	}

	p.setUpdated()
	return nil
}

// Suggest song to suggestion queue
func (p *Party) Suggest(uid UserUUID, sid SongUID) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	if can, err := p.canUserPerformAction(uid, UserCanSuggestSongPermission); err != nil {
		return err
	} else if !can {
		return fmt.Errorf("user can't suggest")
	}

	err := p.suggestionQueue.AddSong(uid, sid)
	if err != nil {
		return err
	}

	// check if there is a song currently playing
	if !p.nowPlaying.CurrentlyPlaying() {
		nextSid, err := p.getNextSong()

		// TODO: bigger check, this shouldn't happen
		if err != nil {
			return err
		}

		// set the currently playing song
		p.nowPlaying.ChangeSong(nextSid)
	}

	p.setUpdated()
	return nil
}

// Seek to a position in the song.
// Error if there isn't anything playing or the user doesn't
// have permission.
func (p *Party) Seek(uid UserUUID, position uint32) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	// check if teh user can seek
	can, err := p.canUserPerformAction(uid, UserCanSeekPermission)
	if err != nil {
		return err
	}

	if !can {
		return fmt.Errorf("user can not seek")
	}

	p.nowPlaying.Seek(position)
	p.setUpdated()

	return nil
}

// SongFinished is called when a song has finished playing.
func (p *Party) SongFinished(uid UserUUID, sid SongUID) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	// TODO: check that user is owner
	// TODO: check that current song actually ended

	nextSid, err := p.getNextSong()
	if err != nil {

		if nextSid == "" {
			// if there are no more songs to get
			p.nowPlaying.SetNonePlaying()
			p.setUpdated()

			return err
		}

		return err
	}

	// insert song
	p.nowPlaying.ChangeSong(nextSid)
	p.setUpdated()

	return nil
}

// Skip the currently playing song.
func (p *Party) Skip(uid UserUUID, sid SongUID) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	// TODO: check that they are on the actual current end song
	nextSid, err := p.getNextSong()
	if err != nil {

		if nextSid == "" {
			// if there are no more songs to get
			p.nowPlaying.SetNonePlaying()
			p.setUpdated()

			return err
		}

		return err
	}

	// insert song
	p.nowPlaying.ChangeSong(nextSid)
	p.setUpdated()

	return nil
}

// Pause the song
func (p *Party) Pause(uid UserUUID, pos uint32) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	if can, err := p.canUserPerformAction(uid, UserCanPlayPausePermission); err != nil {
		return err
	} else if !can {
		return fmt.Errorf("user can't play/pause")
	}

	if err := p.nowPlaying.SetPaused(pos); err != nil {
		return err
	}

	p.setUpdated()
	return nil
}

// Play the song
func (p *Party) Play(uid UserUUID) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	if can, err := p.canUserPerformAction(uid, UserCanPlayPausePermission); err != nil {
		return err
	} else if !can {
		return fmt.Errorf("user can't play/pause")
	}

	if err := p.nowPlaying.SetPlaying(); err != nil {
		return err
	}

	p.setUpdated()
	return nil
}

// SetVolume sets the volume for the player
func (p *Party) SetVolume(uid UserUUID, level uint32) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	// check that the user can perform this action
	if can, err := p.canUserPerformAction(uid, UserCanChangeVolumePermission); err != nil {
		return err
	} else if !can {
		return fmt.Errorf("user can't change volume")
	}

	// check for error, could be on bounds
	if err := p.nowPlaying.SetVolume(level); err != nil {
		return err
	}

	p.setUpdated()
	return nil
}

// playNextSong from the queues
func (p *Party) getNextSong() (SongUID, error) {
	nextSong, err := p.suggestionQueue.Pop()
	return nextSong, err
}

// updated increments the update tracker.
// Should call this whenever there's an update everyone should know about.
func (p *Party) setUpdated() {
	p.changeID++
	p.lastChangeT = time.Now()
}

// TimeSinceLastChange in duration
func (p *Party) TimeSinceLastChange() time.Duration {
	p.mux.Lock()
	defer p.mux.Unlock()

	return time.Since(p.lastChangeT)
}

// consts for pull
const (
	PullChangeKey     = "change"
	PullPlayingKey    = "playing"
	PullSuggestKey    = "suggest"
	PullPermissionKey = "permissions"
)

// Pull returns the user data in a serializable format.
// NOTE: this checks for changes before checking uid.
func (p *Party) Pull(userUUID UserUUID, clientChangeID uint64) (interface{}, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	// if the client's change is larger than our current change
	if p.changeID < clientChangeID {
		return nil, fmt.Errorf("bad pull id")
	}

	// up to date
	if p.changeID == clientChangeID {
		return nil, nil
	}

	_, err := p.getUser(userUUID)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data[PullChangeKey] = p.changeID
	data[PullPermissionKey] = p.permMap
	data[PullPlayingKey] = p.nowPlaying.Data()
	data[PullSuggestKey] = p.suggestionQueue.Pull(userUUID)

	return data, nil
}
