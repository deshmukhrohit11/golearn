package main


import (
"fmt"
    "code.google.com/p/goauth2/oauth"
    "net/http"
    "html/template"
     "encoding/json"
      // "io/ioutil"
)

var notAuthenticatedTemplate = template.Must(template.New("").Parse(`
<html><body background="http://goo.gl/eYQDMZ">
<font face="verdana" color="green" size=6> Please authenticate this app with the Google.
<form action="/authorize" method="POST"><input type="image" src="http://goo.gl/beHwhM" alt="Google" 
border="2" width="150" height="48"></form>
<form action="/authorizeFacebook" method="POST"><input type="image" src="http://goo.gl/fm3wgs" alt="Faecbook" border="2" 
width="150" height="48"></form>
</font></body></html>
`));

var userInfoTemplate = template.Must(template.New("").Parse(`
<html><html><body background="http://goo.gl/eYQDMZ">
<font face="verdana" color="green" size=6>This app is now authenticated to access your Google user info.  Your details are:<br />
{{.}}</font>
</body></html>
`));

// variables used during oauth protocol flow of authentication
var (
    code = ""
    token = ""
)

var oauthCfg = &oauth.Config {
         ClientId: TO-DO,
     ClientSecret: TO-DO,
     AuthURL: "https://accounts.google.com/o/oauth2/auth",
         TokenURL: "https://accounts.google.com/o/oauth2/token",
     RedirectURL: "http://localhost:8080/oauth2callback",
     Scope: "https://www.googleapis.com/auth/userinfo.profile",
    }

var oauthFacebookCfg = &oauth.Config {
         ClientId: "TO-DO, 
         ClientSecret:TO-DO,
    AuthURL: "https://www.facebook.com/dialog/oauth",
    TokenURL: "https://graph.facebook.com/oauth/access_token",
    RedirectURL: "http://localhost:8080/oauth2callbackFacebook",
    Scope: "",
    }
//This is the URL that Google has defined so that an authenticated application may obtain the user's info in json format
const profileInfoURL = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"

func main() {
    http.HandleFunc("/", handleRoot)
    http.HandleFunc("/authorize", handleAuthorize)
    http.HandleFunc("/authorizeFacebook", handleAuthorizeFacebook)

    //Google will redirect to this page to return your code, so handle it appropriately
    http.HandleFunc("/oauth2callback", handleOAuth2Callback)
    http.HandleFunc("/oauth2callbackFacebook", handleOAuth2CallbackFacebook)

    http.ListenAndServe("localhost:8080", nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
    notAuthenticatedTemplate.Execute(w, nil)
}

func handleAuthorizeFacebook(w http.ResponseWriter, r *http.Request) {
    //Get the Google URL which shows the Authentication page to the user
    url := oauthFacebookCfg.AuthCodeURL("")

    //redirect user to that page
    http.Redirect(w, r, url, http.StatusFound)
}
func handleAuthorize(w http.ResponseWriter, r *http.Request) {
    //Get the Google URL which shows the Authentication page to the user
    url := oauthCfg.AuthCodeURL("")

    //redirect user to that page
    http.Redirect(w, r, url, http.StatusFound)
}

type Message struct {
    id uint64 `json:",string"`
    name string `json:"name"`
    given_name string `json:"given_name"`
    family_name string `json:"family_name"`
    link string `json:"link"`
    picture string `json:"picture"`
    gender string `json:"gender"`
    locale string `json:"locale"`
}

func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
    //Get the code from the response
    code := r.FormValue("code")

    t := &oauth.Transport{Config: oauthCfg}

    // Exchange the received code for a token
    t.Exchange(code)

    //now get user data based on the Transport which has the token
    resp, _ := t.Client().Get(profileInfoURL)

    buf := make([]byte, 1024)
    var  value Message
    resp.Body.Read(buf)
  
    json.Unmarshal(buf, &value)
    fmt.Println( value)
    userInfoTemplate.Execute(w, string(buf))
}

func handleOAuth2CallbackFacebook(w http.ResponseWriter, r *http.Request) {
    //Get the code from the response
    code := r.FormValue("code")

    t := &oauth.Transport{Config: oauthCfg}

    // Exchange the received code for a token
    t.Exchange(code)

    //now get user data based on the Transport which has the token
    // resp, _ := t.Client().Get(profileInfoURL)
    // c := t.Client()
    // buf := make([]byte, 1024)
    // var  value Message
    // resp.Body.Read(buf)
  
    // json.Unmarshal(buf, &value)
    // fmt.Println( value)
    userInfoTemplate.Execute(w, "Login Done!!")
}
