package controllers

import (
    "github.com/klenin/orc/db"
    "github.com/klenin/orc/mvc/models"
    "github.com/klenin/orc/sessions"
    "github.com/klenin/orc/utils"
    "log"
    "net/http"
    "strconv"
    "time"
)

func (*BaseController) BlankController() *BlankController {
    return new(BlankController)
}

type BlankController struct {
    Controller
}

func (this *BlankController) GetPersonBlankFromGroup() {
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

    var personalForm bool
    switch request["personal"].(string) {
    case "true":
        personalForm = true
        break
    case "false":
        personalForm = false
        break
    default:
        panic("Invalid bool value")
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
            db.QueryRow(db.QueryGetCaptFaceIdAndRegId, []interface{}{groupRegId}).Scan(&faceId, &regId)
        } else {
            if err := this.GetModel("faces").
                LoadWherePart(map[string]interface{}{"user_id": userId}).
                SelectRow([]string{"id"}).
                Scan(&faceId);
                err != nil {
                utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)

                return
            }
            db.QueryRow(db.QueryGetCaptRegIdByGroupRegIdAndFaceId, []interface{}{groupRegId, faceId}).Scan(&regId)
        }
    } else {
        db.QueryRow(db.QueryGetRegIdByGroupRegIdAndFaceId, []interface{}{groupRegId, faceId}).Scan(&regId)
    }

    log.Println("faceId: ", faceId, ", groupRegId: ", groupRegId, ", regId: ", regId, ", formType: ", personalForm)

    blank := new(models.BlankManager).NewGroupBlank(personalForm)
    blank.SetGroupRegId(groupRegId).SetFaceId(faceId)

    utils.SendJSReply(
        map[string]interface{}{
            "result": "ok",
            "data": blank.GetBlank(),
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

    blank := new(models.BlankManager).NewPersonalBlank(true).SetRegId(regId)
    result := blank.GetBlank()

    if len(result) == 0 {
        result = blank.SetPersonal(false).GetBlank()
    }

    utils.SendJSReply(
        map[string]interface{}{
            "result": "ok",
            "data": result,
            "role": this.isAdmin()},
        this.Response)
}

func (this *BlankController) GetGroupBlank() {
    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    groupRegId, err := strconv.Atoi(request["group_reg_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    utils.SendJSReply(
        map[string]interface{}{
            "result": "ok",
            "data": new(models.BlankManager).NewGroupBlank(false).SetGroupRegId(groupRegId).GetTeamBlank()},
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

    this.Render(
        []string{"mvc/views/item.html"},
        "item",
        map[string]interface{}{"data": new(models.BlankManager).NewPersonalBlank(true).GetEmptyBlank(eventId)})
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
    events := eventsForms.
        LoadWherePart(ids).
        SetCondition(models.OR).
        Select_([]string{"event_id"})

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
        this.GetModel("param_values").Update(this.isAdmin(), userId, params, where)
    }

    utils.SendJSReply(map[string]interface{}{"result": "Изменения сохранены"}, this.Response)
}
