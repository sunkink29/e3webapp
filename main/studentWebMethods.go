package main

import (
	"encoding/json"
	"net/http"
	"fmt"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/sunkink29/e3SelectionWebApp/student"
	"github.com/sunkink29/e3SelectionWebApp/user"
	"github.com/sunkink29/e3SelectionWebApp/teacher"
	"github.com/sunkink29/e3SelectionWebApp/errors"
)

func addStudentMethods() {
	addToWebMethods("setTeacher", setTeacher)
	addToWebMethods("getCurrentStudentBlocks", getCurrentStudentBlocks)
	addToWebMethods("newStudent", newStudent)
	addToWebMethods("editStudent", editStudent)
	addToWebMethods("deleteStudent", deleteStudent)
	addToWebMethods("getAllStudents", getAllStudents)
}

func setTeacher(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	var variables struct { Teacher string; Block int}
	if err := dec.Decode(&variables); err != nil {
		return errors.New(err.Error())
	}
	
	usr, err := student.GetCurrent(ctx, false, debug)
	if err != nil {
		return err
	}
	newTeacher, err := teacher.GetWithEmail(ctx, variables.Teacher, false, debug)
	if err != nil {
		return err
	}
	if variables.Block == 0 {
		prevTeacher, err := teacher.GetWithEmail(ctx, usr.Teacher1, false, debug)
		prevOpen := true
		if err != nil && err.(errors.Error).Message != teacher.TeacherNotFound {
			return err
		} else if err == nil {
			prevOpen = prevTeacher.Block1.BlockOpen
		}
		newOpen := newTeacher.Block1.BlockOpen
		newFull := newTeacher.Block1.CurSize >= newTeacher.Block1.MaxSize
		
		if prevOpen && newOpen && !newFull {
			usr.Teacher1 = variables.Teacher
		}
	} else {
		prevTeacher, err := teacher.GetWithEmail(ctx, usr.Teacher2, false, debug)
		prevOpen := true
		if err != nil && err.(errors.Error).Message != teacher.TeacherNotFound {
			return err
		} else if err == nil {
			prevOpen = prevTeacher.Block2.BlockOpen
		}
		newOpen := newTeacher.Block2.BlockOpen
		newFull := newTeacher.Block2.CurSize >= newTeacher.Block2.MaxSize
		
		if prevOpen && newOpen && !newFull {
			usr.Teacher2 = variables.Teacher
		}
	}
	
	err = student.Edit(ctx, usr)
	return err
}

func getCurrentStudentBlocks(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	var current bool
	if err := dec.Decode(&current); err != nil {
		return errors.New(err.Error())
	}
	
	usr, err := student.GetCurrent(ctx, current, debug)
	if err != nil {
		return err
	}
	block1, err := teacher.GetWithEmail(ctx, usr.Teacher1, current, debug)
	if err != nil && err.(errors.Error).Message != teacher.TeacherNotFound {
		return err
	}
	block2, err := teacher.GetWithEmail(ctx, usr.Teacher2, current, debug)
	if err != nil && err.(errors.Error).Message != teacher.TeacherNotFound {
		return err
	}
	blocks := []*teacher.Teacher{block1, block2}
	jBlocks, err := json.Marshal(blocks)
	if err != nil {
		return errors.New(err.Error())
	}
	s := string(jBlocks[:])
	fmt.Fprintln(w, s)
	return nil
}

func newStudent(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.GetCurrent(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		newUsr := new(student.Student)
		if err := dec.Decode(newUsr); err != nil {
			return errors.New(err.Error())
		}
		err := student.New(ctx, newUsr, debug)
		return err
	}
	return errors.New(errors.AccessDenied)
}

func editStudent(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.GetCurrent(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		usr := new(student.Student)
		if err := dec.Decode(usr); err != nil {
			return errors.New(err.Error())
		}
		err := student.Edit(ctx, usr)
		return err
	}
	return errors.New(errors.AccessDenied)
}

func deleteStudent(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
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
		err = student.Delete(ctx, k)
		return err
	}
	return errors.New(errors.AccessDenied)
}

func getAllStudents(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	var current bool
	if err := dec.Decode(&current); err != nil {
			return errors.New(err.Error())
	}
	students, err := student.GetAll(ctx, current, debug)
	if err != nil {
		return err
	}

	jStudents, err := json.Marshal(students)
	if err != nil {
		return errors.New(err.Error())
	}
	s := string(jStudents[:])

	fmt.Fprintln(w, s)
	return nil
}
