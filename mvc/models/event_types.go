package models

func (c *ModelManager) EventTypes() *EventTypesModel {
    model := new(EventTypesModel)

    model.TableName = "event_types"
    model.Caption = "Типы мероприятий"

    model.Columns = []string{"id", "name", "description", "topicality"}
    model.ColNames = []string{"ID", "Тип", "Описание", "Актуальность"}

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
            "field": "description",
            "type":  "text",
            "null":  "NOT NULL",
            "extra": ""},
        {
            "field": "topicality",
            "type":  "boolean",
            "null":  "NOT NULL",
            "extra": ""},
    }

    model.Ref = false
    model.RefData = nil
    model.RefFields = nil

    model.Sub = true
    model.SubTable = []string{"forms_types", "events_types"}
    model.SubField = "type_id"

    return model
}

type EventTypesModel struct {
    Entity
}
