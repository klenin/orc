package models

type Event struct {
    id         int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    name       string `name:"name" type:"text" null:"NOT NULL" extra:"UNIQUE"`
    dateStart  string `name:"date_start" type:"date" null:"NOT NULL" extra:""`
    dateFinish string `name:"date_finish" type:"date" null:"NOT NULL" extra:""`
    time       string `name:"time" type:"time" null:"NOT NULL" extra:""`
    team       bool   `name:"team" type:"boolean" null:"NOT NULL" extra:""`
    url        string `name:"url" type:"text" null:"NULL" extra:""`
}

func (this *Event) GetId() int {
    return this.id
}

func (this *Event) SetName(name string) {
    this.name = name
}

func (this *Event) GetName() string {
    return this.name
}

func (this *Event) SetDateStart(dateStart string) {
    this.dateStart = dateStart
}

func (this *Event) GetDateStart() string {
    return this.dateStart
}

func (this *Event) SetDateFinish(dateFinish string) {
    this.dateFinish = dateFinish
}

func (this *Event) GetDateFinish() string {
    return this.dateFinish
}

func (this *Event) SetTime(time string) {
    this.time = time
}

func (this *Event) GetTime() string {
    return this.time
}

func (this *Event) SetTeam(team bool) {
    this.team = team
}

func (this *Event) GetTeam() bool {
    return this.team
}

func (this *Event) SetUrl(url string) {
    this.url = url
}

func (this *Event) GetUrl() string {
    return this.url
}

type EventsModel struct {
    Entity
}

func (*ModelManager) Events() *EventsModel {
    model := new(EventsModel)
    model.SetTableName("events").
        SetCaption("Мероприятия").
        SetColumns([]string{"id", "name", "date_start", "date_finish", "time", "team", "url"}).
        SetColNames([]string{"ID", "Название", "Дата начала", "Дата окончания", "Время", "Командное", "Сайт"}).
        SetFields(new(Event)).
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

func (*EventsModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
    return []map[string]interface{} {
        0: map[string]interface{} {
            "index": "id",
            "name": "id",
            "editable": false,
            "width": 20,
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
            "index": "team",
            "name": "team",
            "editable": true,
            "editrules": map[string]interface{}{"required": true},
            "formatter": "checkbox",
            "formatoptions": map[string]interface{}{"disabled": true},
            "edittype": "checkbox",
            "editoptions": map[string]interface{}{"value": "true:false"},
            "width": 70,
        },
        6: map[string]interface{} {
            "index": "url",
            "name": "url",
            "editable": true,
            "formatter": "link",
            "editrules": map[string]interface{}{"url": true, "required": false},
            "width": 250,
        },
    }
}
