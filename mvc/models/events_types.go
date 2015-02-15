package models

type EventsTypes struct {
    Id      string `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    EventId string `name:"event_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"events" refField:"id" refFieldShow:"name"`
    TypeId  string `name:"type_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"event_types" refField:"id" refFieldShow:"name"`
}

func (c *ModelManager) EventsTypes() *EventsTypesModel {
    model := new(EventsTypesModel)

    model.TableName = "events_types"
    model.Caption = "Мероприятия - Типы"

    model.Columns = []string{"id", "event_id", "type_id"}
    model.ColNames = []string{"ID", "Мероприятие", "Тип"}

    model.Fields = new(EventsTypes)

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

type EventsTypesModel struct {
    Entity
}
