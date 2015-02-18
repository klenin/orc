package models

type EventsTypes struct {
    Id      int `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    EventId int `name:"event_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"events" refField:"id" refFieldShow:"name"`
    TypeId  int `name:"type_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"event_types" refField:"id" refFieldShow:"name"`
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
