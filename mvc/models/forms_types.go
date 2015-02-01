package models

import (
    "github.com/orc/db"
)

func (c *ModelManager) FormsTypes() *FormsTypesModel {
    model := new(FormsTypesModel)

    model.TableName = "forms_types"
    model.Caption = "Формы - Типы мероприятий"

    model.Columns = []string{"id", "form_id", "type_id", "serial_number"}
    model.ColNames = []string{"ID", "Форма", "Тип", "Порядковый номер"}

    model.Fields = []map[string]string{
        {
            "field": "id",
            "type":  "int",
            "null":  "NOT NULL",
            "extra": "PRIMARY"},
        {
            "field":    "form_id",
            "type":     "int",
            "null":     "NOT NULL",
            "extra":    "REFERENCES",
            "refTable": "forms",
            "refField": "id"},
        {
            "field":    "type_id",
            "type":     "int",
            "null":     "NOT NULL",
            "extra":    "REFERENCES",
            "refTable": "event_types",
            "refField": "id"},
        {
            "field": "serial_number",
            "type":  "integer",
            "null":  "NOT NULL",
            "extra": "UNIQUE"},
    }

    model.Ref = true
    model.RefFields = []string{"name"}
    model.RefData = make(map[string]interface{}, 2)

    result := db.Select("forms", nil, "", []string{"id", "name"})
    model.RefData["form_id"] = make([]interface{}, len(result))
    model.RefData["form_id"] = result

    result = db.Select("event_types", nil, "", []string{"id", "name"})
    model.RefData["type_id"] = make([]interface{}, len(result))
    model.RefData["type_id"] = result

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

type FormsTypesModel struct {
    Entity
}
