package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/sunkink29/e3webapp/errors"
	"github.com/sunkink29/e3webapp/messaging"
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

			teachers := []teacher.Teacher{*cTchr}
			jTeachers, err := json.Marshal(teachers)
			if err != nil {
				return errors.New(err.Error())
			}
			message := string(jTeachers[:])

			err = messaging.SendEvent(ctx, messaging.EventTypes.ClassEdit, message, messaging.Topics.Student)
			if err != nil {
				return err
			}
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

		var stdntTchr *string
		if variables.Block == 0 {
			stdntTchr = &stdnt.Teacher1
		} else {
			stdntTchr = &stdnt.Teacher2
		}
		prevTchr, err := teacher.WithEmail(ctx, *stdntTchr, false, debug)
		prevOpen := true
		if err != nil && err.(errors.Error).Message != teacher.TeacherNotFound {
			return err
		} else if err == nil {
			if variables.Block == 0 {
				prevOpen = prevTchr.Block1.BlockOpen
			} else {
				prevOpen = prevTchr.Block2.BlockOpen
			}
		}

		var newFull bool
		if variables.Block == 0 {
			newFull = curT.Block1.CurSize >= curT.Block1.MaxSize
		} else {
			newFull = curT.Block2.CurSize >= curT.Block2.MaxSize
		}

		if prevOpen && !newFull {
			*stdntTchr = curT.Email
		} else {
			if !prevOpen {
				return errors.New("Student Current class closed")
			}
			return errors.New("Current class full")
		}

		stdnt.Edit(ctx)

		var block *teacher.Block
		var prevBlock *teacher.Block
		if variables.Block == 0 {
			block = &curT.Block1
			if prevTchr != nil {
				prevBlock = &prevTchr.Block1
			}
		} else {
			block = &curT.Block2
			if prevTchr != nil {
				prevBlock = &prevTchr.Block2
			}
		}

		block.CurSize++
		teachers := []teacher.Teacher{*curT}

		if prevBlock != nil {
			prevBlock.CurSize--
			teachers = append(teachers, *prevTchr)
		}

		jTeachers, err := json.Marshal(teachers)
		if err != nil {
			return errors.New(err.Error())
		}
		message := string(jTeachers[:])
		err = messaging.SendEvent(ctx, messaging.EventTypes.ClassEdit, message, messaging.Topics.Student)
		if err != nil {
			return err
		}

		stdntUsr, err := user.WithEmail(ctx, stdnt.Email, false)
		if err != nil {
			return err
		}
		changeTeacher := struct {
			Block   int
			Teacher teacher.Teacher
		}{variables.Block, *curT}
		jTeacher, err := json.Marshal(changeTeacher)
		if err != nil {
			return errors.New(err.Error())
		}
		message2 := string(jTeacher[:])
		err = messaging.SendUserEvent(ctx, messaging.EventTypes.CurrentChange, message2, stdntUsr.RToken)
		if err != nil {
			return err
		}

		if prevTchr != nil {
			tchrUsr, err := user.WithEmail(ctx, prevTchr.Email, false)
			if err != nil {
				return err
			}
			changeStudent := struct {
				Block   int
				Student *student.Student
				Method  string
			}{variables.Block, stdnt, "remove"}
			jStudent, err := json.Marshal(changeStudent)
			if err != nil {
				return errors.New(err.Error())
			}
			message3 := string(jStudent[:])
			err = messaging.SendUserEvent(ctx, messaging.EventTypes.StudentUpdate, message3, tchrUsr.RToken)
			if err != nil {
				return err
			}
		}

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

		var tchr *string
		if variables.Block == 0 {
			tchr = &stdnt.Teacher1
		} else {
			tchr = &stdnt.Teacher2
		}

		if *tchr == usr.Email {
			*tchr = ""
		} else {
			return errors.New(errors.AccessDenied)
		}

		curT, err := teacher.Current(ctx, false, debug)
		if err != nil {
			return err
		}

		var block *teacher.Block
		if variables.Block == 0 {
			block = &curT.Block1
		} else {
			block = &curT.Block2
		}
		block.CurSize--
		teacherList := []teacher.Teacher{*curT}
		jTeachers, err := json.Marshal(teacherList)
		if err != nil {
			return errors.New(err.Error())
		}
		message := string(jTeachers[:])
		err = messaging.SendEvent(ctx, messaging.EventTypes.ClassEdit, message, messaging.Topics.Student)
		if err != nil {
			return err
		}

		stdntUsr, err := user.WithEmail(ctx, stdnt.Email, false)
		if err != nil {
			return err
		}

		changeTeacher := struct {
			Block   int
			Teacher *teacher.Teacher
		}{variables.Block, nil}
		jTeacher, err := json.Marshal(changeTeacher)
		if err != nil {
			return errors.New(err.Error())
		}
		message2 := string(jTeacher[:])
		err = messaging.SendUserEvent(ctx, messaging.EventTypes.CurrentChange, message2, stdntUsr.RToken)

		err = stdnt.Edit(ctx)
		if err != nil {
			return err
		}

		return nil
	}
	return errors.New(errors.AccessDenied)
}
