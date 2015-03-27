package controllers

import (
    "encoding/json"
    "github.com/orc/db"
    "github.com/orc/utils"
    "html/template"
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

    tmp, err := template.ParseFiles(
        "mvc/views/index.html",
        "mvc/views/header.html",
        "mvc/views/footer.html",
        "mvc/views/login.html")
    if utils.HandleErr("[IndexController::Index] ParseFiles: ", err, this.Response) {
        return
    }

    err = tmp.ExecuteTemplate(this.Response, "index", nil)
    utils.HandleErr("[IndexController::Index] ExecuteTemplate: ", err, this.Response)
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
        var id int
        err = db.QueryInsert_(event, "RETURNING id").Scan(&id)
        utils.HandleErr("[loadContestsFromCats] db.QueryInsert_().Scan(&id): ", err, this.Response)
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

    var form_id2 int
    forms.LoadModelData(map[string]interface{}{"name": "Личные данные"})
    db.QueryInsert_(forms, "RETURNING id").Scan(&form_id2)

    eventsForms := GetModel("events_forms")
    eventsForms.LoadModelData(map[string]interface{}{"form_id": form_id1, "event_id": event_id})
    db.QueryInsert_(eventsForms, "")

    eventsForms.LoadModelData(map[string]interface{}{"form_id": form_id2, "event_id": event_id})
    db.QueryInsert_(eventsForms, "")

    var param_text_type_id int
    paramTypes := GetModel("param_types")
    paramTypes.LoadModelData(map[string]interface{}{"name": "text"})
    db.QueryInsert_(paramTypes, "RETURNING id").Scan(&param_text_type_id)

    var param_pass_type_id int
    paramTypes.LoadModelData(map[string]interface{}{"name": "password"})
    db.QueryInsert_(paramTypes, "RETURNING id").Scan(&param_pass_type_id)

    params := GetModel("params")
    params.LoadModelData(map[string]interface{}{
        "name":          "Имя",
        "form_id":       form_id2,
        "param_type_id": param_text_type_id,
        "identifier":    1})
    db.QueryInsert_(params, "RETURNING id")

    params.LoadModelData(map[string]interface{}{
        "name":          "Логин",
        "form_id":       form_id1,
        "param_type_id": param_text_type_id,
        "identifier":    2})
    db.QueryInsert_(params, "RETURNING id")

    params.LoadModelData(map[string]interface{}{
        "name":          "Пароль",
        "form_id":       form_id1,
        "param_type_id": param_pass_type_id,
        "identifier":    3})
    db.QueryInsert_(params, "RETURNING id")

    params.LoadModelData(map[string]interface{}{
        "name":          "Подтвердите пароль",
        "form_id":       form_id1,
        "param_type_id": param_pass_type_id,
        "identifier":    4})
    db.QueryInsert_(params, "RETURNING id")
}
