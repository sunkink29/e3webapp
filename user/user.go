package user

import (
	"encoding/json"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/appengine/datastore"
	appUser "google.golang.org/appengine/user"

	"github.com/sunkink29/e3webapp/errors"
)

// User is a struct that stores the information and permissions for a user
type User struct {
	ID             string `datastore:"-"`
	Email, Name    string
	Teacher, Admin bool
	AuthState      string        `json:"-"`
	Token          *oauth2.Token `json:"-"`
	RToken         string        `json:"-"`
}

func (u *User) Load(ps []datastore.Property) error {
	for _, p := range ps {
		switch p.Name {
		case "Email":
			u.Email = p.Value.(string)
		case "Name":
			u.Name = p.Value.(string)
		case "Teacher":
			u.Teacher = p.Value.(bool)
		case "Admin":
			u.Admin = p.Value.(bool)
		case "AuthState":
			u.AuthState = p.Value.(string)
		case "Token":
			if p.Value.(string) != "null" {
				u.Token = new(oauth2.Token)
				tByte := []byte(p.Value.(string))
				err := json.Unmarshal(tByte, u.Token)
				if err != nil {
					u.Token = nil
				}
			}
		case "RToken":
			u.RToken = p.Value.(string)
		}
	}
	return nil
}

func (u *User) Save() ([]datastore.Property, error) {
	bToken, err := json.Marshal(u.Token)
	sToken := string(bToken[:])
	return []datastore.Property{
		{
			Name:  "Email",
			Value: u.Email,
		}, {
			Name:  "Name",
			Value: u.Name,
		}, {
			Name:  "Teacher",
			Value: u.Teacher,
		}, {
			Name:  "Admin",
			Value: u.Admin,
		}, {
			Name:  "AuthState",
			Value: u.AuthState,
		}, {
			Name:  "Token",
			Value: sToken,
		}, {
			Name:  "RToken",
			Value: u.RToken,
		},
	}, err
}

// New stores the given user as a new user
func (usr *User) New(ctx context.Context) error {
	pKey := parentKey(ctx)
	k := datastore.NewIncompleteKey(ctx, "User", pKey)
	newK, err := datastore.Put(ctx, k, usr)
	if err != nil {
		return errors.New(err.Error())
	}
	usr.ID = newK.Encode()
	return nil
}

// Edit changes the user with the given ID to the values given
func (usr *User) Edit(ctx context.Context) error {
	key, err := datastore.DecodeKey(usr.ID)
	if err != nil {
		return errors.New(err.Error())
	}
	_, err = datastore.Put(ctx, key, usr)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// Delete removes the user with the given key
func (usr *User) Delete(ctx context.Context) error {
	key, err := datastore.DecodeKey(usr.ID)
	if err != nil {
		return errors.New(err.Error())
	}

	err = datastore.Delete(ctx, key)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// Current returns the current User
func Current(ctx context.Context) (*User, error) {
	u := appUser.Current(ctx)
	user, err := WithEmail(ctx, u.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// WithKey returns the user with the given key
func WithKey(ctx context.Context, k *datastore.Key) (*User, error) {
	var usr *User
	err := datastore.Get(ctx, k, usr)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	usr.ID = k.Encode()
	return usr, nil
}

// WithEmail reterns the first user with matching email
func WithEmail(ctx context.Context, email string) (*User, error) {
	ancestor := parentKey(ctx)
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
func All(ctx context.Context) ([]*User, error) {
	ancestor := parentKey(ctx)
	q := datastore.NewQuery("User").Ancestor(ancestor).Order("Name")
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

func parentKey(ctx context.Context) *datastore.Key {
	return datastore.NewKey(ctx, "User", "Release", 0, nil)
}
