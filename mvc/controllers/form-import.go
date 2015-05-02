package controllers

import (
    "database/sql"
    "github.com/orc/db"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "net/http"
    "strconv"
)

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
        return
    }

    event_id, err := strconv.Atoi(request["event_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    query := `SELECT event_types.id, event_types.name FROM events_types
        INNER JOIN events ON events.id = events_types.event_id
        INNER JOIN event_types ON event_types.id = events_types.type_id
        WHERE events.id = $1 ORDER BY event_types.id;`
    result := db.Query(query, []interface{}{event_id})

    utils.SendJSReply(map[string]interface{}{"result": "ok", "data": result}, this.Response)

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
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
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
        formsResult := db.Query(query, []interface{}{int(eventResult[0].(map[string]interface{})["id"].(int))})

        for i := 0; i < len(formsResult); i++ {
            form_id := int(formsResult[i].(map[string]interface{})["id"].(int))
            eventsForms := this.GetModel("events_forms")
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
