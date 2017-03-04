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
