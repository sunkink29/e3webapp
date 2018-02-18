package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/sunkink29/e3SelectionWebApp/errors"
	"github.com/sunkink29/e3SelectionWebApp/student"
	"github.com/sunkink29/e3SelectionWebApp/teacher"
	"github.com/sunkink29/e3SelectionWebApp/user"
)

func addTeacherMethods() {
	addToWebMethods("newTeacher", newTeacher)
	addToWebMethods("editTeacher", editTeacher)
	addToWebMethods("deleteTeacher", deleteTeacher)
	addToWebMethods("getAllTeachers", getAllTeachers)
	addToWebMethods("getCurrentStudents", getCurrentStudents)
	addToWebMethods("getBlocks", getBlocks)
	addToWebMethods("setBlocks", setBlock)
	addToWebMethods("addStudentToClass", addStudentToClass)
	addToWebMethods("removeFromClass", removeFromClass)
}

func newTeacher(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		newTchr := new(teacher.Teacher)
		if err := dec.Decode(newTchr); err != nil {
			return errors.New(err.Error())
		}
		err := newTchr.New(ctx, debug)
		return err
	}
	return errors.New(errors.AccessDenied)
}

func editTeacher(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		tchr := new(teacher.Teacher)
		if err := dec.Decode(tchr); err != nil {
			return errors.New(err.Error())
		}
		err := tchr.Edit(ctx)
		return err
	}
	return errors.New(errors.AccessDenied)
}

func deleteTeacher(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		sKey := new(string)
		if err := dec.Decode(sKey); err != nil {
			return errors.New(err.Error())
		}
		tchr := new(teacher.Teacher)
		tchr.ID = *sKey
		err = tchr.Delete(ctx)
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
	teachers, err := teacher.All(ctx, current, debug)
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

func getCurrentStudents(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Teacher {
		var current bool
		if err := dec.Decode(&current); err != nil {
			return errors.New(err.Error())
		}
		tchr := teacher.Teacher{Email: curU.Email, Current: current}
		
		var block1 []*student.Student
		var block2 []*student.Student
		block1, err = tchr.StudentList(ctx, 0, debug)
		if err != nil && err.(errors.Error).Message != student.StudentNotFound {
			return err
		}
		
		block2, err = tchr.StudentList(ctx, 1, debug)
		if err != nil && err.(errors.Error).Message != student.StudentNotFound {
			return err
		}
		blocks := [][]*student.Student{block1, block2}
		jBlock, err := json.Marshal(blocks)
		if err != nil {
			return errors.New(err.Error())
		}
		s := string(jBlock[:])

		fmt.Fprintln(w, s)
		return nil
	}
	return errors.New(errors.AccessDenied)
}

func getBlocks(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Teacher {
		cTeacher, err := teacher.Current(ctx, false, debug)
		if err != nil && err.(errors.Error).Message != teacher.TeacherNotFound {
			return err
		} else if err != nil && err.(errors.Error).Message == teacher.TeacherNotFound {
			cTeacher = new(teacher.Teacher)
		}
		blocks := []teacher.Block{cTeacher.Block1, cTeacher.Block2}
		jBlock, err := json.Marshal(blocks)
		if err != nil {
			return errors.New(err.Error())
		}
		s := string(jBlock[:])

		fmt.Fprintln(w, s)
		return nil
	}
	return errors.New(errors.AccessDenied)
}

func setBlock(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Teacher {
		var blocks []teacher.Block
		if err := dec.Decode(&blocks); err != nil {
			return errors.New(err.Error())
		}
		cTchr, err := teacher.Current(ctx, false, debug)
		if err != nil && err.(errors.Error).Message == teacher.TeacherNotFound {
			newTchr := new(teacher.Teacher)
			newTchr.Email = curU.Email
			newTchr.Name = curU.Name
			newTchr.Current = false
			newTchr.Block1 = blocks[0]
			newTchr.Block2 = blocks[1]
			err = newTchr.New(ctx, debug)
			if err != nil {
				return err
			}
			return nil
		} else if err != nil {
			return err
		} else {
			cTchr.Block1 = blocks[0]
			cTchr.Block2 = blocks[1]
			cTchr.Edit(ctx)
			return nil
		}

	}
	return errors.New(errors.AccessDenied)
}

func addStudentToClass(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Teacher {
		var variables struct {
			Key   string
			Block int
		}
		if err := dec.Decode(&variables); err != nil {
			return errors.New(err.Error())
		}

		k, err := datastore.DecodeKey(variables.Key)
		if err != nil {
			return errors.New(err.Error())
		}

		curT, err := teacher.Current(ctx, false, debug)
		if err != nil {
			return err
		}

		stdnt, err := student.WithKey(ctx, k)
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
			newFull := curT.Block1.CurSize >= curT.Block1.MaxSize

			if prevOpen && !newFull {
				stdnt.Teacher1 = curT.Email
			} else {
				if !prevOpen {
					return errors.New("Current student class closed")
				} else {
					return errors.New("Current class full")
				}
			}
		} else {
			prevTeacher, err := teacher.WithEmail(ctx, stdnt.Teacher2, false, debug)
			prevOpen := true
			if err != nil && err.(errors.Error).Message != teacher.TeacherNotFound {
				return err
			} else if err == nil {
				prevOpen = prevTeacher.Block2.BlockOpen
			}
			newOpen := curT.Block2.BlockOpen
			newFull := curT.Block2.CurSize >= curT.Block2.MaxSize

			if prevOpen && newOpen && !newFull {
				stdnt.Teacher2 = curT.Email
			} else {
				return errors.New(errors.AccessDenied)
			}
		}

		stdnt.Edit(ctx)
		return nil
	}
	return errors.New(errors.AccessDenied)
}

func removeFromClass(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Teacher {
		var variables struct {
			Key   string
			Block int
		}
		if err := dec.Decode(&variables); err != nil {
			return errors.New(err.Error())
		}

		k, err := datastore.DecodeKey(variables.Key)
		if err != nil {
			return errors.New(err.Error())
		}

		stdnt, err := student.WithKey(ctx, k)
		if err != nil {
			return err
		}

		usr, err := user.Current(ctx, debug)
		if err != nil {
			return err
		}

		if variables.Block == 0 {
			if stdnt.Teacher1 == usr.Email {
				stdnt.Teacher1 = ""
			} else {
				return errors.New(errors.AccessDenied)
			}
		} else {
			if stdnt.Teacher2 == usr.Email {
				stdnt.Teacher2 = ""
			} else {
				return errors.New(errors.AccessDenied)
			}
		}
		err = stdnt.Edit(ctx)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New(errors.AccessDenied)
}
