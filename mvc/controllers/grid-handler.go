package controllers

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/orc/db"
    "github.com/orc/mailer"
    "github.com/orc/mvc/models"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "net/http"
    "strconv"
    "strings"
    "errors"
)

func (c *BaseController) GridHandler() *GridHandler {
    return new(GridHandler)
}

type GridHandler struct {
    Controller
}

func (this *GridHandler) GetSubTable() {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(err.Error(), this.Response)
        return
    }

    model := GetModel(request["table"].(string))
    index, _ := strconv.Atoi(request["index"].(string))
    subModel := GetModel(model.GetSubTable(index))
    subModel.LoadWherePart(map[string]interface{}{model.GetSubField(): request["id"]})
    refFields, refData := GetModelRefDate(subModel)

    response, err := json.Marshal(map[string]interface{}{
        "name":      subModel.GetTableName(),
        "caption":   subModel.GetCaption(),
        "colnames":  subModel.GetColNames(),
        "columns":   subModel.GetColumns(),
        "reffields": refFields,
        "refdata":   refData})
    if utils.HandleErr("[GridHandler::GetSubTable] Marshal: ", err, this.Response) {
        return
    }

    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *GridHandler) Select(tableName string) {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        http.Redirect(this.Response, this.Request, "/", http.StatusForbidden)
        return
    }

    model := GetModel(tableName)
    refFields, refData := GetModelRefDate(model)

    this.Render([]string{"mvc/views/table.html"}, "table", Model{
        RefData:   refData,
        RefFields: refFields,
        TableName: model.GetTableName(),
        ColNames:  model.GetColNames(),
        Columns:   model.GetColumns(),
        Caption:   model.GetCaption(),
        Sub:       model.GetSub()})
}

func (this *GridHandler) Edit(tableName string) {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "", http.StatusUnauthorized)
        return
    }

    model := GetModel(tableName)
    if model == nil {
        utils.HandleErr("[GridHandler::Edit] GetModel: invalid model", nil, this.Response)
        return
    }

    params := make(map[string]interface{}, len(model.GetColumns()))
    for i := 0; i < len(model.GetColumns()); i++ {
        params[model.GetColumnByIdx(i)] = this.Request.PostFormValue(model.GetColumnByIdx(i))
    }

    oper := this.Request.PostFormValue("oper")

    switch oper {
    case "edit":
        id, err := strconv.Atoi(this.Request.PostFormValue("id"))
        if utils.HandleErr("[GridHandler::Edit] strconv.Atoi id: ", err, this.Response) {
            return
        }

        if tableName == "groups" && !this.isAdmin() {
            query := `SELECT groups.face_id FROM groups
                INNER JOIN faces ON faces.id = groups.face_id
                INNER JOIN users ON users.id = faces.user_id
                WHERE users.id = $1 AND groups.id = $2;`

            face_id := -1
            err := db.QueryRow(query, []interface{}{user_id, id}).Scan(&face_id)

            if err != nil {
                utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
                return
            } else if face_id == -1 {
                utils.SendJSReply(map[string]interface{}{"result": "Нет прав редактировать эту группу."}, this.Response)
                return
            }
            params["face_id"] = face_id
        }
        model.LoadModelData(params)
        model.LoadWherePart(map[string]interface{}{"id": id})
        db.QueryUpdate_(model).Scan()
        break
    case "add":
        if tableName == "groups" {
            var face_id int
            face := GetModel("faces")
            face.LoadModelData(map[string]interface{}{"user_id": user_id})
            db.QueryInsert_(face, "RETURNING id").Scan(&face_id)
            params["face_id"] = face_id

        } else if tableName == "persons" {
            to := params["name"].(string)
            address := params["email"].(string)
            token := utils.GetRandSeq(HASH_SIZE)
            params["token"] = token

            query := `SELECT param_values.value
                FROM reg_param_vals
                INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
                INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
                INNER JOIN params ON params.id = param_values.param_id
                INNER JOIN events ON events.id = registrations.event_id
                INNER JOIN faces ON faces.id = registrations.face_id
                INNER JOIN users ON users.id = faces.user_id
                WHERE params.id in (5, 6, 7) AND users.id = $1 AND events.id = 1 ORDER BY params.id;`

            data := db.Query(query, []interface{}{user_id})
            headName := data[0].(map[string]interface{})["value"].(string)
            headName += " " + data[1].(map[string]interface{})["value"].(string)
            headName += " " + data[2].(map[string]interface{})["value"].(string)

            group_id, err := strconv.Atoi(params["group_id"].(string))
            if utils.HandleErr("[GridHandler::Edit] group_id Atoi: ", err, this.Response) {
                return
            }

            var groupName string
            db.QueryRow("SELECT name FROM groups WHERE id = $1;", []interface{}{group_id}).Scan(&groupName)

            if !mailer.InviteToGroup(to, address, token, headName, groupName) {
                utils.HandleErr("Mailer: ", errors.New("Письмо с приглашением в группу не отправлено."), this.Response)
            }
        }
        model.LoadModelData(params)
        db.QueryInsert_(model, "").Scan()
        break
    case "del":
        db.QueryDeleteByIds(tableName, this.Request.PostFormValue("id"))
        break
    }
}

func (this *GridHandler) ResetPassword() {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        http.Redirect(this.Response, this.Request, "/", http.StatusForbidden)
        return
    }

    this.Response.Header().Set("Access-Control-Allow-Origin", "*")
    this.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    this.Response.Header().Set("Content-type", "application/json")

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(err.Error(), this.Response)
        return
    }

    pass1 := request["pass1"].(string)
    pass2 := request["pass2"].(string)

    if !utils.MatchRegexp("^.{6,36}$", pass1) || !utils.MatchRegexp("^.{6,36}$", pass2) {
        utils.SendJSReply(map[string]interface{}{"result": "badPassword"}, this.Response)
        return
    } else if pass1 != pass2 {
        utils.SendJSReply(map[string]interface{}{"result": "differentPasswords"}, this.Response)
        return
    }

    id, err :=  strconv.Atoi(request["id"].(string))
    if utils.HandleErr("[Grid-Handler::ResetPassword] strconv.Atoi: ", err, this.Response) {
        return
    }

    user := GetModel("users")
    user.LoadWherePart(map[string]interface{}{"id": id})

    var salt string
    var enabled bool
    db.SelectRow(user, []string{"salt", "enabled"}).Scan(&salt, &enabled)

    user.GetFields().(*models.User).Enabled = enabled

    user.LoadModelData(map[string]interface{}{"pass": utils.GetMD5Hash(pass1 + salt)})
    db.QueryUpdate_(user).Scan()

    utils.SendJSReply(map[string]interface{}{"result": "ok"}, this.Response)
}

func (this *GridHandler) GetEventTypesByEventId() {
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
    } else {
        event_id, err := strconv.Atoi(request["event_id"].(string))
        if utils.HandleErr("[GridHandler::GetEventTypesByEventId] event_id Atoi: ", err, this.Response) {
            return
        }

        query := `SELECT event_types.id, event_types.name FROM events_types
            INNER JOIN events ON events.id = events_types.event_id
            INNER JOIN event_types ON event_types.id = events_types.type_id
            WHERE events.id = $1 ORDER BY event_types.id;`
        result := db.Query(query, []interface{}{event_id})

        utils.SendJSReply(map[string]interface{}{"result": "ok", "data": result}, this.Response)
    }
}

func (this *GridHandler) ImportForms() {
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
    if utils.HandleErr("[GridHandler::ImportForms] event_id Atoi: ", err, this.Response) {
        return
    }

    for _, v := range request["event_types_ids"].([]interface{}) {
        type_id, err := strconv.Atoi(v.(string))
        if err != nil {
            utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
            return
        }
        query := `SELECT events.id FROM events
            INNER JOIN events_types ON events_types.event_id = events.id
            INNER JOIN event_types ON event_types.id = events_types.type_id
            WHERE event_types.id = $1 AND events.id <> $2
            ORDER BY id DESC LIMIT 1;`

        eventResult := db.Query(query, []interface{}{type_id, event_id})

        query = `SELECT forms.id FROM forms
            INNER JOIN events_forms ON events_forms.form_id = forms.id
            INNER JOIN events ON events.id = events_forms.event_id
            WHERE events.id = $1 ORDER BY forms.id;`

        formsResult := db.Query(query, []interface{}{int(eventResult[0].(map[string]interface{})["id"].(int64))})

        for i := 0; i < len(formsResult); i++ {
            form_id := int(formsResult[i].(map[string]interface{})["id"].(int64))
            eventsForms := GetModel("events_forms")
            eventsForms.LoadWherePart(map[string]interface{}{"event_id":  event_id, "form_id": form_id})
            var p int
            err := db.SelectRow(eventsForms, []string{"id"}).Scan(&p)
            if err != sql.ErrNoRows {
                continue
            }
            eventsForms.LoadModelData(map[string]interface{}{"event_id":  event_id, "form_id": form_id})
            db.QueryInsert_(eventsForms, "").Scan()
        }
    }

    utils.SendJSReply(map[string]interface{}{"result": "ok"}, this.Response)
}

func (this *GridHandler) GetPersonsByEventId() {
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
    } else {
        event_id, err := strconv.Atoi(request["event_id"].(string))
        if utils.HandleErr("[GridHandler::GetPersonsByEventId] event_id Atoi: ", err, this.Response) {
            return
        }

        params := request["params_ids"].([]interface{})

        if len(params) == 0 {
            utils.SendJSReply(map[string]interface{}{"result": "Выберите параметры."}, this.Response)
            return
        }

        q := "SELECT params.name FROM params WHERE params.id in ("+strings.Join(db.MakeParams(len(params)), ", ")+") ORDER BY id;"

        var caption []string
        for _, v := range db.Query(q, params) {
            caption = append(caption, v.(map[string]interface{})["name"].(string))
        }

        result := []interface{}{0: map[string]interface{}{"id": -1, "name": strings.Join(caption, " ")}}

        query := `SELECT reg_param_vals.reg_id as id, array_to_string(array_agg(param_values.value), ' ') as name
            FROM reg_param_vals
            INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
            INNER JOIN events ON events.id = registrations.event_id
            INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
            INNER JOIN params ON params.id = param_values.param_id
            WHERE params.id in (` + strings.Join(db.MakeParams(len(params)), ", ")
        query += ") AND events.id = $" + strconv.Itoa(len(params)+1) + " GROUP BY reg_param_vals.reg_id ORDER BY reg_param_vals.reg_id;"

        data := db.Query(query, append(params, event_id))
        utils.SendJSReply(map[string]interface{}{"result": "ok", "data": append(result, data...)}, this.Response)
    }
}

func (this *GridHandler) GetParamsByEventId() {
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
    } else {
        event_id, err := strconv.Atoi(request["event_id"].(string))
        if utils.HandleErr("[GridHandler::GetParamsByEventId] event_id Atoi: ", err, this.Response) {
            return
        }

        query := `SELECT DISTINCT params.id, params.name
            FROM reg_param_vals
            INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
            INNER JOIN params ON params.id = param_values.param_id
            INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
            INNER JOIN events ON events.id = registrations.event_id
            WHERE events.id = $1 ORDER BY params.id;`

        result := db.Query(query, []interface{}{event_id})

        utils.SendJSReply(map[string]interface{}{"result": "ok", "data": result}, this.Response)
    }
}

func (this *GridHandler) GetRequest() {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    } else {
        person_id, err := strconv.Atoi(request["person_id"].(string))
        if utils.HandleErr("[GridHandler::GetRequest] person_id Atoi: ", err, this.Response) {
            return
        }

        group_reg_id, err := strconv.Atoi(request["group_reg_id"].(string))
        if utils.HandleErr("[GridHandler::GetRequest] group_reg_id Atoi: ", err, this.Response) {
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
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN persons ON persons.face_id = faces.id
            INNER JOIN groups ON groups.id = persons.group_id
            INNER JOIN group_registrations ON group_registrations.group_id = groups.id
                AND group_registrations.event_id = events.id

            WHERE group_registrations.id = $1 AND persons.id = $2 ORDER BY forms.id;`

        res := db.Query(query, []interface{}{group_reg_id, person_id})

        utils.SendJSReply(map[string]interface{}{"data": res}, this.Response)
    }
}

func (this *GridHandler) GetPersonRequest() {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    } else {
        reg_id, err := strconv.Atoi(request["reg_id"].(string))
        if utils.HandleErr("[GridHandler::GetPersonRequest] reg_id Atoi: ", err, this.Response) {
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
            WHERE registrations.id=$1 ORDER BY forms.id;`

        data := db.Query(query, []interface{}{reg_id})

        utils.SendJSReply(map[string]interface{}{"result": "ok", "data": data}, this.Response)
    }
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

    } else {
        event_id, err := strconv.Atoi(request["event_id"].(string))
        if utils.HandleErr("[GridHandler::ConfirmOrRejectPersonRequest] event_id Atoi: ", err, this.Response) {
            return
        }
        reg_id, err := strconv.Atoi(request["reg_id"].(string))
        if utils.HandleErr("[GridHandler::ConfirmOrRejectPersonRequest] reg_id Atoi: ", err, this.Response) {
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
            WHERE params.id in (1, 4) AND users.id in (
                SELECT users.id FROM registrations INNER JOIN events ON events.id = registrations.event_id
                INNER JOIN faces ON faces.id = registrations.face_id
                INNER JOIN users ON users.id = faces.user_id
                WHERE registrations.id = $1
            ) ORDER BY params.id;`

        data := db.Query(query, []interface{}{reg_id})

        if len(data) < 2 {
            utils.SendJSReply(map[string]interface{}{"result": "Нет данных о логине или e-mail пользователя."}, this.Response)
            return
        }

        to := data[0].(map[string]interface{})["value"].(string)
        email := data[1].(map[string]interface{})["value"].(string)
        event := db.Query("SELECT name FROM events WHERE id=$1;", []interface{}{event_id})[0].(map[string]interface{})["name"].(string)

        if request["confirm"].(bool) {
            if event_id == 1 {
                utils.SendJSReply(map[string]interface{}{"result": "Эту заявку нельзя подтвердить письмом."}, this.Response)

            } else {
                if mailer.SendEmailToConfirmRejectPersonRequest(to, email, event, true) {
                    utils.SendJSReply(map[string]interface{}{"result": "Письмо с подтверждением заявки отправлено."}, this.Response)
                } else {
                    utils.SendJSReply(map[string]interface{}{"result": "Ошибка. Письмо с подтверждением заявки не отправлено."}, this.Response)
                }
            }

        } else {
            if event_id == 1 {
                utils.SendJSReply(map[string]interface{}{"result": "Эту заявку нельзя отклонить письмом."}, this.Response)

            } else {
                query := `DELETE
                    FROM param_values USING reg_param_vals
                    WHERE param_values.id in (SELECT reg_param_vals.param_val_id WHERE reg_param_vals.reg_id = $1);`
                db.Query(query, []interface{}{reg_id})

                query = `DELETE FROM registrations WHERE id = $1;`
                db.Query(query, []interface{}{reg_id})

                if mailer.SendEmailToConfirmRejectPersonRequest(to, email, event, false) {
                    utils.SendJSReply(map[string]interface{}{"result": "Письмо с отклонением заявки отправлено."}, this.Response)
                } else {
                    utils.SendJSReply(map[string]interface{}{"result": "Ошибка. Письмо с отклонением заявки не отправлено."}, this.Response)
                }
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

        param_value := GetModel("param_values")
        param_value.LoadModelData(map[string]interface{}{"value": value})
        param_value.LoadWherePart(map[string]interface{}{"id": param_val_id})
        db.QueryUpdate_(param_value).Scan()
    }

    utils.SendJSReply(map[string]interface{}{"result": "Изменения сохранены."}, this.Response)
}

func (this *GridHandler) RegGroup() {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    } else {
        group_id, err := strconv.Atoi(request["group_id"].(string))
        if utils.HandleErr("[GridHandler::RegGroup] group_id Atoi: ", err, this.Response) {
            return
        }

        event_id, err := strconv.Atoi(request["event_id"].(string))
        if utils.HandleErr("[GridHandler::RegGroup] event_id Atoi: ", err, this.Response) {
            return
        }

        event := GetModel("events")
        event.LoadWherePart(map[string]interface{}{"id": event_id})

        var eventName string
        err = db.SelectRow(event, []string{"name"}).Scan(&eventName)
        if err != nil {
            utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
            return
        }

        query := `SELECT groups.face_id, groups.name FROM groups
            INNER JOIN faces ON faces.id = groups.face_id
            INNER JOIN users ON users.id = faces.user_id
            WHERE users.id = $1 AND groups.id = $2;`

        var face_id int
        var groupName string
        err = db.QueryRow(query, []interface{}{user_id, group_id}).Scan(&face_id, &groupName)

        if err != nil {
            utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
            return
        }

        group_reg := GetModel("group_registrations")
        group_reg.LoadModelData(map[string]interface{}{"event_id": event_id, "group_id": group_id})
        db.QueryInsert_(group_reg, "").Scan()

        query = `SELECT persons.name, persons.email, persons.face_id, persons.status FROM persons
            INNER JOIN groups ON groups.id = persons.group_id
            INNER JOIN faces ON faces.id = persons.face_id
            WHERE groups.id = $1;`
        data := db.Query(query, []interface{}{group_id})

        query = `SELECT params.id FROM events_forms
            INNER JOIN events ON events.id = events_forms.event_id
            INNER JOIN forms ON forms.id = events_forms.form_id
            INNER JOIN params ON forms.id = params.form_id
            WHERE events.id = $1 ORDER BY forms.id;`
        params := db.Query(query, []interface{}{event_id})

        for _, v := range data {
            face_id := int(v.(map[string]interface{})["face_id"].(int64))
            status := v.(map[string]interface{})["status"].(bool)

            if !status {
                continue
            }

            var reg_id int
            regs := GetModel("registrations")
            regs.LoadModelData(map[string]interface{}{"face_id": face_id, "event_id": event_id})
            db.QueryInsert_(regs, "RETURNING id").Scan(&reg_id)

            to := v.(map[string]interface{})["name"].(string)
            address := v.(map[string]interface{})["email"].(string)
            if !mailer.AttendAnEvent(to, address, eventName, groupName) {
                utils.SendJSReply(map[string]interface{}{"result": "Ошибка. Письмо с уведомлением не отправлено."}, this.Response)
            }

            for _, p := range params {
                param_id := int(p.(map[string]interface{})["id"].(int64))

                var param_val_id int
                paramValues := GetModel("param_values")
                paramValues.LoadModelData(map[string]interface{}{"param_id": param_id, "value": "  "})
                db.QueryInsert_(paramValues, "RETURNING id").Scan(&param_val_id)

                regParamValue := GetModel("reg_param_vals")
                regParamValue.LoadModelData(map[string]interface{}{
                    "reg_id":        reg_id,
                    "param_val_id":  param_val_id})
                db.QueryInsert_(regParamValue, "").Scan()
            }

        }

        utils.SendJSReply(map[string]interface{}{"result": "ok"}, this.Response)
    }
}
