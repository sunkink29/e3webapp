package user

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"

	"github.com/sunkink29/e3SelectionWebApp/errors"
)

// User is a struct that stores the information and permissions for a user
type User struct {
	ID             string `datastore:"-"`
	Email, Name    string
	Teacher, Admin bool
}

// GetCurrent returns the current User
func GetCurrent(ctx context.Context, debug bool) (*User, error) {
	u := user.Current(ctx)
	user, err := GetWithEmail(ctx, u.Email, debug)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetWithEmail reterns the first user with matching email
func GetWithEmail(ctx context.Context, email string, debug bool) (*User, error) {
	ancestor := parentKey(ctx, debug)
	q := datastore.NewQuery("User").Ancestor(ancestor).Filter("Email =", email)
	t := q.Run(ctx)
	var user User
	key, err := t.Next(&user)
	if err == datastore.Done {
		return nil, errors.New("User not found")
	}
	if err != nil {
		return nil, errors.New(err.Error())
	}
	user.ID = key.Encode()
	return &user, nil
}

// GetAll returns all users
func GetAll(ctx context.Context, debug bool) ([]*User, error) {
	ancestor := parentKey(ctx, debug)
	q := datastore.NewQuery("User").Ancestor(ancestor)
	var users []*User
	keys, err := q.GetAll(ctx, &users)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	for i := 0; i < len(users); i++ {
		users[i].ID = keys[i].Encode()
	}
	return users, nil
}

// Get returns the user with the given key
func Get(ctx context.Context, k *datastore.Key) (*User, error) {
	var usr *User
	err := datastore.Get(ctx, k, usr)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return usr,nil
}

// New stores the given user as a new user
func New(ctx context.Context, newUsr *User, debug bool) error {
	pKey := parentKey(ctx, debug)
	k := datastore.NewIncompleteKey(ctx, "User", pKey)
	_, err := datastore.Put(ctx, k, newUsr)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// Edit changes the user with the given ID to the values given
func Edit(ctx context.Context, user *User) error {
	key, err := datastore.DecodeKey(user.ID)
	if err != nil {
		return errors.New(err.Error())
	}
	_, err = datastore.Put(ctx, key, user)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// Delete removes the user with the given key
func Delete(ctx context.Context,  k *datastore.Key) error{
	err := datastore.Delete(ctx, k)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func parentKey(ctx context.Context, debug bool) *datastore.Key {
	var keyLiteral string
	if debug {
		keyLiteral = "Debug"
	} else {
		keyLiteral = "Release"
	}
	return datastore.NewKey(ctx, "User", keyLiteral, 0, nil)
}
