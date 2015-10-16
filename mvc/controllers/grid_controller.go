package controllers

import (
    "database/sql"
    "errors"
    "encoding/csv"
    "encoding/json"
    "fmt"
    "github.com/klenin/orc/db"
    "github.com/klenin/orc/sessions"
    "github.com/klenin/orc/utils"
    "net/http"
    "reflect"
    "strconv"
    "strings"
    "time"
)

func (*BaseController) GridController() *GridController {
    return new(GridController)
}

type GridController struct {
    Controller
}

func (this *GridController) GetSubTable() {
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
    if utils.HandleErr("[GridController::GetSubTable] Marshal: ", err, this.Response) {
        return
    }

    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *GridController) CreateGrid(tableName string) {
    userId, err := this.CheckSid()
    if err != nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)

        return
    }

    if !this.isAdmin() {
        http.Redirect(this.Response, this.Request, "/", http.StatusForbidden)

        return
    }

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

    if tableName == "search" {
        this.Render([]string{"mvc/views/search.html"}, "search", map[string]interface{}{"params": params, "faces": faces})

        return
    }

    model = this.GetModel(tableName)
    obj := Model{
        ColModel:  model.GetColModel(true, userId),
        TableName: model.GetTableName(),
        ColNames:  model.GetColNames(),
        Caption:   model.GetCaption(),
        Sub:       model.GetSub()}

    if tableName == "groups" || tableName == "group_registrations" {
        model = this.GetModel("events")
        events := Model{
            ColModel:  model.GetColModel(true, userId),
            TableName: model.GetTableName(),
            ColNames:  model.GetColNames(),
            Caption:   model.GetCaption(),
            Sub:       model.GetSub()}
        this.Render(
            []string{"mvc/views/table.html"},
            "table",
            map[string]interface{}{"model": obj, "events": events, "faces": faces, "params": params})
    } else {
        this.Render(
            []string{"mvc/views/table.html"},
            "table",
            map[string]interface{}{"model": obj})
    }
}

func (this *GridController) EditGridRow(tableName string) {
    userId, err := this.CheckSid()
    if err != nil{
        http.Redirect(this.Response, this.Request, "", http.StatusUnauthorized)

        return
    }

    model := this.GetModel(tableName)
    if model == nil {
        utils.HandleErr("[GridController::Edit] GetModel: ", errors.New("Unexpected table name"), this.Response)
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
        model.Update(this.isAdmin(), userId, params, map[string]interface{}{"id": rowId})
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

//-----------------------------------------------------------------------------
func (this *GridController) GetEventTypesByEventId() {
    if !sessions.CheckSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)

        return
    }

    if !this.isAdmin() {
        http.Redirect(this.Response, this.Request, "/", http.StatusForbidden)

        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)

        return
    }

    eventId, err := strconv.Atoi(request["event_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)

        return
    }

    query := `SELECT event_types.id, event_types.name FROM events_types
        INNER JOIN events ON events.id = events_types.event_id
        INNER JOIN event_types ON event_types.id = events_types.type_id
        WHERE events.id = $1 ORDER BY event_types.id;`
    result := db.Query(query, []interface{}{eventId})

    utils.SendJSReply(map[string]interface{}{"result": "ok", "data": result}, this.Response)
}

func (this *GridController) ImportForms() {
    if !sessions.CheckSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)

        return
    }

    if !this.isAdmin() {
        http.Redirect(this.Response, this.Request, "/", http.StatusForbidden)

        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)

        return
    }

    eventId, err := strconv.Atoi(request["event_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)

        return
    }

    for _, v := range request["event_types_ids"].([]interface{}) {
        typeId, err := strconv.Atoi(v.(string))
        if err != nil {
            utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)

            return
        }

        var lastEventId int
        query := `SELECT events.id FROM events
            INNER JOIN events_types ON events_types.event_id = events.id
            INNER JOIN event_types ON event_types.id = events_types.type_id
            WHERE event_types.id = $1 AND events.id <> $2
            ORDER BY id DESC LIMIT 1;`
        db.QueryRow(query, []interface{}{typeId, eventId}).Scan(&lastEventId)

        query = `SELECT forms.id FROM forms
            INNER JOIN events_forms ON events_forms.form_id = forms.id
            INNER JOIN events ON events.id = events_forms.event_id
            WHERE events.id = $1 ORDER BY forms.id;`
        formsResult := db.Query(query, []interface{}{lastEventId})

        for i := 0; i < len(formsResult); i++ {
            formId := int(formsResult[i].(map[string]interface{})["id"].(int))

            eventsForms := this.GetModel("events_forms")

            var eventFormId int
            if err := eventsForms.
                LoadWherePart(map[string]interface{}{"event_id": eventId, "form_id": formId}).
                SelectRow([]string{"id"}).
                Scan(&eventFormId);
                err != sql.ErrNoRows {
                continue
            }

            eventsForms.
                LoadModelData(map[string]interface{}{"event_id":  eventId, "form_id": formId}).
                QueryInsert("").
                Scan()
        }
    }

    utils.SendJSReply(map[string]interface{}{"result": "ok"}, this.Response)
}

//-----------------------------------------------------------------------------
func (this *GridController) GetPersonsByEventId() {
    if !sessions.CheckSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)

        return
    }

    if !this.isAdmin() {
        http.Redirect(this.Response, this.Request, "/", http.StatusForbidden)

        return
    }

    if this.Request.URL.Query().Get("event") == "" || this.Request.URL.Query().Get("params") == "" {

        return
    }

    eventId, err := strconv.Atoi(this.Request.URL.Query().Get("event"))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)

        return
    }

    paramsIds := strings.Split(this.Request.URL.Query().Get("params"), ",")

    if len(paramsIds) == 0 {
        utils.SendJSReply(map[string]interface{}{"result": "Выберите параметры."}, this.Response)

        return
    }

    var queryParams []interface{}
    query := "SELECT params.name FROM params WHERE params.id in ("

    for k, v := range paramsIds {
        paramId, err := strconv.Atoi(v)
        if err != nil {
            utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)

            return
        }
        query += "$"+strconv.Itoa(k+1)+", "
        queryParams = append(queryParams, paramId)
    }
    query = query[:len(query)-2]
    query+=") ORDER BY id;"

    var caption []string
    for _, v := range db.Query(query, queryParams) {
        caption = append(caption, v.(map[string]interface{})["name"].(string))
    }

    result := []interface{}{0: map[string]interface{}{"id": -1, "data": caption}}

    query = `SELECT
        reg.id as id,
        ARRAY(
            SELECT param_values.value
            FROM param_values
            INNER JOIN registrations ON registrations.id = param_values.reg_id
            INNER JOIN events ON events.id = registrations.event_id
            INNER JOIN params ON params.id = param_values.param_id
            WHERE param_values.param_id IN (` + strings.Join(db.MakeParams(len(queryParams)), ", ")
    query += `) AND events.id = $` + strconv.Itoa(len(queryParams)+1) + ` AND registrations.id = reg.id ORDER BY param_values.param_id
        ) as data

        FROM param_values
        INNER JOIN registrations as reg ON reg.id = param_values.reg_id
        INNER JOIN events as ev ON ev.id = reg.event_id
        INNER JOIN params ON params.id = param_values.param_id
        WHERE ev.id = $` + strconv.Itoa(len(queryParams)+1) + ` GROUP BY reg.id ORDER BY reg.id;`

    data := db.Query(query, append(queryParams, eventId))

    this.Render([]string{"mvc/views/list.html"}, "list", append(result, data...))
}

func (this *GridController) GetParamsByEventId() {
    if !sessions.CheckSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)

        return
    }

    if !this.isAdmin() {
        utils.SendJSReply(map[string]interface{}{"result": errors.New("Forbidden")}, this.Response)
        http.Redirect(this.Response, this.Request, "/", http.StatusForbidden)

        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)

        return
    }

    eventId, err := strconv.Atoi(request["event_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)

        return
    }

    query := `SELECT params.id, params.name
        FROM events_forms
        INNER JOIN forms ON forms.id = events_forms.form_id
        INNER JOIN params ON params.form_id = forms.id
        INNER JOIN events ON events.id = events_forms.event_id
        WHERE events.id = $1 ORDER BY params.id;`
    result := db.Query(query, []interface{}{eventId})

    utils.SendJSReply(map[string]interface{}{"result": "ok", "data": result}, this.Response)
}

//-----------------------------------------------------------------------------
func (this *GridController) JsonToExcel(tableName string) {
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
    data := this.GetModel(tableName).
        SetSorting(request["sord"].(string)).
        SetOrder(request["sidx"].(string)).
        Select(fields, filters)

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
