package sessions

import (
    "github.com/gorilla/securecookie"
    "log"
    "net/http"
    "time"
)

var lifetime = 300 //5 min

var CookieHandler = securecookie.New(
    securecookie.GenerateRandomKey(64),
    securecookie.GenerateRandomKey(32))

func SetSession(response http.ResponseWriter, values map[string]interface{}) {
    values["time"] = int(time.Now().Unix())
    if encoded, err := CookieHandler.Encode("session", values); err == nil {
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
    if cookie, err := request.Cookie("session"); err == nil {
        cookieValue := make(map[string]interface{})
        if err = CookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
            return cookieValue[field]
        } else {
            log.Println("session.GetValue [CookieHandler.Decode]: ", err)
            return nil
        }
    } else {
        log.Println("session.GetValue [request.Cookie]: ", err)
        return nil
    }
}

func setValue(field string, value interface{}, request *http.Request) {
    if cookie, err := request.Cookie("session"); err == nil {
        cookieValue := make(map[string]interface{})
        if err = CookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
            cookieValue[field] = value
            cookie.MaxAge = int(time.Now().Unix())
        } else {
            log.Println("session.setValue ErrorSetValue", err)
        }
    } else {
        log.Println("session.setValue Error", err)
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
    oldTime := GetValue("time", request)
    if oldTime == nil {
        return false
    } else {
        newTime := int(time.Now().Unix())
        if oldTime.(int)+lifetime < newTime {
            ClearSession(response)
            return false
        } else {
            setValue("time", newTime, request)
            return true
        }
    }
}
