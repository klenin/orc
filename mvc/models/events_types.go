package models

import (
    "github.com/orc/db"
)

func (c *ModelManager) EventsTypes() *EventsTypesModel {
    model := new(EventsTypesModel)

    model.TableName = "events_types"
    model.Caption = "Мероприятия - Типы"

    model.Columns = []string{"id", "event_id", "type_id"}
    model.ColNames = []string{"ID", "Мероприятие", "Тип"}

    model.Fields = []map[string]string{
        {
            "field": "id",
            "type":  "int",
            "null":  "NOT NULL",
            "extra": "PRIMARY"},
        {
            "field":    "event_id",
            "type":     "int",
            "null":     "NOT NULL",
            "extra":    "REFERENCES",
            "refTable": "events",
            "refField": "id"},
        {
            "field":    "type_id",
            "type":     "int",
            "null":     "NOT NULL",
            "extra":    "REFERENCES",
            "refTable": "event_types",
            "refField": "id"},
    }

    model.Ref = true
    model.RefFields = []string{"name"}
    model.RefData = make(map[string]interface{}, 2)

    result := db.Select("events", nil, "", []string{"id", "name"})
    model.RefData["event_id"] = make([]interface{}, len(result))
    model.RefData["event_id"] = result

    result = db.Select("event_types", nil, "", []string{"id", "name"})
    model.RefData["type_id"] = make([]interface{}, len(result))
    model.RefData["type_id"] = result

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

type EventsTypesModel struct {
    Entity
}
