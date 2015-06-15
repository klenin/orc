package controllers

import (
    "github.com/lib/pq"
    "github.com/orc/db"
    "github.com/orc/mailer"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "net/http"
    "strconv"
    "strings"
    "time"
)

func (this *GridHandler) GetPersonRequestFromGroup() {
    if !sessions.CheckSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    faceId, err := strconv.Atoi(request["face_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    groupRegId, err := strconv.Atoi(request["group_reg_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    query := `SELECT forms.id as form_id, forms.name as form_name,
            params.id as param_id, params.name as param_name, params.required, params.editable,
            events.name as event_name, events.id as event_id,
            param_types.name as type, param_values.id as param_val_id, param_values.value
        FROM events_forms
        INNER JOIN events ON events.id = events_forms.event_id
        INNER JOIN forms ON forms.id = events_forms.form_id
        INNER JOIN params ON forms.id = params.form_id
        INNER JOIN param_types ON param_types.id = params.param_type_id
        INNER JOIN param_values ON params.id = param_values.param_id
        INNER JOIN reg_param_vals ON reg_param_vals.param_val_id = param_values.id
        INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN group_registrations ON group_registrations.event_id = events.id
        INNER JOIN groups ON group_registrations.group_id = groups.id
        INNER JOIN regs_groupregs ON regs_groupregs.reg_id = registrations.id
            AND regs_groupregs.groupreg_id = group_registrations.id
        WHERE group_registrations.id = $1 AND faces.id = $2 ORDER BY forms.id, params.id;`

    utils.SendJSReply(
        map[string]interface{}{
            "result": "ok",
            "data": db.Query(query, []interface{}{groupRegId, faceId}),
            "role": this.isAdmin()},
        this.Response)
}

func (this *GridHandler) GetPersonRequest() {
    if !sessions.CheckSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    regId, err := strconv.Atoi(request["reg_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    query := `SELECT forms.id as form_id, forms.name as form_name,
            params.id as param_id, params.name as param_name, params.required, params.editable,
            events.name as event_name, events.id as event_id,
            param_types.name as type, param_values.id as param_val_id, param_values.value
        FROM events_forms
        INNER JOIN events ON events.id = events_forms.event_id
        INNER JOIN forms ON forms.id = events_forms.form_id
        INNER JOIN params ON forms.id = params.form_id
        INNER JOIN param_types ON param_types.id = params.param_type_id
        INNER JOIN param_values ON params.id = param_values.param_id
        INNER JOIN reg_param_vals ON reg_param_vals.param_val_id = param_values.id
        INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
            AND events.id = registrations.event_id
        WHERE registrations.id = $1 ORDER BY forms.id, params.id;`

    utils.SendJSReply(
        map[string]interface{}{
            "result": "ok",
            "data": db.Query(query, []interface{}{regId}),
            "role": this.isAdmin()},
        this.Response)
}

func (this *GridHandler) ConfirmOrRejectPersonRequest() {
    if !sessions.CheckSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        http.Redirect(this.Response, this.Request, "/", http.StatusForbidden)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    eventId, err := strconv.Atoi(request["event_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    regId, err := strconv.Atoi(request["reg_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    query := `SELECT param_values.value, users.id as user_id
        FROM reg_param_vals
        INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
        INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
        INNER JOIN params ON params.id = param_values.param_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE params.id in (4, 5, 6, 7) AND users.id in (
            SELECT users.id FROM registrations INNER JOIN events ON events.id = registrations.event_id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN users ON users.id = faces.user_id
            WHERE registrations.id = $1
        ) ORDER BY params.id;`

    data := db.Query(query, []interface{}{regId})

    if len(data) < 2 {
        utils.SendJSReply(
            map[string]interface{}{"result": "Нет регистрационных данных пользователя"},
            this.Response)
        return
    }

    email := data[0].(map[string]interface{})["value"].(string)

    to := data[1].(map[string]interface{})["value"].(string)
    to += " " + data[2].(map[string]interface{})["value"].(string)
    to += " " + data[3].(map[string]interface{})["value"].(string)

    event := db.Query(
        "SELECT name FROM events WHERE id=$1;",
        []interface{}{eventId})[0].(map[string]interface{})["name"].(string)

    if request["confirm"].(bool) {
        if eventId == 1 {
            utils.SendJSReply(map[string]interface{}{"result": "Эту заявку нельзя подтвердить письмом"}, this.Response)
        } else {
            if mailer.SendEmailToConfirmRejectPersonRequest(to, email, event, true) {
                utils.SendJSReply(map[string]interface{}{"result": "Письмо с подтверждением заявки отправлено"}, this.Response)
            } else {
                utils.SendJSReply(map[string]interface{}{"result": "Ошибка. Письмо с подтверждением заявки не отправлено"}, this.Response)
            }
        }

    } else {
        if eventId == 1 {
            utils.SendJSReply(map[string]interface{}{"result": "Эту заявку нельзя отклонить письмом"}, this.Response)
        } else {
            query := `DELETE
                FROM param_values USING reg_param_vals
                WHERE param_values.id in (SELECT reg_param_vals.param_val_id WHERE reg_param_vals.reg_id = $1);`
            db.Query(query, []interface{}{regId})

            query = `DELETE FROM registrations WHERE id = $1;`
            db.Query(query, []interface{}{regId})

            if mailer.SendEmailToConfirmRejectPersonRequest(to, email, event, false) {
                utils.SendJSReply(map[string]interface{}{"result": "Письмо с отклонением заявки отправлено"}, this.Response)
            } else {
                utils.SendJSReply(map[string]interface{}{"result": "Ошибка. Письмо с отклонением заявки не отправлено"}, this.Response)
            }
        }
    }
}

func (this *GridHandler) EditParams() {
    userId, err := this.CheckSid()
    if err != nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    date := time.Now().Format("2006-01-02T15:04:05Z00:00")

    for _, v := range request["data"].([]interface{}) {
        paramValId, err := strconv.Atoi(v.(map[string]interface{})["param_val_id"].(string))
        if err != nil {
            utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
            return
        }

        query := `SELECT params.name, params.required, params.editable
            FROM params
            INNER JOIN param_values ON param_values.param_id = params.id
            WHERE param_values.id = $1;`
        result := db.Query(query, []interface{}{paramValId})

        name := result[0].(map[string]interface{})["name"].(string)
        required := result[0].(map[string]interface{})["required"].(bool)
        editable := result[0].(map[string]interface{})["editable"].(bool)
        value := v.(map[string]interface{})["value"].(string)

        if required && utils.MatchRegexp("^[ \t\v\r\n\f]{0,}$", value) {
            utils.SendJSReply(map[string]interface{}{"result": "Заполните параметр '"+name+"'"}, this.Response)
            return
        }

        if !this.isAdmin() && !editable {
            continue
        }

        if value == "" {
            value = " "
        }

        paramValue := this.GetModel("param_values")
        paramValue.LoadModelData(map[string]interface{}{"value": value, "date": date, "user_id": userId})
        paramValue.LoadWherePart(map[string]interface{}{"id": paramValId})
        db.QueryUpdate_(paramValue).Scan()
    }

    utils.SendJSReply(map[string]interface{}{"result": "Изменения сохранены"}, this.Response)
}

func (this *Handler) AddPerson() {
    userId, err := this.CheckSid()
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": "Unauthorized"}, this.Response)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    groupId, err := strconv.Atoi(request["group_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    var groupName string
    db.QueryRow("SELECT name FROM groups WHERE id = $1;", []interface{}{groupId}).Scan(&groupName)

    date := time.Now().Format("2006-01-02T15:04:05Z00:00")
    token := utils.GetRandSeq(HASH_SIZE)
    to, address, headName := "", "", ""

    query := `SELECT param_values.value
        FROM reg_param_vals
        INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
        INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
        INNER JOIN params ON params.id = param_values.param_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE params.id in (5, 6, 7) AND users.id = $1 AND events.id = 1 ORDER BY params.id;`
    data := db.Query(query, []interface{}{userId})

    if len(data) < 3 {
        utils.SendJSReply(map[string]interface{}{"result": "Данные о руководителе группы отсутсвуют"}, this.Response)
        return

    } else {
        headName = data[0].(map[string]interface{})["value"].(string)
        headName += " " + data[1].(map[string]interface{})["value"].(string)
        headName += " " + data[2].(map[string]interface{})["value"].(string)
    }

    var faceId int
    face := this.GetModel("faces")
    db.QueryInsert_(face, "RETURNING id").Scan(&faceId)

    persons := this.GetModel("persons")
    persons.LoadModelData(map[string]interface{}{"face_id": faceId, "group_id": groupId, "status": false, "token": token})
    db.QueryInsert_(persons, "").Scan()

    var regId int
    registration := this.GetModel("registrations")
    registration.LoadModelData(map[string]interface{}{"face_id": faceId, "event_id": 1})
    db.QueryInsert_(registration, "RETURNING id").Scan(&regId)

    var paramValueIds []string

    for _, element := range request["data"].([]interface{}) {
        paramId, err := strconv.Atoi(element.(map[string]interface{})["id"].(string))
        if err != nil {
            continue
        }

        query := `SELECT params.name FROM params WHERE params.id = $1;`
        res := db.Query(query, []interface{}{paramId})

        name := res[0].(map[string]interface{})["name"].(string)
        value := element.(map[string]interface{})["value"].(string)

        if utils.MatchRegexp("^[ \t\v\r\n\f]{0,}$", value) {
            db.QueryDeleteByIds("param_vals", strings.Join(paramValueIds, ", "))
            db.QueryDeleteByIds("registrations", strconv.Itoa(regId))
            db.QueryDeleteByIds("faces", strconv.Itoa(faceId))
            utils.SendJSReply(map[string]interface{}{"result": "Заполните параметр '"+name+"'."}, this.Response)
            return
        }

        var paramValId int
        paramValues := this.GetModel("param_values")
        paramValues.LoadModelData(map[string]interface{}{"param_id": paramId, "value": value, "date": date, "user_id": userId})
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

        if paramId == 4 {
            address = value
        } else if paramId == 5 || paramId == 6 || paramId == 7 {
            to += value + " "
        }
    }

    if !mailer.InviteToGroup(to, address, token, headName, groupName) {
        utils.SendJSReply(
            map[string]interface{}{
                "result": "Вы указали неправильный email, отправить письмо-приглашенине невозможно"},
            this.Response)
        return
    }

    utils.SendJSReply(map[string]interface{}{"result": "ok"}, this.Response)
}
