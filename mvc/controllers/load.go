package controllers

import (
    "github.com/orc/db"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "math"
    "net/http"
    "strconv"
    "encoding/json"
    "log"
    "strings"
)

func (this *GridHandler) Load(tableName string) {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        http.Redirect(this.Response, this.Request, "/", http.StatusForbidden)
        return
    }

    qWhere := ""

    if this.Request.PostFormValue("_search") == "true" {
        qWhere = " WHERE "
        var filters map[string]interface{}
        err := json.NewDecoder(strings.NewReader(this.Request.PostFormValue("filters"))).Decode(&filters)
        if utils.HandleErr("[GridHandler::Load]: ", err, this.Response) {
            return
        }

        groupOp := filters["groupOp"].(string)
        rules := filters["rules"].([]interface{})

        if len(rules) > 10 {
            log.Println("More 10 rules for serching!")
        }

        firstElem := true

        for _, v := range rules {
            if !firstElem {
                if groupOp != "AND" && groupOp != "OR" {
                    log.Println("`groupOp` parameter is not allowed!")
                    continue
                }
                qWhere += " " + groupOp + " "
            } else {
                firstElem = false
            }

            rule := v.(map[string]interface{})

            switch rule["op"].(string) {
            case "eq":
                qWhere += rule["field"].(string) + " = " + rule["data"].(string)
                break
            case "ne":
                qWhere += rule["field"].(string) + " <> " + rule["data"].(string)
                break
            case "bw":
                qWhere += rule["field"].(string) + " LIKE " + rule["data"].(string) + "%"
                break
            case "cn":
                qWhere += rule["field"].(string) + " LIKE %" + rule["data"].(string) + "%"
                break
            default:
                panic("`op` parameter is not allowed!")
            }
        }
    }

    limit, err := strconv.Atoi(this.Request.PostFormValue("rows"))
    if utils.HandleErr("[GridHandler::Load] limit Atoi: ", err, this.Response) {
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if utils.HandleErr("[GridHandler::Load] page Atoi: ", err, this.Response) {
        return
    }

    sord := this.Request.PostFormValue("sord")
    sidx := this.Request.FormValue("sidx")
    start := limit*page - limit

    model := GetModel(tableName)

    query := `SELECT ` + strings.Join(model.GetColumns(), ", ") + ` FROM ` + model.GetTableName()
    query += qWhere + " ORDER BY $1 " + sord + " LIMIT $2 OFFSET $3;"
    rows := db.Query(query, []interface{}{sidx, limit, start})

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

func (this *Handler) GroupsLoad() {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    limit, err := strconv.Atoi(this.Request.PostFormValue("rows"))
    if utils.HandleErr("[GridHandler::GroupsLoad] limit Atoi: ", err, this.Response) {
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if utils.HandleErr("[GridHandler::GroupsLoad] page Atoi: ", err, this.Response) {
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

    count := int(db.Query(query, []interface{}{user_id})[0].(map[string]interface{})["count"].(int64))

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

func (this *Handler) RegistrationsLoad() {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    limit, err := strconv.Atoi(this.Request.PostFormValue("rows"))
    if utils.HandleErr("[GridHandler::RegistrationsLoad] limit Atoi: ", err, this.Response) {
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if utils.HandleErr("[GridHandler::RegistrationsLoad] page Atoi: ", err, this.Response) {
        return
    }

    sidx := this.Request.FormValue("sidx")
    start := limit*page - limit

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

    count := int(db.Query(query, []interface{}{user_id})[0].(map[string]interface{})["count"].(int64))

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
    if utils.HandleErr("[GridHandler::GroupRegistrationsLoad] limit Atoi: ", err, this.Response) {
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if utils.HandleErr("[GridHandler::GroupRegistrationsLoad] page Atoi: ", err, this.Response) {
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

    count := int(db.Query(query, []interface{}{user_id})[0].(map[string]interface{})["count"].(int64))

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
    if utils.HandleErr("[GridHandler::PersonsLoad] limit Atoi: ", err, this.Response) {
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if utils.HandleErr("[GridHandler::PersonsLoad] page Atoi: ", err, this.Response) {
        return
    }

    id, err := strconv.Atoi(group_id)
    if utils.HandleErr("[GridHandler::PersonsLoad] id Atoi: ", err, this.Response) {
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

    count := int(db.Query(query, []interface{}{user_id})[0].(map[string]interface{})["count"].(int64))

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