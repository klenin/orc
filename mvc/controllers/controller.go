package controllers

import (
    "github.com/orc/db"
    "github.com/orc/mvc/models"
    "net/http"
    "reflect"
)

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
}

type RequestModel struct {
    E []interface{} //events
    T []interface{} //event_types
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
    //case "teams":
    //    return base.Teams()
    case "persons":
        return base.Persons()
    case "users":
        return base.Users()
    //case "teams_persons":
    //    return base.TeamsPersons()
    case "forms":
        return base.Forms()
    case "params":
        return base.Params()
    case "forms_types":
        return base.FormsTypes()
    case "param_values":
        return base.ParamValues()
    case "persons_events":
        return base.PersonsEvents()
    case "param_types":
        return base.ParamTypes()
    }
    panic("controller.GetModel: have no such table name")
}

func GetModelRefDate(model models.VirtEntity) (fields []string, result map[string]interface{}) {
    result = make(map[string]interface{})
    rt := reflect.TypeOf(model.GetFields())

    for i := 0; i < rt.Elem().NumField(); i++ {
        refFieldShow := rt.Elem().Field(i).Tag.Get("refFieldShow")
        if refFieldShow != "" {
            fields = append(fields, refFieldShow)
            refField := rt.Elem().Field(i).Tag.Get("refField")
            data := db.Select(rt.Elem().Field(i).Tag.Get("refTable"), []string{refField, refFieldShow}, nil, "")
            result[rt.Elem().Field(i).Tag.Get("name")] = make([]interface{}, len(data))
            result[rt.Elem().Field(i).Tag.Get("name")] = data
        }
    }

    return fields, result
}
