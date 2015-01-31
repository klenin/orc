package models

import (
    "github.com/orc/db"
)

func (c *ModelManager) PersonsEvents() *PersonsEventsModel {
    model := new(PersonsEventsModel)

    model.TableName = "persons_events"
    model.Caption = "Персоны - Мероприятия"

    model.Columns = []string{"id", "person_id", "event_id", "reg_date", "last_date"}
    model.ColNames = []string{"ID", "Персона", "Мероприятие", "Дата регистрации", "Дата последних изменений"}

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
            "field": "reg_date",
            "type":  "date",
            "null":  "NOT NULL",
            "extra": ""},
        {
            "field": "last_date",
            "type":  "date",
            "null":  "NOT NULL",
            "extra": ""},
    }

    model.Ref = true
    model.RefFields = []string{"name", "fname"}
    model.RefData = make(map[string]interface{}, 2)

    result := db.Select("persons", nil, "", []string{"id", "fname"})
    model.RefData["person_id"] = make([]interface{}, len(result))
    model.RefData["person_id"] = result

    result = db.Select("events", nil, "", []string{"id", "name"})
    model.RefData["event_id"] = make([]interface{}, len(result))
    model.RefData["event_id"] = result

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

type PersonsEventsModel struct {
    Entity
}
