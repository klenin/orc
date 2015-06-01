package controllers

import (
    "github.com/orc/db"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "net/http"
    "strconv"
    "errors"
    "strings"
)

func (this *GridHandler) GetPersonsByEventId() {
    if !sessions.CheckSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        http.Redirect(this.Response, this.Request, "/", http.StatusForbidden)
        return
    }

    if this.Request.URL.Query().Get("event") == "" || this.Request.URL.Query().Get("params") == "" {
        return
    }

    eventId, err := strconv.Atoi(this.Request.URL.Query().Get("event"))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    paramsIds := strings.Split(this.Request.URL.Query().Get("params"), ",")

    if len(paramsIds) == 0 {
        utils.SendJSReply(map[string]interface{}{"result": "Выберите параметры."}, this.Response)
        return
    }

    var queryParams []interface{}
    query := "SELECT params.name FROM params WHERE params.id in ("

    for k, v := range paramsIds {
        param_id, err := strconv.Atoi(v)
        if err != nil {
            utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
            return
        }
        query += "$"+strconv.Itoa(k+1)+", "
        queryParams = append(queryParams, param_id)
    }
    query = query[:len(query)-2]
    query+=") ORDER BY id;"

    var caption []string
    for _, v := range db.Query(query, queryParams) {
        caption = append(caption, v.(map[string]interface{})["name"].(string))
    }

    result := []interface{}{0: map[string]interface{}{"id": -1, "data": caption}}

    query = `SELECT reg_param_vals.reg_id as id, array_agg(param_values.value) as data
        FROM reg_param_vals
        INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
        INNER JOIN params ON params.id = param_values.param_id
        WHERE params.id in (` + strings.Join(db.MakeParams(len(queryParams)), ", ")
    query += ") AND events.id = $" + strconv.Itoa(len(queryParams)+1) + " GROUP BY reg_param_vals.reg_id ORDER BY reg_param_vals.reg_id;"

    data := db.Query(query, append(queryParams, eventId))

    this.Render([]string{"mvc/views/list.html"}, "list", append(result, data...))
}

func (this *GridHandler) GetParamsByEventId() {
    if !sessions.CheckSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        utils.SendJSReply(map[string]interface{}{"result": errors.New("Forbidden")}, this.Response)
        http.Redirect(this.Response, this.Request, "/", http.StatusForbidden)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    eventId, err := strconv.Atoi(request["event_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    query := `SELECT params.id, params.name
        FROM events_forms
        INNER JOIN forms ON forms.id = events_forms.form_id
        INNER JOIN params ON params.form_id = forms.id
        INNER JOIN events ON events.id = events_forms.event_id
        WHERE events.id = $1 ORDER BY params.id;`
    result := db.Query(query, []interface{}{eventId})

    utils.SendJSReply(map[string]interface{}{"result": "ok", "data": result}, this.Response)
}
