package controllers

import (
    "database/sql"
    "github.com/orc/db"
    "github.com/orc/mailer"
    "github.com/orc/mvc/models"
    "github.com/orc/sessions"
    "github.com/orc/utils"
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

func (this *Handler) GetList() {
    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    } else {
        fields := request["fields"].([]interface{})
        result := db.Select(this.GetModel(request["table"].(string)), utils.ArrayInterfaceToString(fields))
        utils.SendJSReply(map[string]interface{}{"result": "ok", "data": result}, this.Response)
    }
}

func (this *Handler) Index() {
    var response interface{}

    data, err := utils.ParseJS(this.Request, this.Response)
    if utils.HandleErr("[Handle::Index]: ", err, this.Response) {
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

        sid := sessions.GetValue("sid", this.Request)
        if sid == nil {
            result = map[string]interface{}{"result": "no"}

        } else {
            user := this.GetModel("users")
            user.LoadWherePart(map[string]interface{}{"sid": sid})
            err := db.SelectRow(user, []string{"sid"}).Scan(&userHash)
            if err != sql.ErrNoRows {
                result = map[string]interface{}{"result": "ok"}
            } else {
                result = map[string]interface{}{"result": "no"}
            }
        }

        utils.SendJSReply(result, this.Response)
        break

    case "sendEmailWellcomeToProfile":
        userId, err := strconv.Atoi(data["user_id"].(string))
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
        data := db.Query(query, []interface{}{userId})

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
            utils.SendJSReply(map[string]interface{}{"result": "Проверьте правильность email."}, this.Response)
            break
        }

        user := this.GetModel("users")
        user.GetFields().(*models.User).Token = token
        user.GetFields().(*models.User).Enabled = true
        user.LoadWherePart(map[string]interface{}{"id": userId})
        db.QueryUpdate(user).Scan()

        utils.SendJSReply(map[string]interface{}{"result": "Письмо отправлено"}, this.Response)
        break
    }
}

func (this *Handler) ShowCabinet() {
    userId, err := this.CheckSid()
    if err != nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    user := this.GetModel("users")
    user.LoadWherePart(map[string]interface{}{"id": userId})

    var role string
    err = db.SelectRow(user, []string{"role"}).Scan(&role)
    if err != nil {
        utils.HandleErr("[Handle::ShowCabinet]: ", err, this.Response)
        return
    }

    if role == "admin" {
        model := Model{Columns: db.Tables, ColNames: db.TableNames}
        this.Render([]string{"mvc/views/"+role+".html"}, role, model)

    } else {
        groups := this.GetModel("groups")
        persons := this.GetModel("persons")
        groupsModel := Model{
            TableName:    groups.GetTableName(),
            ColNames:     groups.GetColNames(),
            ColModel:     groups.GetColModel(false, userId),
            Caption:      groups.GetCaption(),
            Sub:          groups.GetSub(),
            SubTableName: persons.GetTableName(),
            SubCaption:   persons.GetCaption(),
            SubColModel:  persons.GetColModel(false, userId),
            SubColNames:  persons.GetColNames()}

        regs := this.GetModel("registrations")
        regsModel := Model{
            TableName: regs.GetTableName(),
            ColNames:  regs.GetColNames(),
            ColModel:  regs.GetColModel(false, userId),
            Caption:   regs.GetCaption()}

        groupRegs := this.GetModel("group_registrations")
        groupRegsModel := Model{
            TableName:    groupRegs.GetTableName(),
            ColNames:     groupRegs.GetColNames(),
            ColModel:     groupRegs.GetColModel(false, userId),
            Caption:      groupRegs.GetCaption(),
            Sub:          groupRegs.GetSub(),
            SubTableName: persons.GetTableName(),
            SubCaption:   persons.GetCaption(),
            SubColModel:  persons.GetColModel(false, userId),
            SubColNames:  persons.GetColNames()}

        query := `SELECT params.name, param_values.value, users.login
            FROM reg_param_vals
            INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
            INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
            INNER JOIN params ON params.id = param_values.param_id
            INNER JOIN events ON events.id = registrations.event_id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN users ON users.id = faces.user_id
            WHERE params.id = 4 AND users.id = $1;`

        data := db.Query(query, []interface{}{userId})

        faces := this.GetModel("faces")
        facesModel := Model{
            ColModel:     faces.GetColModel(false, userId),
            TableName:    faces.GetTableName(),
            ColNames:     faces.GetColNames(),
            Caption:      faces.GetCaption()}

        params := this.GetModel("param_values")
        paramsModel := Model{
            ColModel:  params.GetColModel(false, userId),
            TableName: params.GetTableName(),
            ColNames:  params.GetColNames(),
            Caption:   params.GetCaption()}

        events := this.GetModel("events")
        eventsModel := Model{
            ColModel:  events.GetColModel(false, userId),
            TableName: events.GetTableName(),
            ColNames:  events.GetColNames(),
            Caption:   events.GetCaption()}

        this.Render(
            []string{"mvc/views/"+role+".html"},
            role,
            map[string]interface{}{
                "group": groupsModel,
                "reg": regsModel,
                "groupreg": groupRegsModel,
                "faces": facesModel,
                "params": paramsModel,
                "events": eventsModel,
                "userData": data})
    }
}

func WellcomeToProfile(w http.ResponseWriter, r *http.Request) {
    newContreoller := new(BaseController).Handler()
    newContreoller.Request = r
    newContreoller.Response = w

    parts := strings.Split(r.URL.Path, "/")
    token := parts[len(parts)-1]

    user := newContreoller.GetModel("users")
    user.LoadWherePart(map[string]interface{}{"token": token})

    var id int
    err := db.SelectRow(user, []string{"id"}).Scan(&id)
    if utils.HandleErr("[WellcomeToProfile]: ", err, newContreoller.Response) || id == 0 {
        return
    }

    sid := utils.GetRandSeq(HASH_SIZE)
    user = newContreoller.GetModel("users")
    user.GetFields().(*models.User).Sid = sid
    user.GetFields().(*models.User).Enabled = true
    user.LoadWherePart(map[string]interface{}{"id": id})
    db.QueryUpdate(user).Scan()

    sessions.SetSession(newContreoller.Response, map[string]interface{}{"sid": sid})

    http.Redirect(newContreoller.Response, newContreoller.Request, "/handler/showcabinet/users", 200)
}

func (this *Handler) Login(userId string) {
    if !sessions.CheckSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        http.Redirect(this.Response, this.Request, "/", http.StatusForbidden)
        return
    }

    id, err := strconv.Atoi(userId)
    if utils.HandleErr("[GridHandler::login] user_id Atoi: ", err, this.Response) {
        return
    }

    if !db.IsExists("users", []string{"id"}, []interface{}{id}) {
        http.Error(this.Response, "Have not such user with the id", http.StatusInternalServerError)
        return
    }

    sid := utils.GetRandSeq(HASH_SIZE)

    user := this.GetModel("users")
    user.GetFields().(*models.User).Sid = sid
    user.GetFields().(*models.User).Enabled = true
    user.LoadWherePart(map[string]interface{}{"id": id})
    db.QueryUpdate(user).Scan()

    sessions.SetSession(this.Response, map[string]interface{}{"sid": sid})

    http.Redirect(this.Response, this.Request, "/handler/showcabinet/users", 200)
}

func (this *Handler) CheckEnableOfUser(id string) {
    eventId, err := strconv.Atoi(id)
    if utils.HandleErr("[Handler::CheckEnableOfUser] event_id Atoi: ", err, this.Response) {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    if eventId == 1 {
        if sessions.CheckSession(this.Response, this.Request) {
            utils.SendJSReply(map[string]interface{}{"result": "authorized"}, this.Response)
            return
        }
        utils.SendJSReply(map[string]interface{}{"result": "ok"}, this.Response)
        return
    }

    userId, err := this.CheckSid()
    if err != nil && eventId != 1 {
        utils.SendJSReply(map[string]interface{}{"result": "Unauthorized"}, this.Response)
        return
    }

    query := `SELECT registrations.id
        FROM registrations
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 AND events.id = $2;`

    var regId int
    err = db.QueryRow(query, []interface{}{userId, eventId}).Scan(&regId)
    if err != sql.ErrNoRows {
        utils.SendJSReply(map[string]interface{}{"result": "regExists", "regId": strconv.Itoa(regId)}, this.Response)
    } else {
        utils.SendJSReply(map[string]interface{}{"result": "ok"}, this.Response)
    }
}
