package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"strconv"
	"net/http"

	"github.com/sunkink29/e3SelectionWebApp/errors"
	"github.com/sunkink29/e3SelectionWebApp/student"
	"github.com/sunkink29/e3SelectionWebApp/teacher"
	"github.com/sunkink29/e3SelectionWebApp/user"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/api/sheets/v4"
	appUser "google.golang.org/appengine/user"
)

func addAdminMethods() {
	addToWebMethods("print", returnInput)
	addToWebMethods("addFirstUser", addFirstUser)
	addToWebMethods("newUser", addNewUser)
	addToWebMethods("editUser", editUser)
	addToWebMethods("deleteUser", deleteUser)
	addToWebMethods("getAllUsers", getAllUsers)
	addToWebMethods("getStudentsInClass", getStudentsInClass)
	addToWebMethods("importUsers", importUsers)
	addToWebMethods("getClientID", getClientID)
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

func importUsers(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		id := ""
		if err := dec.Decode(&id); err != nil {
			return errors.New(err.Error())
		}
		
		client, err := user.Client(ctx)
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
				    sEmail : &email,
				    sName: &name,
				    sTeacher: &teacher,
				    sAdmin: &admin,
				    sGrade: &grade,
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
			
			c := make(chan map[string]*user.User)
			for i := 0; i < numRoutines; i++ {
				go func (cells [][]interface{}, c chan map[string]*user.User) {
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
					c <- newUsers
				}(divided[i], c)
			}
			
			newUsers := make(map[string]*user.User)
			var usrs map[string]*user.User
			for i := 0; i < numRoutines; i++ {
				select {
				case usrs = <- c:
					for key, value := range usrs {
					    newUsers[key] = value
					}
				}
			}
			
		  	for _, usr := range users {
		  		if _, ok := newUsers[usr.Email]; !ok {
		  			usr.Delete(ctx)
		  		}
		  	}
		  	
		  	for _, stdnt := range stdnts {
		  		if usr, ok := newUsers[stdnt.Email]; !ok ||  usr.Teacher {
		  			stdnt.Delete(ctx)
		  		}
		  	}
		} 
		return nil
	}
	return errors.New(errors.AccessDenied)
}

func getClientID(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	debug := r.Form.Get("debug") == "true"
	curU, err := user.Current(ctx, debug)
	if err != nil {
		return err
	}
	if curU.Admin {
		cid := user.ClientID()
		
		jCid, err := json.Marshal(cid)
		if err != nil {
			return errors.New(err.Error())
		}
		s := string(jCid[:])
	
		fmt.Fprintln(w, s)
		return nil
	}
	return errors.New(errors.AccessDenied)
}

func setAuthInfo(dec *json.Decoder, w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	usr := appUser.Current(ctx)
	
	if usr.Admin {
		var creds user.Credentials
		if err := dec.Decode(creds); err != nil {
			return errors.New(err.Error())
		}
		
		key := datastore.NewKey(ctx, "Auth", "Auth", 0, nil)
		_, err := datastore.Put(ctx, key, creds)
		if err != nil {
			return errors.New(err.Error())
		}
		return nil
	}
	return errors.New(errors.AccessDenied)
}