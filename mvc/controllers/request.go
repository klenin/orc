package controllers

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/orc/db"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "html/template"
    "net/http"
    "strconv"
)

func (this *Handler) GetHistoryRequest() {
    if flag := sessions.CheackSession(this.Response, this.Request); !flag {
        return
    }

    this.Response.Header().Set("Access-Control-Allow-Origin", "*")
    this.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    this.Response.Header().Set("Content-type", "application/json")

    result := make(map[string]interface{}, 2)
    result["result"] = "ok"

    var data map[string]string
    decoder := json.NewDecoder(this.Request.Body)
    err := decoder.Decode(&data)
    if utils.HandleErr("[Handler::GetHistoryRequest] Decode: ", err, this.Response) {
        return
    }

    user_id := sessions.GetValue("id", this.Request)
    event_id := data["event_id"]

    query := `select param_id, params.name, event_types.id as event_type_id, param_types.name as type, param_values.value, forms.id from events
            inner join events_types on events_types.event_id = events.id
            inner join event_types on events_types.type_id = event_types.id
            inner join events_regs on events_regs.event_id = events.id
            inner join registrations on registrations.id = events_regs.reg_id
            inner join reg_param_vals on reg_param_vals.reg_id = registrations.id
                                     and reg_param_vals.event_id = events.id
                                     and reg_param_vals.event_type_id = event_types.id
            inner join faces on faces.id = registrations.face_id
            inner join users on users.id = faces.user_id
            inner join forms_types on forms_types.type_id = event_types.id
            inner join forms on forms.id = forms_types.form_id
            inner join params on params.form_id = forms.id
            inner join param_types on param_types.id = params.param_type_id
            inner join param_values on param_values.param_id = params.id and reg_param_vals.param_val_id = param_values.id
            where users.id = $1 and events.id = $2;`

    result["data"] = db.Query(query, []interface{}{user_id, event_id})
    response, err := json.Marshal(result)
    if utils.HandleErr("[Handle::GetHistoryRequest] Marshal: ", err, this.Response) {
        return
    }

    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *Handler) GetListHistoryEvents() {
    if !sessions.CheackSession(this.Response, this.Request) {
        response, err := json.Marshal(map[string]interface{}{"result": "no"})
        if utils.HandleErr("[Handle::GetListHistoryEvents] Marshal: ", err, this.Response) {
            return
        }

        fmt.Fprintf(this.Response, "%s", string(response))
        return
    }

    result := make(map[string]interface{}, 2)
    result["result"] = "ok"

    var data map[string]interface{}
    decoder := json.NewDecoder(this.Request.Body)
    err := decoder.Decode(&data)
    if utils.HandleErr("[Handler::GetListHistoryEvents] Decode: ", err, this.Response) {
        return
    }

    user_id := sessions.GetValue("id", this.Request)
    if user_id == nil {
        http.Redirect(this.Response, this.Request, "/", 401)
        return
    }

    ids := make(map[string]interface{}, 1)
    ids["form_id"] = make([]interface{}, 0)
    if data["form_ids"].(map[string]interface{})["form_id"] == nil {
        result["result"] = "no"
    } else {
        for _, v := range data["form_ids"].(map[string]interface{})["form_id"].([]interface{}) {
            ids["form_id"] = append(ids["form_id"].([]interface{}), int(v.(float64)))
        }

        formsTypes := GetModel("forms_types")
        formsTypes.LoadWherePart(ids)
        types := db.Select(formsTypes, []string{"type_id"}, "OR")

        if len(types) != 0 {

            query := `SELECT DISTINCT events.id, events.name FROM events
                inner join events_types on events_types.event_id = events.id
                inner join event_types on events_types.type_id = event_types.id
                inner join events_regs on events_regs.event_id = events.id
                inner join registrations on registrations.id = events_regs.reg_id
                inner join faces on faces.id = registrations.face_id
                inner join users on users.id = faces.user_id
                WHERE users.id=$1 AND events.id IN (SELECT DISTINCT event_id FROM events_types WHERE `

            var i int
            var params []interface{}

            params = append(params, user_id)

            for i = 2; i < len(types); i++ {
                query += "type_id=$" + strconv.Itoa(i) + " OR "
                params = append(params, int(types[i-2].(map[string]interface{})["type_id"].(int64)))
            }

            query += "type_id=$" + strconv.Itoa(i) + ")"
            params = append(params, int(types[i-2].(map[string]interface{})["type_id"].(int64)))
            result["data"] = db.Query(query, params)
        }
    }

    response, err := json.Marshal(result)
    if utils.HandleErr("[Handle::GetListHistoryEvents] Marshal: ", err, this.Response) {
        return
    }

    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *Handler) SaveUserRequest() {
    var param_val_ids []interface{}
    var result string
    var reg_id int

    var data map[string]interface{}
    decoder := json.NewDecoder(this.Request.Body)
    err := decoder.Decode(&data)
    if utils.HandleErr("[Handler::SaveUserRequest] Decode :", err, this.Response) {
        return
    }

    event_id := int(data["event_id"].(float64))

    if event_id == 1 && sessions.CheackSession(this.Response, this.Request) {
        return
    }

    if sessions.CheackSession(this.Response, this.Request) {
        user_id := sessions.GetValue("id", this.Request)
        if user_id == nil {
            http.Redirect(this.Response, this.Request, "/", 401)
            return
        }

        var face_id int
        face := GetModel("faces")
        face.LoadWherePart(map[string]interface{}{"user_id": user_id})
        err = db.SelectRow(face, []string{"id"}, "").Scan(&face_id)
        if err != nil {
            response, err := json.Marshal(map[string]interface{}{"result": err.Error()})
            if utils.HandleErr("[Handle::SaveUserRequest] Marshal: ", err, this.Response) {
                return
            }
            fmt.Fprintf(this.Response, "%s", string(response))
            return
        }

        reg := GetModel("registrations")
        reg.LoadWherePart(map[string]interface{}{"face_id": face_id})
        err = db.SelectRow(reg, []string{"id"}, "").Scan(&reg_id)
        if err != nil {
            response, err := json.Marshal(map[string]interface{}{"result": err.Error()})
            if utils.HandleErr("[Handle::SaveUserRequest] Marshal: ", err, this.Response) {
                return
            }
            fmt.Fprintf(this.Response, "%s", string(response))
            return
        }

        var event_reg_id int
        eventsRegs := GetModel("events_regs")
        eventsRegs.LoadModelData(map[string]interface{}{"reg_id": reg_id, "event_id": event_id})
        err := db.SelectRow(eventsRegs, []string{"id"}, "AND").Scan(&event_reg_id)

        if err != sql.ErrNoRows {
            response, err := json.Marshal(map[string]interface{}{"result": "Вы уже заполняли эту анкету."})
            if utils.HandleErr("[Handle::SaveUserRequest] Marshal: ", err, this.Response) {
                return
            }
            fmt.Fprintf(this.Response, "%s", string(response))
            return
        } else {
            db.QueryInsert_(eventsRegs, "")
        }

        param_val_ids, _, _ = InsertUserParams(data["data"].([]interface{}))

    } else if event_id == 1 {
        var userLogin, userPass string
        param_val_ids, userLogin, userPass = InsertUserParams(data["data"].([]interface{}))

        result, reg_id = this.HandleRegister_(userLogin, userPass, "user")
        if result != "ok" {
            response, err := json.Marshal(map[string]interface{}{"result": result})
            if utils.HandleErr("[Handle::SaveUserRequest] Marshal: ", err, this.Response) {
                return
            }
            fmt.Fprintf(this.Response, "%s", string(response))
            return
        }
        eventsRegs := GetModel("events_regs")
        eventsRegs.LoadModelData(map[string]interface{}{"reg_id": reg_id, "event_id": event_id})
        db.QueryInsert_(eventsRegs, "")
    }

    for _, v := range param_val_ids {
        regParamValue := GetModel("reg_param_vals")
        regParamValue.LoadModelData(map[string]interface{}{
            "reg_id":        reg_id,
            "event_id":      event_id,
            "event_type_id": v.(map[string]int)["event_type_id"],
            "param_val_id":  v.(map[string]int)["param_val_id"]})
        db.QueryInsert_(regParamValue, "")
    }

    response, err := json.Marshal(map[string]interface{}{"result": "ok"})
    if utils.HandleErr("[Handle::SaveUserRequest] Marshal: ", err, this.Response) {
        return
    }
    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *Handler) GetRequest(tableName, id string) {
    tmp, err := template.ParseFiles(
        "mvc/views/item.html",
        "mvc/views/header.html",
        "mvc/views/footer.html")
    if utils.HandleErr("[Handler::GetRequest] ParseFiles: ", err, this.Response) {
        return
    }

    reaponse, err := json.Marshal(MegoJoin(tableName, id))
    if utils.HandleErr("[Handler::GetRequest] Marshal: ", err, this.Response) {
        return
    }

    err = tmp.ExecuteTemplate(this.Response, "item", template.JS(reaponse))
    utils.HandleErr("[Handler::GetRequest] ExecuteTemplate: ", err, this.Response)
}

func MegoJoin(tableName, id string) RequestModel {
    var E []interface{}
    var T []interface{}
    var F []interface{}
    var P []interface{}

    event := GetModel("events")
    event.LoadWherePart(map[string]interface{}{"id": id})
    E = db.Select(event, []string{"id", "name"}, "")

    query := db.InnerJoin(
        []string{"t.id", "t.name"},
        "events_types",
        "e_t",
        []string{"event_id", "type_id"},
        []string{"events", "event_types"},
        []string{"e", "t"},
        []string{"id", "id"},
        "where e.id=$1")
    T = db.Query(query, []interface{}{id})

    for i := 0; i < len(T); i++ {
        id := T[i].(map[string]interface{})["id"]
        query := db.InnerJoin(
            []string{"f.id", "f.name"},
            "forms_types",
            "f_t",
            []string{"form_id", "type_id"},
            []string{"forms", "event_types"},
            []string{"f", "t"},
            []string{"id", "id"},
            "where t.id=$1")
        F = append(F, db.Query(query, []interface{}{id}))
    }

    for i := 0; i < len(F); i++ {
        var PP []interface{}
        for j := 0; j < len(F[i].([]interface{})); j++ {
            item := F[i].([]interface{})[j]
            id := item.(map[string]interface{})["id"]
            query := db.InnerJoin(
                []string{"p.id", "p.name", "p_t.name as type"},
                "params",
                "p",
                []string{"form_id", "param_type_id"},
                []string{"forms", "param_types"},
                []string{"f", "p_t"},
                []string{"id", "id"},
                "where f.id=$1")
            PP = append(PP, db.Query(query, []interface{}{id}))
        }
        P = append(P, PP)
    }

    return RequestModel{E: E, T: T, F: F, P: P}
}

func InsertUserParams(data []interface{}) ([]interface{}, string, string) {
    param_val_ids := make([]interface{}, 0)
    userLogin := ""
    userPass := ""

    for _, element := range data {
        param_id, err := strconv.Atoi(element.(map[string]interface{})["id"].(string))
        if err != nil {

        }
        event_type_id, err := strconv.Atoi(element.(map[string]interface{})["event_type_id"].(string))
        if err != nil {

        }
        value := element.(map[string]interface{})["value"].(string)

        if param_id == 2 {
            userLogin = value
            continue
        } else if param_id == 3 || param_id == 4 {
            userPass = value
            continue
        }

        var param_val_id int
        paramValues := GetModel("param_values")
        paramValues.LoadModelData(map[string]interface{}{"param_id": param_id, "value": value})
        db.QueryInsert_(paramValues, "RETURNING id").Scan(&param_val_id)

        item := make(map[string]int, 2)
        item["param_val_id"] = param_val_id
        item["event_type_id"] = event_type_id
        param_val_ids = append(param_val_ids, item)
    }

    return param_val_ids, userLogin, userPass
}
