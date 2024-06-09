// Package users contains utilities for working with users.
package users

// UserType defines a class of user.
type UserType int

const (
	// UserTypeUser is a normal user.
	UserTypeUser UserType = iota

	// UserTypeAdmin is an administrative user.
	UserTypeAdmin
)

// User defines a user.
type User struct {
	ID       int
	Username string
	Type     UserType
}

// UserRepository provides a way of accessing users.
type UserRepository interface {
	FindUserByUsername(username string) (*User, error)
	GetAllUsersOfType(t UserType) ([]User, error)
}
