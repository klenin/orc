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

    var model Model
    if role == "admin" {
        model = Model{Columns: db.Tables, ColNames: db.TableNames}
    } else if role == "user" {
        var face_id int
        face := GetModel("faces")
        face.LoadWherePart(map[string]interface{}{"user_id": user_id})
        err = db.SelectRow(face, []string{"id"}).Scan(&face_id)
        if err != nil {
            utils.HandleErr("[Handle::ShowCabinet]: ", err, this.Response)
            return
        }

        query := `SELECT params.name, param_values.value FROM param_values
            INNER JOIN params ON params.id = param_values.param_id
            INNER JOIN reg_param_vals ON reg_param_vals.param_val_id = param_values.id
            INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
            INNER JOIN events ON events.id = reg_param_vals.event_id
            INNER JOIN events_regs ON events_regs.event_id = events.id AND events_regs.reg_id = registrations.id
            INNER JOIN faces ON faces.id = registrations.face_id
            WHERE events.id=$1 AND faces.id=$2 ORDER BY params.id`

        regParamVals := GetModel("reg_param_vals")
        data := db.Query(query, []interface{}{1, face_id})
        model = Model{Table: data, Columns: regParamVals.GetColumns(), ColNames: regParamVals.GetColNames()}
    }

    this.Render([]string{"mvc/views/"+role+".html"}, role, model)
}
