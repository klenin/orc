package controllers

import (
    "encoding/json"
    "fmt"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "html/template"
    "strconv"
    "strings"
)

func (c *BaseController) GridHandler() *GridHandler {
    return new(GridHandler)
}

type GridHandler struct {
    Controller
}

func (this *GridHandler) GetSubTable() {
    if flag := sessions.CheackSession(this.Response, this.Request); !flag {
        return
    }
    this.Response.Header().Set("Access-Control-Allow-Origin", "*")
    this.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    this.Response.Header().Set("Content-type", "application/json")

    var request map[string]string
    decoder := json.NewDecoder(this.Request.Body)
    err := decoder.Decode(&request)
    utils.HandleErr("[GridHandler::GetSubTable] Decode :", err, this.Response)

    id := request["id"]
    tableName := request["table"]
    model := GetModel(tableName)
    index, _ := strconv.Atoi(request["index"])

    subTableName := model.GetSubTable(index)
    subModel := GetModel(subTableName)

    result, refdata := subModel.Select([]string{model.GetSubField(), id}, "", subModel.GetColumns())
    response, err := json.Marshal(map[string]interface{}{
        "data":      result,
        "name":      subModel.GetTableName(),
        "caption":   subModel.GetCaption(),
        "colnames":  subModel.GetColNames(),
        "columns":   subModel.GetColumns(),
        "refdata":   refdata,
        "reffields": subModel.GetRefFields()})
    utils.HandleErr("[GridHandler::GetSubTable] json.Marshal: ", err, nil)
    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *GridHandler) Load(tableName string) {
    if flag := sessions.CheackSession(this.Response, this.Request); !flag {
        return
    }
    this.Response.Header().Set("Access-Control-Allow-Origin", "*")
    this.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    this.Response.Header().Set("Content-type", "application/json")

    model := GetModel(tableName)
    answer, _ := model.Select(nil, "", model.GetColumns())

    response, err := json.Marshal(answer)
    utils.HandleErr("[GridHandler::Load] json.Marshal: ", err, nil)
    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *GridHandler) Select(tableName string) {
    if flag := sessions.CheackSession(this.Response, this.Request); !flag {
        return
    }
    model := GetModel(tableName)
    _, refdata := model.Select(nil, "", model.GetColumns())
    tmp, err := template.ParseFiles(
        "mvc/views/table.html",
        "mvc/views/header.html",
        "mvc/views/footer.html")
    utils.HandleErr("[GridHandler::Select] template.ParseFiles: ", err, nil)
    err = tmp.ExecuteTemplate(this.Response, "table", Model{
        RefData:   refdata,
        RefFields: model.GetRefFields(),
        TableName: model.GetTableName(),
        ColNames:  model.GetColNames(),
        Columns:   model.GetColumns(),
        Caption:   model.GetCaption(),
        Sub:       model.GetSub()})
    utils.HandleErr("[GridHandler::Select] tmp.Execute: ", err, nil)
}

func (this *GridHandler) Edit(tableName string) {
    if flag := sessions.CheackSession(this.Response, this.Request); !flag {
        return
    }
    this.Response.Header().Set("Access-Control-Allow-Origin", "*")
    this.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    this.Response.Header().Set("Content-type", "application/json")

    var i int
    oper := this.Request.FormValue("oper")
    model := GetModel(tableName)
    params := make([]interface{}, len(model.GetColumns())-1)

    for i = 0; i < len(model.GetColumns())-1; i++ {
        if model.GetColumnByIdx(i+1) == "date" {
            params[i] = this.Request.FormValue(model.GetColumnByIdx(i + 1))[0:10]
        } else {
            params[i] = this.Request.FormValue(model.GetColumnByIdx(i + 1))
        }
    }

    switch oper {
    case "edit":
        params = append(params, this.Request.FormValue("id"))
        model.Update(model.GetColumnSlice(1), params, "id=$"+strconv.Itoa(i+1))
        break
    case "add":
        model.Insert(model.GetColumnSlice(1), params)
        break
    case "del":
        ids := strings.Split(this.Request.FormValue("id"), ",")
        tmp := make([]interface{}, len(ids))
        for i, v := range ids {
            tmp[i] = interface{}(v)
        }
        model.Delete("id", tmp)
        break
    }
}
