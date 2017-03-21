package event

import (
	"net/http"
)
// User at a party
type User struct {
	name        string
	permissions map[string]bool
}

// UserUID is the uuid for a user
type UserUID string

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

func (s *Server) CreateEvent(w http.ResponseWriter, r *http.Request) {
// TODO: Implement
}

func (s *Server) EndEvent(w http.ResponseWriter, r *http.Request) {
// TODO: Implement
}

func (u *User) JoinEvent(w http.ResponseWriter, r *http.Request) {
// TODO: Implement
}

func (u *User) LeaveEvent(w http.ResponseWriter, r *http.Request) {
// TODO: Implement
}
