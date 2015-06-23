package controllers

import (
    "github.com/orc/db"
    "github.com/orc/utils"
    "net/http"
    "strconv"
    "time"
    "github.com/orc/mailer"
)

func (this *GridHandler) RegGroup() {
    userId, err := this.CheckSid()
    if err != nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    }

    groupId, err := strconv.Atoi(request["group_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    }

    eventId, err := strconv.Atoi(request["event_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    }

    event := this.GetModel("events")
    event.LoadWherePart(map[string]interface{}{"id": eventId})

    var eventName string
    err = db.SelectRow(event, []string{"name"}).Scan(&eventName)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    faceId, groupName := -1, ""
    query := `SELECT groups.face_id, groups.name FROM groups
        INNER JOIN faces ON faces.id = groups.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 AND groups.id = $2;`
    err = db.QueryRow(query, []interface{}{userId, groupId}).Scan(&faceId, &groupName)

    if (err != nil || faceId == -1 || groupName == "") && !this.isAdmin() {
        utils.SendJSReply(map[string]interface{}{"result": "Вы не являетесь владельцем группы"}, this.Response)
        return
    }

    if db.IsExists_("group_registrations", []string{"group_id", "event_id"}, []interface{}{groupId, eventId}) {
        utils.SendJSReply(map[string]interface{}{"result": "Группа уже зарегистрированна в этом мероприятии"}, this.Response)
        return
    }

    var groupregId int
    groupReg := this.GetModel("group_registrations")
    groupReg.LoadModelData(map[string]interface{}{"event_id": eventId, "group_id": groupId})
    db.QueryInsert(groupReg, "RETURNING id").Scan(&groupregId)

    query = `SELECT persons.status, faces.id FROM persons
        INNER JOIN groups ON groups.id = persons.group_id
        INNER JOIN faces ON faces.id = persons.face_id
        WHERE groups.id = $1;`
    data := db.Query(query, []interface{}{groupId})

    query = `SELECT params.id FROM events_forms
        INNER JOIN events ON events.id = events_forms.event_id
        INNER JOIN forms ON forms.id = events_forms.form_id
        INNER JOIN params ON forms.id = params.form_id
        WHERE events.id = $1 ORDER BY forms.id;`
    params := db.Query(query, []interface{}{eventId})

    date := time.Now().Format("2006-01-02T15:04:05Z00:00")

    for _, v := range data {
        status := v.(map[string]interface{})["status"].(bool)
        personfaceId := v.(map[string]interface{})["id"]

        if !status {
            continue
        }

        var regId int
        regs := this.GetModel("registrations")
        regs.LoadModelData(map[string]interface{}{"face_id": personfaceId, "event_id": eventId})
        db.QueryInsert(regs, "RETURNING id").Scan(&regId)

        to := v.(map[string]interface{})["name"].(string)
        address := v.(map[string]interface{})["email"].(string)
        if !mailer.AttendAnEvent(to, address, eventName, groupName) {
            utils.SendJSReply(map[string]interface{}{"result": "Ошибка. Письмо с уведомлением не отправлено."}, this.Response)
        }
        regsGroupRegs := this.GetModel("regs_groupregs")
        regsGroupRegs.LoadModelData(map[string]interface{}{"groupreg_id": groupregId, "reg_id": regId})
        db.QueryInsert(regsGroupRegs, "").Scan()


        for _, elem := range params {
            param_id := int(elem.(map[string]interface{})["id"].(int))

            var paramValId int
            paramValues := this.GetModel("param_values")
            paramValues.LoadModelData(map[string]interface{}{"param_id": param_id, "value": " ", "date": date, "user_id": userId})
            db.QueryInsert(paramValues, "RETURNING id").Scan(&paramValId)

            regParamValue := this.GetModel("reg_param_vals")
            regParamValue.LoadModelData(map[string]interface{}{
                "reg_id":        regId,
                "param_val_id":  paramValId})
            db.QueryInsert(regParamValue, "").Scan()
        }

    }

    utils.SendJSReply(map[string]interface{}{"result": "ok"}, this.Response)
}

func (this *Handler) ConfirmInvitationToGroup(token string) {
    person := this.GetModel("persons")
    person.LoadWherePart(map[string]interface{}{"token": token})

    var faceId, groupId int
    err = db.SelectRow(person, []string{"face_id", "group_id",}).Scan(&faceId, &groupId)

    if err != nil {
        if this.Response != nil {
            this.Render([]string{"mvc/views/msg.html"}, "msg", "Неверный токен.")
        }
        return
    }

    person = this.GetModel("persons")
    person.LoadModelData(map[string]interface{}{"status": true, "token": " "})
    person.LoadWherePart(map[string]interface{}{"token": token})
    db.QueryUpdate(person).Scan()

    if this.Response != nil {
        this.Render([]string{"mvc/views/msg.html"}, "msg", "Вы успешно присоединены к группе.")
    }
}

func (this *Handler) RejectInvitationToGroup(token string) {
    db.Exec("DELETE FROM persons WHERE token = $1;", []interface{}{token})

    if this.Response != nil {
        this.Render([]string{"mvc/views/msg.html"}, "msg", "Запрос о присоединении к группе успешно отклонен.")
    }
}

func (this *GridHandler) IsRegGroup() {
    _, err := this.CheckSid()
    if err != nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    }

    groupId, err := strconv.Atoi(request["group_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    }

    addDelFlag := !db.IsExists_("group_registrations", []string{"group_id"}, []interface{}{groupId})
    utils.SendJSReply(map[string]interface{}{"result": "ok", "addDelFlag": addDelFlag}, this.Response)
}
