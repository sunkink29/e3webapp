package main

import (
	"bytes"
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
	admin   bool
	teacher bool
	url     string
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
	err := indexTemplate.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func include(filename string) string {
	mux.Lock()
	file, err := ioutil.ReadFile(filename)
	mux.Unlock()
	if err != nil {
		return err.Error()
	}
	n := bytes.IndexByte(file, 0)
	s := string(file[:n])
	return s
}

func includeHTML(filename string) template.HTML {
	return template.HTML(include(filename))
}

func includeJS(filename string) template.JS {
	return template.JS(include(filename))
}

func includeCSS(filename string) template.CSS {
	return template.CSS(include(filename))
}
