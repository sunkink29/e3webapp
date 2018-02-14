package main

import (
	"encoding/json"
	"net/http"
	"fmt"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/sunkink29/e3SelectionWebApp/student"
	"github.com/sunkink29/e3SelectionWebApp/teacher"
	"github.com/sunkink29/e3SelectionWebApp/user"
	"github.com/sunkink29/e3SelectionWebApp/errors"
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

func getCurrentStudents(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.GetCurrent(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Teacher {
		var current bool
		if err := dec.Decode(&current); err != nil {
			return errors.New(err.Error())
		}
		block1, err := teacher.GetStudentList(ctx, 0, current, debug)
		if err != nil && err.(errors.Error).Message != student.StudentNotFound && 
				err.(errors.Error).Message != teacher.TeacherNotFound {
			return err
		}
		block2, err := teacher.GetStudentList(ctx, 1, current, debug)
		if err != nil && err.(errors.Error).Message != student.StudentNotFound && 
				err.(errors.Error).Message != teacher.TeacherNotFound {
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
	curU, err := user.GetCurrent(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Teacher {
		cTeacher, err := teacher.GetCurrent(ctx, false, debug)
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
	curU, err := user.GetCurrent(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Teacher {
		var blocks []teacher.Block
		if err := dec.Decode(&blocks); err != nil {
			return errors.New(err.Error())
		}
		cTeacher, err := teacher.GetCurrent(ctx, false, debug)
		if err != nil && err.(errors.Error).Message == teacher.TeacherNotFound {
			newT := new(teacher.Teacher)
			newT.Email = curU.Email
			newT.Name = curU.Name
			newT.Current = false;
			newT.Block1 = blocks[0]
			newT.Block2 = blocks[1]
			err = teacher.New(ctx, newT, debug)
			if err != nil {
				return err
			}
			if err != nil {
				return err
			}
			return nil
		} else if err != nil {
			return err
		} else {
			cTeacher.Block1 = blocks[0]
			cTeacher.Block2 = blocks[1]
			teacher.Edit(ctx, cTeacher)
			return nil
		}

	}
	return errors.New(errors.AccessDenied)
}

func addStudentToClass(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.GetCurrent(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Teacher {
		var variables struct{Key string; Block int}
		if err := dec.Decode(&variables); err != nil {
			return errors.New(err.Error())
		}
		
		k, err := datastore.DecodeKey(variables.Key)
		if err != nil {
			return errors.New(err.Error())
		}
		
		curT, err := teacher.GetCurrent(ctx, false, debug)
		if err != nil {
			return err
		}
		
		stdnt, err := student.Get(ctx, k)
		if err != nil {
			return err
		}
		
		if variables.Block == 0 {
			prevTeacher, err := teacher.GetWithEmail(ctx, stdnt.Teacher1, false, debug)
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
			prevTeacher, err := teacher.GetWithEmail(ctx, stdnt.Teacher2, false, debug)
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
		
		student.Edit(ctx, stdnt)
		return nil
	}
	return errors.New(errors.AccessDenied)
}

func removeFromClass(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.GetCurrent(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Teacher {
		var variables struct{Key string; Block int}
		if err := dec.Decode(&variables); err != nil {
			return errors.New(err.Error())
		}
		
		k, err := datastore.DecodeKey(variables.Key)
		if err != nil {
			return errors.New(err.Error())
		}
		
		stdnt, err := student.Get(ctx, k)
		if err != nil {
			return err
		}
		
		usr, err := user.GetCurrent(ctx, debug)
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
		err = student.Edit(ctx, stdnt)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New(errors.AccessDenied)
}