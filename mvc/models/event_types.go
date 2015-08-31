package models

type EventTypes struct {
    id          int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    name        string `name:"name" type:"text" null:"NOT NULL" extra:"UNIQUE"`
    description string `name:"description" type:"text" null:"NOT NULL" extra:""`
}

func (this *EventTypes) GetId() int {
    return this.id
}

func (this *EventTypes) SetName(name string) {
    this.name = name
}

func (this *EventTypes) GetName() string {
    return this.name
}

func (this *EventTypes) SetDescription(description string) {
    this.description = description
}

func (this *EventTypes) GetDescription() string {
    return this.description
}

type EventTypesModel struct {
    Entity
}

func (*ModelManager) EventTypes() *EventTypesModel {
    model := new(EventTypesModel)
    model.SetTableName("event_types").
        SetCaption("Типы мероприятий").
        SetColumns([]string{"id", "name", "description"}).
        SetColNames([]string{"ID", "Тип", "Описание"}).
        SetFields(new(EventTypes)).
        SetCondition(AND).
        SetOrder("id").
        SetLimit("ALL").
        SetOffset(0).
        SetSorting("ASC").
        SetWherePart(make(map[string]interface{}, 0)).
        SetSub(false).
        SetSubTables(nil).
        SetSubField("")

    return model
}

func (*EventTypesModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
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
