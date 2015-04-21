package controllers

import (
    "github.com/orc/db"
    "github.com/orc/mvc/models"
    "github.com/orc/sessions"
    "net/http"
    "reflect"
    "html/template"
)

const HASH_SIZE = 32

type BaseController struct{}

type Controller struct {
    Request  *http.Request
    Response http.ResponseWriter
}

type Model struct {
    Id        string
    TableName string
    Caption   string
    Table     []interface{}
    RefData   map[string]interface{}
    RefFields []string
    Columns   []string
    ColNames  []string
    Sub       bool
    SubTableName string
    SubCaption   string
    SubRefData   map[string]interface{}
    SubRefFields []string
    SubColumns   []string
    SubColNames  []string
}

type RequestModel struct {
    E []interface{} //events
    F []interface{} //forms
    P []interface{} //params
}

func GetModel(tableName string) models.VirtEntity {
    base := new(models.ModelManager)
    switch tableName {
    case "events":
        return base.Events()
    case "event_types":
        return base.EventTypes()
    case "events_types":
        return base.EventsTypes()
    case "persons":
        return base.Persons()
    case "users":
        return base.Users()
    case "forms":
        return base.Forms()
    case "params":
        return base.Params()
    case "events_forms":
        return base.EventsForms()
    case "param_values":
        return base.ParamValues()
    case "param_types":
        return base.ParamTypes()
    case "registrations":
        return base.Registrations()
    case "faces":
        return base.Faces()
    case "reg_param_vals":
        return base.RegParamVals()
    case "groups":
        return base.Groups()
    case "group_registrations":
        return base.GroupRegistrations()
    }
    return nil
}

func GetModelRefDate(model models.VirtEntity) (fields []string, result map[string]interface{}) {
    result = make(map[string]interface{})
    rt := reflect.TypeOf(model.GetFields())

    for i := 0; i < rt.Elem().NumField(); i++ {
        refFieldShow := rt.Elem().Field(i).Tag.Get("refFieldShow")
        if refFieldShow != "" {
            fields = append(fields, refFieldShow)
            refField := rt.Elem().Field(i).Tag.Get("refField")
            data := db.Select(GetModel(rt.Elem().Field(i).Tag.Get("refTable")), []string{refField, refFieldShow})
            result[rt.Elem().Field(i).Tag.Get("name")] = make([]interface{}, len(data))
            result[rt.Elem().Field(i).Tag.Get("name")] = data
        }
    }

    return fields, result
}

func (this *Controller) Render(filenames []string, tmpname string, data interface{}) {
    filenames = append(filenames, "mvc/views/header.html")
    filenames = append(filenames, "mvc/views/footer.html")
    tmpl, err := template.ParseFiles(filenames...)
    if err != nil {
        http.Error(this.Response, err.Error(), http.StatusInternalServerError)
    }
    if err := tmpl.ExecuteTemplate(this.Response, tmpname, data); err != nil {
        http.Error(this.Response, err.Error(), http.StatusInternalServerError)
    }
}

type VirtController interface {
    Render(filename string, data interface{})
}

func (this *Controller) isAdmin() bool {
    var role string

    user_id := sessions.GetValue("id", this.Request)
    if user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return false
    }

    user := GetModel("users")
    user.LoadWherePart(map[string]interface{}{"id": user_id})
    err := db.SelectRow(user, []string{"role"}).Scan(&role)
    if err != nil || role == "user" {
        return false
    }

    return role == "admin"
}
