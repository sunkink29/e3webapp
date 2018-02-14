package student

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	
	"github.com/sunkink29/e3SelectionWebApp/user"
	"github.com/sunkink29/e3SelectionWebApp/errors"
)

const StudentNotFound = "Student not found"

type Student struct {
	ID string `datastore:"-"`
	Email, Name string
	Teacher1, Teacher2 string
	Current bool
}

func New(ctx context.Context, student *Student, debug bool) error {
	pKey := ParentKey(ctx, debug)
	k := datastore.NewIncompleteKey(ctx, "Student", pKey)
	_, err := datastore.Put(ctx, k, student)
	if err != nil {
		return errors.New(err.Error())
	}
	return err
}

func Get(ctx context.Context, k *datastore.Key) (*Student, error) {
	var usr Student
	err := datastore.Get(ctx, k, &usr)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	usr.ID = k.Encode()
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &usr, nil
}

func GetCurrent(ctx context.Context, current bool, debug bool) (*Student, error) {
	curU, err := user.GetCurrent(ctx, debug)
	if err != nil {
		return nil, err
	}
	usr, err := GetWithEmail(ctx, curU.Email, current, debug)
	if err != nil && err.(errors.Error).Message == StudentNotFound {
		newS := Student{"",curU.Email,curU.Name,"","",current}
		usr = &newS
		err = New(ctx, &newS, debug)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return usr, nil
}

func GetWithEmail(ctx context.Context, email string, current bool, debug bool) (*Student, error) {
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

func GetAll(ctx context.Context, currentWeek bool, debug bool) ([]*Student, error) {
	ancestor := ParentKey(ctx, debug)
	q := datastore.NewQuery("Student").Ancestor(ancestor).Filter("Current =", currentWeek)
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

func Edit(ctx context.Context, student *Student) error {
	key, err := datastore.DecodeKey(student.ID)
	if err != nil {
		return errors.New(err.Error())
	}
	_, err = datastore.Put(ctx, key, student)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func Delete(ctx context.Context, k *datastore.Key) error {
	err := datastore.Delete(ctx, k)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
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

