package controllers

import (
    "database/sql"
    "github.com/orc/db"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "github.com/orc/mailer"
    "math"
    "net/http"
    "strconv"
)

func (c *BaseController) Handler() *Handler {
    return new(Handler)
}

type Handler struct {
    Controller
}

func (this *Handler) GetList() {
    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    } else {
        fields := request["fields"].([]interface{})
        result := db.Select(GetModel(request["table"].(string)), utils.ArrayInterfaceToString(fields))
        utils.SendJSReply(map[string]interface{}{"result": "ok", "data": result}, this.Response)
    }
}

func (this *Handler) Index() {
    var response interface{}

    data, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    switch data["action"].(string) {
    case "login":
        response = this.HandleLogin(data["login"].(string), data["password"].(string))
        utils.SendJSReply(response, this.Response)
        break

    case "logout":
        utils.SendJSReply(this.HandleLogout(), this.Response)
        break

    case "checkSession":
        var userHash string
        var result interface{}

        hash := sessions.GetValue("hash", this.Request)

        if hash == nil {
            result = map[string]interface{}{"result": "no"}
        } else {
            user := GetModel("users")
            user.LoadWherePart(map[string]interface{}{"hash": hash})
            err := db.SelectRow(user, []string{"hash"}).Scan(&userHash)
            if err != sql.ErrNoRows {
                result = map[string]interface{}{"result": "ok"}
            } else {
                result = map[string]interface{}{"result": "no"}
            }
        }

        utils.SendJSReply(result, this.Response)
        break
    }
}

func (this *Handler) ShowCabinet(tableName string) {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    user := GetModel("users")
    user.LoadWherePart(map[string]interface{}{"id": user_id})

    var role string
    err := db.SelectRow(user, []string{"role"}).Scan(&role)
    if err != nil {
        utils.HandleErr("[Handle::ShowCabinet]: ", err, this.Response)
        return
    }

    var model Model
    if role == "admin" {
        model = Model{Columns: db.Tables, ColNames: db.TableNames}
    } else {

        query := `SELECT params.name, param_values.value FROM param_values
            INNER JOIN params ON params.id = param_values.param_id
            INNER JOIN reg_param_vals ON reg_param_vals.param_val_id = param_values.id
            INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
            INNER JOIN events ON events.id = registrations.event_id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN users ON users.id = faces.user_id
            WHERE events.id=$1 AND users.id=$2 ORDER BY params.id`

        regParamVals := GetModel("reg_param_vals")
        data := db.Query(query, []interface{}{1, user_id})
        model = Model{Table: data, Columns: regParamVals.GetColumns(), ColNames: regParamVals.GetColNames()}
    }

    this.Render([]string{"mvc/views/"+role+".html"}, role, model)
}

func (this *Handler) CreateGroup() {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    } else {
        event_id, err := strconv.Atoi(request["event_id"].(string))
        if utils.HandleErr("[GridHandler::GetParamsByEventId] event_id Atoi: ", err, this.Response) {
            return
        }

        if request["group-name"] == nil {
            utils.SendJSReply(map[string]interface{}{"result": "Заполните поле \"Название группы\""}, this.Response)
            return
        }

        var eventName string
        err = db.QueryRow("SELECT events.name FROM events WHERE events.id = $1", []interface{}{event_id}).Scan(&eventName)
        if err != nil {
            utils.SendJSReply(map[string]interface{}{"result": "Нет информации о таком мероприятии."}, this.Response)
            return
        }

        groupName := request["group-name"].(string)

        var face_id int
        face := GetModel("faces")
        face.LoadModelData(map[string]interface{}{"user_id": user_id})
        db.QueryInsert_(face, "RETURNING id").Scan(&face_id)

        var group_id int
        group := GetModel("groups")
        group.LoadModelData(map[string]interface{}{"name": groupName, "face_id": face_id})
        db.QueryInsert_(group, "RETURNING id").Scan(&group_id)

        group_reg := GetModel("group_registrations")
        group_reg.LoadModelData(map[string]interface{}{"event_id": event_id, "group_id": group_id})
        db.QueryInsert_(group_reg, "").Scan()

        query := `SELECT param_values.value
            FROM reg_param_vals
            INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
            INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
            INNER JOIN params ON params.id = param_values.param_id
            INNER JOIN events ON events.id = registrations.event_id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN users ON users.id = faces.user_id
            WHERE params.id in (5, 6, 7) AND users.id = $1 AND events.id = 2 ORDER BY params.id;`

        data := db.Query(query, []interface{}{user_id})
        headName := data[0].(map[string]interface{})["value"].(string)
        headName += " " + data[1].(map[string]interface{})["value"].(string)
        headName += " " + data[2].(map[string]interface{})["value"].(string)

        if request["persons"] != nil {
            for _, v := range request["persons"].([]interface{}) {
                to := v.(map[string]interface{})["name"].(string)
                address := v.(map[string]interface{})["email"].(string)

                token := utils.GetRandSeq(HASH_SIZE)
                eventUrl := "/handler/getrequest/event/"+strconv.Itoa(event_id)+"/"+token

                mailer.InviteToGroup(to, address, eventName, eventUrl, headName, groupName)

                person := GetModel("persons")
                person.LoadModelData(map[string]interface{}{"group_id": group_id, "token": token, "name": to})
                db.QueryInsert_(person, "").Scan()
            }
        }

        utils.SendJSReply(map[string]interface{}{"result": "Группа создана."}, this.Response)
    }
}

func (this *Handler) GetGroups() {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    query := `SELECT groups.id, groups.name FROM groups
        INNER JOIN faces ON faces.id = groups.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1;`

    result := db.Query(query, []interface{}{user_id})
    utils.SendJSReply(map[string]interface{}{"result": "ok", "data": result}, this.Response)
}

func (this *Handler) GetGroup() {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    // request, err := utils.ParseJS(this.Request, this.Response)
    // if err != nil {
    //     utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    // } else {
        // group_id := int(request["group_id"].(float64))

        // query := `SELECT groups.face_id, group_registrations.event_id FROM group_registrations
        //     INNER JOIN groups ON groups.id = group_registrations.group_id
        //     INNER JOIN events ON events.id = group_registrations.event_id
        //     INNER JOIN faces ON faces.id = groups.face_id
        //     INNER JOIN users ON users.id = faces.user_id
        //     WHERE users.id = $1 AND groups.id = $2;`

        // data := db.Query(query, []interface{}{user_id, group_id})

        // if len(data) < 1 {
        //     utils.SendJSReply(map[string]interface{}{"result": "Такой группы вы не создавали."}, this.Response)
        //     return
        // }

        // face_id := data[0].(map[string]interface{})["face_id"].(int)
        // event_id := data[0].(map[string]interface{})["event_id"].(int)

        model := GetModel("groups")
        refFields, refData := GetModelRefDate(model)

        ans := map[string]interface{}{
            "RefData":   refData,
            "RefFields": refFields,
            "TableName": model.GetTableName(),
            "ColNames":  model.GetColNames(),
            "Columns":   model.GetColumns(),
            "Caption":   model.GetCaption(),
            "Sub":       model.GetSub()}

        utils.SendJSReply(map[string]interface{}{"result": "ok", "data": ans}, this.Response)
    // }
}

func (this *Handler) GroupLoad(event_id, group_id string) {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    limit, err := strconv.Atoi(this.Request.PostFormValue("rows"))
    if utils.HandleErr("[GridHandler::Load]  limit Atoi: ", err, this.Response) {
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if utils.HandleErr("[GridHandler::Load] page Atoi: ", err, this.Response) {
        return
    }

    sidx := this.Request.FormValue("sidx")
    start := limit*page - limit

    query := `SELECT groups.id, groups.name FROM group_registrations
        INNER JOIN groups ON groups.id = group_registrations.group_id
        INNER JOIN events ON events.id = group_registrations.event_id
        INNER JOIN faces ON faces.id = groups.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 ORDER BY $2 LIMIT $3 OFFSET $4;`

    rows := db.Query(query, []interface{}{user_id, sidx, limit, start})

    if len(rows) < 1 {
        utils.SendJSReply(map[string]interface{}{"result": "Эмм"}, this.Response)
        return
    }

    count := len(rows)

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
