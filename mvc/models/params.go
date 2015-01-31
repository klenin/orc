package models

import (
    "github.com/orc/db"
)

func (c *ModelManager) Params() *ParamsModel {
    model := new(ParamsModel)

    model.TableName = "params"
    model.Caption = "Параметры"

    model.Columns = []string{"id", "name", "type", "form_id", "identifier"}
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
            "field": "type",
            "type":  "text",
            "null":  "NOT NULL",
            "extra": ""},
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
    model.RefData = make(map[string]interface{}, 1)

    result := db.Select("forms", nil, "", []string{"id", "name"})
    model.RefData["form_id"] = make([]interface{}, len(result))
    model.RefData["form_id"] = result

    model.RefData["type"] = []interface{}{
        map[string]string{"id": "0", "name": "date"},
        map[string]string{"id": "1", "name": "region"},
        map[string]string{"id": "2", "name": "district"},
        map[string]string{"id": "3", "name": "city"},
        map[string]string{"id": "4", "name": "street"},
        map[string]string{"id": "5", "name": "building"},
        map[string]string{"id": "6", "name": "input"},
        map[string]string{"id": "7", "name": "textarea"}}

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

type ParamsModel struct {
    Entity
}
