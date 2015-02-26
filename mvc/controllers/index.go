package controllers

import (
    "github.com/orc/db"
    "github.com/orc/utils"
    "html/template"
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
