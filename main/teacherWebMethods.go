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

func addTeacherHandle(url string, handle appHandler) {
	http.Handle("/api/teacher/"+url, handle)
}

func addTeacherMethods() {
	addTeacherHandle("new", appHandler(newTeacher))
	addTeacherHandle("edit", appHandler(editTeacher))
	addTeacherHandle("delete", appHandler(deleteTeacher))
	addTeacherHandle("getall", appHandler(getAllTeachers))
	addTeacherHandle("getstudents", appHandler(getCurrentStudents))
	addTeacherHandle("getblocks", appHandler(getBlocks))
	addTeacherHandle("setblocks", appHandler(setBlocks))
	addTeacherHandle("addstudent", appHandler(addStudentToClass))
	addTeacherHandle("removestudent", appHandler(removeFromClass))
}

func newTeacher(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		decoder := json.NewDecoder(r.Body)
		newTchr := new(teacher.Teacher)
		if err := decoder.Decode(newTchr); err != nil {
			return errors.New(err.Error())
		}
		err := newTchr.New(ctx, debug)
		return err
	}
	return errors.New(errors.AccessDenied)
}

func editTeacher(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		decoder := json.NewDecoder(r.Body)
		tchr := new(teacher.Teacher)
		if err := decoder.Decode(tchr); err != nil {
			return errors.New(err.Error())
		}
		err := tchr.Edit(ctx)
		return err
	}
	return errors.New(errors.AccessDenied)
}

func deleteTeacher(w http.ResponseWriter, r *http.Request) error {
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
		tchr := new(teacher.Teacher)
		tchr.ID = *sKey
		err = tchr.Delete(ctx)
		return err
	}
	return errors.New(errors.AccessDenied)
}

func getAllTeachers(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	current, _ := strconv.ParseBool(r.Form.Get("current"))
	
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

func getCurrentStudents(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Teacher {
		current, _ := strconv.ParseBool(r.Form.Get("current"))
		
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

func getBlocks(w http.ResponseWriter, r *http.Request) error {
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

func setBlocks(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Teacher {
		decoder := json.NewDecoder(r.Body)
		var blocks []teacher.Block
		if err := decoder.Decode(&blocks); err != nil {
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

func addStudentToClass(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Teacher {
		decoder := json.NewDecoder(r.Body)
		var variables struct {
			Key   string
			Block int
		}
		if err := decoder.Decode(&variables); err != nil {
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

func removeFromClass(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Teacher {
		decoder := json.NewDecoder(r.Body)
		var variables struct {
			Key   string
			Block int
		}
		if err := decoder.Decode(&variables); err != nil {
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
