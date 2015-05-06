package controllers

import (
    "github.com/orc/db"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "math"
    "net/http"
    "strconv"
    "encoding/json"
    "strings"
    "errors"
)

func (this *GridHandler) Load(tableName string) {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        utils.SendJSReply(map[string]interface{}{"result": errors.New("Forbidden")}, this.Response)
        http.Redirect(this.Response, this.Request, "/", http.StatusForbidden)
        return
    }

    var filters map[string]interface{}

    if this.Request.PostFormValue("_search") == "true" {
        err := json.NewDecoder(strings.NewReader(this.Request.PostFormValue("filters"))).Decode(&filters)
        if err != nil {
            utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
            return
        }
    }

    limit, err := strconv.Atoi(this.Request.PostFormValue("rows"))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    sord := this.Request.PostFormValue("sord")
    sidx := this.Request.FormValue("sidx")
    start := limit*page - limit

    if tableName == "search" {
        model := this.GetModel("param_values")

        var filters map[string]interface{}

        err := json.NewDecoder(strings.NewReader(this.Request.PostFormValue("filters"))).Decode(&filters)
        if err != nil {
            utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
            return
        }

        where, params := model.Where(filters)

        query := `SELECT faces.id, faces.user_id
            FROM reg_param_vals
            INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN events ON events.id = registrations.event_id
            INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
            INNER JOIN params ON params.id = param_values.param_id
            INNER JOIN users ON users.id = faces.user_id` + where + ` ORDER BY params.id ` + sord+` LIMIT $`+strconv.Itoa(len(params)+1)+` OFFSET $`+strconv.Itoa(len(params)+2)+`;`

        rows := db.Query(query, append(params, []interface{}{limit, start}...))

        query = `SELECT COUNT(*)
            FROM reg_param_vals
            INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN events ON events.id = registrations.event_id
            INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
            INNER JOIN params ON params.id = param_values.param_id
            INNER JOIN users ON users.id = faces.user_id` + where + `;`

        count := int(db.Query(query, params)[0].(map[string]interface{})["count"].(int))

        var totalPages int
        if count > 0 {
            totalPages = int(math.Ceil(float64(count) / float64(limit)))
        } else {
            totalPages = 0
        }

        result := make(map[string]interface{}, 4)
        result["rows"] = rows
        result["page"] = page
        result["total"] = totalPages
        result["records"] = count

        utils.SendJSReply(result, this.Response)

        return
    }

    model := this.GetModel(tableName)
    where, params := model.Where(filters)
    if len(where) < 8 {
        where = ""
    }
    query := `SELECT `+strings.Join(model.GetColumns(), ", ")+` FROM `+model.GetTableName()+where+` ORDER BY `+sidx+` `+ sord+` LIMIT $`+strconv.Itoa(len(params)+1)+` OFFSET $`+strconv.Itoa(len(params)+2)+`;`
    rows := db.Query(query, append(params, []interface{}{limit, start}...))
    count := db.SelectCount(tableName)

    var totalPages int
    if count > 0 {
        totalPages = int(math.Ceil(float64(count) / float64(limit)))
    } else {
        totalPages = 0
    }

    result := make(map[string]interface{}, 4)
    result["rows"] = rows
    result["page"] = page
    result["total"] = totalPages
    result["records"] = count

    utils.SendJSReply(result, this.Response)
}

func (this *Handler) UserGroupsLoad() {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    limit, err := strconv.Atoi(this.Request.PostFormValue("rows"))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    sidx := this.Request.FormValue("sidx")
    start := limit*page - limit

    query := `SELECT groups.id, groups.name FROM groups
        INNER JOIN faces ON faces.id = groups.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 ORDER BY $2 LIMIT $3 OFFSET $4;`
    rows := db.Query(query, []interface{}{user_id, sidx, limit, start})

    query = `SELECT COUNT(*) FROM groups
        INNER JOIN faces ON faces.id = groups.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1;`
    count := int(db.Query(query, []interface{}{user_id})[0].(map[string]interface{})["count"].(int))

    var totalPages int
    if count > 0 {
        totalPages = int(math.Ceil(float64(count) / float64(limit)))
    } else {
        totalPages = 0
    }

    result := make(map[string]interface{}, 2)
    result["rows"] = rows
    result["page"] = page
    result["total"] = totalPages
    result["records"] = count

    utils.SendJSReply(result, this.Response)
}

func (this *Handler) GroupsLoad() {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    limit, err := strconv.Atoi(this.Request.PostFormValue("rows"))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    sidx := this.Request.FormValue("sidx")
    start := limit*page - limit

    query := `SELECT groups.id, groups.name FROM groups
        INNER JOIN persons ON persons.group_id = groups.id
        INNER JOIN faces ON faces.id = persons.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 ORDER BY $2 LIMIT $3 OFFSET $4;`
    rows := db.Query(query, []interface{}{user_id, sidx, limit, start})

    query = `SELECT COUNT(*) FROM groups
        INNER JOIN persons ON persons.group_id = groups.id
        INNER JOIN faces ON faces.id = persons.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1;`
    count := int(db.Query(query, []interface{}{user_id})[0].(map[string]interface{})["count"].(int))

    var totalPages int
    if count > 0 {
        totalPages = int(math.Ceil(float64(count) / float64(limit)))
    } else {
        totalPages = 0
    }

    result := make(map[string]interface{}, 2)
    result["rows"] = rows
    result["page"] = page
    result["total"] = totalPages
    result["records"] = count

    utils.SendJSReply(result, this.Response)
}

func (this *Handler) RegistrationsLoad(userId string) {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    limit, err := strconv.Atoi(this.Request.PostFormValue("rows"))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    sidx := this.Request.FormValue("sidx")
    start := limit*page - limit

    if this.isAdmin() {
        user_id, err = strconv.Atoi(userId)
        if err != nil {
            utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
            return
        }

    }
    query := `SELECT registrations.id, registrations.event_id, registrations.status FROM registrations
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 ORDER BY $2 LIMIT $3 OFFSET $4;`
    rows := db.Query(query, []interface{}{user_id, sidx, limit, start})

    query = `SELECT COUNT(*) FROM registrations
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1;`
    count := int(db.Query(query, []interface{}{user_id})[0].(map[string]interface{})["count"].(int))

    var totalPages int
    if count > 0 {
        totalPages = int(math.Ceil(float64(count) / float64(limit)))
    } else {
        totalPages = 0
    }

    result := make(map[string]interface{}, 2)
    result["rows"] = rows
    result["page"] = page
    result["total"] = totalPages
    result["records"] = count

    utils.SendJSReply(result, this.Response)
}

func (this *Handler) GroupRegistrationsLoad() {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    limit, err := strconv.Atoi(this.Request.PostFormValue("rows"))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    sidx := this.Request.FormValue("sidx")
    start := limit*page - limit

    query := `SELECT group_registrations.id, group_registrations.event_id, group_registrations.group_id FROM group_registrations
        INNER JOIN events ON events.id = group_registrations.event_id
        INNER JOIN groups ON groups.id = group_registrations.group_id
        INNER JOIN faces ON faces.id = groups.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 ORDER BY $2 LIMIT $3 OFFSET $4;`
    rows := db.Query(query, []interface{}{user_id, sidx, limit, start})

    query = `SELECT COUNT(*) FROM group_registrations
        INNER JOIN events ON events.id = group_registrations.event_id
        INNER JOIN groups ON groups.id = group_registrations.group_id
        INNER JOIN faces ON faces.id = groups.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1;`
    count := int(db.Query(query, []interface{}{user_id})[0].(map[string]interface{})["count"].(int))

    var totalPages int
    if count > 0 {
        totalPages = int(math.Ceil(float64(count) / float64(limit)))
    } else {
        totalPages = 0
    }

    result := make(map[string]interface{}, 2)
    result["rows"] = rows
    result["page"] = page
    result["total"] = totalPages
    result["records"] = count

    utils.SendJSReply(result, this.Response)
}

func (this *Handler) PersonsLoad(group_id string) {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    limit, err := strconv.Atoi(this.Request.PostFormValue("rows"))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    id, err := strconv.Atoi(group_id)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    sidx := this.Request.FormValue("sidx")
    start := limit*page - limit

    var rows []interface{}

    if this.isAdmin() {
        query := `SELECT persons.id, persons.name, persons.email, persons.group_id, persons.face_id, persons.status
            FROM persons
            INNER JOIN groups ON groups.id = persons.group_id
            INNER JOIN faces ON faces.id = groups.face_id
            WHERE groups.id = $1 ORDER BY $2 LIMIT $3 OFFSET $4;`
        rows = db.Query(query, []interface{}{id, sidx, limit, start})

    } else {
        query := `SELECT persons.id, persons.name, persons.email, persons.group_id, persons.face_id, persons.status
        FROM persons
            INNER JOIN groups ON groups.id = persons.group_id
            INNER JOIN faces ON faces.id = groups.face_id
            INNER JOIN users ON users.id = faces.user_id
            WHERE users.id = $1 AND groups.id = $2 ORDER BY $3 LIMIT $4 OFFSET $5;`
        rows = db.Query(query, []interface{}{user_id, id, sidx, limit, start})
    }

    query := `SELECT COUNT(*) FROM persons
        INNER JOIN groups ON groups.id = persons.group_id
        INNER JOIN faces ON faces.id = groups.face_id
        WHERE groups.id = $1;`
    count := int(db.Query(query, []interface{}{user_id})[0].(map[string]interface{})["count"].(int))

    var totalPages int
    if count > 0 {
        totalPages = int(math.Ceil(float64(count) / float64(limit)))
    } else {
        totalPages = 0
    }

    result := make(map[string]interface{}, 2)
    result["rows"] = rows
    result["page"] = page
    result["total"] = totalPages
    result["records"] = count

    utils.SendJSReply(result, this.Response)
}
