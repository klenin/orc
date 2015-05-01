package controllers

import (
    "database/sql"
    "github.com/orc/db"
    "github.com/orc/sessions"
    "github.com/orc/mvc/models"
    "github.com/orc/utils"
    "net/http"
    "strconv"
    "github.com/orc/mailer"
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
    case "sendEmailWellcomeToProfile":
        user_id, err := strconv.Atoi(data["user_id"].(string))
        if err != nil {
            utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
            return
        }
        query := `SELECT param_values.value
            FROM reg_param_vals
            INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
            INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
            INNER JOIN params ON params.id = param_values.param_id
            INNER JOIN events ON events.id = registrations.event_id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN users ON users.id = faces.user_id
            WHERE params.id in (4, 5, 6, 7) AND users.id = $1 ORDER BY params.id;`

        data := db.Query(query, []interface{}{user_id})

        if len(data) < 4 {
            utils.SendJSReply(map[string]interface{}{"result": "Нет регистрационных данных пользователя."}, this.Response)
            break
        }

        to := data[1].(map[string]interface{})["value"].(string)+" "
        to += data[2].(map[string]interface{})["value"].(string)+" "
        to += data[3].(map[string]interface{})["value"].(string)
        email := data[0].(map[string]interface{})["value"].(string)

        token := utils.GetRandSeq(HASH_SIZE)
        if !mailer.SendEmailWellcomeToProfile(to, email, token) {
            utils.SendJSReply(map[string]interface{}{"result": "Проверьте правильность введенного Вами email."}, this.Response)
            break
        }
        user := GetModel("users")
        user.LoadModelData(map[string]interface{}{"token": token})
        user.GetFields().(*models.User).Enabled = true
        user.LoadWherePart(map[string]interface{}{"id": user_id})
        db.QueryUpdate_(user).Scan()
        utils.SendJSReply(map[string]interface{}{"result": "ok"}, this.Response)
        break
    }
}

func (this *Handler) ShowCabinet() {
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

        query = `SELECT params.name, param_values.value
            FROM reg_param_vals
            INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
            INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
            INNER JOIN params ON params.id = param_values.param_id
            INNER JOIN events ON events.id = registrations.event_id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN users ON users.id = faces.user_id
            WHERE params.id in (1, 4) AND users.id = $1 ORDER BY params.id;`

        data := db.Query(query, []interface{}{user_id})

        this.Render(
            []string{"mvc/views/"+role+".html"},
            role,
            map[string]interface{}{"group": groupsModel, "reg": regsModel, "groupreg": groupRegsModel, "userData": data})
    }
}

func (this *Handler) ConfirmInvitationToGroup(token string) {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    var face_id int
    // face := GetModel("faces")
    // face.LoadModelData(map[string]interface{}{"user_id": user_id})
    // db.QueryInsert_(face, "RETURNING id").Scan(&face_id)

    query := `SELECT faces.id
        FROM registrations
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN users ON faces.user_id = users.id
        WHERE users.id = $1 AND events.id = $2;`

    db.QueryRow(query, []interface{}{user_id, 1}).Scan(&face_id)

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

func (this *Handler) WellcomeToProfile(token string) {
    user := GetModel("users")
    user.LoadWherePart(map[string]interface{}{"token": token})

    var id int
    err := db.SelectRow(user, []string{"id"}).Scan(&id)
    if utils.HandleErr("[Handle::WellcomeToProfile]: ", err, this.Response) && id != 0 {
        return
    }

    user = GetModel("users")
    user.LoadModelData(map[string]interface{}{"token": " ", "hash": utils.GetRandSeq(HASH_SIZE)})
    user.GetFields().(*models.User).Enabled = true
    user.LoadWherePart(map[string]interface{}{"id": id})
    db.QueryUpdate_(user).Scan()
    hash := utils.GetRandSeq(HASH_SIZE)
    sessions.SetSession(this.Response, map[string]interface{}{"id": id, "hash": hash})

    this.ShowCabinet()
}

func (this *Handler) Login(user_id string) {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        http.Redirect(this.Response, this.Request, "/", http.StatusForbidden)
        return
    }

    id, err := strconv.Atoi(user_id)
    if utils.HandleErr("[GridHandler::login] user_id Atoi: ", err, this.Response) {
        return
    }

    if !db.IsExists_("users", []string{"id"}, []interface{}{id}) {
        http.Error(this.Response, "Have not such user with the id", http.StatusInternalServerError)
        return
    }

    hash := utils.GetRandSeq(HASH_SIZE)

    user := GetModel("users")
    user.LoadModelData(map[string]interface{}{"hash": hash})
    user.GetFields().(*models.User).Enabled = true
    user.LoadWherePart(map[string]interface{}{"id": id})
    err = db.QueryUpdate_(user).Scan()

    sessions.SetSession(this.Response, map[string]interface{}{"id": id, "hash": hash})

    this.ShowCabinet()
}
