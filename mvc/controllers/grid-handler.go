package controllers

import (
    "errors"
    "encoding/csv"
    "encoding/json"
    "fmt"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "net/http"
    "reflect"
    "strconv"
    "strings"
    "time"
)

func (c *BaseController) GridHandler() *GridHandler {
    return new(GridHandler)
}

type GridHandler struct {
    Controller
}

func (this *GridHandler) GetSubTable() {
    userId, err := this.CheckSid()
    if err != nil {
        http.Error(this.Response, "Unauthorized", 400)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        http.Error(this.Response, err.Error(), 400)
        return
    }

    model := this.GetModel(request["table"].(string))
    index, _ := strconv.Atoi(request["index"].(string))
    subModel := this.GetModel(model.GetSubTable(index))
    subModel.LoadWherePart(map[string]interface{}{model.GetSubField(): request["id"]})

    response, err := json.Marshal(map[string]interface{}{
        "name":     subModel.GetTableName(),
        "caption":  subModel.GetCaption(),
        "colnames": subModel.GetColNames(),
        "columns":  subModel.GetColumns(),
        "colmodel": subModel.GetColModel(this.isAdmin(), userId)})
    if utils.HandleErr("[GridHandler::GetSubTable] Marshal: ", err, this.Response) {
        return
    }

    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *GridHandler) CreateGrid(tableName string) {
    userId, err := this.CheckSid()
    if err != nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        http.Redirect(this.Response, this.Request, "/", http.StatusForbidden)
        return
    }

    if tableName == "search" {
        model := this.GetModel("faces")
        regs := this.GetModel("registrations")
        faces := Model{
            ColModel:     model.GetColModel(true, userId),
            TableName:    model.GetTableName(),
            ColNames:     model.GetColNames(),
            Caption:      model.GetCaption(),
            Sub:          true,
            SubTableName: regs.GetTableName(),
            SubCaption:   regs.GetCaption(),
            SubColModel:  regs.GetColModel(true, userId),
            SubColNames:  regs.GetColNames()}

        model = this.GetModel("param_values")
        params := Model{
            ColModel:  model.GetColModel(true, userId),
            TableName: model.GetTableName(),
            ColNames:  model.GetColNames(),
            Caption:   model.GetCaption()}

        this.Render([]string{"mvc/views/search.html"}, "search", map[string]interface{}{"params": params, "faces": faces})
        return
    }

    model := this.GetModel(tableName)
    this.Render([]string{"mvc/views/table.html"}, "table", Model{
        ColModel:  model.GetColModel(true, userId),
        TableName: model.GetTableName(),
        ColNames:  model.GetColNames(),
        Caption:   model.GetCaption(),
        Sub:       model.GetSub()})
}

func (this *GridHandler) EditGridRow(tableName string) {
    userId, err := this.CheckSid()
    if err != nil{
        http.Redirect(this.Response, this.Request, "", http.StatusUnauthorized)
        return
    }

    model := this.GetModel(tableName)
    if model == nil {
        utils.HandleErr("[GridHandler::Edit] GetModel: ", errors.New("Unexpected table name"), this.Response)
        http.Error(this.Response, "Unexpected table name", 400)
        return
    }

    params := make(map[string]interface{}, len(model.GetColumns()))
    for i := 0; i < len(model.GetColumns()); i++ {
        params[model.GetColumnByIdx(i)] = this.Request.PostFormValue(model.GetColumnByIdx(i))
    }

    switch this.Request.PostFormValue("oper") {
    case "edit":
        rowId, err := strconv.Atoi(this.Request.PostFormValue("id"))
        if err != nil {
            http.Error(this.Response, err.Error(), 400)
            return
        }
        model.Update(userId, rowId, params)
        break

    case "add":
        err := model.Add(userId, params)
        if err != nil {
            http.Error(this.Response, err.Error(), 400)
        }
        break

    case "del":
        for _, v := range strings.Split(this.Request.PostFormValue("id"), ",") {
            id, err := strconv.Atoi(v)
            if err != nil {
                http.Error(this.Response, err.Error(), 400)
                return
            }
            model.Delete(id)
        }
        break
    }
}

func (this *GridHandler) JsonToExcel(tableName string) {
    if !sessions.CheckSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        http.Error(this.Response, fmt.Sprintf(err.Error()), 400)
        return
    }

    var filters map[string]interface{}
    if request["filters"] == nil {
        filters = nil
    } else {
        filters = request["filters"].(map[string]interface{})
    }

    fields := utils.ArrayInterfaceToString(request["fields"].([]interface{}))
    sord := request["sord"].(string)
    sidx := request["sidx"].(string)
    data := this.GetModel(tableName).Select(fields, filters, -1, -1, sord, sidx)

    this.Response.Header().Set("Content-type", "text/csv")
    w := csv.NewWriter(this.Response)

    for _, obj := range data {
        var record []string

        for _, col := range obj.(map[string]interface{}) {
            fmt.Printf("type=%s\n", reflect.TypeOf(col))
            switch col.(type) {
            case int:
                record = append(record, strconv.Itoa(col.(int)))
                break
            case int64:
                record = append(record, strconv.Itoa(int(col.(int64))))
                break
            case string:
                record = append(record, col.(string))
                break
            case bool:
                record = append(record, strconv.FormatBool(col.(bool)))
                break
            case []string:
                record = append(record, col.([]string)[0])
                break
            case time.Time:
                record = append(record, col.(time.Time).Format("2006-01-02 15:04:05 07:00"))
            default:
                panic("JsonToExcel: Unexpected type.")
            }
        }

        w.Write(record)
    }

    w.Flush()
}
