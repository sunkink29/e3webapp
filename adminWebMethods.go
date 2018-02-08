package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sunkink29/e3SelectionWebApp/user"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"errors"
)

func addAdminMethods() {
	addToWebMethods("print", returnInput)
	addToWebMethods("addFirstUser", addFirstUser)
	addToWebMethods("newUser", addNewUser)
	addToWebMethods("editUser", editUser)
	addToWebMethods("deleteUser", deleteUser)
	addToWebMethods("getAllUsers", getAllUsers)
}

func returnInput(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	str := r.Form.Get("data")
	fmt.Fprintln(w, str)
	return nil
}

func addFirstUser(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	if users, _ := user.GetAll(ctx, debug); len(users) <= 0 {
		usr := new(user.User)
		if err := dec.Decode(usr); err != nil {
			return err
		}
		err := user.New(ctx, usr, debug)
		return err
	}
	return errors.New("Access Denied")
}

func addNewUser(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.GetCurrent(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		usr := new(user.User)
		if err := dec.Decode(usr); err != nil {
			return err
		}
		err := user.New(ctx, usr, debug)
		return err
	}
	return errors.New("Access Denied")

}

func editUser(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.GetCurrent(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		usr := new(user.User)
		if err := dec.Decode(usr); err != nil {
			return err
		}
		err := user.Edit(ctx, usr)
		return err
	}
	return errors.New("Access Denied")
}

func deleteUser(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.GetCurrent(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		sKey := new(string)
		if err := dec.Decode(sKey); err != nil {
			return err
		}
		k, err := datastore.DecodeKey(*sKey)
		if err != nil {
			return err
		}
		err = user.Delete(ctx, k)
		return err
	}
	return errors.New("Access Denied")
}

func getAllUsers(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.GetCurrent(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		users, err := user.GetAll(ctx, debug)
		if err != nil {
			return err
		}

		jUsers, err := json.Marshal(users)
		if err != nil {
			return err
		}
		s := string(jUsers[:])

		fmt.Fprintln(w, s)
		return nil
	}
	return errors.New("Access Denied")
}
