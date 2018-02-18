package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sunkink29/e3SelectionWebApp/errors"
	"github.com/sunkink29/e3SelectionWebApp/student"
	"github.com/sunkink29/e3SelectionWebApp/teacher"
	"github.com/sunkink29/e3SelectionWebApp/user"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func addAdminMethods() {
	addToWebMethods("print", returnInput)
	addToWebMethods("addFirstUser", addFirstUser)
	addToWebMethods("newUser", addNewUser)
	addToWebMethods("editUser", editUser)
	addToWebMethods("deleteUser", deleteUser)
	addToWebMethods("getAllUsers", getAllUsers)
	addToWebMethods("getStudentsInClass", getStudentsInClass)
	//	addToWebMethods("switchNextToCurrent", switchNextToCurrent)
}

func returnInput(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	str := r.Form.Get("data")
	fmt.Fprintln(w, str)
	return nil
}

func addFirstUser(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	if users, _ := user.All(ctx, debug); len(users) <= 0 {
		usr := new(user.User)
		if err := dec.Decode(usr); err != nil {
			return errors.New(err.Error())
		}
		err := usr.New(ctx, debug)
		return err
	}
	return errors.New(errors.AccessDenied)
}

func addNewUser(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		usr := new(user.User)
		if err := dec.Decode(usr); err != nil {
			return errors.New(err.Error())
		}
		err := usr.New(ctx, debug)
		return err
	}
	return errors.New("Access Denied")

}

func editUser(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		usr := new(user.User)
		if err := dec.Decode(usr); err != nil {
			return errors.New(err.Error())
		}
		err := usr.Edit(ctx)
		return err
	}
	return errors.New("Access Denied")
}

func deleteUser(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
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
		usr := new(user.User)
		usr.ID = *sKey
		err = usr.Delete(ctx)
		return err
	}
	return errors.New(errors.AccessDenied)
}

func getAllUsers(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		users, err := user.All(ctx, debug)
		if err != nil {
			return err
		}

		jUsers, err := json.Marshal(users)
		if err != nil {
			return errors.New(err.Error())
		}
		s := string(jUsers[:])

		fmt.Fprintln(w, s)
		return nil
	}
	return errors.New(errors.AccessDenied)
}

func getStudentsInClass(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		key := new(string)
		if err := dec.Decode(key); err != nil {
			return errors.New(err.Error())
		}
		k, err := datastore.DecodeKey(*key)
		if err != nil {
			return errors.New(err.Error())
		}
		tchr, err := teacher.WithKey(ctx, k, debug)
		if err != nil {
			return err
		}

		block1, err := tchr.StudentList(ctx, 0, debug)
		if err != nil {
			return err
		}

		block2, err := tchr.StudentList(ctx, 1, debug)
		if err != nil {
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