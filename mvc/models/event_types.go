package models

type EventTypes struct {
    Id          int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    Name        string `name:"name" type:"text" null:"NOT NULL" extra:"UNIQUE"`
    Description string `name:"description" type:"text" null:"NOT NULL" extra:""`
    Topicality  bool   `name:"topicality" type:"boolean" null:"NOT NULL" extra:""`
}

func (c *ModelManager) EventTypes() *EventTypesModel {
    model := new(EventTypesModel)

    model.TableName = "event_types"
    model.Caption = "Типы мероприятий"

    model.Columns = []string{"id", "name", "description", "topicality"}
    model.ColNames = []string{"ID", "Тип", "Описание", "Актуальность"}

    model.Fields = new(EventTypes)
    model.WherePart = make(map[string]interface{}, 0)
    model.OrderBy = "id"
    model.Limit = "ALL"
    model.Offset = 0

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

type EventTypesModel struct {
    Entity
}
