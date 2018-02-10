package student

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	
	"github.com/sunkink29/e3SelectionWebApp/user"
)
type Error string

func (e Error) Error() string { return string(e) }

const StudentNotFound = Error("Student not found")

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

	return err
}

func GetCurrent(ctx context.Context, current bool, debug bool) (*Student, error) {
	curU, err := user.GetCurrent(ctx, debug)
	if err != nil {
		return nil, err
	}
	usr, err := GetWithEmail(ctx, curU.Email, current, debug)
	if err == StudentNotFound {
		newS := Student{"",curU.Email,curU.Name,"","",current}
		New(ctx, &newS, debug)
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
		return nil, StudentNotFound
	}
	if err != nil {
		return nil, err
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
		return nil, err
	}
	for i := 0; i < len(students); i++ {
		students[i].ID = keys[i].Encode()
	}
	return students, nil
}

func Edit(ctx context.Context, student *Student) error {
	key, err := datastore.DecodeKey(student.ID)
	if err != nil {
		return err
	}
	_, err = datastore.Put(ctx, key, student)
	return err
}

func Delete(ctx context.Context, k *datastore.Key) error {
	err := datastore.Delete(ctx, k)
	return err
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

