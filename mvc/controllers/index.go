package controllers

import (
    "encoding/json"
    "github.com/orc/db"
    "github.com/orc/utils"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
    "fmt"
)

func (c *BaseController) Index() *IndexController {
    return new(IndexController)
}

type IndexController struct {
    Controller
}

func (this *IndexController) Index() {
    this.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
    model := this.GetModel("events")
    this.Render([]string{"mvc/views/login.html", "mvc/views/index.html"}, "index", map[string]interface{}{"events": Model{
        ColModel:  model.GetColModel(),
        TableName: model.GetTableName(),
        ColNames:  model.GetColNames(),
        Caption:   model.GetCaption()}})
}

func (this *IndexController) Init(runTest bool) {
    if !runTest {
        return
    }

    for k, v := range db.Tables {
        db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", v), nil)
        db.Exec(fmt.Sprintf("DROP SEQUENCE IF EXISTS %s_id_seq;", v), nil)
        db.QueryCreateTable_(this.GetModel(db.Tables[k]))
    }
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
        event := this.GetModel("events")
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
        event.LoadModelData(map[string]interface{}{
            "name":        contest["name"],
            "date_start":  startDate.Format("2006-01-02 15:04:05"),
            "date_finish": finishDate.Format("2006-01-02 15:04:05"),
            "time":        time_.Format("15:04:05"),
            "url":         "http://imcs.dvfu.ru/cats/main.pl?f=contests;cid="+strconv.Itoa(int(contest["id"].(float64))),
        })
        db.QueryInsert_(event, "").Scan()
    }
}

func (this *IndexController) CreateRegistrationEvent() {
    var eventId int
    events := this.GetModel("events")
    events.LoadModelData(map[string]interface{}{"name": "Регистрация для входа в систему", "date_start": "2006-01-02", "date_finish": "2006-01-02", "time": "00:00:00"})
    db.QueryInsert_(events, "RETURNING id").Scan(&eventId)

    var formId1 int
    forms := this.GetModel("forms")
    forms.LoadModelData(map[string]interface{}{"name": "Регистрационные данные"})
    db.QueryInsert_(forms, "RETURNING id").Scan(&formId1)

    eventsForms := this.GetModel("events_forms")
    eventsForms.LoadModelData(map[string]interface{}{"form_id": formId1, "event_id": eventId})
    db.QueryInsert_(eventsForms, "").Scan()

    var paramTextTypeId int
    paramTypes := this.GetModel("param_types")
    paramTypes.LoadModelData(map[string]interface{}{"name": "text"})
    db.QueryInsert_(paramTypes, "RETURNING id").Scan(&paramTextTypeId)

    var paramPassTypeId int
    paramTypes.LoadModelData(map[string]interface{}{"name": "password"})
    db.QueryInsert_(paramTypes, "RETURNING id").Scan(&paramPassTypeId)

    params := this.GetModel("params")
    params.LoadModelData(map[string]interface{}{
        "name":          "Логин",
        "form_id":       formId1,
        "param_type_id": paramTextTypeId,
        "identifier":    2})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Пароль",
        "form_id":       formId1,
        "param_type_id": paramPassTypeId,
        "identifier":    3})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Подтвердите пароль",
        "form_id":       formId1,
        "param_type_id": paramPassTypeId,
        "identifier":    4})
    db.QueryInsert_(params, "").Scan()

    var paramEmailTypeId int
    paramTypes.LoadModelData(map[string]interface{}{"name": "email"})
    db.QueryInsert_(paramTypes, "RETURNING id").Scan(&paramEmailTypeId)

    params.LoadModelData(map[string]interface{}{
        "name":          "E-mail",
        "form_id":       formId1,
        "param_type_id": paramTextTypeId,
        "identifier":    5})
    db.QueryInsert_(params, "").Scan()

    var formId3 int
    forms.LoadModelData(map[string]interface{}{"name": "Общие сведения"})
    db.QueryInsert_(forms, "RETURNING id").Scan(&formId3)

    params.LoadModelData(map[string]interface{}{
        "name":          "Фамилия",
        "form_id":       formId3,
        "param_type_id": paramTextTypeId,
        "identifier":    6})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Имя",
        "form_id":       formId3,
        "param_type_id": paramTextTypeId,
        "identifier":    7})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Отчество",
        "form_id":       formId3,
        "param_type_id": paramTextTypeId,
        "identifier":    8})
    db.QueryInsert_(params, "").Scan()

    eventsForms.LoadModelData(map[string]interface{}{"form_id": formId3, "event_id": eventId})
    db.QueryInsert_(eventsForms, "").Scan()

    /* Турнир юных программистов */

    events.LoadModelData(map[string]interface{}{
        "name": "Турнир юных программистов",
        "date_start": "2015-04-25",
        "date_finish": "2015-04-25",
        "time": "10:00:00",
        "url": "http://imcs.dvfu.ru/cats/main.pl?f=problems;cid=990998"})
    db.QueryInsert_(events, "RETURNING id").Scan(&eventId)

    eventsForms.LoadModelData(map[string]interface{}{"form_id": formId3, "event_id": eventId})
    db.QueryInsert_(eventsForms, "").Scan()

    var formId4 int
    forms.LoadModelData(map[string]interface{}{"name": "Домашний адрес и телефоны"})
    db.QueryInsert_(forms, "RETURNING id").Scan(&formId4)

    eventsForms.LoadModelData(map[string]interface{}{"form_id": formId4, "event_id": eventId})
    db.QueryInsert_(eventsForms, "").Scan()

    var param_region_type_id int
    paramTypes.LoadModelData(map[string]interface{}{"name": "region"})
    db.QueryInsert_(paramTypes, "RETURNING id").Scan(&param_region_type_id)

    var param_city_type_id int
    paramTypes.LoadModelData(map[string]interface{}{"name": "city"})
    db.QueryInsert_(paramTypes, "RETURNING id").Scan(&param_city_type_id)

    var param_street_type_id int
    paramTypes.LoadModelData(map[string]interface{}{"name": "street"})
    db.QueryInsert_(paramTypes, "RETURNING id").Scan(&param_street_type_id)

    var param_building_type_id int
    paramTypes.LoadModelData(map[string]interface{}{"name": "building"})
    db.QueryInsert_(paramTypes, "RETURNING id").Scan(&param_building_type_id)

    var param_phon_type_id int
    paramTypes.LoadModelData(map[string]interface{}{"name": "phon"})
    db.QueryInsert_(paramTypes, "RETURNING id").Scan(&param_phon_type_id)

    params.LoadModelData(map[string]interface{}{
        "name":          "Регион",
        "form_id":       formId4,
        "param_type_id": param_region_type_id,
        "identifier":    9})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Город",
        "form_id":       formId4,
        "param_type_id": param_city_type_id,
        "identifier":    10})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Улица",
        "form_id":       formId4,
        "param_type_id": param_street_type_id,
        "identifier":    11})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Дом",
        "form_id":       formId4,
        "param_type_id": param_building_type_id,
        "identifier":    12})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Квартира",
        "form_id":       formId4,
        "param_type_id": param_building_type_id,
        "identifier":    13})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Контактный телефон",
        "form_id":       formId4,
        "param_type_id": param_phon_type_id,
        "identifier":    14})
    db.QueryInsert_(params, "").Scan()

    var formId5 int
    forms.LoadModelData(map[string]interface{}{"name": "Образование"})
    db.QueryInsert_(forms, "RETURNING id").Scan(&formId5)

    eventsForms.LoadModelData(map[string]interface{}{"form_id": formId5, "event_id": eventId})
    db.QueryInsert_(eventsForms, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Учебное заведение",
        "form_id":       formId5,
        "param_type_id": paramTextTypeId,
        "identifier":    15})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Класс",
        "form_id":       formId5,
        "param_type_id": paramTextTypeId,
        "identifier":    16})
    db.QueryInsert_(params, "").Scan()

    var formId6 int
    forms.LoadModelData(map[string]interface{}{"name": "Участие в мероприятии"})
    db.QueryInsert_(forms, "RETURNING id").Scan(&formId6)

    eventsForms.LoadModelData(map[string]interface{}{"form_id": formId6, "event_id": eventId})
    db.QueryInsert_(eventsForms, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Тип участия (очное/дистанционное)",
        "form_id":       formId6,
        "param_type_id": paramTextTypeId,
        "identifier":    17})
    db.QueryInsert_(params, "").Scan()

    var formId7 int
    forms.LoadModelData(map[string]interface{}{"name": "Руководитель"})
    db.QueryInsert_(forms, "RETURNING id").Scan(&formId7)

    eventsForms.LoadModelData(map[string]interface{}{"form_id": formId7, "event_id": eventId})
    db.QueryInsert_(eventsForms, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Фамилия",
        "form_id":       formId7,
        "param_type_id": paramTextTypeId,
        "identifier":    18})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Имя",
        "form_id":       formId7,
        "param_type_id": paramTextTypeId,
        "identifier":    19})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Отчество",
        "form_id":       formId7,
        "param_type_id": paramTextTypeId,
        "identifier":    20})
    db.QueryInsert_(params, "").Scan()
}
