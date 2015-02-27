package controllers

import (
    "encoding/json"
    "fmt"
    "github.com/orc/db"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "html/template"
    "math"
    "net/http"
    "strconv"
)

func (c *BaseController) GridHandler() *GridHandler {
    return new(GridHandler)
}

type GridHandler struct {
    Controller
}

func (this *GridHandler) GetSubTable() {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", 401)
        return
    }

    if !this.isAdmin() {
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(err.Error(), this.Response)
        return
    }

    model := GetModel(request["table"].(string))
    index, _ := strconv.Atoi(request["index"].(string))
    subModel := GetModel(model.GetSubTable(index))
    subModel.LoadWherePart(map[string]interface{}{model.GetSubField(): request["id"]})
    result := db.Select(subModel, subModel.GetColumns(), "")
    refFields, refData := GetModelRefDate(subModel)

    response, err := json.Marshal(map[string]interface{}{
        "data":      result,
        "name":      subModel.GetTableName(),
        "caption":   subModel.GetCaption(),
        "colnames":  subModel.GetColNames(),
        "columns":   subModel.GetColumns(),
        "reffields": refFields,
        "refdata":   refData})
    if utils.HandleErr("[GridHandler::GetSubTable] Marshal: ", err, this.Response) {
        return
    }

    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *GridHandler) Load(tableName string) {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", 401)
        return
    }

    if !this.isAdmin() {
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

    model := GetModel(tableName)
    model.SetOrder(sidx)
    model.SetLimit(limit)
    model.SetOffset(start)

    rows := db.Select(model, model.GetColumns(), "")
    count := db.SelectCount(tableName)

    var totalPages int
    if count > 0 {
        totalPages = int(math.Ceil(float64(count) / float64(limit)))
    } else {
        totalPages = 0
    }

    result := make(map[string]interface{}, 4)
    result["rows"] = rows
    result["page"] = page
    result["total"] = totalPages
    result["records"] = count

    utils.SendJSReply(result, this.Response)
}

func (this *GridHandler) Select(tableName string) {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", 401)
        return
    }

    if !this.isAdmin() {
        return
    }

    tmp, err := template.ParseFiles(
        "mvc/views/table.html",
        "mvc/views/header.html",
        "mvc/views/footer.html")
    if utils.HandleErr("[GridHandler::Select] ParseFiles: ", err, this.Response) {
        return
    }

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
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", 401)
        return
    }

    if !this.isAdmin() {
        return
    }

    model := GetModel(tableName)
    if model == nil {
        utils.HandleErr("[Grid-Handler::Edit] GetModel: invalid model", nil, this.Response)
        return
    }

    params := make(map[string]interface{}, len(model.GetColumns()))
    for i := 0; i < len(model.GetColumns()); i++ {
        params[model.GetColumnByIdx(i)] = this.Request.PostFormValue(model.GetColumnByIdx(i))
    }

    oper := this.Request.PostFormValue("oper")
    switch oper {
    case "edit":
        id, err := strconv.Atoi(this.Request.PostFormValue("id"))
        if utils.HandleErr("[Grid-Handler::Edit] strconv.Atoi: ", err, this.Response) {
            return
        }
        model.LoadModelData(params)
        model.LoadWherePart(map[string]interface{}{"id": id})
        db.QueryUpdate_(model, "")
        break
    case "add":
        model.LoadModelData(params)
        db.QueryInsert_(model, "")
        break
    case "del":
        db.QueryDeleteByIds(tableName, this.Request.PostFormValue("id"))
        break
    }
}

func (this *GridHandler) ResetPassword() {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", 401)
        return
    }

    if !this.isAdmin() {
        return
    }

    this.Response.Header().Set("Access-Control-Allow-Origin", "*")
    this.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    this.Response.Header().Set("Content-type", "application/json")

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(err.Error(), this.Response)
        return
    }

    pass := request["pass"].(string)

    id, err :=  strconv.Atoi(request["id"].(string))
    if utils.HandleErr("[Grid-Handler::ResetPassword] strconv.Atoi: ", err, this.Response) {
        return
    }

    user := GetModel("users")
    user.LoadWherePart(map[string]interface{}{"id": id})

    var salt string
    db.SelectRow(user, []string{"salt"}, "").Scan(&salt)

    user.LoadModelData(map[string]interface{}{"pass": utils.GetMD5Hash(pass + salt)})
    db.QueryUpdate_(user, "")

    utils.SendJSReply(map[string]interface{}{"result": "ok"}, this.Response)
}

func (this *GridHandler) isAdmin() bool {
    var role string

    user_id := sessions.GetValue("id", this.Request)
    if user_id == nil {
        http.Redirect(this.Response, this.Request, "/", 401)
        return false
    }

    user := GetModel("users")
    user.LoadWherePart(map[string]interface{}{"id": user_id})
    err := db.SelectRow(user, []string{"role"}, "").Scan(&role)
    if err != nil || role == "user" {
        http.Redirect(this.Response, this.Request, "/", 403)
        return false
    }

    return role == "admin"
}
