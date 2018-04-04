package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/sunkink29/e3webapp/errors"
	"github.com/sunkink29/e3webapp/messaging"
	"github.com/sunkink29/e3webapp/student"
	"github.com/sunkink29/e3webapp/teacher"
	"github.com/sunkink29/e3webapp/user"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var funcMap = map[string]interface{}{
	"includeHTML":    includeHTML,
	"include":        include,
	"clientID":       user.ClientID,
	"apiKey":         user.ApiKey,
	"isDevServer":    isDevServer,
	"firebaseKey":    messaging.FirebaseKey,
	"projectID":      messaging.ID,
	"firebaseApiKey": messaging.APIKey,
	"senderID":       messaging.SenderID}
var indexTemplate = template.Must(template.New("index").Funcs(funcMap).ParseFiles("html/index.html"))
var fileMux sync.Mutex

func include(filename string) (string, error) {
	fileMux.Lock()
	file, err := ioutil.ReadFile(filename)
	fileMux.Unlock()
	if err != nil {
		return "", errors.New(err.Error())
	}
	s := fmt.Sprintf("%s", file)
	return s, nil
}

func isDevServer() bool {
	return appengine.IsDevAppServer()
}

func includeHTML(filename string) (template.HTML, error) {
	text, err := include(filename)
	return template.HTML(text), err
}

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	user.InitAuth(ctx)
	messaging.InitAuth(ctx)
	k := datastore.NewKey(ctx, "lock", "lock", 0, nil)
	lock := new(struct{ lock bool })
	err := datastore.Get(ctx, k, lock)
	if true || err != nil {
		r.ParseForm()
		if err := fn(w, r); err != nil {
			http.Error(w, err.(errors.Error).HttpError(ctx), 500)
		}
	} else {
		http.Error(w, "Database is locked currently\n check back in 5 minutes", 500)
	}
}

func main() {
	appengine.Main()
}

func init() {
	//	validMethods = make(map[string]webMethod)
	addAdminMethods()
	addTeacherMethods()
	addStudentMethods()

	http.Handle("/", appHandler(root))
	http.Handle("/api/registertoken", appHandler(messaging.RegisterTopicHandler))
	//	http.HandleFunc("/async", async)
	http.Handle("/worker/usrswitch", appHandler(usrswitch))
	http.Handle("/worker/importusers", appHandler(importUsers))
	http.HandleFunc("/auth", user.AuthHandle)
}

func root(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	ctx := appengine.NewContext(r)
	r.ParseForm()
	debug := r.Form.Get("debug") == "true"
	usr, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}

	err = indexTemplate.ExecuteTemplate(w, "index.html", usr)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

/*
func async(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	methodStr := r.Form.Get("method")
	method, ok := validMethods[methodStr]
	if !ok {
		http.Error(w, "method "+methodStr+" not found", http.StatusUnprocessableEntity)
		return
	}
	strReader := strings.NewReader(r.Form.Get("data"))
	dec := json.NewDecoder(strReader)
	ctx := appengine.NewContext(r)
	k := datastore.NewKey(ctx, "lock", "lock", 0, nil)
	lock := new(struct{ lock bool })
	err := datastore.Get(ctx, k, lock)
	if err != nil {
		err := method(dec, w, r)
		if err != nil {
			if err.Error() == "redirect" {
				redirect := err.(errors.Redirect)
				http.Redirect(w , r, redirect.URL, redirect.Code)
				return
			}
			http.Error(w, err.(errors.Error).HttpError(ctx), http.StatusInternalServerError)
			return
		}
		return
	}
	return
}
*/

func usrswitch(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := false

	lockKey := datastore.NewKey(ctx, "lock", "lock", 0, nil)
	lock := new(struct{ lock bool })
	_, err := datastore.Put(ctx, lockKey, lock)
	if err != nil {
		return errors.New(err.Error())
	}

	teachers, err := teacher.All(ctx, true, debug)
	if err != nil {
		return err
	}
	for _, tchr := range teachers {
		err = tchr.Delete(ctx)
		if err != nil {
			return err
		}
	}

	teachers, err = teacher.All(ctx, false, debug)
	if err != nil {
		return err
	}
	for _, tchr := range teachers {
		tchr.Current = true
		err = tchr.Edit(ctx)
		if err != nil {
			return err
		}
	}

	students, err := student.All(ctx, true, debug)
	if err != nil {
		return err
	}
	for _, stdnt := range students {
		err = stdnt.Delete(ctx)
		if err != nil {
			return err
		}
	}

	students, err = student.All(ctx, false, debug)
	if err != nil {
		return err
	}
	for _, stdnt := range students {
		stdnt.Current = true
		err = stdnt.Edit(ctx)
		if err != nil {
			return err
		}
		newS := student.Student{Name: stdnt.Name, Email: stdnt.Email, Grade: stdnt.Grade}
		err = newS.New(ctx, debug)
		if err != nil {
			return err
		}
	}
	err = datastore.Delete(ctx, lockKey)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}
