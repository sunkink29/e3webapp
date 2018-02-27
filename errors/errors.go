package errors

import (
//	"fmt"
//	"os"
//	"encoding/json"
	"runtime/debug"
	"net/http"

	"golang.org/x/net/context"
//	"cloud.google.com/go/errorreporting"
//	"google.golang.org/api/option"
//	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
	"google.golang.org/appengine/log"
)

const AccessDenied = "Access Denied"

type Error struct {
	Message string
	Stack   []byte
}

type ErrorLog struct {
	ServiceContent struct {
		Service string `json:"service"`
	} `json:"serviceContent"`
	Message string `json:"message"`
	Content struct {
		HttpRequest struct {
			Method             string `json:"method"`
			Url                string `json:"url"`
			ResponseStatusCode int    `json:"responseStatusCode"`
		} `json:"httpRequest"`
		User string `json:"user"`
	} `json:"Content"`
}

func (this Error) Error() string {
	return this.Message
}

func (this Error) HttpError(ctx context.Context, usr string, url string, r *http.Request) string {
	output := this.Message
	logging := string(this.Stack[:])
//	errLog := new(ErrorLog)
//	errLog.ServiceContent.Service = "frontend"
//	errLog.Message = fmt.Sprintf("%v\n%v", output, logging)
//	errLog.Content.HttpRequest.Method = "Get"
//	errLog.Content.HttpRequest.Url = url
//	errLog.Content.HttpRequest.ResponseStatusCode = 500
//	errLog.Content.User = usr
//	jLog, _ := json.Marshal(errLog)
//	sLog := string(jLog[:])
//	fmt.Fprint(os.Stderr, sLog)

//	key := datastore.NewKey(ctx, "Auth", "Auth", 0, nil)
//	var Auth struct {ID string; APIKey string}
//	datastore.Get(ctx, key, &Auth)
//	ctf := errorreporting.Config{}
//	client, _ := errorreporting.NewClient(ctx, Auth.ID, ctf, option.WithAPIKey(Auth.APIKey))
//	entry := errorreporting.Entry{this, r, this.Stack}
//	err = client.ReportSync(ctx,entry)
//	defer client.Close()
//	defer client.Flush()

	if user.IsAdmin(ctx) {
		output += "\n"+logging
//		if err != nil {
//			output += "\n" + err.Error()
//		}
//		output = sLog

	}
	log.Errorf(ctx, "%s\n%s", output, logging)
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
