package controllers

import (
    "encoding/json"
    "github.com/klenin/orc/db"
    "github.com/klenin/orc/sessions"
    "github.com/klenin/orc/utils"
    "math"
    "log"
    "net/http"
    "strconv"
    "strings"
)

func (c *BaseController) Handler() *Handler {
    return new(Handler)
}

type Handler struct {
    Controller
}

func (this *GridController) Load(tableName string) {
    if tableName != "events" && !sessions.CheckSession(this.Response, this.Request) {
        http.Error(this.Response, "Unauthorized", 400)
        return
    }

    isAdmin := this.isAdmin()

    var filters map[string]interface{}
    if this.Request.PostFormValue("_search") == "true" {
        err := json.NewDecoder(strings.NewReader(this.Request.PostFormValue("filters"))).Decode(&filters)
        if err != nil {
            http.Error(this.Response, err.Error(), 400)
            return
        }
    }

    limit, err := strconv.Atoi(this.Request.PostFormValue("rows"))
    if err != nil {
        http.Error(this.Response, err.Error(), 400)
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if err != nil {
        http.Error(this.Response, err.Error(), 400)
        return
    }

    sord := this.Request.PostFormValue("sord")
    sidx := this.Request.FormValue("sidx")
    start := limit*page - limit

    if tableName == "search" {
        var filters map[string]interface{}
        err := json.NewDecoder(strings.NewReader(this.Request.PostFormValue("filters"))).Decode(&filters)
        if err != nil {
            utils.SendJSReply(nil, this.Response)
            return
        }

        model := this.GetModel("faces")
        query := `SELECT DISTINCT faces.id, faces.user_id
            FROM param_values
            INNER JOIN registrations ON registrations.id = param_values.reg_id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN events ON events.id = registrations.event_id
            INNER JOIN params ON params.id = param_values.param_id
            INNER JOIN users ON users.id = faces.user_id`

        where, params, _ := model.WhereByParams(filters, 1)

        if !isAdmin {
            where = ` WHERE events.id = 1 AND users.enabled = true AND `+where
        } else {
            if where != "" {
                where = " WHERE "+where
            }
        }
        where += ` ORDER BY faces.id `+sord
        query += where+` LIMIT $`+strconv.Itoa(len(params)+1)+` OFFSET $`+strconv.Itoa(len(params)+2)+`;`
        rows := db.Query(query, append(params, []interface{}{limit, start}...))

        query = `SELECT COUNT(*)
            FROM (SELECT DISTINCT faces.id, faces.user_id
            FROM param_values
            INNER JOIN registrations ON registrations.id = param_values.reg_id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN events ON events.id = registrations.event_id
            INNER JOIN params ON params.id = param_values.param_id
            INNER JOIN users ON users.id = faces.user_id`
        query += where+") as count;"
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
    where, params, _ := model.Where(filters, 1)

    if tableName == "param_values" && !isAdmin {
        w := " WHERE param_values.param_id in (4, 5, 6, 7)"
        if where != "" {
            where = w+" AND "+where
        } else {
            where = w
        }
    } else {
        if where != "" {
            where = " WHERE "+where
        }
    }

    query := `SELECT `+strings.Join(model.GetColumns(), ", ")+` FROM `+model.GetTableName()+where+` ORDER BY `+sidx+` `+sord+` LIMIT $`+strconv.Itoa(len(params)+1)+` OFFSET $`+strconv.Itoa(len(params)+2)+`;`
    rows := db.Query(query, append(params, []interface{}{limit, start}...))

    query = `SELECT COUNT(*) FROM (SELECT `+model.GetTableName()+`.id FROM `+model.GetTableName()
    query += where+`) as count;`
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
}

func (this *Handler) UserGroupsLoad() {
    userId, err := this.CheckSid()
    if err != nil {
        http.Error(this.Response, "Unauthorized", 400)
        return
    }

    limit, err := strconv.Atoi(this.Request.PostFormValue("rows"))
    if err != nil {
        http.Error(this.Response, err.Error(), 400)
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if err != nil {
        http.Error(this.Response, err.Error(), 400)
        return
    }

    sidx := this.Request.FormValue("sidx")
    start := limit * page - limit

    query := `SELECT groups.id, groups.name FROM groups
        INNER JOIN faces ON faces.id = groups.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 ORDER BY $2 LIMIT $3 OFFSET $4;`
    rows := db.Query(query, []interface{}{userId, sidx, limit, start})

    query = `SELECT COUNT(*) FROM (SELECT groups.id FROM groups
        INNER JOIN faces ON faces.id = groups.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1) as count;`
    count := int(db.Query(query, []interface{}{userId})[0].(map[string]interface{})["count"].(int))

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
    userId, err := this.CheckSid()
    if err != nil {
        http.Error(this.Response, "Unauthorized", 400)
        return
    }

    limit, err := strconv.Atoi(this.Request.PostFormValue("rows"))
    if err != nil {
        http.Error(this.Response, err.Error(), 400)
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if err != nil {
        http.Error(this.Response, err.Error(), 400)
        return
    }

    sidx := this.Request.FormValue("sidx")
    start := limit * page - limit

    query := `SELECT groups.id, groups.name, groups.face_id
        FROM groups
        INNER JOIN persons ON persons.group_id = groups.id
        INNER JOIN faces ON faces.id = persons.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 ORDER BY $2 LIMIT $3 OFFSET $4;`
    rows := db.Query(query, []interface{}{userId, sidx, limit, start})

    query = `SELECT COUNT(*) FROM (SELECT groups.id FROM groups
        INNER JOIN persons ON persons.group_id = groups.id
        INNER JOIN faces ON faces.id = persons.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1) as count;`
    count := int(db.Query(query, []interface{}{userId})[0].(map[string]interface{})["count"].(int))

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

func (this *Handler) RegistrationsLoad(userId_ string) {
    userId, err := this.CheckSid()
    if err != nil {
        http.Error(this.Response, "Unauthorized", 400)
        return
    }

    limit, err := strconv.Atoi(this.Request.PostFormValue("rows"))
    if err != nil {
        http.Error(this.Response, err.Error(), 400)
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if err != nil {
        http.Error(this.Response, err.Error(), 400)
        return
    }

    sidx := this.Request.FormValue("sidx")
    start := limit * page - limit

    var id int
    if this.isAdmin() {
        id, err = strconv.Atoi(userId_)
        if err != nil {
            http.Error(this.Response, err.Error(), 400)
            return
        }
    } else {
        id = userId
    }

    query := `SELECT registrations.id, registrations.event_id, registrations.status FROM registrations
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN events_forms ON events_forms.event_id = events.id
        INNER JOIN forms ON forms.id = events_forms.form_id
        INNER JOIN params ON forms.id = params.form_id
        INNER JOIN param_types ON param_types.id = params.param_type_id
        INNER JOIN param_values ON params.id = param_values.param_id AND param_values.reg_id = registrations.id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 AND forms.personal = true
        GROUP BY registrations.id
        ORDER BY $2 LIMIT $3 OFFSET $4;`
    rows := db.Query(query, []interface{}{id, sidx, limit, start})

    query = `SELECT COUNT(*) FROM (SELECT registrations.id FROM registrations
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN events_forms ON events_forms.event_id = events.id
        INNER JOIN forms ON forms.id = events_forms.form_id
        INNER JOIN params ON forms.id = params.form_id
        INNER JOIN param_types ON param_types.id = params.param_type_id
        INNER JOIN param_values ON params.id = param_values.param_id AND param_values.reg_id = registrations.id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 AND forms.personal = true
        GROUP BY registrations.id) as count;`
    count := int(db.Query(query, []interface{}{id})[0].(map[string]interface{})["count"].(int))

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

func (this *Handler) UserGroupRegistrationsLoad(isTeam string) {
    userId, err := this.CheckSid()
    if err != nil {
        http.Error(this.Response, "Unauthorized", 400)
        return
    }

    limit, err := strconv.Atoi(this.Request.PostFormValue("rows"))
    if err != nil {
        http.Error(this.Response, err.Error(), 400)
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if err != nil {
        http.Error(this.Response, err.Error(), 400)
        return
    }

    sidx := this.Request.FormValue("sidx")
    start := limit * page - limit

    isTeam_, err := strconv.ParseBool(isTeam)
    if err != nil {
        http.Error(this.Response, err.Error(), 400)
        return
    }

    query := `SELECT group_registrations.id, group_registrations.event_id,
            group_registrations.group_id
        FROM group_registrations
        INNER JOIN events ON events.id = group_registrations.event_id
        INNER JOIN groups ON groups.id = group_registrations.group_id
        INNER JOIN faces ON faces.id = groups.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 AND events.team = $2 ORDER BY $3 LIMIT $4 OFFSET $5;`
    rows := db.Query(query, []interface{}{userId, isTeam_, sidx, limit, start})

    query = `SELECT COUNT(*) FROM (SELECT group_registrations.id FROM group_registrations
        INNER JOIN events ON events.id = group_registrations.event_id
        INNER JOIN groups ON groups.id = group_registrations.group_id
        INNER JOIN faces ON faces.id = groups.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 AND events.team = $2) as count;`

    var count int
    db.QueryRow(query, []interface{}{userId, isTeam_}).Scan(&count)

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
    userId, err := this.CheckSid()
    if err != nil {
        http.Error(this.Response, "Unauthorized", 400)
        return
    }

    limit, err := strconv.Atoi(this.Request.PostFormValue("rows"))
    if err != nil {
        http.Error(this.Response, err.Error(), 400)
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if err != nil {
        http.Error(this.Response, err.Error(), 400)
        return
    }

    sidx := this.Request.FormValue("sidx")
    start := limit * page - limit

    query := `SELECT group_registrations.id, group_registrations.event_id,
            group_registrations.group_id
        FROM group_registrations
        INNER JOIN events ON events.id = group_registrations.event_id
        INNER JOIN groups ON groups.id = group_registrations.group_id
        INNER JOIN persons ON persons.group_id = groups.id
        INNER JOIN faces ON faces.id = persons.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 AND events.team = $2 ORDER BY $3 LIMIT $4 OFFSET $5;`
    rows := db.Query(query, []interface{}{userId, true, sidx, limit, start})

    query = `SELECT COUNT(*) FROM (SELECT group_registrations.id
        FROM group_registrations
        INNER JOIN events ON events.id = group_registrations.event_id
        INNER JOIN groups ON groups.id = group_registrations.group_id
        INNER JOIN persons ON persons.group_id = groups.id
        INNER JOIN faces ON faces.id = persons.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 AND events.team = $2 GROUP BY group_registrations.id) as count;`

    var count int
    db.QueryRow(query, []interface{}{userId, true}).Scan(&count)

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

func (this *Handler) PersonsLoad(groupId string) {
    userId, err := this.CheckSid()
    if err != nil {
        http.Error(this.Response, "Unauthorized", 400)
        return
    }

    limit, err := strconv.Atoi(this.Request.PostFormValue("rows"))
    if err != nil {
        http.Error(this.Response, err.Error(), 400)
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if err != nil {
        http.Error(this.Response, err.Error(), 400)
        return
    }

    id, err := strconv.Atoi(groupId)
    if err != nil {
        http.Error(this.Response, err.Error(), 400)
        return
    }

    faceId := -1
    query := `SELECT groups.face_id FROM groups
        INNER JOIN faces ON faces.id = groups.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 AND groups.id = $2;`
    err = db.QueryRow(query, []interface{}{userId, id}).Scan(&faceId)

    if (err != nil || faceId == -1) && !this.isAdmin() {
        http.Error(this.Response, "Вы не являетесь владельцем группы", 400)
        return
    }

    sidx := this.Request.FormValue("sidx")
    start := limit * page - limit

    query = `SELECT persons.id, persons.group_id, persons.face_id, persons.status
        FROM persons
        INNER JOIN groups ON groups.id = persons.group_id
        WHERE groups.id = $1 ORDER BY $2 LIMIT $3 OFFSET $4;`
    rows := db.Query(query, []interface{}{id, sidx, limit, start})

    query = `SELECT COUNT(*) FROM (SELECT persons.id FROM persons
        INNER JOIN groups ON groups.id = persons.group_id
        WHERE groups.id = $1) as count;`
    count := int(db.Query(query, []interface{}{id})[0].(map[string]interface{})["count"].(int))

    log.Println("count: ", count)

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
