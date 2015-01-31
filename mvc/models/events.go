package models

func (c *ModelManager) Events() *EventModel {
    model := new(EventModel)

    model.TableName = "events"
    model.Caption = "Мероприятия"

    model.Columns = []string{"id", "name", "date_start", "date_finish", "time", "url"}
    model.ColNames = []string{"ID", "Название", "Дата начала", "Дата окончания", "Время", "Сайт"}

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
            "extra": "UNIQUE"},
        {
            "field": "date_start",
            "type":  "date",
            "null":  "NOT NULL",
            "extra": ""},
        {
            "field": "date_finish",
            "type":  "date",
            "null":  "NOT NULL",
            "extra": ""},
        {
            "field": "time",
            "type":  "time",
            "null":  "NOT NULL",
            "extra": ""},
        {
            "field": "url",
            "type":  "text",
            "null":  "NOT NULL",
            "extra": ""},
    }

    model.Ref = false
    model.RefData = nil
    model.RefFields = nil

    model.Sub = true
    model.SubTable = []string{"events_types"}
    model.SubField = "event_id"
    return model
}

type EventModel struct {
    Entity
}
