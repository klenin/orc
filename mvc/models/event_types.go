package models

type EventTypesModel struct {
    Entity
}

type EventTypes struct {
    Id          int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    Name        string `name:"name" type:"text" null:"NOT NULL" extra:"UNIQUE"`
    Description string `name:"description" type:"text" null:"NOT NULL" extra:""`
}

func (c *ModelManager) EventTypes() *EventTypesModel {
    model := new(EventTypesModel)

    model.TableName = "event_types"
    model.Caption = "Типы мероприятий"

    model.Columns = []string{"id", "name", "description"}
    model.ColNames = []string{"ID", "Тип", "Описание"}

    model.Fields = new(EventTypes)
    model.WherePart = make(map[string]interface{}, 0)
    model.Condition = AND
    model.OrderBy = "id"
    model.Limit = "ALL"
    model.Offset = 0

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

func (this *EventTypesModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
    return []map[string]interface{} {
        0: map[string]interface{} {
            "index": "id",
            "name": "id",
            "editable": false,
        },
        1: map[string]interface{} {
            "index": "name",
            "name": "name",
            "editable": true,
            "editrules": map[string]interface{}{"required": true},
            "edittype": "text",
        },
        2: map[string]interface{} {
            "index": "description",
            "name": "description",
            "editable": true,
            "editrules": map[string]interface{}{"required": true},
            "edittype": "textarea",
        },
    }
}
