package party

// User at a party
type User struct {
	name        string
	permissions map[string]bool
}

// UserUUID is the uuid for a user
type UserUUID string

// NewUser with name
func NewUser(name string) *User {
	return &User{
		name:        name,
		permissions: make(map[string]bool),
	}
}

// CanPerform checks if a user can perform an action
func (u User) CanPerform(action string) bool {
	canPerform, has := u.permissions[action]
	return canPerform && has
}

// SetPermission sets a permission for an action
func (u *User) SetPermission(action string, canPerform bool) {
	u.permissions[action] = canPerform
}

// Data satisfies the serializable interface
func (u User) Data() interface{} {
	return u.permissions
}

// user permissions
const (
	UserCanSeekPermission           = "Seek"
	UserCanVoteSuggestionPermission = "SuggestVote"
	UserCanSuggestSongPermission    = "Suggest"
	UserCanChangeVolumePermission   = "Volume"
	UserCanPlayPausePermission      = "PlayPause"
)

// maps can't be const in go
var (
	PermissionDescriptionMap = map[string]string{
		UserCanSeekPermission:           "Users may seek",
		UserCanSuggestSongPermission:    "Add songs to the suggestion queue",
		UserCanVoteSuggestionPermission: "Users can vote on songs in the suggestion queue",
		UserCanChangeVolumePermission:   "Users can change the music volume",
		UserCanPlayPausePermission:      "Users can play and pause music",
	}
)
