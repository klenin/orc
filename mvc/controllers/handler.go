package controllers

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/orc/db"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "html/template"
    "net/http"
)

func (c *BaseController) Handler() *Handler {
    return new(Handler)
}

type Handler struct {
    Controller
}

func (this *Handler) GetEventList() {
    var request map[string]interface{}
    decoder := json.NewDecoder(this.Request.Body)
    err := decoder.Decode(&request)
    if utils.HandleErr("[Handler::GetEventList] Decode :", err, this.Response) {
        return
    }

    fields := request["fields"].([]interface{})
    result := db.Select(GetModel(request["table"].(string)), utils.ArrayInterfaceToString(fields), "")

    response, err := json.Marshal(map[string]interface{}{"data": result})
    if utils.HandleErr("[Handle::GetEventList] Marshal: ", err, this.Response) {
        return
    }

    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *Handler) Index() {
    var data map[string]interface{}
    var response string

    decoder := json.NewDecoder(this.Request.Body)
    err := decoder.Decode(&data)
    if utils.HandleErr("[Handler::Index] Decode :", err, this.Response) {
        return
    }

    switch data["action"] {
    case "login":
        response = this.HandleLogin(data["login"].(string), data["password"].(string))
        fmt.Fprintf(this.Response, "%s", response)
        break

    case "logout":
        response = this.HandleLogout()
        fmt.Fprintf(this.Response, "%s", response)
        break

    case "editProfile":
        params := make(map[string]interface{}, 0)

        for _, element := range data["data"].([]interface{}) {
            elem := element.(map[string]interface{})
            params[elem["name"].(string)] = elem["value"]
        }

        model := GetModel(data["table"].(string))
        model.LoadModelData(params)
        model.LoadWherePart(map[string]interface{}{"id":  int(data["id"].(float64))})
        db.QueryUpdate_(model, "")

        response, err := json.Marshal(map[string]interface{}{"result": "ok"})
        if utils.HandleErr("[Handle::Index] Marshal: ", err, this.Response) {
            return
        }

        fmt.Fprintf(this.Response, "%s", string(response))
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
            err := db.SelectRow(user, []string{"hash"}, "").Scan(&userHash)
            if err != sql.ErrNoRows {
                result = map[string]interface{}{"result": "ok"}
            } else {
                result = map[string]interface{}{"result": "no"}
            }
        }

        response, err := json.Marshal(result)
        if utils.HandleErr("[Handle::Index] Marshal: ", err, this.Response) {
            return
        }

        fmt.Fprintf(this.Response, "%s", string(response))
        break
    }
}

func (this *Handler) ShowCabinet(tableName string) {
    if flag := sessions.CheackSession(this.Response, this.Request); !flag {
        return
    }

    user_id := sessions.GetValue("id", this.Request)
    if user_id == nil {
        http.Redirect(this.Response, this.Request, "/", 401)
        return
    }
    user := GetModel("users")
    user.LoadWherePart(map[string]interface{}{"id": user_id})

    var role string
    err := db.SelectRow(user, []string{"role"}, "").Scan(&role)
    if err != nil {
        panic("ShowCabinet: " + err.Error())
    }

    var model Model
    if role == "admin" {
        model = Model{Columns: db.Tables, ColNames: db.TableNames}
    } else if role == "user" {
        var face_id int
        face := GetModel("faces")
        face.LoadWherePart(map[string]interface{}{"user_id": user_id})
        err = db.SelectRow(face, []string{"id"}, "").Scan(&face_id)
        if err != nil {
            response, err := json.Marshal(map[string]interface{}{"result": err.Error()})
            if utils.HandleErr("[Handle::SaveUserRequest] Marshal: ", err, this.Response) {
                return
            }
            fmt.Fprintf(this.Response, "%s", string(response))
            return
        }
        var reg_id int
        reg := GetModel("registrations")
        reg.LoadWherePart(map[string]interface{}{"face_id": face_id})
        err = db.SelectRow(reg, []string{"id"}, "").Scan(&reg_id)
        if err != nil {
            response, err := json.Marshal(map[string]interface{}{"result": err.Error()})
            if utils.HandleErr("[Handle::SaveUserRequest] Marshal: ", err, this.Response) {
                return
            }
            fmt.Fprintf(this.Response, "%s", string(response))
            return
        }

        query := `SELECT params.name, param_values.value from param_values
            inner join params on params.id = param_values.param_id
            inner join reg_param_vals on reg_param_vals.param_val_id = param_values.id
            inner join registrations on registrations.id = reg_param_vals.reg_id
            inner join events on events.id = reg_param_vals.event_id WHERE registrations.id=$1`

        regParamVals := GetModel("reg_param_vals")
        data := db.Query(query, []interface{}{reg_id})
        model = Model{Table: data, Columns: regParamVals.GetColumns(), ColNames: regParamVals.GetColNames()}
    }

    tmp, err := template.ParseFiles(
        "mvc/views/"+role+".html",
        "mvc/views/header.html",
        "mvc/views/footer.html")
    if utils.HandleErr("[Handler::ShowCabinet] ParseFiles: ", err, this.Response) {
        return
    }

    err = tmp.ExecuteTemplate(this.Response, role, model)
    utils.HandleErr("[Handler::ShowCabinet] ExecuteTemplate: ", err, this.Response)
}
