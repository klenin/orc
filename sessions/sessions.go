package sessions

import (
	"github.com/gorilla/securecookie"
	"net/http"
	"time"
)

var lifetime = 300 //5 min

var CookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func SetSession(id, login string, response http.ResponseWriter) {
	value := map[string]interface{}{
		"id":   id,
		"name": login,
		"time": int(time.Now().Unix()),
	}
	if encoded, err := CookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:   "session",
			Value:  encoded,
			Path:   "/",
			MaxAge: int(time.Now().Unix()) + lifetime,
		}
		http.SetCookie(response, cookie)
	}
}

func GetValue(field string, request *http.Request) interface{} {
	var value interface{}
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]interface{})
		if err = CookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			value = cookieValue[field]
		} else {
			println("session.GetValue Error", err)
			return nil
		}
	} else {
		println("session.GetValue Error", err)
		return nil
	}
	return value
}

func setValue(field string, value interface{}, request *http.Request) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]interface{})
		if err = CookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			cookieValue[field] = value
			cookie.MaxAge = int(time.Now().Unix())
		} else {
			println("session.setValue ErrorSetValue", err)
		}
	} else {
		println("session.setValue Error", err)
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
	oldTime, ok := GetValue("time", request).(int)
	if ok != true || oldTime == 0 {
		http.Redirect(response, request, "/", 302)
		return false
	} else {
		newTime := int(time.Now().Unix())
		if oldTime+lifetime < newTime {
			ClearSession(response)
			http.Redirect(response, request, "/", 302)
			return false
		} else {
			setValue("time", newTime, request)
			return true
		}
	}
}
