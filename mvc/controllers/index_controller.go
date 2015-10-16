package controllers

import (
    "encoding/json"
    "github.com/klenin/orc/db"
    "github.com/klenin/orc/utils"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
    "fmt"
)

func (c *BaseController) IndexController() *IndexController {
    return new(IndexController)
}

type IndexController struct {
    Controller
}

func (this *IndexController) Index() {
    this.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
    model := this.GetModel("events")
    this.Render([]string{"mvc/views/login.html", "mvc/views/index.html"}, "index", map[string]interface{}{"events": Model{
        ColModel:  model.GetColModel(false, -1),
        TableName: model.GetTableName(),
        ColNames:  model.GetColNames(),
        Caption:   model.GetCaption()}})
}

func (this *IndexController) Init(runTest bool) {
    for k, v := range db.Tables {
        db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", v), nil)
        db.Exec(fmt.Sprintf("DROP SEQUENCE IF EXISTS %s_id_seq;", v), nil)
        db.QueryCreateTable(this.GetModel(db.Tables[k]))
    }

    this.CreateRegistrationEvent()
}

func (this *IndexController) LoadContestsFromCats() {
    url := "http://imcs.dvfu.ru/cats/main.pl?f=contests;filter=unfinished;json=1"
    result, err := http.Get(url)
    if utils.HandleErr("[loadContestsFromCats] http.Get(url): ", err, this.Response) {
        return
    }
    defer result.Body.Close()

    body, err := ioutil.ReadAll(result.Body)
    if utils.HandleErr("[loadContestsFromCats] ioutil.ReadAll(data.Body): ", err, this.Response) {
        return
    }

    var data map[string]interface{}
    err = json.Unmarshal(body, &data)
    if utils.HandleErr("[loadContestsFromCats] json.Unmarshal(body, &data): ", err, this.Response) {
        return
    }

    for _, v := range data["contests"].([]interface{}) {
        contest := v.(map[string]interface{})
        time_, err := time.Parse("20060102T150405", contest["start_time"].(string))
        if utils.HandleErr("[loadContestsFromCats] time.Parse: ", err, this.Response) {
            continue
        }
        startDate, err := time.Parse("02.01.2006 15:04", contest["start_date"].(string))
        if utils.HandleErr("[loadContestsFromCats] time.Parse: ", err, this.Response) {
            continue
        }
        finishDate, err := time.Parse("02.01.2006 15:04", contest["finish_date"].(string))
        if utils.HandleErr("[loadContestsFromCats] time.Parse: ", err, this.Response) {
            continue
        }
        this.GetModel("events").
            LoadModelData(map[string]interface{}{
                "name":        contest["name"],
                "date_start":  startDate.Format("2006-01-02 15:04:05"),
                "date_finish": finishDate.Format("2006-01-02 15:04:05"),
                "time":        time_.Format("15:04:05"),
                "url":         "http://imcs.dvfu.ru/cats/main.pl?f=contests;cid="+strconv.Itoa(int(contest["id"].(float64))),
            }).
            QueryInsert("").
            Scan()
    }
}

func (this *IndexController) CreateRegistrationEvent() {
    var eventId int
    this.GetModel("events").
        LoadModelData(map[string]interface{}{
            "name": "Регистрация для входа в систему",
            "date_start": "2006-01-02",
            "date_finish": "2006-01-02",
            "time": "00:00:00"}).
        QueryInsert("RETURNING id").
        Scan(&eventId)

    var formId1 int
    this.GetModel("forms").
        LoadModelData(map[string]interface{}{"name": "Регистрационные данные", "personal": true}).
        QueryInsert("RETURNING id").
        Scan(&formId1)

    this.GetModel("events_forms").
        LoadModelData(map[string]interface{}{"form_id": formId1, "event_id": eventId}).
        QueryInsert("").
        Scan()

    var paramTextTypeId int
    paramTypes := this.GetModel("param_types")
    paramTypes.LoadModelData(map[string]interface{}{"name": "text"}).
        QueryInsert("RETURNING id").
        Scan(&paramTextTypeId)

    var paramPassTypeId int
    paramTypes.LoadModelData(map[string]interface{}{"name": "password"}).
        QueryInsert("RETURNING id").
        Scan(&paramPassTypeId)

    params := this.GetModel("params")
    params.LoadModelData(map[string]interface{}{
        "name":          "Логин",
        "form_id":       formId1,
        "param_type_id": paramTextTypeId,
        "identifier":    2,
        "required":      true,
        "editable":      true}).
        QueryInsert("").
        Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Пароль",
        "form_id":       formId1,
        "param_type_id": paramPassTypeId,
        "identifier":    3,
        "required":      true,
        "editable":      true}).
        QueryInsert("").
        Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Подтвердите пароль",
        "form_id":       formId1,
        "param_type_id": paramPassTypeId,
        "identifier":    4,
        "required":      true,
        "editable":      true}).
        QueryInsert("").
        Scan()

    var paramEmailTypeId int
    paramTypes.LoadModelData(map[string]interface{}{"name": "email"}).
        QueryInsert("RETURNING id").
        Scan(&paramEmailTypeId)

    params.LoadModelData(map[string]interface{}{
        "name":          "E-mail",
        "form_id":       formId1,
        "param_type_id": paramTextTypeId,
        "identifier":    5,
        "required":      true,
        "editable":      true}).
        QueryInsert("").
        Scan()

    var formId2 int
    this.GetModel("forms").
        LoadModelData(map[string]interface{}{"name": "Общие сведения", "personal": true}).
        QueryInsert("RETURNING id").
        Scan(&formId2)

    this.GetModel("events_forms").
        LoadModelData(map[string]interface{}{"form_id": formId2, "event_id": eventId}).
        QueryInsert("").
        Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Фамилия",
        "form_id":       formId2,
        "param_type_id": paramTextTypeId,
        "identifier":    6,
        "required":      true,
        "editable":      true}).
        QueryInsert("").
        Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Имя",
        "form_id":       formId2,
        "param_type_id": paramTextTypeId,
        "identifier":    7,
        "required":      true,
        "editable":      true}).
        QueryInsert("").
        Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Отчество",
        "form_id":       formId2,
        "param_type_id": paramTextTypeId,
        "identifier":    8,
        "required":      true,
        "editable":      true}).
        QueryInsert("").
        Scan()
}
