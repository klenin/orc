package controllers

import (
    "encoding/json"
    "github.com/orc/db"
    "github.com/orc/utils"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
)

func (c *BaseController) Index() *IndexController {
    return new(IndexController)
}

type IndexController struct {
    Controller
}

func (this *IndexController) Index() {
    this.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
    this.Render([]string{"mvc/views/login.html", "mvc/views/index.html"}, "index", nil)
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
        event := GetModel("events")
        time_, err := time.Parse("20060102T150405", contest["start_time"].(string))
        if utils.HandleErr("[loadContestsFromCats] time.Parse: ", err, this.Response) {
            continue
        }
        start_date, err := time.Parse("02.01.2006 15:04", contest["start_date"].(string))
        if utils.HandleErr("[loadContestsFromCats] time.Parse: ", err, this.Response) {
            continue
        }
        finish_date, err := time.Parse("02.01.2006 15:04", contest["finish_date"].(string))
        if utils.HandleErr("[loadContestsFromCats] time.Parse: ", err, this.Response) {
            continue
        }
        event.LoadModelData(map[string]interface{}{
            "name":        contest["name"],
            "date_start":  start_date.Format("2006-01-02 15:04:05"),
            "date_finish": finish_date.Format("2006-01-02 15:04:05"),
            "time":        time_.Format("15:04:05"),
            "url":         "http://imcs.dvfu.ru/cats/main.pl?f=contests;cid="+strconv.Itoa(int(contest["id"].(float64))),
        })
        err = db.QueryInsert_(event, "").Scan()
    }

}

func CreateRegistrationEvent() {

    var event_id int
    events := GetModel("events")
    events.LoadModelData(map[string]interface{}{"name": "Регистрация", "date_start": "2006-01-02", "date_finish": "2006-01-02", "time": "00:00:00"})
    db.QueryInsert_(events, "RETURNING id").Scan(&event_id)

    var form_id1 int
    forms := GetModel("forms")
    forms.LoadModelData(map[string]interface{}{"name": "Регистрационные данные"})
    db.QueryInsert_(forms, "RETURNING id").Scan(&form_id1)

    eventsForms := GetModel("events_forms")
    eventsForms.LoadModelData(map[string]interface{}{"form_id": form_id1, "event_id": event_id})
    db.QueryInsert_(eventsForms, "").Scan()

    var param_text_type_id int
    paramTypes := GetModel("param_types")
    paramTypes.LoadModelData(map[string]interface{}{"name": "text"})
    db.QueryInsert_(paramTypes, "RETURNING id").Scan(&param_text_type_id)

    var param_pass_type_id int
    paramTypes.LoadModelData(map[string]interface{}{"name": "password"})
    db.QueryInsert_(paramTypes, "RETURNING id").Scan(&param_pass_type_id)

    params := GetModel("params")
    params.LoadModelData(map[string]interface{}{
        "name":          "Логин",
        "form_id":       form_id1,
        "param_type_id": param_text_type_id,
        "identifier":    2})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Пароль",
        "form_id":       form_id1,
        "param_type_id": param_pass_type_id,
        "identifier":    3})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Подтвердите пароль",
        "form_id":       form_id1,
        "param_type_id": param_pass_type_id,
        "identifier":    4})
    db.QueryInsert_(params, "").Scan()

    var param_email_type_id int
    paramTypes.LoadModelData(map[string]interface{}{"name": "email"})
    db.QueryInsert_(paramTypes, "RETURNING id").Scan(&param_email_type_id)

    params.LoadModelData(map[string]interface{}{
        "name":          "E-mail",
        "form_id":       form_id1,
        "param_type_id": param_text_type_id,
        "identifier":    5})
    db.QueryInsert_(params, "").Scan()

    /* Турнир юных программистов */

    events.LoadModelData(map[string]interface{}{
        "name": "Турнир юных программистов",
        "date_start": "2015-04-25",
        "date_finish": "2015-04-25",
        "time": "10:00:00"})
    db.QueryInsert_(events, "RETURNING id").Scan(&event_id)

    var form_id3 int
    forms.LoadModelData(map[string]interface{}{"name": "Общие сведения"})
    db.QueryInsert_(forms, "RETURNING id").Scan(&form_id3)

    eventsForms.LoadModelData(map[string]interface{}{"form_id": form_id3, "event_id": event_id})
    db.QueryInsert_(eventsForms, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Фамилия",
        "form_id":       form_id3,
        "param_type_id": param_text_type_id,
        "identifier":    6})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Имя",
        "form_id":       form_id3,
        "param_type_id": param_text_type_id,
        "identifier":    7})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Отчество",
        "form_id":       form_id3,
        "param_type_id": param_text_type_id,
        "identifier":    8})
    db.QueryInsert_(params, "").Scan()

    var form_id4 int
    forms.LoadModelData(map[string]interface{}{"name": "Домашний адрес и телефоны"})
    db.QueryInsert_(forms, "RETURNING id").Scan(&form_id4)

    eventsForms.LoadModelData(map[string]interface{}{"form_id": form_id4, "event_id": event_id})
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
        "form_id":       form_id4,
        "param_type_id": param_region_type_id,
        "identifier":    9})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Город",
        "form_id":       form_id4,
        "param_type_id": param_city_type_id,
        "identifier":    10})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Улица",
        "form_id":       form_id4,
        "param_type_id": param_street_type_id,
        "identifier":    11})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Дом",
        "form_id":       form_id4,
        "param_type_id": param_building_type_id,
        "identifier":    12})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Контактный телефон",
        "form_id":       form_id4,
        "param_type_id": param_phon_type_id,
        "identifier":    13})
    db.QueryInsert_(params, "").Scan()

    var form_id5 int
    forms.LoadModelData(map[string]interface{}{"name": "Образование"})
    db.QueryInsert_(forms, "RETURNING id").Scan(&form_id5)

    eventsForms.LoadModelData(map[string]interface{}{"form_id": form_id5, "event_id": event_id})
    db.QueryInsert_(eventsForms, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Учебное заведение",
        "form_id":       form_id5,
        "param_type_id": param_text_type_id,
        "identifier":    14})
    db.QueryInsert_(params, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Класс",
        "form_id":       form_id5,
        "param_type_id": param_text_type_id,
        "identifier":    15})
    db.QueryInsert_(params, "").Scan()

    var form_id6 int
    forms.LoadModelData(map[string]interface{}{"name": "Участие в мероприятии"})
    db.QueryInsert_(forms, "RETURNING id").Scan(&form_id6)

    eventsForms.LoadModelData(map[string]interface{}{"form_id": form_id6, "event_id": event_id})
    db.QueryInsert_(eventsForms, "").Scan()

    params.LoadModelData(map[string]interface{}{
        "name":          "Тип участия (очное/дистанционное)",
        "form_id":       form_id6,
        "param_type_id": param_text_type_id,
        "identifier":    16})
    db.QueryInsert_(params, "").Scan()
}
