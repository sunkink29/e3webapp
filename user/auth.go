package user

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "net/http"

	"golang.org/x/net/context"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    gOauth2 "google.golang.org/api/oauth2/v2"
    "google.golang.org/appengine"
    "google.golang.org/appengine/datastore"
    
    "github.com/sunkink29/e3webapp/errors"
)

// Credentials which stores google ids.
type Credentials struct {
	APIKey 	string
    Cid     string `datastore:"clientID"`
    Csecret string `datastore:"clientSecret"`
    ID		string
    URL	string `datastore:"RedirectURL"`
}

var cred Credentials
var conf *oauth2.Config

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func ClientID() (string) {
	return cred.Cid
}

func InitAuth(ctx context.Context) error {
	if conf == nil {
	    key := datastore.NewKey(ctx, "Auth", "Auth", 0, nil)
	    err := datastore.Get(ctx, key, &cred)
		if err != nil {
			return errors.New(err.Error())
		}
		
	    conf = &oauth2.Config{
	        ClientID:     cred.Cid,
	        ClientSecret: cred.Csecret,
	        RedirectURL:  cred.URL,
	        Scopes: []string{
	            "https://www.googleapis.com/auth/spreadsheets.readonly",
	        },
	        Endpoint: google.Endpoint,
	    }
    }
    return nil
}

func Client(ctx context.Context) (*http.Client, error) {
	usr, err := Current(ctx, false)
	if err != nil {
		return nil, err
	}
	if usr.Token == nil || !usr.Token.Valid() {
		err = requestToken(ctx)
		return nil, err
	}
	
	client := conf.Client(ctx, usr.Token)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	
	oauth2Service, err := gOauth2.New(client)
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.AccessToken(usr.Token.AccessToken)
	_, err = tokenInfoCall.Do()
    if err != nil {
        err = requestToken(ctx)
		return nil, err
    }
	return client, nil
}

func requestToken(ctx context.Context) error {
	usr, err := Current(ctx, false)
	if err != nil {
		return err
	}
	state := randToken()
	usr.AuthState = state
	err = usr.Edit(ctx)
	if err != nil {
		return err
	}
//	setEmail = oauth2.SetAuthURLParam(key, value string)
	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return errors.Redirect{URL: url, Code: 308}
}

func AuthHandle(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	usr, err := Current(ctx, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	retrievedState := r.Form.Get("state")
	if retrievedState != usr.AuthState {
		http.Error(w, fmt.Sprintf("Invalid session state: %s", retrievedState), http.StatusUnauthorized)
		return
	}
	
	tok, err := conf.Exchange(ctx, r.Form.Get("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	usr.Token = tok
	err = usr.Edit(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	
//	bToken, err := json.Marshal(tok)
//	sToken := string(bToken[:])
	
	fmt.Fprintln(w, "<script type='text/javascript'> window.opener.finishImport(); window.close() </script>")
}