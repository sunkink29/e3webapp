package errors

import (
	"runtime/debug"

	"golang.org/x/net/context"
	"google.golang.org/appengine/user"
	"google.golang.org/appengine/log"
)

const AccessDenied = "Access Denied"

type Error struct {
	Message string
	Stack   []byte
}

func (this Error) Error() string {
	return this.Message
}

func (this Error) HttpError(ctx context.Context) string {
	output := this.Message
	logging := string(this.Stack[:])
	log.Errorf(ctx, "%s\n%s", output, logging)
	if user.IsAdmin(ctx) {
		output += "\n"+logging
//		if err != nil {
//			output += "\n" + err.Error()
//		}
//		output = sLog

	}
	return output
}

type Redirect struct {
	URL string `json:"url"`
	Code int `json:"code"`
}

func (this Redirect) Error() string {
	return"redirect"
}

func New(input string) error {
	var output Error
	output.Message = input
	output.Stack = debug.Stack()
	return output
}
