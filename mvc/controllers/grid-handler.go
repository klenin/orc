package controllers

import (
    "encoding/json"
    "fmt"
    "github.com/orc/db"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "html/template"
    "strconv"
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

    var request map[string]string
    decoder := json.NewDecoder(this.Request.Body)
    err := decoder.Decode(&request)
    utils.HandleErr("[GridHandler::GetSubTable] Decode :", err, this.Response)

    model := GetModel(request["table"])
    index, _ := strconv.Atoi(request["index"])
    subModel := GetModel(model.GetSubTable(index))
    result := db.Select(model.GetSubTable(index), []string{model.GetSubField(), request["id"]}, "", subModel.GetColumns())
    refFields, refData := GetModelRefDate(subModel)

    response, err := json.Marshal(map[string]interface{}{
        "data":      result,
        "name":      subModel.GetTableName(),
        "caption":   subModel.GetCaption(),
        "colnames":  subModel.GetColNames(),
        "columns":   subModel.GetColumns(),
        "reffields": refFields,
        "refdata":   refData})
    utils.HandleErr("[GridHandler::GetSubTable] Marshal: ", err, this.Response)

    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *GridHandler) Load(tableName string) {
    if flag := sessions.CheackSession(this.Response, this.Request); !flag {
        return
    }

    model := GetModel(tableName)
    response, err := json.Marshal(db.Select(tableName, nil, "", model.GetColumns()))
    utils.HandleErr("[GridHandler::Load] Marshal: ", err, this.Response)

    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *GridHandler) Select(tableName string) {
    if flag := sessions.CheackSession(this.Response, this.Request); !flag {
        return
    }

    tmp, err := template.ParseFiles(
        "mvc/views/table.html",
        "mvc/views/header.html",
        "mvc/views/footer.html")
    utils.HandleErr("[GridHandler::Select] ParseFiles: ", err, this.Response)

    model := GetModel(tableName)
    refFields, refData := GetModelRefDate(model)

    err = tmp.ExecuteTemplate(this.Response, "table", Model{
        RefData:   refData,
        RefFields: refFields,
        TableName: model.GetTableName(),
        ColNames:  model.GetColNames(),
        Columns:   model.GetColumns(),
        Caption:   model.GetCaption(),
        Sub:       model.GetSub()})
    utils.HandleErr("[GridHandler::Select] ExecuteTemplate: ", err, this.Response)
}

func (this *GridHandler) Edit(tableName string) {
    if flag := sessions.CheackSession(this.Response, this.Request); !flag {
        return
    }

    model := GetModel(tableName)
    if model == nil {
        return
    }

    params := make(map[string]interface{}, len(model.GetColumns()))
    for i := 0; i < len(model.GetColumns()); i++ {
        params[model.GetColumnByIdx(i)] = this.Request.PostFormValue(model.GetColumnByIdx(i))
    }
    model.LoadModelData(params)

    oper := this.Request.PostFormValue("oper")
    switch oper {
    case "edit":
        db.QueryUpdate_(model)
        break
    case "add":
        db.QueryInsert_(model, "")
        break
    case "del":
        db.QueryDeleteByIds(tableName, this.Request.PostFormValue("id"))
        break
    }
}
