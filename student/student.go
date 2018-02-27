package student

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	"github.com/sunkink29/e3SelectionWebApp/errors"
	"github.com/sunkink29/e3SelectionWebApp/user"
)

const StudentNotFound = "Student not found"

type Student struct {
	ID                 string `datastore:"-"`
	Email, Name        string
	Grade int
	Teacher1, Teacher2 string
	Current            bool
}

// New stores the given student as a new student
func (stdnt *Student) New(ctx context.Context, debug bool) error {
	pKey := ParentKey(ctx, debug)
	k := datastore.NewIncompleteKey(ctx, "Student", pKey)
	k, err := datastore.Put(ctx, k, stdnt)
	if err != nil {
		return errors.New(err.Error())
	}
	stdnt.ID = k.Encode()
	return err
}

// Edit changes the student with the given ID to the values given
func (stdnt *Student) Edit(ctx context.Context) error {
	key, err := datastore.DecodeKey(stdnt.ID)
	if err != nil {
		return errors.New(err.Error())
	}
	_, err = datastore.Put(ctx, key, stdnt)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// Delete removes the student with the given key
func (stdnt *Student) Delete(ctx context.Context) error {
	key, err := datastore.DecodeKey(stdnt.ID)
	if err != nil {
		return errors.New(err.Error())
	}
	err = datastore.Delete(ctx, key)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// Current returns the current student
func Current(ctx context.Context, current bool, debug bool) (*Student, error) {
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return nil, err
	}
	stdnt, err := WithEmail(ctx, curU.Email, current, debug)
	if err != nil && err.(errors.Error).Message == StudentNotFound {
		newS := Student{Email: curU.Email, Name: curU.Name,Current: current}
		stdnt = &newS
		err = newS.New(ctx, debug)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return stdnt, nil
}

// WithKey returns the student with the given key
func WithKey(ctx context.Context, k *datastore.Key) (*Student, error) {
	stdnt := new(Student)
	err := datastore.Get(ctx, k, stdnt)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	stdnt.ID = k.Encode()
	return stdnt, nil
}

// WithEmail reterns the first student with matching email
func WithEmail(ctx context.Context, email string, current bool, debug bool) (*Student, error) {
	ancestor := ParentKey(ctx, debug)
	q := datastore.NewQuery("Student").Ancestor(ancestor).Filter("Email =", email).Filter("Current =", current)
	t := q.Run(ctx)
	var user Student
	key, err := t.Next(&user)
	if err == datastore.Done {
		return nil, errors.New(StudentNotFound)
	}
	if err != nil {
		return nil, errors.New(err.Error())
	}
	user.ID = key.Encode()
	return &user, nil
}

// All returns all of the students
func All(ctx context.Context, current bool, debug bool) ([]*Student, error) {
	ancestor := ParentKey(ctx, debug)
	q := datastore.NewQuery("Student").Ancestor(ancestor).Filter("Current =", current).Order("Name")
	var students []*Student
	keys, err := q.GetAll(ctx, &students)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	for i := 0; i < len(students); i++ {
		students[i].ID = keys[i].Encode()
	}
	return students, nil
}

func ParentKey(ctx context.Context, debug bool) *datastore.Key {
	var keyLiteral string
	if debug {
		keyLiteral = "Debug"
	} else {
		keyLiteral = "Release"
	}
	return datastore.NewKey(ctx, "Student", keyLiteral, 0, nil)
}
