package controllers

import (
    "github.com/orc/db"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "github.com/orc/mvc/models"
    "github.com/orc/mailer"
    "strconv"
    "time"
)

func (this *Handler) HandleLogin(login, pass string) interface{} {
    var id int
    var enabled bool
    var passHash, salt string
    result := make(map[string]interface{}, 1)

    user := GetModel("users")
    user.LoadWherePart(map[string]interface{}{"login": login})
    err := db.SelectRow(user, []string{"id", "pass", "salt", "enabled"}).Scan(&id, &passHash, &salt, &enabled)

    if err != nil {
        result["result"] = "invalidCredentials"
    } else if enabled == false {
        result["result"] = "notEnabled"
    } else if passHash == utils.GetMD5Hash(pass+salt) {
        result["result"] = "ok"

        hash := utils.GetRandSeq(HASH_SIZE)

        user := GetModel("users")
        user.LoadModelData(map[string]interface{}{"hash": hash})
        user.GetFields().(*models.User).Enabled = true
        user.LoadWherePart(map[string]interface{}{"id": id})
        db.QueryUpdate_(user).Scan()

        sessions.SetSession(this.Response, map[string]interface{}{"id": id, "hash": hash})
    } else {
        result["result"] = "badPassword"
    }

    return result
}

func (this *Handler) HandleLogout() interface{} {
    result := map[string]string{"result": "ok"}
    sessions.ClearSession(this.Response)

    return result
}

func (this *Handler) HandleRegister_(login, password, email, role string) (result string, reg_id int) {
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
    // } else if bad email {
    } else {
        token := utils.GetRandSeq(HASH_SIZE)

        if !mailer.SendConfirmEmail(login, email, token) {
            return "badEmail", -1
        }

        var user_id int
        user := GetModel("users")
        user.LoadModelData(map[string]interface{}{"login": login, "pass": pass, "salt": salt, "role": role, "token": token, "enabled": false})
        user.GetFields().(*models.User).Enabled = false
        db.QueryInsert_(user, "RETURNING id").Scan(&user_id)

        var face_id int
        face := GetModel("faces")
        face.LoadModelData(map[string]interface{}{"user_id": user_id})
        db.QueryInsert_(face, "RETURNING id").Scan(&face_id)

        registration := GetModel("registrations")
        registration.LoadModelData(map[string]interface{}{"face_id": face_id, "event_id": 1})
        db.QueryInsert_(registration, "RETURNING id").Scan(&reg_id)

        return result, reg_id
    }

    return result, -1
}

func (this *Handler) ConfirmUser(token string) {
    user := GetModel("users")
    user.LoadWherePart(map[string]interface{}{"token": token})

    var id string
    err := db.SelectRow(user, []string{"id"}).Scan(&id)
    if utils.HandleErr("[Handle::ConfirmUser]: ", err, this.Response) {
        return
    }

    user = GetModel("users")
    user.LoadModelData(map[string]interface{}{"token": " ", "enabled": true})
    user.GetFields().(*models.User).Enabled = true
    user.LoadWherePart(map[string]interface{}{"id": id})
    db.QueryUpdate_(user).Scan()

    if this.Response != nil {
        this.Render([]string{"mvc/views/msg.html"}, "msg", "Регестрация подтверждена.")
    }
}

func (this *Handler) RejectUser(token string) {
    user := GetModel("users")
    user.LoadWherePart(map[string]interface{}{"token": token})

    var id string
    err := db.SelectRow(user, []string{"id"}).Scan(&id)
    if err != nil {
        utils.HandleErr("[Handle::RejectUser]: ", err, this.Response)
        return
    }

    db.QueryDeleteByIds("users", id)

    if this.Response != nil {
        this.Render([]string{"mvc/views/msg.html"}, "msg", "Вы успешно отписаны от рассылок Secret Oasis.")
    }
}
