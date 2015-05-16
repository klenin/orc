package controllers

import (
    "encoding/json"
    "encoding/csv"
    "fmt"
    "github.com/orc/db"
    "github.com/orc/mailer"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "net/http"
    "strconv"
    "errors"
    "reflect"
    "time"
)

func (c *BaseController) GridHandler() *GridHandler {
    return new(GridHandler)
}

type GridHandler struct {
    Controller
}

func (this *GridHandler) GetSubTable() {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(err.Error(), this.Response)
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
        "colmodel": subModel.GetColModel()})
    if utils.HandleErr("[GridHandler::GetSubTable] Marshal: ", err, this.Response) {
        return
    }

    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *GridHandler) CreateGrid(tableName string) {
    if !sessions.CheackSession(this.Response, this.Request) {
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
            ColModel:     model.GetColModel(),
            TableName:    model.GetTableName(),
            ColNames:     model.GetColNames(),
            Caption:      model.GetCaption(),
            Sub:          true,
            SubTableName: regs.GetTableName(),
            SubCaption:   regs.GetCaption(),
            SubColModel:  regs.GetColModel(),
            SubColNames:  regs.GetColNames()}

        model = this.GetModel("param_values")
        params := Model{
            ColModel:  model.GetColModel(),
            TableName: model.GetTableName(),
            ColNames:  model.GetColNames(),
            Caption:   model.GetCaption()}

        this.Render([]string{"mvc/views/search.html"}, "search", map[string]interface{}{"params": params, "faces": faces})
        return
    }

    model := this.GetModel(tableName)
    this.Render([]string{"mvc/views/table.html"}, "table", Model{
        ColModel:  model.GetColModel(),
        TableName: model.GetTableName(),
        ColNames:  model.GetColNames(),
        Caption:   model.GetCaption(),
        Sub:       model.GetSub()})
}

func (this *GridHandler) EditGridRow(tableName string) {
    user_id := sessions.GetValue("id", this.Request)

    if !sessions.CheackSession(this.Response, this.Request) || user_id == nil {
        http.Redirect(this.Response, this.Request, "", http.StatusUnauthorized)
        return
    }

    model := this.GetModel(tableName)
    if model == nil {
        utils.HandleErr("[GridHandler::Edit] GetModel: ", errors.New("Unexpected table name"), this.Response)
        http.Error(this.Response, fmt.Sprintf("Unexpected table name"), 400)
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
        if err != nil {
            http.Error(this.Response, fmt.Sprintf(err.Error()), 400)
            return
        }

        if tableName == "groups" && !this.isAdmin() {
            face_id, err := db.IsUserGroup(user_id.(int), id)
            if err != nil {
                http.Error(this.Response, fmt.Sprintf(err.Error()), 400)
                return
            }
            params["face_id"] = face_id
        }

        model.LoadModelData(params)
        model.LoadWherePart(map[string]interface{}{"id": id})
        db.QueryUpdate_(model).Scan()
        break

    case "add":
        if tableName == "groups" {
            var face_id int
            query := `SELECT faces.id
                FROM registrations
                INNER JOIN faces ON faces.id = registrations.face_id
                INNER JOIN events ON events.id = registrations.event_id
                INNER JOIN users ON faces.user_id = users.id
                WHERE users.id = $1 AND events.id = $2;`
            db.QueryRow(query, []interface{}{user_id, 1}).Scan(&face_id)
            params["face_id"] = face_id

        } else if tableName == "persons" {
            to := params["name"].(string)
            address := params["email"].(string)
            token := utils.GetRandSeq(HASH_SIZE)
            params["token"] = token

            query := `SELECT param_values.value
                FROM reg_param_vals
                INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
                INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
                INNER JOIN params ON params.id = param_values.param_id
                INNER JOIN events ON events.id = registrations.event_id
                INNER JOIN faces ON faces.id = registrations.face_id
                INNER JOIN users ON users.id = faces.user_id
                WHERE params.id in (5, 6, 7) AND users.id = $1 AND events.id = 1 ORDER BY params.id;`
            data := db.Query(query, []interface{}{user_id})
            headName := data[0].(map[string]interface{})["value"].(string)
            headName += " " + data[1].(map[string]interface{})["value"].(string)
            headName += " " + data[2].(map[string]interface{})["value"].(string)

            group_id, err := strconv.Atoi(params["group_id"].(string))
            if utils.HandleErr("[GridHandler::Edit] group_id Atoi: ", err, this.Response) {
                http.Error(this.Response, fmt.Sprintf(err.Error()), 400)
                return
            }

            var groupName string
            db.QueryRow("SELECT name FROM groups WHERE id = $1;", []interface{}{group_id}).Scan(&groupName)

            if !mailer.InviteToGroup(to, address, token, headName, groupName) {
                http.Error(this.Response, fmt.Sprintf("Проверьте правильность введенного Вами email"), 400)
                return
            }
        }
        model.LoadModelData(params)
        db.QueryInsert_(model, "").Scan()
        break

    case "del":
        db.QueryDeleteByIds(tableName, this.Request.PostFormValue("id"))
        break
    }
}

func (this *GridHandler) JsonToExcel(tableName string) {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
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
