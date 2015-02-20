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
    "time"
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

    user := GetModel("users")
    user.LoadWherePart(map[string]interface{}{"id": user_id})

    var person_id int
    err = db.SelectRow(user, []string{"person_id"}, "").Scan(&person_id)
    if err != nil {
        result["result"] = err.Error()
    } else {
        query := `select param_id, p.name, event_type_id, p_t.name as type, value, form_id from param_values
                inner join params p on param_values.param_id = p.id
                inner join forms on forms.id = p.form_id
                inner join param_types p_t on p_t.id = p.param_type_id
                where person_id = $1 and event_id = $2;`

        result["data"] = db.Query(query, []interface{}{person_id, event_id})
    }

    response, err := json.Marshal(result)
    if utils.HandleErr("[Handle::GetHistoryRequest] Marshal: ", err, this.Response) {
        return
    }

    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *Handler) GetListHistoryEvents() {
    if flag := sessions.CheackSession(this.Response, this.Request); !flag {
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
    user := GetModel("users")
    user.LoadWherePart(map[string]interface{}{"id": user_id})

    var person_id int
    err = db.SelectRow(user, []string{"person_id"}, "").Scan(&person_id)

    if err == sql.ErrNoRows {
        result["result"] = err.Error()
    } else {
        ids := make(map[string]interface{}, 1)
        ids["form_id"] = make([]interface{}, 0)
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
                inner join persons_events on persons_events.event_id = events.id
                inner join persons on persons.id = persons_events.person_id
                WHERE persons.id=$1 AND events.id IN (SELECT DISTINCT event_id FROM events_types WHERE `

            var i int
            var params []interface{}

            params = append(params, person_id)

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
    if flag := sessions.CheackSession(this.Response, this.Request); !flag {
        return
    }

    var response interface{}

    var data map[string]interface{}
    decoder := json.NewDecoder(this.Request.Body)
    err := decoder.Decode(&data)
    if utils.HandleErr("[Handler::SaveUserRequest] Decode :", err, this.Response) {
        return
    }

    var person_id int
    user_id := sessions.GetValue("id", this.Request)
    if user_id == nil {
        http.Redirect(this.Response, this.Request, "/", 401)
        return
    }
    user := GetModel("users")
    user.LoadWherePart(map[string]interface{}{"id": user_id})
    err = db.SelectRow(user, []string{"person_id"}, "").Scan(&person_id)
    if err != nil {
        response = map[string]interface{}{"result": err.Error()}
    } else {
        event_id := int(data["event_id"].(float64))
        personsEvents := GetModel("persons_events")
        personsEvents.LoadWherePart(map[string]interface{}{"person_id": person_id, "event_id": event_id})
        person := db.Select(personsEvents, []string{"id"}, "AND")

        model := GetModel("persons_events")
        if len(person) == 0 {
            model.LoadModelData(map[string]interface{}{
                "person_id": person_id,
                "event_id":  event_id,
                "reg_date":  time.Now().Format("2006-01-02"),
                "last_date": time.Now().Format("2006-01-02"),
            })
            db.QueryInsert_(model, "")
            response = map[string]interface{}{"result": "ok"}

        } else {
            model.LoadModelData(map[string]interface{}{
                "id":        strconv.Itoa(int(person[0].(map[string]interface{})["id"].(int64))),
                "last_date": time.Now().Format("2006-01-02"),
            })
            db.QueryUpdate_(model, "")
            response = map[string]interface{}{"result": "ok"}
        }

        for _, element := range data["data"].([]interface{}) {
            param_id := element.(map[string]interface{})["id"]
            event_type_id := element.(map[string]interface{})["event_type_id"]
            value := element.(map[string]interface{})["value"]

            paramValues := GetModel("param_values")

            if db.IsExists_(
                "param_values",
                []string{"person_id", "event_id", "param_id", "event_type_id"},
                []interface{}{person_id, event_id, param_id, event_type_id}) {

                paramValues.LoadModelData(map[string]interface{}{"value": value})
                paramValues.LoadWherePart(map[string]interface{}{
                    "person_id":     person_id,
                    "event_id":      event_id,
                    "param_id":      param_id,
                    "event_type_id": event_type_id,
                })
                db.QueryUpdate_(paramValues, "AND")

            } else {
                paramValues.LoadModelData(map[string]interface{}{
                    "person_id":     person_id,
                    "event_id":      event_id,
                    "param_id":      param_id,
                    "event_type_id": event_type_id,
                    "value":         value,
                })
                db.QueryInsert_(paramValues, "")
            }
        }
    }

    result, err := json.Marshal(response)
    if utils.HandleErr("[Handle::SaveUserRequest] Marshal: ", err, this.Response) {
        return
    }

    fmt.Fprintf(this.Response, "%s", string(result))
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
