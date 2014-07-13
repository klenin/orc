package sessions

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/orc/utils"
	"net/http"
	"strconv"
	"time"
)

var CookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func SetSession(id, login string, response http.ResponseWriter) {
	value := map[string]string{
		"id":   id,
		"name": login,
		"time": strconv.Itoa(int(time.Now().Unix()) + 900),
	}
	if encoded, err := CookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:   "session",
			Value:  encoded,
			Path:   "/",
			MaxAge: int(time.Now().Unix()) + 900,
		}
		http.SetCookie(response, cookie)
	}
}

func GetValue(field string, request *http.Request) string {
	value := ""
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = CookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			value = cookieValue[field]
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
	return value
}

func SetValue(field, value string, request *http.Request) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = CookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			cookieValue[field] = value
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}

func ClearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func CheackSession(response http.ResponseWriter, request *http.Request) bool {
	t := GetValue("time", request)
	if t == "" {
		fmt.Println("WHERE")
		http.Redirect(response, request, "/", 302)
		return false
	} else {
		oldTime, err := strconv.Atoi(t)
		utils.HandleErr("CheackSession: ", err, nil)
		newTime := int(time.Now().Unix())
		if oldTime < newTime-900 {
			ClearSession(response)
			http.Redirect(response, request, "/", 302)
			return false
		} else {
			SetValue("time", strconv.Itoa(int(time.Now().Unix())+900), request)
			return true
		}

	}
}
