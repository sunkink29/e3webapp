package main

import (
	"encoding/json"
	"net/http"
	"fmt"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/sunkink29/e3SelectionWebApp/teacher"
	"github.com/sunkink29/e3SelectionWebApp/user"
	"github.com/sunkink29/e3SelectionWebApp/errors"
)

func addTeacherMethods() {
	addToWebMethods("newTeacher", newTeacher)
	addToWebMethods("editTeacher", editTeacher)
	addToWebMethods("deleteTeacher", deleteTeacher)
	addToWebMethods("getAllTeachers", getAllTeachers)
}

func newTeacher(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.GetCurrent(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		newUsr := new(teacher.Teacher)
		if err := dec.Decode(newUsr); err != nil {
			return errors.New(err.Error())
		}
		err := teacher.New(ctx, newUsr, debug)
		return err
	}
	return errors.New(errors.AccessDenied)
}

func editTeacher(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.GetCurrent(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		usr := new(teacher.Teacher)
		if err := dec.Decode(usr); err != nil {
			return errors.New(err.Error())
		}
		err := teacher.Edit(ctx, usr)
		return err
	}
	return errors.New(errors.AccessDenied)
}

func deleteTeacher(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.GetCurrent(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		sKey := new(string)
		if err := dec.Decode(sKey); err != nil {
			return errors.New(err.Error())
		}
		k, err := datastore.DecodeKey(*sKey)
		if err != nil {
			return errors.New(err.Error())
		}
		err = teacher.Delete(ctx, k)
		return err
	}
	return errors.New(errors.AccessDenied)
}

func getAllTeachers(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	var current bool
	if err := dec.Decode(&current); err != nil {
			return errors.New(err.Error())
	}
	teachers, err := teacher.GetAll(ctx, current, debug)
	if err != nil {
		return err
	}

	jTeachers, err := json.Marshal(teachers)
	if err != nil {
		return errors.New(err.Error())
	}
	s := string(jTeachers[:])

	fmt.Fprintln(w, s)
	return nil
}
