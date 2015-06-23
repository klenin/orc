package controllers

import (
    "github.com/orc/db"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "github.com/orc/mvc/models"
    // "github.com/orc/mailer"
    "strconv"
    "time"
    "net/http"
)

func (this *Handler) HandleLogin(login, pass string) interface{} {
    var id int
    var enabled bool
    var passHash, salt string
    result := make(map[string]interface{}, 1)

    user := this.GetModel("users")
    user.LoadWherePart(map[string]interface{}{"login": login})
    err := db.SelectRow(user, []string{"id", "pass", "salt", "enabled"}).Scan(&id, &passHash, &salt, &enabled)

    if err != nil {
        result["result"] = "invalidCredentials"

    } else if enabled == false {
        result["result"] = "notEnabled"

    } else if passHash != utils.GetMD5Hash(pass+salt) {
        result["result"] = "badPassword"

    } else {
        result["result"] = "ok"

        sid := utils.GetRandSeq(HASH_SIZE)

        user := this.GetModel("users")
        user.GetFields().(*models.User).Enabled = true
        user.GetFields().(*models.User).Sid = sid
        user.LoadWherePart(map[string]interface{}{"id": id})
        db.QueryUpdate_(user).Scan()

        sessions.SetSession(this.Response, map[string]interface{}{"sid": sid})
    }

    return result
}

func (this *Handler) HandleLogout() interface{} {
    userId, err := this.CheckSid()
    if err != nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return map[string]string{"result": "badSid"}
    }

    var enabled bool
    user := this.GetModel("users")
    user.LoadWherePart(map[string]interface{}{"id": userId})
    err = db.SelectRow(user, []string{"enabled"}).Scan(&enabled)
    if err != nil {
        utils.HandleErr("[Handle::HandleLogout]: ", err, this.Response)
        return map[string]string{"result": err.Error()}
    }

    user = this.GetModel("users")
    user.GetFields().(*models.User).Enabled = enabled
    user.GetFields().(*models.User).Sid = " "
    user.LoadWherePart(map[string]interface{}{"id": userId})
    db.QueryUpdate_(user).Scan()

    sessions.ClearSession(this.Response)

    return map[string]string{"result": "ok"}
}

func (this *Handler) HandleRegister(login, password, email, role string) (result string, regId int) {
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

        // if !mailer.SendConfirmEmail(login, email, token) {
        //     return "badEmail", -1
        // }

        var userId int
        user := this.GetModel("users")
        user.LoadModelData(map[string]interface{}{
            "login": login,
            "pass":  pass,
            "salt":  salt,
            "role":  role,
            "token": token})
        user.GetFields().(*models.User).Enabled = false
        db.QueryInsert(user, "RETURNING id").Scan(&userId)

        var faceId int
        face := this.GetModel("faces")
        face.LoadModelData(map[string]interface{}{"user_id": userId})
        db.QueryInsert(face, "RETURNING id").Scan(&faceId)

        registration := this.GetModel("registrations")
        registration.LoadModelData(map[string]interface{}{"face_id": faceId, "event_id": 1})
        db.QueryInsert(registration, "RETURNING id").Scan(&regId)

        return result, regId
    }

    return result, -1
}

func (this *Handler) ConfirmUser(token string) {
    var userId int
    user := this.GetModel("users")
    user.GetFields().(*models.User).Token = token
    err := db.SelectRow(user, []string{"id"}).Scan(&userId)

    if utils.HandleErr("[Handle::ConfirmUser]: ", err, this.Response) {
        if this.Response != nil {
            this.Render([]string{"mvc/views/msg.html"}, "msg", err.Error())
        }
        return
    }

    user = this.GetModel("users")
    user.GetFields().(*models.User).Enabled = true
    user.GetFields().(*models.User).Token = " "
    user.LoadWherePart(map[string]interface{}{"id": userId})
    db.QueryUpdate_(user).Scan()

    if this.Response != nil {
        this.Render([]string{"mvc/views/msg.html"}, "msg", "Регистрация подтверждена.")
    }
}

func (this *Handler) RejectUser(token string) {
    var userId int
    user := this.GetModel("users")
    user.GetFields().(*models.User).Token = token
    err := db.SelectRow(user, []string{"id"}).Scan(&userId)
    if err != nil {
        utils.HandleErr("[Handle::RejectUser]: ", err, this.Response)
        if this.Response != nil {
            this.Render([]string{"mvc/views/msg.html"}, "msg", err.Error())
        }
        return
    }

    db.QueryDeleteByIds("users", strconv.Itoa(userId))

    if this.Response != nil {
        this.Render([]string{"mvc/views/msg.html"}, "msg", "Вы успешно отписаны от рассылок Secret Oasis.")
    }
}

func (this *Handler) ResetPassword() {
    userId, err := this.CheckSid()
    if err != nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    // if !this.isAdmin() {
    //     http.Redirect(this.Response, this.Request, "/", http.StatusForbidden)
    //     return
    // }

    this.Response.Header().Set("Access-Control-Allow-Origin", "*")
    this.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    this.Response.Header().Set("Content-type", "application/json")

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(err.Error(), this.Response)
        return
    }

    pass := request["pass"].(string)
    if !utils.MatchRegexp("^.{6,36}$", pass) {
        utils.SendJSReply(map[string]interface{}{"result": "badPassword"}, this.Response)
        return
    }

    var id int
    if request["id"] == nil {
        id = userId

    } else {
        id, err =  strconv.Atoi(request["id"].(string))
        if utils.HandleErr("[Grid-Handler::ResetPassword] strconv.Atoi: ", err, this.Response) {
            utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
            return
        }
    }

    var enabled bool
    salt := strconv.Itoa(int(time.Now().Unix()))
    user := this.GetModel("users")
    user.LoadWherePart(map[string]interface{}{"id": id})
    db.SelectRow(user, []string{"enabled"}).Scan(&enabled)
    user.GetFields().(*models.User).Enabled = enabled
    user.GetFields().(*models.User).Salt = salt
    user.GetFields().(*models.User).Pass = utils.GetMD5Hash(pass + salt)
    db.QueryUpdate_(user).Scan()

    utils.SendJSReply(map[string]interface{}{"result": "ok"}, this.Response)
}
