package controllers

import (
    "github.com/orc/db"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "net/http"
    "strconv"
    "github.com/orc/mailer"
)

func (this *GridHandler) RegGroup() {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    }

    group_id, err := strconv.Atoi(request["group_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    }

    event_id, err := strconv.Atoi(request["event_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    }

    event := this.GetModel("events")
    event.LoadWherePart(map[string]interface{}{"id": event_id})

    var eventName string
    err = db.SelectRow(event, []string{"name"}).Scan(&eventName)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    face_id, groupName := -1, ""
    query := `SELECT groups.face_id, groups.name FROM groups
        INNER JOIN faces ON faces.id = groups.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 AND groups.id = $2;`
    err = db.QueryRow(query, []interface{}{user_id, group_id}).Scan(&face_id, &groupName)

    if err != nil || face_id == -1 || groupName == "" {
        utils.SendJSReply(map[string]interface{}{"result": "Вы не являетесь владельцем группы"}, this.Response)
        return
    }

    var groupreg_id int
    group_reg := this.GetModel("group_registrations")
    group_reg.LoadModelData(map[string]interface{}{"event_id": event_id, "group_id": group_id})
    db.QueryInsert_(group_reg, "RETURNING id").Scan(&groupreg_id)

    query = `SELECT persons.name, persons.email, persons.status, users.id as user_id FROM persons
        INNER JOIN groups ON groups.id = persons.group_id
        INNER JOIN faces ON faces.id = persons.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE groups.id = $1;`
    data := db.Query(query, []interface{}{group_id})

    query = `SELECT params.id FROM events_forms
        INNER JOIN events ON events.id = events_forms.event_id
        INNER JOIN forms ON forms.id = events_forms.form_id
        INNER JOIN params ON forms.id = params.form_id
        WHERE events.id = $1 ORDER BY forms.id;`
    params := db.Query(query, []interface{}{event_id})

    for _, v := range data {
        status := v.(map[string]interface{})["status"].(bool)
        p_user_id := v.(map[string]interface{})["user_id"].(int)

        if !status {
            continue
        }

        var face_id int
        face := this.GetModel("faces")
        face.LoadModelData(map[string]interface{}{"user_id": p_user_id})
        db.QueryInsert_(face, "RETURNING id").Scan(&face_id)

        var reg_id int
        regs := this.GetModel("registrations")
        regs.LoadModelData(map[string]interface{}{"face_id": face_id, "event_id": event_id})
        db.QueryInsert_(regs, "RETURNING id").Scan(&reg_id)

        to := v.(map[string]interface{})["name"].(string)
        address := v.(map[string]interface{})["email"].(string)
        if !mailer.AttendAnEvent(to, address, eventName, groupName) {
            utils.SendJSReply(map[string]interface{}{"result": "Ошибка. Письмо с уведомлением не отправлено."}, this.Response)
        }
        regs_groupregs := this.GetModel("regs_groupregs")
        regs_groupregs.LoadModelData(map[string]interface{}{"groupreg_id": groupreg_id, "reg_id": reg_id})
        db.QueryInsert_(regs_groupregs, "").Scan()


        for _, p := range params {
            param_id := int(p.(map[string]interface{})["id"].(int))

            var param_val_id int
            paramValues := this.GetModel("param_values")
            paramValues.LoadModelData(map[string]interface{}{"param_id": param_id, "value": "  "})
            db.QueryInsert_(paramValues, "RETURNING id").Scan(&param_val_id)

            regParamValue := this.GetModel("reg_param_vals")
            regParamValue.LoadModelData(map[string]interface{}{
                "reg_id":        reg_id,
                "param_val_id":  param_val_id})
            db.QueryInsert_(regParamValue, "").Scan()
        }

    }

    utils.SendJSReply(map[string]interface{}{"result": "ok"}, this.Response)
}

func (this *Handler) ConfirmInvitationToGroup(token string) {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    var face_id int
    query := `SELECT faces.id
        FROM registrations
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN users ON faces.user_id = users.id
        WHERE users.id = $1 AND events.id = $2;`
    db.QueryRow(query, []interface{}{user_id, 1}).Scan(&face_id)

    person := this.GetModel("persons")
    person.LoadWherePart(map[string]interface{}{"token": token})

    var group_id int
    err := db.SelectRow(person, []string{"group_id",}).Scan(&group_id)

    if err != nil {
        if this.Response != nil {
            this.Render([]string{"mvc/views/msg.html"}, "msg", "Неверный токен.")
        }
        return
    }


    if db.IsExists_("persons", []string{"face_id", "group_id"}, []interface{}{face_id, group_id}) {
        if this.Response != nil {
            this.Render([]string{"mvc/views/msg.html"}, "msg", "Вы уже состоите в группе.")
        }
        return
    }

    person = this.GetModel("persons")
    person.LoadModelData(map[string]interface{}{"face_id": face_id, "status": true, "token": " "})
    person.LoadWherePart(map[string]interface{}{"token": token})
    db.QueryUpdate_(person).Scan()

    if this.Response != nil {
        this.Render([]string{"mvc/views/msg.html"}, "msg", "Вы успешно присоединены к группе.")
    }
}

func (this *Handler) RejectInvitationToGroup(token string) {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    query := "DELETE FROM persons WHERE token = $1;"
    db.Exec(query, []interface{}{token})

    if this.Response != nil {
        this.Render([]string{"mvc/views/msg.html"}, "msg", "Запрос о присоединении к группе успешно отклонен.")
    }
}
