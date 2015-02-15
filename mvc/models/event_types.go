package models

type EventTypes struct {
    Id          string `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    Name        string `name:"name" type:"text" null:"NOT NULL" extra:"UNIQUE"`
    Description string `name:"description" type:"text" null:"NOT NULL" extra:""`
    Topicality  string `name:"topicality" type:"boolean" null:"NOT NULL" extra:""`
}

func (c *ModelManager) EventTypes() *EventTypesModel {
    model := new(EventTypesModel)

    model.TableName = "event_types"
    model.Caption = "Типы мероприятий"

    model.Columns = []string{"id", "name", "description", "topicality"}
    model.ColNames = []string{"ID", "Тип", "Описание", "Актуальность"}

    model.Fields = new(EventTypes)

    model.Sub = true
    model.SubTable = []string{"forms_types", "events_types"}
    model.SubField = "type_id"

    return model
}

type EventTypesModel struct {
    Entity
}
