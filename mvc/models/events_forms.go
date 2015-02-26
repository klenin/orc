package models

type EventForm struct {
    Id      int `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    EventId int `name:"event_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"events" refField:"id" refFieldShow:"name"`
    FormId  int `name:"form_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"forms" refField:"id" refFieldShow:"name"`
}

func (c *ModelManager) EventsForms() *EventsFormsModel {
    model := new(EventsFormsModel)

    model.TableName = "events_forms"
    model.Caption = "Мероприятия - Формы"

    model.Columns = []string{"id", "event_id", "form_id"}
    model.ColNames = []string{"ID", "Мероприятие", "Форма"}

    model.Fields = new(EventForm)
    model.WherePart = make(map[string]interface{}, 0)
    model.OrderBy = "id"
    model.Limit = "ALL"
    model.Offset = 0

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

type EventsFormsModel struct {
    Entity
}
