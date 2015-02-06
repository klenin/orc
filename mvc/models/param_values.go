package models

import (
    "github.com/orc/db"
)

func (c *ModelManager) ParamValues() *ParamValuesModel {
    model := new(ParamValuesModel)

    model.TableName = "param_values"
    model.Caption = "Знвчение параметров"

    model.Columns = []string{"id", "person_id", "event_id", "param_id", "value"}
    model.ColNames = []string{"ID", "Персона", "Мероприятие", "Параметр", "Значение"}

    model.Fields = []map[string]string{
        {
            "field": "id",
            "type":  "int",
            "null":  "NOT NULL",
            "extra": "PRIMARY"},
        {
            "field":    "person_id",
            "type":     "int",
            "null":     "NOT NULL",
            "extra":    "REFERENCES",
            "refTable": "persons",
            "refField": "id"},
        {
            "field":    "event_id",
            "type":     "int",
            "null":     "NOT NULL",
            "extra":    "REFERENCES",
            "refTable": "events",
            "refField": "id"},
        {
            "field":    "event_type_id",
            "type":     "int",
            "null":     "NOT NULL",
            "extra":    "REFERENCES",
            "refTable": "event_types",
            "refField": "id"},
        {
            "field":    "param_id",
            "type":     "int",
            "null":     "NOT NULL",
            "extra":    "REFERENCES",
            "refTable": "params",
            "refField": "id"},
        {
            "field": "value",
            "type":  "text",
            "null":  "NOT NULL",
            "extra": ""},
    }

    model.Ref = true
    model.RefFields = []string{"fname", "name"}
    model.RefData = make(map[string]interface{}, 3)

    result := db.Select("persons", nil, "", []string{"id", "fname"})
    model.RefData["person_id"] = make([]interface{}, len(result))
    model.RefData["person_id"] = result

    result = db.Select("events", nil, "", []string{"id", "name"})
    model.RefData["event_id"] = make([]interface{}, len(result))
    model.RefData["event_id"] = result

    result = db.Select("params", nil, "", []string{"id", "name"})
    model.RefData["param_id"] = make([]interface{}, len(result))
    model.RefData["param_id"] = result

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

type ParamValuesModel struct {
    Entity
}
