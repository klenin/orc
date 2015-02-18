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
    utils.HandleErr("[Handler::GetEventList] Decode :", err, this.Response)

    fields := request["fields"].([]interface{})
    result := db.Select(request["table"].(string), utils.ArrayInterfaceToString(fields), nil, "")

    response, err := json.Marshal(map[string]interface{}{"data": result})
    utils.HandleErr("[Handle::GetEventList] Marshal: ", err, this.Response)

    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *Handler) ResetPassword() {
    if flag := sessions.CheackSession(this.Response, this.Request); !flag {
        return
    }

    this.Response.Header().Set("Access-Control-Allow-Origin", "*")
    this.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    this.Response.Header().Set("Content-type", "application/json")

    var request map[string]interface{}
    decoder := json.NewDecoder(this.Request.Body)
    err := decoder.Decode(&request)
    utils.HandleErr("[Handler::ResetPassword] Decode :", err, this.Response)

    id, pass := request["id"].(int), request["pass"].(string)
    result := db.Select("users", []string{"salt"}, map[string]interface{}{"id": id}, "")
    salt := result[0].(map[string]interface{})["salt"].(string)
    hash := GetMD5Hash(pass + salt)

    user := GetModel("users")
    user.LoadModelData(map[string]interface{}{"id": id, "pass": hash})
    db.QueryUpdate_(user)

    response, err := json.Marshal(map[string]interface{}{"result": "ok"})
    utils.HandleErr("[Handle::ResetPassword] Marshal: ", err, this.Response)

    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *Handler) Index() {
    var data map[string]interface{}
    var response string

    decoder := json.NewDecoder(this.Request.Body)
    err := decoder.Decode(&data)
    utils.HandleErr("[Handler::Index] Decode :", err, this.Response)

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
        db.QueryUpdate_(model)

        response, err := json.Marshal(map[string]interface{}{"result": "ok"})
        utils.HandleErr("[Handle::Index] Marshal: ", err, this.Response)

        fmt.Fprintf(this.Response, "%s", string(response))
        break

    case "checkSession":
        var userHash string
        var result interface{}

        hash := sessions.GetValue("hash", this.Request)

        if hash == nil {
            result = map[string]interface{}{"result": "no"}
        } else {
            query := db.QuerySelect("users", "hash=$1", []string{"hash"})
            err := db.QueryRow(query, []interface{}{hash.(string)}).Scan(&userHash)
            if err != sql.ErrNoRows {
                result = map[string]interface{}{"result": "ok"}
            } else {
                result = map[string]interface{}{"result": "no"}
            }
        }

        response, err := json.Marshal(result)
        utils.HandleErr("[Handle::Index] Marshal: ", err, this.Response)

        fmt.Fprintf(this.Response, "%s", string(response))
        break
    }
}

func (this *Handler) ShowCabinet(tableName string) {
    if flag := sessions.CheackSession(this.Response, this.Request); !flag {
        return
    }

    id := sessions.GetValue("id", this.Request).(int)
    data := db.Select("users", []string{"role", "person_id"}, map[string]interface{}{"id": id}, "")

    role := data[0].(map[string]interface{})["role"].(string)
    person_id := data[0].(map[string]interface{})["person_id"].(int64)

    var model Model
    if role == "admin" {
        model = Model{Columns: db.Tables, ColNames: db.TableNames}
    } else if role == "user" {
        m := GetModel("persons")
        data := db.Select("persons", m.GetColumns(), map[string]interface{}{"id": person_id}, "")
        model = Model{Table: data, Columns: m.GetColumns(), ColNames: m.GetColNames()}
    }

    tmp, err := template.ParseFiles(
        "mvc/views/"+role+".html",
        "mvc/views/header.html",
        "mvc/views/footer.html")
    utils.HandleErr("[Handler::ShowCabinet] ParseFiles: ", err, this.Response)

    err = tmp.ExecuteTemplate(this.Response, role, model)
    utils.HandleErr("[Handler::ShowCabinet] ExecuteTemplate: ", err, this.Response)
}
