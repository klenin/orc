package controllers

import (
    "errors"
    "github.com/orc/db"
    "github.com/lib/pq"
    "github.com/orc/mvc/models"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "strconv"
    "strings"
    "time"
)

func (this *Handler) GetHistoryRequest() {
    userId, err := this.CheckSid()
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": "Unauthorized"}, this.Response)
        return
    }

    data, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    eventId, err := strconv.Atoi(data["event_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    query := `SELECT params.id as param_id, params.name as param_name,
                param_types.name as type, param_values.value, forms.id as form_id
            FROM events
            INNER JOIN events_forms ON events_forms.event_id = events.id
            INNER JOIN forms ON events_forms.form_id = forms.id
            INNER JOIN registrations ON events.id = registrations.event_id
            INNER JOIN reg_param_vals ON reg_param_vals.reg_id = registrations.id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN users ON users.id = faces.user_id
            INNER JOIN params ON params.form_id = forms.id
            INNER JOIN param_types ON param_types.id = params.param_type_id
            INNER JOIN param_values ON param_values.param_id = params.id
                AND reg_param_vals.param_val_id = param_values.id
            WHERE users.id = $1 AND events.id = $2;`

    utils.SendJSReply(map[string]interface{}{"result": "ok", "data": db.Query(query, []interface{}{userId, eventId})}, this.Response)
}

func (this *Handler) GetListHistoryEvents() {
    userId, err := this.CheckSid()
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": "Unauthorized"}, this.Response)
        return
    }

    data, err := utils.ParseJS(this.Request, this.Response)
    if  err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    ids := make(map[string]interface{}, 1)
    ids["form_id"] = make([]interface{}, 0)
    if data["form_ids"] == nil || len(data["form_ids"].(map[string]interface{})["form_id"].([]interface{})) == 0 {
        utils.SendJSReply(map[string]interface{}{"result": "Нет данных"}, this.Response)
        return
    }

    for _, v := range data["form_ids"].(map[string]interface{})["form_id"].([]interface{}) {
        ids["form_id"] = append(ids["form_id"].([]interface{}), int(v.(float64)))
    }

    eventsForms := this.GetModel("events_forms")
    eventsForms.LoadWherePart(ids)
    eventsForms.SetCondition(models.OR)
    events := db.Select(eventsForms, []string{"event_id"})

    if len(events) == 0 {
        utils.SendJSReply(map[string]interface{}{"result": "Нет данных"}, this.Response)
        return
    }

    query := `SELECT DISTINCT events.id, events.name FROM events
        INNER JOIN events_forms ON events_forms.event_id = events.id
        INNER JOIN forms ON events_forms.form_id = forms.id
        INNER JOIN registrations ON registrations.event_id = events.id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id=$1 AND events.id IN (`

    var i int
    params := []interface{}{userId}

    for i = 2; i < len(events); i++ {
        query += "$" + strconv.Itoa(i) + ", "
        params = append(params, int(events[i-2].(map[string]interface{})["event_id"].(int)))
    }
    query += "$" + strconv.Itoa(i) + ")"
    params = append(params, int(events[i-2].(map[string]interface{})["event_id"].(int)))

    utils.SendJSReply(map[string]interface{}{"result": "ok", "data": db.Query(query, params)}, this.Response)
}

func (this *Handler) RegPerson() {
    var result string
    var regId int

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
        face := this.GetModel("faces")
        face.LoadModelData(map[string]interface{}{"user_id": userId})
        db.QueryInsert_(face, "RETURNING id").Scan(&faceId)

        registration := this.GetModel("registrations")
        registration.LoadModelData(map[string]interface{}{"face_id": faceId, "event_id": eventId})
        db.QueryInsert_(registration, "RETURNING id").Scan(&regId)

        if err = this.InsertUserParams(regId, data["data"].([]interface{})); err != nil {
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

        result, regId = this.HandleRegister_(userLogin, userPass, email, "user")
        if result != "ok" && regId == -1 {
            utils.SendJSReply(map[string]interface{}{"result": result}, this.Response)
            return
        }

        err = this.InsertUserParams(regId, data["data"].([]interface{}))
        if err != nil {
            query := `SELECT users.id
                FROM users
                INNER JOIN faces ON faces.usr_id = users.id
                INNER JOIN registrations ON registrations.face_id = faces.id
                WHERE registrations.id = $1;`
            userId := db.Query(query, []interface{}{regId})[0].(map[string]interface{})["id"].(int)
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

func (this *Handler) GetRequest(id string) {
    eventId, err := strconv.Atoi(id)
    if utils.HandleErr("[Handler::GetRequestGetRequest] event_id Atoi: ", err, this.Response) {
        return
    }

    if !sessions.CheckSession(this.Response, this.Request) && eventId != 1 {
        this.Render([]string{"mvc/views/loginpage.html", "mvc/views/login.html"}, "loginpage", nil)
        return
    }

    query := `SELECT forms.id as form_id, forms.name as form_name,
            params.id as param_id, params.name as param_name, params.required, params.editable,
            param_types.name as type, events.name as event_name, events.id as event_id
        FROM events_forms
        INNER JOIN events ON events.id = events_forms.event_id
        INNER JOIN forms ON forms.id = events_forms.form_id
        INNER JOIN params ON forms.id = params.form_id
        INNER JOIN param_types ON param_types.id = params.param_type_id
        WHERE events.id = $1 ORDER BY forms.id, params.id;`
    res := db.Query(query, []interface{}{eventId})

    this.Render([]string{"mvc/views/item.html"}, "item", map[string]interface{}{"data": res})
}

func (this *Handler) InsertUserParams(regId int, data []interface{}) (err error) {
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
        paramValues := this.GetModel("param_values")
        paramValues.LoadModelData(map[string]interface{}{"param_id": paramId, "value": value, "date": date})
        err = db.QueryInsert_(paramValues, "RETURNING id").Scan(&paramValId)
        if err, ok := err.(*pq.Error); ok {
            println(err.Code.Name())
        }

        regParamValue := this.GetModel("reg_param_vals")
        regParamValue.LoadModelData(map[string]interface{}{
            "reg_id":       regId,
            "param_val_id": paramValId})
        db.QueryInsert_(regParamValue, "").Scan()

        paramValueIds = append(paramValueIds, strconv.Itoa(paramValId))
    }

    return nil
}
