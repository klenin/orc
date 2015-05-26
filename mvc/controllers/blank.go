package controllers

import (

    "github.com/orc/db"
    "github.com/orc/mailer"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "net/http"
    "strconv"
)

func (this *GridHandler) GetPersonRequestFromGroup() {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }
    person_id, err := strconv.Atoi(request["person_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    group_reg_id, err := strconv.Atoi(request["group_reg_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    q := `SELECT users.id FROM users
        INNER JOIN faces ON faces.user_id = users.id
        INNER JOIN persons ON persons.face_id = faces.id
        WHERE persons.id = $1;`

    user_id := db.Query(q, []interface{}{person_id})[0].(map[string]interface{})["id"].(int)

    query := `SELECT forms.id as form_id, forms.name as form_name, params.id as param_id,
        events.name as event_name, events.id as event_id, params.name as param_name,
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
        INNER JOIN users ON users.id = faces.user_id
        INNER JOIN group_registrations ON group_registrations.event_id = events.id
        INNER JOIN groups ON group_registrations.group_id = groups.id
        INNER JOIN regs_groupregs ON regs_groupregs.reg_id = registrations.id AND regs_groupregs.groupreg_id = group_registrations.id
        WHERE group_registrations.id = $1 AND users.id = $2 ORDER BY forms.id, params.id;`

    utils.SendJSReply(
        map[string]interface{}{"result": "ok", "data": db.Query(query, []interface{}{group_reg_id, user_id})},
        this.Response)
}

func (this *GridHandler) GetPersonRequest() {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    reg_id, err := strconv.Atoi(request["reg_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    query := `SELECT forms.id as form_id, forms.name as form_name, params.id as param_id,
        events.name as event_name, events.id as event_id, params.name as param_name,
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
        map[string]interface{}{"result": "ok", "data": db.Query(query, []interface{}{reg_id})},
        this.Response)
}

func (this *GridHandler) ConfirmOrRejectPersonRequest() {
    if !sessions.CheackSession(this.Response, this.Request) {
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

    event_id, err := strconv.Atoi(request["event_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    reg_id, err := strconv.Atoi(request["reg_id"].(string))
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

    data := db.Query(query, []interface{}{reg_id})

    if len(data) < 2 {
        utils.SendJSReply(map[string]interface{}{"result": "Нет регистрационных данных пользователя"}, this.Response)
        return
    }

    email := data[0].(map[string]interface{})["value"].(string)

    to := data[1].(map[string]interface{})["value"].(string)
    to += " " + data[2].(map[string]interface{})["value"].(string)
    to += " " + data[3].(map[string]interface{})["value"].(string)

    event := db.Query("SELECT name FROM events WHERE id=$1;", []interface{}{event_id})[0].(map[string]interface{})["name"].(string)

    if request["confirm"].(bool) {
        if event_id == 1 {
            utils.SendJSReply(map[string]interface{}{"result": "Эту заявку нельзя подтвердить письмом"}, this.Response)
        } else {
            if mailer.SendEmailToConfirmRejectPersonRequest(to, email, event, true) {
                utils.SendJSReply(map[string]interface{}{"result": "Письмо с подтверждением заявки отправлено"}, this.Response)
            } else {
                utils.SendJSReply(map[string]interface{}{"result": "Ошибка. Письмо с подтверждением заявки не отправлено"}, this.Response)
            }
        }

    } else {
        if event_id == 1 {
            utils.SendJSReply(map[string]interface{}{"result": "Эту заявку нельзя отклонить письмом"}, this.Response)
        } else {
            query := `DELETE
                FROM param_values USING reg_param_vals
                WHERE param_values.id in (SELECT reg_param_vals.param_val_id WHERE reg_param_vals.reg_id = $1);`
            db.Query(query, []interface{}{reg_id})

            query = `DELETE FROM registrations WHERE id = $1;`
            db.Query(query, []interface{}{reg_id})

            if mailer.SendEmailToConfirmRejectPersonRequest(to, email, event, false) {
                utils.SendJSReply(map[string]interface{}{"result": "Письмо с отклонением заявки отправлено"}, this.Response)
            } else {
                utils.SendJSReply(map[string]interface{}{"result": "Ошибка. Письмо с отклонением заявки не отправлено"}, this.Response)
            }
        }
    }
}

func (this *GridHandler) EditParams() {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    for _, v := range request["data"].([]interface{}) {

        param_val_id, err := strconv.Atoi(v.(map[string]interface{})["param_val_id"].(string))
        if err != nil {
            utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
            return
        }

        value := v.(map[string]interface{})["value"].(string)

        // !!!
        if value == "" {
            value = " "
        }

        param_value := this.GetModel("param_values")
        param_value.LoadModelData(map[string]interface{}{"value": value})
        param_value.LoadWherePart(map[string]interface{}{"id": param_val_id})
        db.QueryUpdate_(param_value).Scan()
    }

    utils.SendJSReply(map[string]interface{}{"result": "Изменения сохранены"}, this.Response)
}
