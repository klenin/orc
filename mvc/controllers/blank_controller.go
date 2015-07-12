package controllers

import (
    "github.com/orc/db"
    "github.com/orc/mvc/models"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "log"
    "net/http"
    "strconv"
    "time"
)

func (c *BaseController) BlankController() *BlankController {
    return new(BlankController)
}

type BlankController struct {
    Controller
}

func (this *BlankController) GetPersonRequestFromGroup() {
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

    personalForm := request["personal"].(string)
    if personalForm != "true" && personalForm != "false" {
        utils.SendJSReply(map[string]interface{}{"result": "Invalid personalForm"}, this.Response)
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

    var regId int

    if faceId == -1 {
        if this.isAdmin() {
            query := `SELECT faces.id, registrations.id
                FROM regs_groupregs
                INNER JOIN group_registrations ON group_registrations.id = regs_groupregs.groupreg_id
                INNER JOIN registrations ON registrations.id = regs_groupregs.reg_id
                INNER JOIN groups ON groups.id = group_registrations.group_id
                INNER JOIN faces ON faces.id = groups.face_id AND faces.id = registrations.face_id
                INNER JOIN events ON events.id = registrations.event_id
                INNER JOIN events_forms ON events_forms.event_id = events.id
                INNER JOIN forms ON forms.id = events_forms.form_id
                INNER JOIN params ON forms.id = params.form_id
                INNER JOIN param_types ON param_types.id = params.param_type_id
                INNER JOIN param_values ON params.id = param_values.param_id AND param_values.reg_id = registrations.id
                WHERE group_registrations.id = $1 AND forms.personal = false;`
            db.QueryRow(query, []interface{}{groupRegId}).Scan(&faceId, &regId)
        } else {
            if err := this.GetModel("faces").
                LoadWherePart(map[string]interface{}{"user_id": userId}).
                SelectRow([]string{"id"}).
                Scan(&faceId);
                err != nil {
                utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
                return
            }

            query := `SELECT registrations.id
                FROM regs_groupregs
                INNER JOIN group_registrations ON group_registrations.id = regs_groupregs.groupreg_id
                INNER JOIN registrations ON registrations.id = regs_groupregs.reg_id
                INNER JOIN faces ON faces.id = registrations.face_id
                INNER JOIN groups ON groups.face_id = faces.id AND groups.id = group_registrations.group_id
                INNER JOIN events ON events.id = registrations.event_id
                INNER JOIN events_forms ON events_forms.event_id = events.id
                INNER JOIN forms ON forms.id = events_forms.form_id
                INNER JOIN params ON forms.id = params.form_id
                INNER JOIN param_types ON param_types.id = params.param_type_id
                INNER JOIN param_values ON params.id = param_values.param_id AND param_values.reg_id = registrations.id
                WHERE group_registrations.id = $1 AND faces.id = $2 AND forms.personal = false;`
            db.QueryRow(query, []interface{}{groupRegId, faceId}).Scan(&regId)
        }
    } else {
        query := `SELECT registrations.id
                FROM regs_groupregs
                INNER JOIN group_registrations ON group_registrations.id = regs_groupregs.groupreg_id
                INNER JOIN registrations ON registrations.id = regs_groupregs.reg_id
                INNER JOIN faces ON faces.id = registrations.face_id
                INNER JOIN persons ON persons.face_id = faces.id AND persons.group_id = group_registrations.group_id
                INNER JOIN events ON events.id = registrations.event_id
                INNER JOIN events_forms ON events_forms.event_id = events.id
                INNER JOIN forms ON forms.id = events_forms.form_id
                INNER JOIN params ON forms.id = params.form_id
                INNER JOIN param_types ON param_types.id = params.param_type_id
                INNER JOIN param_values ON params.id = param_values.param_id AND param_values.reg_id = registrations.id
                WHERE group_registrations.id = $1 AND faces.id = $2 AND forms.personal = true;`
        db.QueryRow(query, []interface{}{groupRegId, faceId}).Scan(&regId)
    }

    log.Println("faceId: ", faceId)
    log.Println("groupRegId: ", groupRegId)
    log.Println("regId: ", regId)
    log.Println("formType: ", personalForm)

    query := `SELECT forms.id as form_id, forms.name as form_name,
    var query string
    var result []interface{}
            params.id as param_id, params.name as param_name, params.required, params.editable,
            events.name as event_name, events.id as event_id,
            param_types.name as type, param_values.id as param_val_id, param_values.value
        FROM events_forms
        INNER JOIN events ON events.id = events_forms.event_id
        INNER JOIN forms ON forms.id = events_forms.form_id
        INNER JOIN params ON forms.id = params.form_id
        INNER JOIN param_types ON param_types.id = params.param_type_id
        INNER JOIN param_values ON params.id = param_values.param_id
        INNER JOIN registrations ON registrations.id = param_values.reg_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN group_registrations ON group_registrations.event_id = events.id
        INNER JOIN groups ON group_registrations.group_id = groups.id
        INNER JOIN regs_groupregs ON regs_groupregs.reg_id = registrations.id
            AND regs_groupregs.groupreg_id = group_registrations.id
        WHERE group_registrations.id = $1 AND faces.id = $2 AND forms.personal = $3 ORDER BY forms.id, params.id;`
        query = `SELECT forms.id as form_id, forms.name as form_name,
                params.id as param_id, params.name as param_name, params.required, params.editable,
                events.name as event_name, events.id as event_id,
                param_types.name as type, param_values.id as param_val_id, param_values.value
            FROM events_forms
            INNER JOIN events ON events.id = events_forms.event_id
            INNER JOIN forms ON forms.id = events_forms.form_id
            INNER JOIN params ON forms.id = params.form_id
            INNER JOIN param_types ON param_types.id = params.param_type_id
            INNER JOIN param_values ON params.id = param_values.param_id
            INNER JOIN registrations ON registrations.id = param_values.reg_id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN group_registrations ON group_registrations.event_id = events.id
            INNER JOIN groups ON group_registrations.group_id = groups.id
            INNER JOIN regs_groupregs ON regs_groupregs.reg_id = registrations.id
                AND regs_groupregs.groupreg_id = group_registrations.id
            WHERE group_registrations.id = $1 AND faces.id = $2 AND forms.personal = $3 ORDER BY forms.id, params.id;`
        result = db.Query(query, []interface{}{groupRegId, faceId, personalForm})

    utils.SendJSReply(
        map[string]interface{}{
            "result": "ok",
            "data": result,
            "role": this.isAdmin(),
            "regId": regId},
        this.Response)
}

func (this *BlankController) GetBlankByRegId() {
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
        INNER JOIN registrations ON registrations.id = param_values.reg_id
            AND events.id = registrations.event_id
        WHERE registrations.id = $1 AND forms.personal = true ORDER BY forms.id, params.id;`

    utils.SendJSReply(
        map[string]interface{}{
            "result": "ok",
            "data": db.Query(query, []interface{}{regId}),
            "role": this.isAdmin()},
        this.Response)
}

func (this *BlankController) GetBlankByEventId(id string) {
    eventId, err := strconv.Atoi(id)
    if utils.HandleErr("[BlankController::GetBlankByEventId] event_id Atoi: ", err, this.Response) {
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
        WHERE events.id = $1 AND forms.personal = true ORDER BY forms.id, params.id;`
    res := db.Query(query, []interface{}{eventId})

    this.Render([]string{"mvc/views/item.html"}, "item", map[string]interface{}{"data": res})
}

//-----------------------------------------------------------------------------
func (this *BlankController) GetEditHistoryData() {
    data, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    regId, err := strconv.Atoi(data["reg_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    formType := data["personal"].(string)
    if formType != "true" && formType != "false" {
        utils.SendJSReply(map[string]interface{}{"result": "Invalid form type"}, this.Response)
        return
    }

    query := `SELECT params.id as param_id, forms.id as form_id, p.date as edit_date,
        array_to_string(ARRAY(
            SELECT param_values.value
                FROM events
                INNER JOIN events_forms ON events_forms.event_id = events.id
                INNER JOIN forms ON events_forms.form_id = forms.id
                INNER JOIN registrations ON events.id = registrations.event_id
                INNER JOIN faces ON faces.id = registrations.face_id
                INNER JOIN users ON users.id = faces.user_id
                INNER JOIN params ON params.form_id = forms.id
                INNER JOIN param_types ON param_types.id = params.param_type_id
                INNER JOIN param_values ON param_values.param_id = params.id
                    AND registrations.id = param_values.reg_id
                WHERE (params.id in (5, 6, 7) AND events.id = 1) and users.id = p.user_id
        ), ' ') as login
        FROM events
        INNER JOIN events_forms ON events_forms.event_id = events.id
        INNER JOIN forms ON events_forms.form_id = forms.id
        INNER JOIN registrations ON events.id = registrations.event_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN users ON users.id = faces.user_id
        INNER JOIN params ON params.form_id = forms.id
        INNER JOIN param_types ON param_types.id = params.param_type_id
        INNER JOIN param_values as p ON p.param_id = params.id
            AND p.reg_id = registrations.id
        WHERE registrations.id = $1 AND forms.personal = $2;`

    utils.SendJSReply(map[string]interface{}{"result": "ok", "data": db.Query(query, []interface{}{regId, formType})}, this.Response)
}

func (this *BlankController) GetHistoryRequest() {
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
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN users ON users.id = faces.user_id
        INNER JOIN params ON params.form_id = forms.id
        INNER JOIN param_types ON param_types.id = params.param_type_id
        INNER JOIN param_values ON param_values.param_id = params.id
            AND param_values.reg_id = registrations.id
        WHERE users.id = $1 AND events.id = $2 AND forms.personal = true;`

    utils.SendJSReply(map[string]interface{}{"result": "ok", "data": db.Query(query, []interface{}{userId, eventId})}, this.Response)
}

func (this *BlankController) GetListHistoryEvents() {
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

    ids := map[string]interface{}{"form_id": make([]interface{}, 0)}
    if data["form_ids"] == nil || len(data["form_ids"].([]interface{})) == 0 {
        utils.SendJSReply(map[string]interface{}{"result": "Нет данных о формах анкеты"}, this.Response)
        return
    }

    for _, v := range data["form_ids"].([]interface{}) {
        ids["form_id"] = append(ids["form_id"].([]interface{}), int(v.(float64)))
    }

    eventsForms := this.GetModel("events_forms")
    eventsForms.LoadWherePart(ids).SetCondition(models.OR)
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

//-----------------------------------------------------------------------------
func (this *BlankController) EditParams() {
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

        params := map[string]interface{}{"value": value, "date": date, "user_id": userId}
        where := map[string]interface{}{"id": paramValId}
        this.GetModel("param_values").Update(userId, params, where)
    }

    utils.SendJSReply(map[string]interface{}{"result": "Изменения сохранены"}, this.Response)
}
