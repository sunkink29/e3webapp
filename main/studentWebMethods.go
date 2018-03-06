package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/sunkink29/e3webapp/errors"
	"github.com/sunkink29/e3webapp/student"
	"github.com/sunkink29/e3webapp/teacher"
	"github.com/sunkink29/e3webapp/user"
)

func addStudentHandle(url string, handle appHandler) {
	http.Handle("/api/student/"+url, handle)
}

func addStudentMethods() {
	addStudentHandle("setteacher", appHandler(setTeacher))
	addStudentHandle("getteachers", appHandler(getCurrentStudentBlocks))
	addStudentHandle("new", appHandler(newStudent))
	addStudentHandle("edit", appHandler(editStudent))
	addStudentHandle("delete", appHandler(deleteStudent))
	addStudentHandle("getall", appHandler(getAllStudents))
	addStudentHandle("open", appHandler(studentClassOpen))
}

func setTeacher(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	decoder := json.NewDecoder(r.Body)
	var variables struct {
		ID string
		Block   int
	}
	if err := decoder.Decode(&variables); err != nil {
		return errors.New(err.Error())
	}

	stdnt, err := student.Current(ctx, false, debug)
	if err != nil {
		return err
	}
	key, err := datastore.DecodeKey(variables.ID)
	if err != nil {
		return errors.New(err.Error())
	}
	newTeacher, err := teacher.WithKey(ctx, key, debug)
	if err != nil {
		return err
	}
	if variables.Block == 0 {
		prevTeacher, err := teacher.WithEmail(ctx, stdnt.Teacher1, false, debug)
		prevOpen := true
		if err != nil && err.(errors.Error).Message != teacher.TeacherNotFound {
			return err
		} else if err == nil {
			prevOpen = prevTeacher.Block1.BlockOpen
		}
		newOpen := newTeacher.Block1.BlockOpen
		newFull := newTeacher.Block1.CurSize >= newTeacher.Block1.MaxSize

		if prevOpen && newOpen && !newFull {
			stdnt.Teacher1 = newTeacher.Email
		}
	} else {
		prevTeacher, err := teacher.WithEmail(ctx, stdnt.Teacher2, false, debug)
		prevOpen := true
		if err != nil && err.(errors.Error).Message != teacher.TeacherNotFound {
			return err
		} else if err == nil {
			prevOpen = prevTeacher.Block2.BlockOpen
		}
		newOpen := newTeacher.Block2.BlockOpen
		newFull := newTeacher.Block2.CurSize >= newTeacher.Block2.MaxSize

		if prevOpen && newOpen && !newFull {
			stdnt.Teacher2 = newTeacher.Email
		}
	}

	err = stdnt.Edit(ctx)
	return err
}

func getCurrentStudentBlocks(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	current:= r.Form.Get("current") == "true"

	stdnt, err := student.Current(ctx, current, debug)
	if err != nil {
		return err
	}
	block1, err := teacher.WithEmail(ctx, stdnt.Teacher1, current, debug)
	if err != nil && err.(errors.Error).Message != teacher.TeacherNotFound {
		return err
	}
	block2, err := teacher.WithEmail(ctx, stdnt.Teacher2, current, debug)
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

func newStudent(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		decoder := json.NewDecoder(r.Body)
		newS := new(student.Student)
		if err := decoder.Decode(newS); err != nil {
			return errors.New(err.Error())
		}
		err := newS.New(ctx, debug)
		studentList = append(studentList, newS)
		return err
	}
	return errors.New(errors.AccessDenied)
}

func editStudent(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		decoder := json.NewDecoder(r.Body)
		stdnt := new(student.Student)
		if err := decoder.Decode(stdnt); err != nil {
			return errors.New(err.Error())
		}
		for i, j := range studentList {
			if j.ID == stdnt.ID {
				studentList[i] = stdnt
			}
		}
		err := stdnt.Edit(ctx)
		return err
	}
	return errors.New(errors.AccessDenied)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		decoder := json.NewDecoder(r.Body)
		sKey := new(string)
		if err := decoder.Decode(sKey); err != nil {
			return errors.New(err.Error())
		}
		stdnt := student.Student{ID: *sKey}
		for i, j := range studentList {
			if j.ID == stdnt.ID {
				studentList[i] = nil
			}
		}
		err = stdnt.Delete(ctx)
		return err
	}
	return errors.New(errors.AccessDenied)
}

var studentList []*student.Student

func getAllStudents(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	if len(studentList) == 0 {
		var err error
		if studentList, err = student.All(ctx, false, debug); err != nil {
			return err
		}
	}

	jStudents, err := json.Marshal(studentList)
	if err != nil {
		return errors.New(err.Error())
	}
	s := string(jStudents[:])

	fmt.Fprintln(w, s)
	return nil
}

func studentClassOpen(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	stdntID := r.Form.Get("id")
	block, _ := strconv.Atoi(r.Form.Get("Block"))
	
	key, err := datastore.DecodeKey(stdntID)
	if err != nil {
		return errors.New(err.Error())
	}
	stdnt, err := student.WithKey(ctx, key)
	if err != nil {
		return err
	}
	
	open := true
	if block == 0 {
		Teacher, err := teacher.WithEmail(ctx, stdnt.Teacher1, false, false)
		
		if err != nil && err.(errors.Error).Message != teacher.TeacherNotFound {
			return err
		} else if err == nil {
			open = Teacher.Block1.BlockOpen
		}
	} else {
		Teacher, err := teacher.WithEmail(ctx, stdnt.Teacher2, false, false)
		if err != nil && err.(errors.Error).Message != teacher.TeacherNotFound {
			return err
		} else if err == nil {
			open = Teacher.Block2.BlockOpen
		}
	}
	
	jOutput, err := json.Marshal(open)
	if err != nil {
		return errors.New(err.Error())
	}
	s := string(jOutput[:])

	fmt.Fprintln(w, s)
	return nil
}