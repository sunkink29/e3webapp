package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/sunkink29/e3webapp/errors"
	"github.com/sunkink29/e3webapp/student"
	"github.com/sunkink29/e3webapp/teacher"
	"github.com/sunkink29/e3webapp/user"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type webMethod func(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error

var funcMap = map[string]interface{}{
	"includeHTML": includeHTML,
	"include":     include,
	"clientID": user.ClientID}
var indexTemplate = template.Must(template.New("index").Funcs(funcMap).ParseFiles("html/index.html"))
var fileMux sync.Mutex

var validMethods map[string]webMethod

func main() {
	appengine.Main()
}

func init() {
	validMethods = make(map[string]webMethod)
	addAdminMethods()
	addTeacherMethods()
	addStudentMethods()

	http.HandleFunc("/", root)
	http.HandleFunc("/async", async)
	http.HandleFunc("/usrswitch", usrswitch)
	http.HandleFunc("/auth", user.AuthHandle)
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	ctx := appengine.NewContext(r)
	r.ParseForm()
	debug := r.Form.Get("debug") == "true"
	u, err := user.Current(ctx, debug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//	u.Admin = false
	//	u.Teacher = false
	
	user.InitAuth(ctx)
	
	err = indexTemplate.ExecuteTemplate(w, "index.html", u)
	if err != nil {
		err = errors.New(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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
			url := r.URL.String()
//			debug := r.Form.Get("debug") == "true"
//			usr,_ := user.Current(ctx, debug)
			s := "" //usr.ID
			http.Error(w, err.(errors.Error).HttpError(ctx, s, url, r), http.StatusInternalServerError)
			return
		}
		return
	}
	return
}

func usrswitch(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	err := switchNextToCurrent(ctx, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func switchNextToCurrent(ctx context.Context, debug bool) error {
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
		newS := student.Student{Name: stdnt.Name, Email: stdnt.Email}
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

func addToWebMethods(name string, method webMethod) {
	validMethods[name] = method
}

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

func includeHTML(filename string) (template.HTML, error) {
	text, err := include(filename)
	return template.HTML(text), err
}
