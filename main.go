package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/sunkink29/e3SelectionWebApp/user"

	"google.golang.org/appengine"
)

var funcMap = map[string]interface{}{
	"includeHTML": includeHTML,
	"include":     include}
var indexTemplate = template.Must(template.New("index").Funcs(funcMap).ParseFiles("html/index.html"))
var fileMux sync.Mutex

var validMethods = map[string]func(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error{
	"print": func(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
		var str string
		err := dec.Decode(&str)
		if err != nil {
			return err
		}
		fmt.Fprintln(w, str)
		return nil
	},
}

func main() {
	appengine.Main()
}

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/async", async)
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	ctx := appengine.NewContext(r)
	r.ParseForm()
	u, _ := user.Get(ctx)
	err := indexTemplate.ExecuteTemplate(w, "index.html", u)
	if err != nil {
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
	err := method(dec, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func include(filename string) (string, error) {
	fileMux.Lock()
	file, err := ioutil.ReadFile(filename)
	fileMux.Unlock()
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
