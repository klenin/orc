package controllers

import (
    "encoding/json"
    "github.com/orc/db"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "strconv"
    "time"
)

const HASH_SIZE = 32

func (this *Handler) HandleLogin(login, pass string) string {
    var id int
    var passHash, salt string
    result := make(map[string]interface{}, 1)

    model := GetModel("users")
    model.LoadWherePart(map[string]interface{}{"login": login})
    err := db.SelectRow(model, []string{"id", "pass", "salt"}, "").Scan(&id, &passHash, &salt)

    if err != nil {
        result["result"] = "invalidCredentials"
    } else if passHash == utils.GetMD5Hash(pass+salt) {
        result["result"] = "ok"

        hash := utils.GetRandSeq(HASH_SIZE)

        user := GetModel("users")
        user.LoadModelData(map[string]interface{}{"id": id, "hash": hash})
        db.QueryUpdate_(user, "")

        sessions.SetSession(this.Response, map[string]interface{}{"id": id, "hash": hash})
    } else {
        result["result"] = "badPassword"
    }
    response, err := json.Marshal(result)
    if utils.HandleErr("[Handler::HandleLogin] Marshal: ", err, this.Response) {
        return ""
    }

    return string(response)
}

func (this *Handler) HandleLogout() string {
    result := map[string]string{"result": "ok"}
    sessions.ClearSession(this.Response)

    response, err := json.Marshal(result)
    if utils.HandleErr("[Handler::HandleLogout] Marshal: ", err, this.Response) {
        return ""
    }

    return string(response)
}

func (this *Handler) HandleRegister_(login, password, role string) (result string, reg_id int) {
    result = "ok"
    salt := strconv.Itoa(int(time.Now().Unix()))
    pass := utils.GetMD5Hash(password + salt)

    passHasInvalidChars := false
    for i := 0; i < len(password); i++ {
        if strconv.IsPrint(rune(password[i])) == false {
            passHasInvalidChars = true
            break
        }
    }

    if db.IsExists_("users", []string{"login"}, []interface{}{login}) == true {
        result = "loginExists"
    } else if !utils.MatchRegexp("^[a-zA-Z0-9]{2,36}$", login) {
        result = "badLogin"
    } else if !utils.MatchRegexp("^.{6,36}$", password) || passHasInvalidChars {
        result = "badPassword"
    } else {
        var user_id int
        user := GetModel("users")
        user.LoadModelData(map[string]interface{}{"login": login, "pass": pass, "salt": salt, "role": role})
        db.QueryInsert_(user, "RETURNING id").Scan(&user_id)

        var face_id int
        face := GetModel("faces")
        face.LoadModelData(map[string]interface{}{"user_id": user_id})
        db.QueryInsert_(face, "RETURNING id").Scan(&face_id)

        registration := GetModel("registrations")
        registration.LoadModelData(map[string]interface{}{"face_id": face_id})
        db.QueryInsert_(registration, "RETURNING id").Scan(&reg_id)

        return result, reg_id
    }

    return result, -1
}
