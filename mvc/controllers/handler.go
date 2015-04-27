package controllers

import (
    "database/sql"
    "github.com/orc/db"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "net/http"
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

    if role == "admin" {
        model := Model{Columns: db.Tables, ColNames: db.TableNames}
        this.Render([]string{"mvc/views/"+role+".html"}, role, model)
    } else {
        groups := GetModel("groups")
        groupsRefFields, groupsRefData := groups.GetModelRefDate()
        persons := GetModel("persons")

        query := `SELECT groups.id, groups.name FROM groups
            INNER JOIN faces ON faces.id = groups.face_id
            INNER JOIN users ON users.id = faces.user_id
            WHERE users.id = $1 ORDER BY groups.id;`
        personsRefData := map[string]interface{}{"group_id": db.Query(query, []interface{}{user_id})}

        personsRefFields := []string{"name"}

        groupsModel := Model{
            RefData:      groupsRefData,
            RefFields:    groupsRefFields,
            TableName:    groups.GetTableName(),
            ColNames:     groups.GetColNames(),
            Columns:      groups.GetColumns(),
            Caption:      groups.GetCaption(),
            Sub:          groups.GetSub(),
            SubTableName: persons.GetTableName(),
            SubCaption:   persons.GetCaption(),
            SubRefData:   personsRefData,
            SubRefFields: personsRefFields,
            SubColumns:   persons.GetColumns()[:len(persons.GetColumns())-1],
            SubColNames:  persons.GetColNames()[:len(persons.GetColNames())-1]}

        regs := GetModel("registrations")
        regsRefFields, regsRefData := regs.GetModelRefDate()

        regsModel := Model{
            RefData:   regsRefData,
            RefFields: regsRefFields,
            TableName: regs.GetTableName(),
            ColNames:  regs.GetColNames(),
            Columns:   regs.GetColumns(),
            Caption:   regs.GetCaption(),
            Sub:       regs.GetSub()}

        groupRegs := GetModel("group_registrations")
        groupRegsRefFields, groupRegsRefData := groupRegs.GetModelRefDate()

        groupRegsModel := Model{
            RefData:   groupRegsRefData,
            RefFields: groupRegsRefFields,
            TableName: groupRegs.GetTableName(),
            ColNames:  groupRegs.GetColNames(),
            Columns:   groupRegs.GetColumns(),
            Caption:   groupRegs.GetCaption(),
            Sub:          groups.GetSub(),
            SubTableName: persons.GetTableName(),
            SubCaption:   persons.GetCaption(),
            SubRefData:   personsRefData,
            SubRefFields: personsRefFields,
            SubColumns:   persons.GetColumns()[:len(persons.GetColumns())-1],
            SubColNames:  persons.GetColNames()[:len(persons.GetColNames())-1]}

        this.Render(
            []string{"mvc/views/"+role+".html"},
            role,
            map[string]interface{}{"group": groupsModel, "reg": regsModel, "groupreg": groupRegsModel})
    }
}

func (this *Handler) ConfirmInvitationToGroup(token string) {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    var face_id int
    face := GetModel("faces")
    face.LoadModelData(map[string]interface{}{"user_id": user_id})
    db.QueryInsert_(face, "RETURNING id").Scan(&face_id)

    person := GetModel("persons")
    person.LoadModelData(map[string]interface{}{"face_id": face_id, "status": true})
    person.LoadWherePart(map[string]interface{}{"token": token})
    db.QueryUpdate_(person).Scan()

    if this.Response != nil {
        this.Render([]string{"mvc/views/msg.html"}, "msg", "Вы успешно присоединены к группе.")
    }
}

func (this *Handler) RejectInvitationToGroup(token string) {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    query := "DELETE FROM persons WHERE token = $1;"
    db.Exec(query, []interface{}{token})

    if this.Response != nil {
        this.Render([]string{"mvc/views/msg.html"}, "msg", "Запрос о присоединении к группе успешно отклонен.")
    }
}

func (this *Handler) GroupSearchoptions() {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    query := `SELECT groups.id, groups.name FROM groups
        INNER JOIN faces ON faces.id = groups.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 ORDER BY groups.id;`

    rows := db.Query(query, []interface{}{user_id})

    utils.SendJSReply(map[string]interface{}{"data": rows}, this.Response)
}
