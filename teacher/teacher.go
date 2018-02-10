package teacher

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	
	"github.com/sunkink29/e3SelectionWebApp/student"
)
type Error string

func (e Error) Error() string { return string(e) }

const TeacherNotFound = Error("Teacher not found")

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
	Current bool
}

func (t *Teacher) Load(ps []datastore.Property) error {
	for _,p := range ps {
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

func New(ctx context.Context, teacher *Teacher, debug bool) error {
	pKey := parentKey(ctx, debug)
	k := datastore.NewIncompleteKey(ctx, "Teacher", pKey)
	_, err := datastore.Put(ctx, k, teacher)

	return err
}

func Get(ctx context.Context, k *datastore.Key) (*Teacher, error) {
	var usr *Teacher
	err := datastore.Get(ctx, k, usr)
	usr.ID = k.Encode()
	return usr, err
}

func GetWithEmail(ctx context.Context, email string, current bool, debug bool) (*Teacher, error) {
	ancestor := parentKey(ctx, debug)
	q := datastore.NewQuery("Teacher").Ancestor(ancestor).Filter("Email =", email).Filter("Current =", current)
	t := q.Run(ctx)
	var user Teacher
	key, err := t.Next(&user)
	if err == datastore.Done {
		return nil, TeacherNotFound
	}
	if err != nil {
		return nil, err
	}
	user.ID = key.Encode()
	
	count, err := GetStudentCount(ctx, user.Email, 0, current, debug)
	if err != nil {
		return nil, err
	}
	user.Block1.CurSize = count
	
	count, err = GetStudentCount(ctx, user.Email, 1, current, debug)	
	if err != nil {
		return nil, err
	}
	user.Block2.CurSize = count
	return &user, nil
}

func GetAll(ctx context.Context, currentWeek bool, debug bool) ([]*Teacher, error) {
	ancestor := parentKey(ctx, debug)
	q := datastore.NewQuery("Teacher").Ancestor(ancestor).Filter("Current =", currentWeek)
	var teachers []*Teacher
	keys, err := q.GetAll(ctx, &teachers)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(teachers); i++ {
		teachers[i].ID = keys[i].Encode()
		count, err := GetStudentCount(ctx, teachers[i].Email, 0, currentWeek, debug)
		if err != nil {
			return nil, err
		}
		teachers[i].Block1.CurSize = count
		
		count, err = GetStudentCount(ctx, teachers[i].Email, 1, currentWeek, debug)	
		if err != nil {
			return nil, err
		}
		teachers[i].Block2.CurSize = count
	}
	return teachers, nil
}

func GetStudentCount(ctx context.Context, email string, block int, current bool, debug bool) (int, error) {
	ancestor := student.ParentKey(ctx, debug)
	var sBlock string
	if block == 0 {
		sBlock = "Teacher1"
	} else {
		sBlock = "Teacher2"
	}
	q := datastore.NewQuery("Student").Ancestor(ancestor).Filter("Current =", current).Filter(sBlock+" =", email)
	return q.Count(ctx)
}

func Edit(ctx context.Context, teacher *Teacher) error {
	key, err := datastore.DecodeKey(teacher.ID)
	if err != nil {
		return err
	}
	_, err = datastore.Put(ctx, key, teacher)
	return err
}

func Delete(ctx context.Context, k *datastore.Key) error {
	err := datastore.Delete(ctx, k)
	return err
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
