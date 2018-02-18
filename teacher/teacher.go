package teacher

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	appUser "google.golang.org/appengine/user"

	"github.com/sunkink29/e3SelectionWebApp/errors"
	"github.com/sunkink29/e3SelectionWebApp/student"
)

const TeacherNotFound = "Teacher not found"

type Block struct {
	Subject, Description string
	CurSize              int `datastore:"-"`
	MaxSize, RoomNumber  int
	BlockOpen            bool
}

type Teacher struct {
	ID             string `datastore:"-"`
	Email, Name    string
	Block1, Block2 Block
	Current        bool
}

func (t *Teacher) Load(ps []datastore.Property) error {
	for _, p := range ps {
		switch p.Name {
		case "Email":
			t.Email = p.Value.(string)
		case "Name":
			t.Name = p.Value.(string)
		case "Subject1":
			t.Block1.Subject = p.Value.(string)
		case "Description1":
			t.Block1.Description = p.Value.(string)
		case "RoomNumber1":
			t.Block1.RoomNumber = int(p.Value.(int64))
		case "MaxSize1":
			t.Block1.MaxSize = int(p.Value.(int64))
		case "BlockOpen1":
			t.Block1.BlockOpen = p.Value.(bool)
		case "Subject2":
			t.Block2.Subject = p.Value.(string)
		case "Description2":
			t.Block2.Description = p.Value.(string)
		case "RoomNumber2":
			t.Block2.RoomNumber = int(p.Value.(int64))
		case "MaxSize2":
			t.Block2.MaxSize = int(p.Value.(int64))
		case "BlockOpen2":
			t.Block2.BlockOpen = p.Value.(bool)
		case "Current":
			t.Current = p.Value.(bool)
		}
	}
	return nil
}

func (t *Teacher) Save() ([]datastore.Property, error) {
	return []datastore.Property{
		{
			Name:  "Email",
			Value: t.Email,
		}, {
			Name:  "Name",
			Value: t.Name,
		}, {
			Name:  "Subject1",
			Value: t.Block1.Subject,
		}, {
			Name:  "Description1",
			Value: t.Block1.Description,
		}, {
			Name:  "RoomNumber1",
			Value: int64(t.Block1.RoomNumber),
		}, {
			Name:  "MaxSize1",
			Value: int64(t.Block1.MaxSize),
		}, {
			Name:  "BlockOpen1",
			Value: t.Block1.BlockOpen,
		}, {
			Name:  "Subject2",
			Value: t.Block2.Subject,
		}, {
			Name:  "Description2",
			Value: t.Block2.Description,
		}, {
			Name:  "RoomNumber2",
			Value: int64(t.Block2.RoomNumber),
		}, {
			Name:  "MaxSize2",
			Value: int64(t.Block2.MaxSize),
		}, {
			Name:  "BlockOpen2",
			Value: t.Block2.BlockOpen,
		}, {
			Name:  "Current",
			Value: t.Current,
		},
	}, nil
}

// New stores the given teacher as a new teacher
func (tchr *Teacher) New(ctx context.Context, debug bool) error {
	pKey := parentKey(ctx, debug)
	k := datastore.NewIncompleteKey(ctx, "Teacher", pKey)
	k, err := datastore.Put(ctx, k, tchr)
	if err != nil {
		return errors.New(err.Error())
	}
	tchr.ID = k.Encode()
	return nil
}

// Edit changes the teacher with the given ID to the values given
func (tchr *Teacher) Edit(ctx context.Context) error {
	key, err := datastore.DecodeKey(tchr.ID)
	if err != nil {
		return errors.New(err.Error())
	}
	_, err = datastore.Put(ctx, key, tchr)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// StudentCount returns the student count for the given block
func (tchr *Teacher) StudentCount(ctx context.Context, block int, debug bool) (int, error) {
	ancestor := student.ParentKey(ctx, debug)
	var sBlock string
	if block == 0 {
		sBlock = "Teacher1"
	} else {
		sBlock = "Teacher2"
	}
	q := datastore.NewQuery("Student").Ancestor(ancestor).Filter("Current =", tchr.Current).Filter(sBlock+" =", tchr.Email)
	count, err := q.Count(ctx)
	if err != nil {
		return 0, errors.New(err.Error())
	}
	return count, err
}

// StudentList returns a list of the students in the block given of the given teacher
func (tchr *Teacher) StudentList(ctx context.Context, block int, debug bool) ([]*student.Student, error) {
	ancestor := student.ParentKey(ctx, debug)
	var sBlock string
	if block == 0 {
		sBlock = "Teacher1"
	} else {
		sBlock = "Teacher2"
	}
	q := datastore.NewQuery("Student").Ancestor(ancestor).Filter("Current =", tchr.Current).Filter(sBlock+" =", tchr.Email)
	var students = make([]*student.Student, 0)
	keys, err := q.GetAll(ctx, &students)
	if err == datastore.Done {
		return make([]*student.Student, 0), errors.New(student.StudentNotFound)
	}
	if err != nil {
		return make([]*student.Student, 0), errors.New(err.Error())
	}
	for i := 0; i < len(students); i++ {
		students[i].ID = keys[i].Encode()
	}
	return students, nil
}

// Delete removes the teacher with the given key
func (tchr *Teacher) Delete(ctx context.Context) error {
	key, err := datastore.DecodeKey(tchr.ID)
	if err != nil {
		return errors.New(err.Error())
	}
	err = datastore.Delete(ctx, key)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// Current returns the current Teacher
func Current(ctx context.Context, current bool, debug bool) (*Teacher, error) {
	usr := appUser.Current(ctx)
	tchr, err := WithEmail(ctx, usr.Email, current, debug)
	if err != nil {
		return nil, err
	}
	return tchr, nil
}

// WithKey returns the teacher with the given key
func WithKey(ctx context.Context, k *datastore.Key, debug bool) (*Teacher, error) {
	var tchr = new(Teacher)
	err := datastore.Get(ctx, k, tchr)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	tchr.ID = k.Encode()
	count, err := tchr.StudentCount(ctx, 0, debug)
	if err != nil {
		return nil, err
	}
	tchr.Block1.CurSize = count

	count, err = tchr.StudentCount(ctx, 1, debug)
	if err != nil {
		return nil, err
	}
	tchr.Block2.CurSize = count
	return tchr, nil
}

// WithEmail reterns the first teacher with matching email
func WithEmail(ctx context.Context, email string, current bool, debug bool) (*Teacher, error) {
	ancestor := parentKey(ctx, debug)
	q := datastore.NewQuery("Teacher").Ancestor(ancestor).Filter("Email =", email).Filter("Current =", current)
	t := q.Run(ctx)
	var tchr Teacher
	key, err := t.Next(&tchr)
	if err == datastore.Done {
		return nil, errors.New(TeacherNotFound)
	}
	if err != nil {
		return nil, errors.New(err.Error())
	}
	tchr.ID = key.Encode()

	count, err := tchr.StudentCount(ctx, 0, debug)
	if err != nil {
		return nil, err
	}
	tchr.Block1.CurSize = count

	count, err = tchr.StudentCount(ctx, 1, debug)
	if err != nil {
		return nil, err
	}
	tchr.Block2.CurSize = count
	return &tchr, nil
}

// All returns all of the teachers 
func All(ctx context.Context, current bool, debug bool) ([]*Teacher, error) {
	ancestor := parentKey(ctx, debug)
	q := datastore.NewQuery("Teacher").Ancestor(ancestor).Filter("Current =", current)
	var tchrs []*Teacher
	keys, err := q.GetAll(ctx, &tchrs)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	for i, tchr := range tchrs {
		tchr.ID = keys[i].Encode()
		count, err := tchr.StudentCount(ctx, 0, debug)
		if err != nil {
			return nil, err
		}
		tchr.Block1.CurSize = count

		count, err = tchr.StudentCount(ctx, 1, debug)
		if err != nil {
			return nil, err
		}
		tchr.Block2.CurSize = count
	}
	return tchrs, nil
}



func parentKey(ctx context.Context, debug bool) *datastore.Key {
	var keyLiteral string
	if debug {
		keyLiteral = "Debug"
	} else {
		keyLiteral = "Release"
	}
	return datastore.NewKey(ctx, "Teacher", keyLiteral, 0, nil)
}
