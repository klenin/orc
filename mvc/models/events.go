package models

type EventsModel struct {
    Entity
}

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

    // model.Sub = true
    // model.SubTable = []string{"events_types"}
    // model.SubField = "event_id"

    return model
}

func (this *EventsModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
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
            "width": 300,
        },
        2: map[string]interface{} {
            "index": "date_start",
            "name": "date_start",
            "editable": true,
            "formatter": nil,
            "editrules": map[string]interface{}{"date": true, "required": true},
            "editoptions": map[string]interface{}{"dataInit": nil},
            "formatoptions": map[string]string{"srcformat": "Y-m-d", "newformat": "Y-m-d"},
            "searchoptions": map[string]interface{}{"sopt": []string{"eq", "ne"}, "dataInit": nil},
            "type": "date",
            "width": 100,
        },
        3: map[string]interface{} {
            "index": "date_finish",
            "name": "date_finish",
            "editable": true,
            "formatter": nil,
            "editrules": map[string]interface{}{"date": true, "required": true},
            "editoptions": map[string]interface{}{"dataInit": nil},
            "formatoptions": map[string]string{"srcformat": "Y-m-d", "newformat": "Y-m-d"},
            "searchoptions": map[string]interface{}{"sopt": []string{"eq", "ne"}, "dataInit": nil},
            "type": "date",
            "width": 100,
        },
        4: map[string]interface{} {
            "index": "time",
            "name": "time",
            "editable": true,
            "formatter": nil,
            "editrules": map[string]interface{}{"custom": true, "custom_func": nil, "required": true},
            "editoptions": map[string]interface{}{"dataInit": nil},
            "formatoptions": map[string]string{"srcformat": "Y-m-d", "newformat": "Y-m-d"},
            "searchoptions": map[string]interface{}{"sopt": []string{"eq", "ne"}, "dataInit": nil},
            "type": "time",
            "width": 100,
            "fixed": true,
        },
        5: map[string]interface{} {
            "index": "url",
            "name": "url",
            "editable": true,
            "formatter": "link",
            "editrules": map[string]interface{}{"url": true, "required": false},
            "width": 250,
        },
    }
}
