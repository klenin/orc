package models

type PersonEvent struct {
    Id       string `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    EventId  string `name:"event_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"events" refField:"id" refFieldShow:"name"`
    PersonId string `name:"person_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"persons" refField:"id" refFieldShow:"fname"`
    RegDate  string `name:"reg_date" type:"date" null:"NOT NULL" extra:""`
    LastDate string `name:"last_date" type:"date" null:"NOT NULL" extra:""`
}

func (c *ModelManager) PersonsEvents() *PersonsEventsModel {
    model := new(PersonsEventsModel)

    model.TableName = "persons_events"
    model.Caption = "Персоны - Мероприятия"

    model.Columns = []string{"id", "person_id", "event_id", "reg_date", "last_date"}
    model.ColNames = []string{"ID", "Персона", "Мероприятие", "Дата регистрации", "Дата последних изменений"}

    model.Fields = new(PersonEvent)

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

type PersonsEventsModel struct {
    Entity
}
