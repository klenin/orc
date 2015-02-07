package controllers

import (
    "encoding/json"
    "fmt"
    "github.com/orc/db"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "html/template"
    "log"
    "reflect"
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

    var data map[string]string
    decoder := json.NewDecoder(this.Request.Body)
    err := decoder.Decode(&data)
    utils.HandleErr("[Handler::GetHistoryRequest] Decode :", err, this.Response)

    event_id := data["event_id"]
    id := sessions.GetValue("id", this.Request).(string)

    users := GetModel("users")
    person, _ := users.Select([]string{"id", id}, "", []string{"person_id"})
    person_id := int(person[0].(map[string]interface{})["person_id"].(int64))

    query := `select param_id, p.name, event_type_id, p_t.name as type, value, form_id from param_values
        inner join params p on param_values.param_id = p.id
        inner join forms on forms.id = p.form_id
        inner join param_types p_t on p_t.id = p.param_type_id
        where person_id = $1 and event_id = $2;`

    rows := db.Query(query, []interface{}{person_id, event_id})
    rowsInf := db.Exec(query, []interface{}{person_id, event_id})

    size, _ := rowsInf.RowsAffected()
    columns, _ := rows.Columns()
    result := db.ConvertData(columns, size, rows)

    response, err := json.Marshal(result)
    log.Println("history information about blank")
    log.Println(string(response))
    utils.HandleErr("[Handle::GetHistoryRequest] Marshal: ", err, this.Response)
    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *Handler) GetListHistoryEvents() {
    if flag := sessions.CheackSession(this.Response, this.Request); !flag {
        return
    }
    var data map[string]interface{}
    decoder := json.NewDecoder(this.Request.Body)
    err := decoder.Decode(&data)
    utils.HandleErr("[Handler::GetListHistoryEvents] Decode :", err, this.Response)

    id := sessions.GetValue("id", this.Request).(string)
    ids := utils.ArrayInterfaceToString(data["form_ids"].([]interface{}))

    users := GetModel("users")
    person, _ := users.Select([]string{"id", id}, "", []string{"person_id"})
    person_id := int(person[0].(map[string]interface{})["person_id"].(int64))

    model := GetModel("forms_types")
    result, _ := model.Select(ids, "OR", []string{"type_id"})

    if len(result) == 0 {
        return
    }

    query := `SELECT DISTINCT events.id, events.name FROM events
        inner join events_types on events_types.event_id = events.id
        inner join event_types on events_types.type_id = event_types.id
        inner join persons_events on persons_events.event_id = events.id
        inner join persons on persons.id = persons_events.person_id
        WHERE persons.id=$1 AND events.id IN (SELECT DISTINCT event_id FROM events_types WHERE `

    var i int
    var params []interface{}
    params = append(params, person_id)

    for i = 2; i < reflect.ValueOf(result).Len(); i++ {
        query += "type_id=$" + strconv.Itoa(i) + " OR "
        params = append(params, result[i-2].(map[string]interface{})["type_id"])
    }

    query += "type_id=$" + strconv.Itoa(i) + ") AND person_id=$" + strconv.Itoa(i+1)

    params = append(params, result[i-1].(map[string]interface{})["type_id"])
    params = append(params, person_id)

    rows := db.Query(query, params)
    rowsInf := db.Exec(query, params)
    size, _ := rowsInf.RowsAffected()
    columns, _ := rows.Columns()
    events := db.ConvertData(columns, size, rows)

    response, err := json.Marshal(events)
    utils.HandleErr("[Handle::GetListHistoryEvents] Marshal: ", err, this.Response)
    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *Handler) SaveUserRequest() {
    if flag := sessions.CheackSession(this.Response, this.Request); !flag {
        return
    }
    var data map[string]interface{}
    decoder := json.NewDecoder(this.Request.Body)
    err := decoder.Decode(&data)
    utils.HandleErr("[Handler] Decode :", err, this.Response)

    event_id := int(data["event_id"].(float64))
    id := sessions.GetValue("id", this.Request).(string)

    users := GetModel("users")
    person, _ := users.Select([]string{"id", id}, "", []string{"person_id"})
    person_id := int(person[0].(map[string]interface{})["person_id"].(int64))

    persons_events := GetModel("persons_events")
    person, _ = persons_events.Select([]string{"person_id", strconv.Itoa(person_id), "event_id", strconv.Itoa(event_id)}, "AND", []string{"person_id"})

    var response interface{}
    inf := data["data"].([]interface{})
    log.Println("save data from blank")
    log.Println(inf)
    param_values := GetModel("param_values")
    t := time.Now()

    if len(person) == 0 {
        persons_events.Insert(
            []string{"person_id", "event_id", "reg_date", "last_date"},
            []interface{}{person_id, event_id,
                t.Format("2006-01-02"),
                t.Format("2006-01-02")})
        response = map[string]interface{}{"result": "ok"}
    } else if len(person) != 0 {
        persons_events.Update(
            []string{"last_date"},
            []interface{}{t.Format("2006-01-02"),
                person_id, event_id},
            "person_id=$"+strconv.Itoa(2)+" AND event_id=$"+strconv.Itoa(3))
        response = map[string]interface{}{"result": "ok"}
    } else {
        response = map[string]interface{}{"result": "WHAT?!"}//???
    }

    for _, element := range inf {
            param_id := element.(map[string]interface{})["id"]
            event_type_id := element.(map[string]interface{})["event_type_id"]
            value := element.(map[string]interface{})["value"]

        if db.IsExists_("param_values", []string{"person_id", "event_id", "param_id", "event_type_id"}, []interface{}{person_id, event_id, param_id, event_type_id}) {
            param_values.Update(
                []string{"value"},
                []interface{}{value, person_id, event_id, param_id, event_type_id},
                "person_id=$"+strconv.Itoa(2)+" AND event_id=$"+strconv.Itoa(3)+" AND param_id=$"+strconv.Itoa(4)+" AND event_type_id=$"+strconv.Itoa(5))
        } else {
            param_values.Insert(
                []string{"person_id", "event_id", "param_id", "value", "event_type_id"},
                []interface{}{person_id, event_id, param_id, value, event_type_id})
        }
    }

    result, err := json.Marshal(response)
    utils.HandleErr("[Handle::SaveUserRequest] Marshal: ", err, this.Response)
    fmt.Fprintf(this.Response, "%s", string(result))
}

func (this *Handler) GetRequest(tableName, id string) {
    tmp, err := template.ParseFiles(
        "mvc/views/item.html",
        "mvc/views/header.html",
        "mvc/views/footer.html")
    utils.HandleErr("[Handler::GetRequest] ParseFiles: ", err, this.Response)

    reaponse, err := json.Marshal(MegoJoin(tableName, id))
    utils.HandleErr("[Handler::GetRequest] Marshal: ", err, this.Response)

    err = tmp.ExecuteTemplate(this.Response, "item", template.JS(reaponse))
    utils.HandleErr("[Handler::GetRequest] ExecuteTemplate: ", err, this.Response)
}

func MegoJoin(tableName, id string) RequestModel {
    var E []interface{}
    var T []interface{}
    var F []interface{}
    var P []interface{}

    E = db.Select("events", []string{"id", id}, "", []string{"id", "name"})

    query := db.InnerJoin(
        []string{"t.id", "t.name"},
        "events_types",
        "e_t",
        []string{"event_id", "type_id"},
        []string{"events", "event_types"},
        []string{"e", "t"},
        []string{"id", "id"},
        "where e.id=$1")

    rows := db.Query(query, []interface{}{id})
    rowsInf := db.Exec(query, []interface{}{id})
    l, _ := rowsInf.RowsAffected()
    c, _ := rows.Columns()
    T = db.ConvertData(c, l, rows)

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

        rows := db.Query(query, []interface{}{id})
        rowsInf := db.Exec(query, []interface{}{id})
        l, _ := rowsInf.RowsAffected()
        c, _ := rows.Columns()
        F = append(F, db.ConvertData(c, l, rows))
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

            rows := db.Query(query, []interface{}{id})
            rowsInf := db.Exec(query, []interface{}{id})
            l, _ := rowsInf.RowsAffected()
            c, _ := rows.Columns()
            PP = append(PP, db.ConvertData(c, l, rows))
        }
        P = append(P, PP)
    }
    return RequestModel{E: E, T: T, F: F, P: P}
}
