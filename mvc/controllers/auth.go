package controllers

import (
    "crypto/md5"
    "crypto/rand"
    "encoding/base64"
    "encoding/hex"
    "encoding/json"
    "github.com/orc/db"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "log"
    "regexp"
    "strconv"
    "time"
)

const HASH_SIZE = 32

func MatchRegexp(pattern, str string) bool {
    result, _ := regexp.MatchString(pattern, str)
    return result
}

func GetMD5Hash(text string) string {
    hasher := md5.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}

func GetRandSeq(size int) string {
    rb := make([]byte, size)
    _, err := rand.Read(rb)
    if err != nil {
        log.Println(err)
    }
    return base64.URLEncoding.EncodeToString(rb)
}

func (this *Handler) HandleRegister(login, password, role, fname, lname string) string {
    result := map[string]string{"result": "ok"}
    salt := strconv.Itoa(int(time.Now().Unix()))
    pass := GetMD5Hash(password + salt)

    passHasInvalidChars := false
    for i := 0; i < len(password); i++ {
        if strconv.IsPrint(rune(password[i])) == false {
            passHasInvalidChars = true
            break
        }
    }

    if db.IsExists_("users", []string{"login"}, []interface{}{login}) == true {
        result["result"] = "loginExists"
    } else if !MatchRegexp("^[a-zA-Z0-9]{2,36}$", login) {
        result["result"] = "badLogin"
    } else if !MatchRegexp("^.{6,36}$", password) || !passHasInvalidChars {
        result["result"] = "badPassword"
    } else {
        var p_id int
        person := GetModel("persons")
        person.LoadModelData(map[string]interface{}{"fname": fname, "lname": lname})
        db.QueryInsert_(person, "RETURNING id").Scan(&p_id)

        user := GetModel("users")
        user.LoadModelData(map[string]interface{}{"login": login, "pass": pass, "salt": salt, "role": role, "person_id": p_id})
        db.QueryInsert_(user, "")
    }

    response, err := json.Marshal(result)
    utils.HandleErr("[Handler::HandleRegister] Marshal: ", err, this.Response)

    return string(response)
}

func (this *Handler) HandleLogin(login, pass string) string {
    var id int
    var passHash, salt string
    result := make(map[string]interface{}, 1)

    model := GetModel("users")
    model.LoadWherePart(map[string]interface{}{"login": login})
    err := db.SelectRow(model, []string{"id", "pass", "salt"}, "").Scan(&id, &passHash, &salt)

    if err != nil {
        result["result"] = "invalidCredentials"
    } else if passHash == GetMD5Hash(pass+salt) {
        result["result"] = "ok"

        hash := GetRandSeq(HASH_SIZE)

        user := GetModel("users")
        user.LoadModelData(map[string]interface{}{"id": id, "hash": hash})
        db.QueryUpdate_(user, "")

        sessions.SetSession(this.Response, map[string]interface{}{"id": id, "hash": hash})
    }
    response, err := json.Marshal(result)
    utils.HandleErr("[Handler::HandleLogin] Marshal: ", err, this.Response)

    return string(response)
}

func (this *Handler) HandleLogout() string {
    result := map[string]string{"result": "ok"}
    sessions.ClearSession(this.Response)

    response, err := json.Marshal(result)
    utils.HandleErr("[Handler::HandleLogout] Marshal: ", err, this.Response)

    return string(response)
}
