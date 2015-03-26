package models

type Event struct {
    Id         int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    Name       string `name:"name" type:"text" null:"NOT NULL" extra:"UNIQUE"`
    DateStart  string `name:"date_start" type:"date" null:"NOT NULL" extra:""`
    DateFinish string `name:"date_finish" type:"date" null:"NOT NULL" extra:""`
    Time       string `name:"time" type:"time" null:"NOT NULL" extra:""`
    Url        string `name:"url" type:"text" null:"NULL" extra:""`
}

func (c *ModelManager) Events() *EventsModel {
    model := new(EventsModel)

    model.TableName = "events"
    model.Caption = "Мероприятия"

    model.Columns = []string{"id", "name", "date_start", "date_finish", "time", "url"}
    model.ColNames = []string{"ID", "Название", "Дата начала", "Дата окончания", "Время", "Сайт"}

    model.Fields = new(Event)
    model.WherePart = make(map[string]interface{}, 0)
    model.Condition = AND
    model.OrderBy = "id"
    model.Limit = "ALL"
    model.Offset = 0

    model.Sub = true
    model.SubTable = []string{"events_types"}
    model.SubField = "event_id"

    return model
}

type EventsModel struct {
    Entity
}
