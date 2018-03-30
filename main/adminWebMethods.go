package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/oauth2"

	"github.com/sunkink29/e3webapp/errors"
	"github.com/sunkink29/e3webapp/student"
	"github.com/sunkink29/e3webapp/teacher"
	"github.com/sunkink29/e3webapp/user"

	gOauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/sheets/v4"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/taskqueue"
	appUser "google.golang.org/appengine/user"
)

func addAdminHandle(url string, handle appHandler) {
	http.Handle("/api/admin/"+url, handle)
}

func addAdminMethods() {
	addAdminHandle("print", appHandler(returnInput))
	addAdminHandle("addfirstuser", appHandler(addFirstUser))
	addAdminHandle("newuser", appHandler(addNewUser))
	addAdminHandle("current", appHandler(getCurrent))
	addAdminHandle("edituser", appHandler(editUser))
	addAdminHandle("deleteuser", appHandler(deleteUser))
	addAdminHandle("getallusers", appHandler(getAllUsers))
	addAdminHandle("getstudentclass", appHandler(getStudentsInClass))
	addAdminHandle("importusers", appHandler(startImport))
	addAdminHandle("setauthinfo", appHandler(setAuthInfo))
}

func returnInput(w http.ResponseWriter, r *http.Request) error {
	str := r.Form.Get("string")
	fmt.Fprintln(w, str)
	return nil
}

func addFirstUser(w http.ResponseWriter, r *http.Request) error {
	if !appengine.IsDevAppServer() {
		return errors.New(errors.AccessDenied)
	}
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	if users, _ := user.All(ctx, debug); len(users) <= 0 {
		usr := new(user.User)
		uByte := []byte(r.Form.Get("user"))
		err := json.Unmarshal(uByte, usr)
		if err != nil {
			return errors.New(err.Error())
		}

		err = usr.New(ctx, debug)
		return err
	}
	return errors.New(errors.AccessDenied)
}

func addNewUser(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		decoder := json.NewDecoder(r.Body)
		usr := new(user.User)
		if err := decoder.Decode(usr); err != nil {
			return errors.New(err.Error())
		}
		err := usr.New(ctx, debug)
		userList = append(userList, usr)
		return err
	}
	return errors.New("Access Denied")

}

func getCurrent(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}

	jUsers, err := json.Marshal(curU)
	if err != nil {
		return errors.New(err.Error())
	}
	s := string(jUsers[:])

	fmt.Fprintln(w, s)
	return nil
}

func editUser(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	usr := appUser.Current(ctx)

	if curU.Admin || usr.Admin {
		decoder := json.NewDecoder(r.Body)
		usr := new(user.User)
		if err := decoder.Decode(usr); err != nil {
			return errors.New(err.Error())
		}
		for i, j := range userList {
			if j.ID == usr.ID {
				userList[i] = usr
			}
		}
		err := usr.Edit(ctx)
		return err
	}
	return errors.New("Access Denied")
}

func deleteUser(w http.ResponseWriter, r *http.Request) error {
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
		usr := new(user.User)
		usr.ID = *sKey
		for i, j := range userList {
			if j.ID == usr.ID {
				userList[i] = nil
			}
		}
		err = usr.Delete(ctx)
		return err
	}
	return errors.New(errors.AccessDenied)
}

var userList []*user.User

func getAllUsers(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {

		if len(userList) == 0 {
			var err error
			if userList, err = user.All(ctx, debug); err != nil {
				return err
			}
		}

		jUsers, err := json.Marshal(userList)
		if err != nil {
			return errors.New(err.Error())
		}
		s := string(jUsers[:])

		fmt.Fprintln(w, s)
		return nil
	}
	return errors.New(errors.AccessDenied)
}

func getStudentsInClass(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		key, err := datastore.DecodeKey(r.Form.Get("id"))
		if err != nil {
			return errors.New(err.Error())
		}

		tchr, err := teacher.WithKey(ctx, key, debug)
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

func startImport(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}

	if curU.Admin {
		decoder := json.NewDecoder(r.Body)
		id := ""
		if err := decoder.Decode(&id); err != nil {
			return errors.New(err.Error())
		}

		_, err := user.Client(ctx)
		if err != nil && err.Error() == "redirect" {
			jUrl, err := json.Marshal(err.(errors.Redirect))
			if err != nil {
				return errors.New(err.Error())
			}
			s := string(jUrl[:])

			fmt.Fprintln(w, s)
			return nil
		} else if err != nil {
			return err
		}

		bToken, err := json.Marshal(curU.Token)
		if err != nil {
			return errors.New(err.Error())
		}
		sToken := string(bToken[:])

		t := taskqueue.NewPOSTTask("/worker/importusers", url.Values{
			"ID":    {id},
			"token": {sToken},
		})
		_, err = taskqueue.Add(ctx, t, "")
		if err != nil {
			return errors.New(err.Error())
		}

		return nil
	}
	return errors.New(errors.AccessDenied)
}

func importUsers(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"

	id := r.Form.Get("ID")
	sToken := r.Form.Get("token")
	var token *oauth2.Token

	if sToken != "null" {
		token = new(oauth2.Token)
		tByte := []byte(sToken)
		err := json.Unmarshal(tByte, token)
		if err != nil {
			return errors.New(err.Error())
		}
	} else {
		return errors.New("no token")
	}

	client := user.Conf.Client(ctx, token)

	oauth2Service, err := gOauth2.New(client)
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.AccessToken(token.AccessToken)
	_, err = tokenInfoCall.Do()
	if err != nil {
		return errors.New("token not valid\n" + err.Error())
	}

	srv, err := sheets.New(client)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to retrieve Sheets Client %v", err))
	}

	readRange := "user data!A:Z"
	resp, err := srv.Spreadsheets.Values.Get(id, readRange).Do()
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to retrieve data from sheet. %v", err))
	}
	if len(resp.Values) > 0 {
		cells := resp.Values
		var email, name, teacher, admin, grade = -1, -1, -1, -1, -1
		var sEmail, sName, sTeacher, sAdmin, sGrade = "email", "name", "teacher", "admin", "grade"
		var fields = map[string]*int{
			sEmail:   &email,
			sName:    &name,
			sTeacher: &teacher,
			sAdmin:   &admin,
			sGrade:   &grade,
		}
		output := ""
		for index, column := range cells[0] {
			for key, value := range fields {
				if strings.ToLower(column.(string)) == key {
					*value = index
					output += fmt.Sprint(key, ":", *value, " ")
				}
			}

		}
		fmt.Fprintln(w, output)
		var missing = ""
		for key, value := range fields {
			if *value == -1 {
				missing += key + ", "
			}
		}
		if missing != "" {
			return errors.New(fmt.Sprintf("Unable to find %v Fields in sheet", missing))
		}

		var numRoutines = 20

		var divided [][][]interface{}

		chunkSize := (len(cells) + numRoutines - 1) / numRoutines

		for i := 0; i < len(cells); i += chunkSize {
			end := i + chunkSize

			if end > len(cells) {
				end = len(cells)
			}

			divided = append(divided, cells[i:end])
		}

		//			var sum int
		//			for _, block := range divided {
		//				length := len(block)
		//				sum += length
		//				fmt.Fprintln(w, length)
		//			}
		//			fmt.Fprintln(w, sum)
		//			return nil

		users, err := user.All(ctx, false)
		if err != nil {
			return err
		}

		stdnts, err := student.All(ctx, false, false)
		if err != nil {
			return err
		}

		var userMap = make(map[string]*user.User)
		for _, usr := range users {
			userMap[usr.Email] = usr
		}

		var stdntMap = make(map[string]*student.Student)
		for _, stdnt := range stdntMap {
			stdntMap[stdnt.Email] = stdnt
		}

		// c := make(chan map[string]*user.User)
		// for i := 0; i < numRoutines; i++ {
		// 	go func(cells [][]interface{}, c chan map[string]*user.User) {

		// 		c <- newUsers
		// 	}(divided[i], c)
		// }

		// newUsers := make(map[string]*user.User)
		// var usrs map[string]*user.User
		// for i := 0; i < numRoutines; i++ {
		// 	select {
		// 	case usrs = <-c:
		// 		for key, value := range usrs {
		// 			newUsers[key] = value
		// 		}
		// 	}
		// }

		newUsers := make(map[string]*user.User)
		for _, row := range cells {
			if row[name] != "" && row[email] != "" && strings.ToLower(row[name].(string)) != strings.ToLower(sName) {
				fmt.Fprintln(w, fmt.Sprint(row))
				usr := new(user.User)
				usr.Email = row[email].(string)
				usr.Name = row[name].(string)
				usr.Teacher = row[teacher].(string) == "TRUE"
				usr.Admin = row[admin].(string) == "TRUE"
				if !usr.Teacher && !usr.Admin && len(row) >= grade {
					stdnt := new(student.Student)
					stdnt.Email = usr.Email
					stdnt.Name = usr.Name
					stdnt.Grade, err = strconv.Atoi(row[grade].(string))
					stdnt.Current = false
					if _, ok := stdntMap[usr.Email]; !ok {
						stdnt.New(ctx, false)
					} else {
						oldStdnt := stdntMap[usr.Email]
						stdnt.ID = oldStdnt.ID
						stdnt.Teacher1 = oldStdnt.Teacher1
						stdnt.Teacher2 = oldStdnt.Teacher2
						if *stdnt != *oldStdnt {
							stdnt.Edit(ctx)
						}
					}
				}
				if _, ok := userMap[usr.Email]; !ok {
					usr.New(ctx, false)
				} else {
					oldUsr := userMap[usr.Email]
					usr.ID = oldUsr.ID
					if *usr != *oldUsr {
						usr.Edit(ctx)
					}
				}
				newUsers[usr.Email] = usr
			}
		}

		for _, usr := range users {
			if _, ok := newUsers[usr.Email]; !ok {
				usr.Delete(ctx)
			}
		}

		for _, stdnt := range stdnts {
			if usr, ok := newUsers[stdnt.Email]; !ok || usr.Teacher {
				stdnt.Delete(ctx)
			}
		}
		if userList, err = user.All(ctx, debug); err != nil {
			return err
		}
		if studentList, err = student.All(ctx, false, debug); err != nil {
			return err
		}
	}
	return nil
}

func setAuthInfo(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)

	if appengine.IsDevAppServer() {
		var creds user.Credentials

		cByte := []byte(r.Form.Get("auth"))
		err := json.Unmarshal(cByte, &creds)
		if err != nil {
			return errors.New(err.Error())
		}

		key := datastore.NewKey(ctx, "Auth", "Auth", 0, nil)
		_, err = datastore.Put(ctx, key, &creds)
		if err != nil {
			return errors.New(err.Error())
		}
		return nil
	}
	return errors.New(errors.AccessDenied)
}
