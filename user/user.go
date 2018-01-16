package user

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/user"
)

// User is a struct that stores the information and permissions for a user
type User struct {
	ID, Email, Name string
	Teacher, Admin  bool
}

// Get returns the current User
func Get(ctx context.Context) (*User, error) {
	u := user.Current(ctx)
	return &User{ID: u.ID, Email: u.Email, Name: "", Admin: true, Teacher: true}, nil
}

// func NewUser(user *User) (*User, error) {

// }
