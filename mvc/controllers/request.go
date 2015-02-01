package controllers

import (
    "encoding/json"
    "fmt"
    "github.com/orc/db"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "html/template"
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
    utils.HandleErr("[Handler] Decode :", err, this.Response)

    event_id := data["event_id"]
    id := sessions.GetValue("id", this.Request).(string)

    users := GetModel("users")
    person, _ := users.Select([]string{"id", id}, "", []string{"person_id"})
    person_id := int(person[0].(map[string]interface{})["person_id"].(int64))

    query := `select param_id, p.name param_name, p.type, value, form_id, forms.name form_name from param_values 
        inner join params p on param_values.param_id = p.id
        inner join forms on forms.id = p.form_id
        where person_id = $1 and event_id = $2;`

    rows := db.Query(query, []interface{}{person_id, event_id})
    rowsInf := db.Exec(query, []interface{}{person_id, event_id})

    size, _ := rowsInf.RowsAffected()
    columns, _ := rows.Columns()
    result := db.ConvertData(columns, size, rows)

    response, err := json.Marshal(result)
    utils.HandleErr("[Handle select] json.Marshal: ", err, nil)
    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *Handler) GetListHistoryEvents() {
    if flag := sessions.CheackSession(this.Response, this.Request); !flag {
        return
    }
    this.Response.Header().Set("Access-Control-Allow-Origin", "*")
    this.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    this.Response.Header().Set("Content-type", "application/json")

    var data map[string]interface{}
    decoder := json.NewDecoder(this.Request.Body)
    err := decoder.Decode(&data)
    utils.HandleErr("[Handler] Decode :", err, this.Response)

    id := sessions.GetValue("id", this.Request).(string)
    ids := utils.ArrayInterfaceToString(data["form_ids"].([]interface{}))

    users := GetModel("users")
    person, _ := users.Select([]string{"id", id}, "", []string{"person_id"})
    person_id := int(person[0].(map[string]interface{})["person_id"].(int64))

    model := GetModel("forms_types")
    result, _ := model.Select(ids, "OR", []string{"type_id"})
    //fmt.Println("result: ", result)

    query := `SELECT DISTINCT event_id, name FROM param_values 
    inner join events on events.id = param_values.event_id
    WHERE event_id IN (SELECT DISTINCT event_id FROM events_types WHERE `

    var i int
    var params []interface{}

    for i = 1; i < reflect.ValueOf(result).Len(); i++ {
        query += "type_id=$" + strconv.Itoa(i) + " OR "
        params = append(params, result[i-1].(map[string]interface{})["type_id"])
    }

    query += "type_id=$" + strconv.Itoa(i) + ") AND person_id=$" + strconv.Itoa(i+1)

    params = append(params, result[i-1].(map[string]interface{})["type_id"])
    params = append(params, person_id)
    //fmt.Println("params: ", params)

    rows := db.Query(query, params)
    rowsInf := db.Exec(query, params)
    size, _ := rowsInf.RowsAffected()
    columns, _ := rows.Columns()
    events := db.ConvertData(columns, size, rows)

    response, err := json.Marshal(events)
    utils.HandleErr("[Handle GetListHistoryEvents] json.Marshal: ", err, nil)
    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *Handler) SaveUserRequest() {
    if flag := sessions.CheackSession(this.Response, this.Request); !flag {
        return
    }
    this.Response.Header().Set("Access-Control-Allow-Origin", "*")
    this.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    this.Response.Header().Set("Content-type", "application/json")

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
    param_values := GetModel("param_values")
    t := time.Now()

    if len(person) == 0 {
        persons_events.Insert(
            []string{"person_id", "event_id", "reg_date", "last_date"},
            []interface{}{person_id, event_id,
                t.Format("2006-01-02"),
                t.Format("2006-01-02")})

        for _, element := range inf {
            param_id := element.(map[string]interface{})["name"]
            value := element.(map[string]interface{})["value"]
            param_values.Insert(
                []string{"person_id", "event_id", "param_id", "value"},
                []interface{}{person_id, event_id, param_id, value})
        }
        response = map[string]interface{}{"result": "ok"}
    } else if len(person) != 0 {
        for _, element := range inf {
            param_id := element.(map[string]interface{})["name"]
            value := element.(map[string]interface{})["value"]
            param_values.Update(
                []string{"value"},
                []interface{}{value, person_id, event_id, param_id},
                "person_id=$"+strconv.Itoa(2)+" AND event_id=$"+strconv.Itoa(3)+" AND param_id=$"+strconv.Itoa(4))
        }
        persons_events.Update(
            []string{"last_date"},
            []interface{}{t.Format("2006-01-02"),
                person_id, event_id},
            "person_id=$"+strconv.Itoa(2)+" AND event_id=$"+strconv.Itoa(3))
        response = map[string]interface{}{"result": "ok"}
    } else {
        response = map[string]interface{}{"result": "exists"}
    }

    result, err := json.Marshal(response)
    utils.HandleErr("[Handle select] json.Marshal: ", err, nil)
    fmt.Fprintf(this.Response, "%s", string(result))
}

func (this *Handler) GetRequest(tableName, id string) {
    tmp, err := template.ParseFiles(
        "mvc/views/item.html",
        "mvc/views/header.html",
        "mvc/views/footer.html")
    utils.HandleErr("[Handler.Show] template.ParseFiles: ", err, nil)

    reaponse, err := json.Marshal(MegoJoin(tableName, id))
    utils.HandleErr("[Handler.Show] template.json.Marshal: ", err, nil)

    err = tmp.ExecuteTemplate(this.Response, "item", template.JS(reaponse))
    utils.HandleErr("[Handler.Show] tmp.Execute: ", err, nil)
}

func MegoJoin(tableName, id string) RequestModel {
    var E []interface{}
    var T []interface{}
    var F []interface{}
    var P []interface{}

    E = db.Select("events", []string{"id", id}, "", []string{"id", "name"})

    query := db.InnerJoin(
        []string{"id", "name"},
        "t",
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
            []string{"id", "name"},
            "f",
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
                []string{"id", "name", "type"},
                "p",
                "params",
                "p",
                []string{"form_id"},
                []string{"forms"},
                []string{"f"},
                []string{"id"},
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
