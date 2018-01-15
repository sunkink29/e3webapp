package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"sync"

	"google.golang.org/appengine"
)

var funcMap = map[string]interface{}{
	"includeHTML": includeHTML,
	"include":     include}
var indexTemplate = template.Must(template.New("index").Funcs(funcMap).ParseFiles("html/index.html"))
var mux sync.Mutex

type options struct {
	Admin, Teacher bool
}

func main() {
	appengine.Main()
}

func init() {
	http.HandleFunc("/", doGet)
}

func doGet(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "text/html; charset=utf-8")
	// ctx := appengine.NewContext(r)
	// u := user.Current(ctx)
	r.ParseForm()
	opt := options{Admin: r.Form.Get("admin") == "true", Teacher: r.Form.Get("teacher") == "true"}
	err := indexTemplate.ExecuteTemplate(w, "index.html", opt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func include(filename string) (string, error) {
	mux.Lock()
	file, err := ioutil.ReadFile(filename)
	mux.Unlock()
	if err != nil {
		return "", err
	}
	s := fmt.Sprintf("%s", file)
	return s, nil
}

func includeHTML(filename string) (template.HTML, error) {
	text, err := include(filename)
	return template.HTML(text), err
}

func includeJS(filename string) (template.JS, error) {
	text, err := include(filename)
	return template.JS(text), err
}

func includeCSS(filename string) (template.CSS, error) {
	text, err := include(filename)
	return template.CSS(text), err
}
