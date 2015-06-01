package controllers

import (
    "github.com/orc/db"
    "github.com/orc/mvc/models"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "strconv"
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

    query := `SELECT param_id, params.name, param_types.name as type, param_values.value, forms.id as form_id FROM events
            INNER JOIN events_forms ON events_forms.event_id = events.id
            INNER JOIN forms ON events_forms.form_id = forms.id

            INNER JOIN registrations ON events.id = registrations.event_id
            INNER JOIN reg_param_vals ON reg_param_vals.reg_id = registrations.id

            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN users ON users.id = faces.user_id

            INNER JOIN params ON params.form_id = forms.id
            INNER JOIN param_types ON param_types.id = params.param_type_id
            INNER JOIN param_values ON param_values.param_id = params.id AND reg_param_vals.param_val_id = param_values.id
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
    if len(data["form_ids"].(map[string]interface{})["form_id"].([]interface{})) == 0 {
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
    var paramValIds []interface{}
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

        paramValIds, _, _, _ = this.InsertUserParams(data["data"].([]interface{}))

    } else if eventId == 1 {
        var userLogin, userPass, email string
        paramValIds, userLogin, userPass, email = this.InsertUserParams(data["data"].([]interface{}))

        result, regId = this.HandleRegister_(userLogin, userPass, email, "user")
        if result != "ok" && regId == -1 {
            utils.SendJSReply(map[string]interface{}{"result": result}, this.Response)
            return
        }

    } else {
        utils.SendJSReply(map[string]interface{}{"result": "Unauthorized"}, this.Response)
        return
    }

    for _, v := range paramValIds {
        regParamValue := this.GetModel("reg_param_vals")
        regParamValue.LoadModelData(map[string]interface{}{
            "reg_id":        regId,
            "param_val_id":  v.(map[string]int)["param_val_id"]})
        db.QueryInsert_(regParamValue, "").Scan()
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

    query := `SELECT forms.id as form_id, forms.name as form_name, params.id as param_id,
            params.name as param_name, param_types.name as type, events.name as event_name,
            events.id as event_id
        FROM events_forms
        INNER JOIN events ON events.id = events_forms.event_id
        INNER JOIN forms ON forms.id = events_forms.form_id
        INNER JOIN params ON forms.id = params.form_id
        INNER JOIN param_types ON param_types.id = params.param_type_id
        WHERE events.id = $1 ORDER BY forms.id, params.id;`
    res := db.Query(query, []interface{}{eventId})

    this.Render([]string{"mvc/views/item.html"}, "item", map[string]interface{}{"data": res})
}

func (this *Handler) InsertUserParams(data []interface{}) ([]interface{}, string, string, string) {
    paramValIds := make([]interface{}, 0)
    userLogin := ""
    userPass := ""
    email := ""

    for _, element := range data {
        paramId, err := strconv.Atoi(element.(map[string]interface{})["id"].(string))
        if err != nil {
            continue
        }

        value := element.(map[string]interface{})["value"].(string)

        if paramId == 1 {
            userLogin = value
            continue
        } else if paramId == 2 || paramId == 3 {
            userPass = value
            continue
        } else if paramId == 4 {
            email = value
        }

        var paramValId int
        paramValues := this.GetModel("param_values")
        paramValues.LoadModelData(map[string]interface{}{"param_id": paramId, "value": value})
        db.QueryInsert_(paramValues, "RETURNING id").Scan(&paramValId)
        paramValIds = append(paramValIds, map[string]int{"param_val_id": paramValId})
    }

    return paramValIds, userLogin, userPass, email
}
