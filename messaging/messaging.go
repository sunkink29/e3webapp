package messaging

import (
	"encoding/json"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"context"

	"firebase.google.com/go"
	"firebase.google.com/go/messaging"

	"github.com/sunkink29/e3webapp/errors"
	"github.com/sunkink29/e3webapp/user"
)

type Credentials struct {
	APIKey      string
	ID          string
	SenderID    string
	FirebaseKey string
}

func APIKey() string {
	return cred.APIKey
}

func ID() string {
	return cred.ID
}

func SenderID() string {
	return cred.SenderID
}

func FirebaseKey() string {
	return cred.FirebaseKey
}

// func

type topic struct {
	value string
}

var Topics = struct {
	Student, Teacher, Admin topic
}{topic{"student"}, topic{"teacher"}, topic{"admin"}}

type event struct {
	value string
}

var EventTypes = struct {
	Popup, ClassEdit, CurrentChange, StudentUpdate event
}{event{"popup"}, event{"classEdit"}, event{"currentChange"}, event{"studentUpdate"}}

var conf *firebase.Config
var cred *Credentials

func InitAuth(ctx context.Context) error {
	if conf == nil {
		key := datastore.NewKey(ctx, "Auth", "firebase", 0, nil)
		cred = new(Credentials)
		err := datastore.Get(ctx, key, cred)
		if err != nil {
			return errors.New(err.Error())
		}

		conf = &firebase.Config{
			ProjectID: cred.ID,
		}

		return nil
	}
	return nil
}

func getClient(ctx context.Context) (*messaging.Client, error) {
	if err := InitAuth(ctx); err != nil {
		return nil, err
	}

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return client, nil
}

func RegisterTopicHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	curU, err := user.Current(ctx)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(r.Body)
	token := new(string)
	if err := decoder.Decode(token); err != nil {
		return errors.New(err.Error())
	}

	curU.RToken = *token
	curU.Edit(ctx)

	if !curU.Teacher && !curU.Admin {
		err = registerTopic(ctx, *token, Topics.Student)
	} else {
		if curU.Teacher {
			err = registerTopic(ctx, *token, Topics.Teacher)
		}
		if curU.Admin {
			err = registerTopic(ctx, *token, Topics.Admin)
		}
	}

	return err
}

func registerTopic(ctx context.Context, token string, group topic) error {
	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	// return errors.New(string(group))
	tokens := []string{token}
	response, err := client.SubscribeToTopic(ctx, tokens, group.value)
	if err != nil {
		return errors.New(err.Error())
	}
	if response.FailureCount > 0 {
		return errors.New(response.Errors[0].Reason)
	}
	return nil
}

func SendEvent(ctx context.Context, event event, data string, group topic) error {
	message := messaging.Message{
		Data: map[string]string{
			"event": event.value,
			"data":  data,
		},
		Topic: group.value,
	}

	return sendMessage(ctx, message)
}

func SendUserEvent(ctx context.Context, event event, data, Token string) error {
	message := messaging.Message{
		Data: map[string]string{
			"event": event.value,
			"data":  data,
		},
		Token: Token,
	}

	return sendMessage(ctx, message)
}

func sendMessage(ctx context.Context, message messaging.Message) error {
	client, err := getClient(ctx)
	if err != nil {
		return errors.New(err.Error())
	}

	if _, err = client.Send(ctx, &message); err != nil {
		return errors.New(err.Error())
	}
	return nil
}
