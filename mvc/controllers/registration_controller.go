package controllers

import (
    "errors"
    "github.com/lib/pq"
    "github.com/klenin/orc/db"
    "github.com/klenin/orc/mailer"
    "github.com/klenin/orc/sessions"
    "github.com/klenin/orc/utils"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"
)

func (c *BaseController) RegistrationController() *RegistrationController {
    return new(RegistrationController)
}

type RegistrationController struct {
    Controller
}

func (this *RegistrationController) EventRegisterAction() {
    var result string; var regId int

    data, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)

        return
    }

    eventId := int(data["event_id"].(float64))

    if eventId == 1 && sessions.CheckSession(this.Response, this.Request) {
        utils.SendJSReply(map[string]interface{}{"result": "authorized"}, this.Response)

        return
    }

    if sessions.CheckSession(this.Response, this.Request) {
        userId, err := this.CheckSid()
        if err != nil {
            utils.SendJSReply(map[string]interface{}{"result": "Unauthorized"}, this.Response)

            return
        }

        var faceId int
        query := `SELECT faces.id FROM faces
            INNER JOIN registrations ON registrations.face_id = faces.id
            INNER JOIN events ON events.id = registrations.event_id
            INNER JOIN users ON users.id = faces.user_id
            WHERE users.id = $1 AND events.id = 1;`
        err = db.QueryRow(query, []interface{}{userId}).Scan(&faceId)

        if err != nil {
            utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)

            return
        }

        this.GetModel("registrations").
            LoadModelData(map[string]interface{}{"face_id": faceId, "event_id": eventId, "status": false}).
            QueryInsert("RETURNING id").
            Scan(&regId)

        if err = this.InsertUserParams(userId, regId, data["data"].([]interface{})); err != nil {
            utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)

            return
        }

    } else if eventId == 1 {
        userLogin, userPass, email, flag := "", "", "", 0

        for _, element := range data["data"].([]interface{}) {
            paramId, err := strconv.Atoi(element.(map[string]interface{})["id"].(string))
            if err != nil {
                continue
            }

            value := element.(map[string]interface{})["value"].(string)

            if paramId == 1 {
                if utils.MatchRegexp("^[ \t\v\r\n\f]{0,}$", value) {
                    utils.SendJSReply(map[string]interface{}{"result": "Заполните параметр 'Логин'."}, this.Response)

                    return
                }
                userLogin = value
                flag += 1
                continue

            } else if paramId == 2 || paramId == 3 {
                if utils.MatchRegexp("^[ \t\v\r\n\f]{0,}$", value) {
                    utils.SendJSReply(map[string]interface{}{"result": "Заполните параметр 'Пароль/Подтвердите пароль'."}, this.Response)

                    return
                }
                userPass = value
                flag += 1
                continue

            } else if paramId == 4 {
                if utils.MatchRegexp("^[ \t\v\r\n\f]{0,}$", value) {
                    utils.SendJSReply(map[string]interface{}{"result": "Заполните параметр 'Email'."}, this.Response)

                    return
                }
                email = value
                flag += 1
                continue

            } else if flag > 3 {
                break
            }
        }

        result, regId = this.Register(userLogin, userPass, email, "user")
        if result != "ok" && regId == -1 {
            utils.SendJSReply(map[string]interface{}{"result": result}, this.Response)

            return
        }

        query := `SELECT users.id
            FROM users
            INNER JOIN faces ON faces.user_id = users.id
            INNER JOIN registrations ON registrations.face_id = faces.id
            WHERE registrations.id = $1;`
        userId := db.Query(query, []interface{}{regId})[0].(map[string]interface{})["id"].(int)

        err = this.InsertUserParams(userId, regId, data["data"].([]interface{}))
        if err != nil {
            db.QueryDeleteByIds("users", strconv.Itoa(userId))
            utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)

            return
        }

    } else {
        utils.SendJSReply(map[string]interface{}{"result": "Unauthorized"}, this.Response)

        return
    }

    utils.SendJSReply(map[string]interface{}{"result": "ok"}, this.Response)
}

func (this *RegistrationController) InsertUserParams(userId, regId int, data []interface{}) (err error) {
    var paramValueIds []string

    date := time.Now().Format("2006-01-02T15:04:05Z00:00")

    for _, element := range data {
        paramId, err := strconv.Atoi(element.(map[string]interface{})["id"].(string))
        if err != nil {
            continue
        }

        if paramId == 1 || paramId == 2 || paramId == 3 {
            continue
        }

        query := `SELECT params.name, params.required, params.editable
            FROM params
            WHERE params.id = $1;`
        result := db.Query(query, []interface{}{paramId})

        name := result[0].(map[string]interface{})["name"].(string)
        required := result[0].(map[string]interface{})["required"].(bool)
        editable := result[0].(map[string]interface{})["editable"].(bool)
        value := element.(map[string]interface{})["value"].(string)

        if required && utils.MatchRegexp("^[ \t\v\r\n\f]{0,}$", value) {
            db.QueryDeleteByIds("param_vals", strings.Join(paramValueIds, ", "))
            db.QueryDeleteByIds("registrations", strconv.Itoa(regId))

            return errors.New("Заполните параметр '"+name+"'.")
        }

        if !editable {
            value = " "
        }

        var paramValId int
        err = this.GetModel("param_values").
            LoadModelData(map[string]interface{}{
                "param_id": paramId,
                "value": value,
                "date": date,
                "user_id": userId,
                "reg_id": regId}).
            QueryInsert("RETURNING id").
            Scan(&paramValId)
        if err, ok := err.(*pq.Error); ok {
            log.Println(err.Code.Name())
        }

        paramValueIds = append(paramValueIds, strconv.Itoa(paramValId))
    }

    return nil
}

func (this *RegistrationController) Register(login, password, email, role string) (result string, regId int) {
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

    if db.IsExists("users", []string{"login"}, []interface{}{login}) == true {
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

        var userId int
        this.GetModel("users").
            LoadModelData(map[string]interface{}{
                "login": login,
                "pass":  pass,
                "salt":  salt,
                "role":  role,
                "token": token,
                "enabled": false}).
            QueryInsert("RETURNING id").
            Scan(&userId)

        var faceId int
        this.GetModel("faces").
            LoadModelData(map[string]interface{}{"user_id": userId}).
            QueryInsert("RETURNING id").
            Scan(&faceId)

        this.GetModel("registrations").
            LoadModelData(map[string]interface{}{"face_id": faceId, "event_id": 1, "status": false}).
            QueryInsert("RETURNING id").
            Scan(&regId)

        return result, regId
    }

    return result, -1
}

func (this *RegistrationController) Login() {
    data, err := utils.ParseJS(this.Request, this.Response)
    if utils.HandleErr("[RegistrationController::Login]: ", err, this.Response) {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)

        return
    }

    login := data["login"].(string)
    pass := data["password"].(string)

    var id int
    var enabled bool
    var passHash, salt string
    result := make(map[string]interface{}, 1)

    if err = this.GetModel("users").
        LoadWherePart(map[string]interface{}{"login": login}).
        SelectRow([]string{"id", "pass", "salt", "enabled"}).
        Scan(&id, &passHash, &salt, &enabled);
        err != nil {
        result["result"] = "invalidCredentials"

    } else if enabled == false {
        result["result"] = "notEnabled"

    } else if passHash != utils.GetMD5Hash(pass+salt) {
        result["result"] = "badPassword"

    } else {
        result["result"] = "ok"

        sid := utils.GetRandSeq(HASH_SIZE)
        params := map[string]interface{}{"sid": sid, "enabled": true}
        where := map[string]interface{}{"id": id}
        this.GetModel("users").Update(this.isAdmin(), id, params, where)
        sessions.SetSession(this.Response, map[string]interface{}{"sid": sid})
    }

    utils.SendJSReply(result, this.Response)
}

func (this *RegistrationController) Logout() {
    userId, err := this.CheckSid()
    if err != nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        utils.SendJSReply(map[string]string{"result": "badSid"}, this.Response)

        return
    }

    var enabled bool
    if err = this.GetModel("users").
        LoadWherePart(map[string]interface{}{"id": userId}).
        SelectRow([]string{"enabled"}).
        Scan(&enabled);
        utils.HandleErr("[RegistrationController::Logout]: ", err, this.Response) {
        utils.SendJSReply(map[string]string{"result": err.Error()}, this.Response)

        return
    }

    params := map[string]interface{}{"enabled": enabled, "sid": " "}
    where := map[string]interface{}{"id": userId}
    this.GetModel("users").Update(this.isAdmin(), userId, params, where)
    sessions.ClearSession(this.Response)
    utils.SendJSReply(map[string]string{"result": "ok"}, this.Response)
}

func (this *RegistrationController) ConfirmUser(token string) {
    var userId int
    err := this.GetModel("users").
        LoadWherePart(map[string]interface{}{"token": token}).
        SelectRow([]string{"id"}).
        Scan(&userId)

    if utils.HandleErr("[RegistrationController::ConfirmUser]: ", err, this.Response) {
        if this.Response != nil {
            this.Render([]string{"mvc/views/msg.html"}, "msg", err.Error())
        }

        return
    }

    params := map[string]interface{}{"enabled": true, "token": " "}
    where := map[string]interface{}{"id": userId}
    this.GetModel("users").Update(this.isAdmin(), userId, params, where)

    if this.Response != nil {
        this.Render([]string{"mvc/views/msg.html"}, "msg", "Регистрация подтверждена.")
    }
}

func (this *RegistrationController) RejectUser(token string) {
    var userId int
    err := this.GetModel("users").
        LoadWherePart(map[string]interface{}{"token": token}).
        SelectRow([]string{"id"}).
        Scan(&userId)

    if utils.HandleErr("[RegistrationController::RejectUser]: ", err, this.Response) {
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
