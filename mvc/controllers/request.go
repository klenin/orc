package controllers

import (
    "encoding/json"
    "github.com/orc/db"
    "github.com/orc/mvc/models"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "html/template"
    "log"
    "strconv"
)

func (this *Handler) GetHistoryRequest() {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        utils.SendJSReply(map[string]interface{}{"result": "notAuthorized"}, this.Response)
        return
    }

    result := make(map[string]interface{}, 2)
    result["result"] = "ok"

    data, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    event_id := data["event_id"]

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

    result["data"] = db.Query(query, []interface{}{user_id, event_id})
    utils.SendJSReply(result, this.Response)
}

func (this *Handler) GetListHistoryEvents() {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        utils.SendJSReply(map[string]interface{}{"result": "notAuthorized"}, this.Response)
        return
    }

    result := make(map[string]interface{}, 2)
    result["result"] = "ok"

    data, err := utils.ParseJS(this.Request, this.Response)
    if  err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    ids := make(map[string]interface{}, 1)
    ids["form_id"] = make([]interface{}, 0)
    if len(data["form_ids"].(map[string]interface{})["form_id"].([]interface{})) == 0 {
        result["result"] = "no"
    } else {
        for _, v := range data["form_ids"].(map[string]interface{})["form_id"].([]interface{}) {
            ids["form_id"] = append(ids["form_id"].([]interface{}), int(v.(float64)))
            log.Println("FORM_ID: ", strconv.Itoa(int(v.(float64))))
        }

        eventsForms := GetModel("events_forms")
        eventsForms.LoadWherePart(ids)
        eventsForms.SetCondition(models.OR)
        events := db.Select(eventsForms, []string{"event_id"})

        if len(events) != 0 {
            query := `SELECT DISTINCT events.id, events.name FROM events
                INNER JOIN events_forms ON events_forms.event_id = events.id
                INNER JOIN forms ON events_forms.form_id = forms.id
                INNER JOIN registrations ON registrations.event_id = events.id
                INNER JOIN faces ON faces.id = registrations.face_id
                INNER JOIN users ON users.id = faces.user_id
                WHERE users.id=$1 AND events.id IN (`

            var i int
            var params []interface{}
            params = append(params, user_id)

            for i = 2; i < len(events); i++ {
                query += "$" + strconv.Itoa(i) + ", "
                params = append(params, int(events[i-2].(map[string]interface{})["event_id"].(int64)))
                log.Println("EVENT_ID: ", strconv.Itoa(int(events[i-2].(map[string]interface{})["event_id"].(int64))))
            }

            query += "$" + strconv.Itoa(i) + ")"
            params = append(params, int(events[i-2].(map[string]interface{})["event_id"].(int64)))
            result["data"] = db.Query(query, params)
        }
    }

    utils.SendJSReply(result, this.Response)
}

func (this *Handler) SaveUserRequest() {
    var param_val_ids []interface{}
    var result string
    var reg_id int

    data, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    event_id := int(data["event_id"].(float64))

    if event_id == 1 && sessions.CheackSession(this.Response, this.Request) {
        utils.SendJSReply(map[string]interface{}{"result": "authorized"}, this.Response)
        return
    }

    if sessions.CheackSession(this.Response, this.Request) {
        user_id := sessions.GetValue("id", this.Request)
        if user_id == nil {
            utils.SendJSReply(map[string]interface{}{"result": "notAuthorized"}, this.Response)
            return
        }

        var face_id int
        face := GetModel("faces")
        face.LoadModelData(map[string]interface{}{"user_id": user_id})
        db.QueryInsert_(face, "RETURNING id").Scan(&face_id)

        registration := GetModel("registrations")
        registration.LoadModelData(map[string]interface{}{"face_id": face_id, "event_id": event_id})
        db.QueryInsert_(registration, "RETURNING id").Scan(&reg_id)

        param_val_ids, _, _, _ = InsertUserParams(data["data"].([]interface{}))

    } else if event_id == 1 {
        var userLogin, userPass, email string
        param_val_ids, userLogin, userPass, email = InsertUserParams(data["data"].([]interface{}))

        result, reg_id = this.HandleRegister_(userLogin, userPass, email, "user")
        if result != "ok" {
            utils.SendJSReply(map[string]interface{}{"result": result}, this.Response)
            return
        }

    } else {
        utils.SendJSReply(map[string]interface{}{"result": "notAuthorized"}, this.Response)
        return
    }

    for _, v := range param_val_ids {
        regParamValue := GetModel("reg_param_vals")
        regParamValue.LoadModelData(map[string]interface{}{
            "reg_id":        reg_id,
            "param_val_id":  v.(map[string]int)["param_val_id"]})
        db.QueryInsert_(regParamValue, "").Scan()
    }

    utils.SendJSReply(map[string]interface{}{"result": "ok"}, this.Response)
}

func (this *Handler) GetRequest(tableName, id string) {
    if !sessions.CheackSession(this.Response, this.Request) && id != "1" {
        this.Render([]string{"mvc/views/loginpage.html", "mvc/views/login.html"}, "loginpage", nil)
        return
    }

    // проверка id - число !!!

    response, err := json.Marshal(MegoJoin(tableName, id))
    if utils.HandleErr("[Handler::GetRequest] Marshal: ", err, this.Response) {
        return
    }

    this.Render([]string{"mvc/views/item.html"}, "item", template.JS(response))
}

func MegoJoin(tableName, id string) RequestModel {
    var E []interface{}
    var F []interface{}
    var P []interface{}

    event := GetModel("events")
    event.LoadWherePart(map[string]interface{}{"id": id})
    E = db.Select(event, []string{"id", "name"})

    query := db.InnerJoin(
        []string{"f.id", "f.name"},
        "events_forms",
        "e_f",
        []string{"event_id", "form_id"},
        []string{"events", "forms"},
        []string{"e", "f"},
        []string{"id", "id"},
        "where e.id=$1 ORDER BY id")
    F = db.Query(query, []interface{}{id})

    for j := 0; j < len(F); j++ {
        f_id := F[j].(map[string]interface{})["id"]
        query := db.InnerJoin(
            []string{"p.id", "p.name", "p_t.name as type"},
            "params",
            "p",
            []string{"form_id", "param_type_id"},
            []string{"forms", "param_types"},
            []string{"f", "p_t"},
            []string{"id", "id"},
            "where f.id=$1 ORDER BY id")
        P = append(P, db.Query(query, []interface{}{f_id}))
    }

    return RequestModel{E: E, F: F, P: P}
}

func InsertUserParams(data []interface{}) ([]interface{}, string, string, string) {
    param_val_ids := make([]interface{}, 0)
    userLogin := ""
    userPass := ""
    email := ""

    for _, element := range data {
        param_id, err := strconv.Atoi(element.(map[string]interface{})["id"].(string))
        if err != nil {
            continue
        }

        value := element.(map[string]interface{})["value"].(string)

        if param_id == 1 {
            userLogin = value
        } else if param_id == 2 || param_id == 3 {
            userPass = value
            continue
        } else if param_id == 4 {
            email = value
        }

        var param_val_id int
        paramValues := GetModel("param_values")
        paramValues.LoadModelData(map[string]interface{}{"param_id": param_id, "value": value})
        db.QueryInsert_(paramValues, "RETURNING id").Scan(&param_val_id)
        param_val_ids = append(param_val_ids, map[string]int{"param_val_id": param_val_id})
    }

    return param_val_ids, userLogin, userPass, email
}
