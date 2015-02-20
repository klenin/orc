package controllers

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/orc/db"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "html/template"
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
    case "register":
        login, password := data["login"].(string), data["password"].(string)
        fname, lname := data["fname"].(string), data["lname"].(string)
        response = this.HandleRegister(login, password, "user", fname, lname)
        fmt.Fprintf(this.Response, "%s", response)
        break

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

        params["id"] = data["id"].(string)
        for _, element := range data["data"].([]interface{}) {
            elem := element.(map[string]interface{})
            params[elem["name"].(string)] = elem["value"]
        }

        model := GetModel(data["table"].(string))
        model.LoadModelData(params)
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

    id := sessions.GetValue("id", this.Request).(int)
    user := GetModel("users")
    user.LoadWherePart(map[string]interface{}{"id": id})

    var role string
    var person_id int
    err := db.SelectRow(user, []string{"role", "person_id"}, "").Scan(&role, &person_id)
    if err != nil {
        panic("ShowCabinet: " + err.Error())
    }

    var model Model
    if role == "admin" {
        model = Model{Columns: db.Tables, ColNames: db.TableNames}
    } else if role == "user" {
        m := GetModel("persons")
        m.LoadWherePart(map[string]interface{}{"id": person_id})
        data := db.Select(m, m.GetColumns(), "")
        model = Model{Table: data, Columns: m.GetColumns(), ColNames: m.GetColNames()}
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
