package models

import (
    "github.com/orc/db"
)

func (c *ModelManager) Params() *ParamsModel {
    model := new(ParamsModel)

    model.TableName = "params"
    model.Caption = "Параметры"

    model.Columns = []string{"id", "name", "param_type_id", "form_id", "identifier"}
    model.ColNames = []string{"ID", "Название", "Тип", "Форма", "Идентификатор"}

    model.Fields = []map[string]string{
        {
            "field": "id",
            "type":  "int",
            "null":  "NOT NULL",
            "extra": "PRIMARY"},
        {
            "field": "name",
            "type":  "text",
            "null":  "NOT NULL",
            "extra": ""},
        {
            "field": "param_type_id",
            "type":  "int",
            "null":  "NOT NULL",
            "extra": "REFERENCES",
            "refTable": "param_types",
            "refField": "id"},
        {
            "field":    "form_id",
            "type":     "int",
            "null":     "NOT NULL",
            "extra":    "REFERENCES",
            "refTable": "forms",
            "refField": "id"},
        {
            "field": "identifier",
            "type":  "text",
            "null":  "NOT NULL",
            "extra": ""},
    }

    model.Ref = true
    model.RefFields = []string{"name"}
    model.RefData = make(map[string]interface{}, 2)

    result := db.Select("forms", nil, "", []string{"id", "name"})
    model.RefData["form_id"] = make([]interface{}, len(result))
    model.RefData["form_id"] = result

    result = db.Select("param_types", nil, "", []string{"id", "name"})
    model.RefData["param_type_id"] = make([]interface{}, len(result))
    model.RefData["param_type_id"] = result

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

type ParamsModel struct {
    Entity
}
