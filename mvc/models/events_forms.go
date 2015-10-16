package models

import (
    "github.com/klenin/orc/db"
    "strconv"
)

type EventForm struct {
    id      int `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    eventId int `name:"event_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"events" refField:"id" refFieldShow:"name"`
    formId  int `name:"form_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"forms" refField:"id" refFieldShow:"name"`
}

func (this *EventForm) GetId() int {
    return this.id
}

func (this *EventForm) GetEventId() int {
    return this.eventId
}

func (this *EventForm) SetEventId(eventId int) {
    this.eventId = eventId
}

func (this *EventForm) GetFormId() int {
    return this.formId
}

func (this *EventForm) SetFormId(formId int) {
    this.formId = formId
}

type EventsFormsModel struct {
    Entity
}

func (*ModelManager) EventsForms() *EventsFormsModel {
    model := new(EventsFormsModel)
    model.SetTableName("events_forms").
        SetCaption("Мероприятия - Формы").
        SetColumns([]string{"id", "event_id", "form_id"}).
        SetColNames([]string{"ID", "Мероприятие", "Форма"}).
        SetFields(new(EventForm)).
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

func (this *EventsFormsModel) Select(fields []string, filters map[string]interface{}) (result []interface{}) {
    if len(fields) == 0 {
        return nil
    }

    query := `SELECT `

    for _, field := range fields {
        switch field {
        case "id":
            query += "events_forms.id, "
            break
        case "event_id":
            query += "events.name as event_name, "
            break
        case "form_id":
            query += "forms.name as form_name, "
            break
        }
    }

    query = query[:len(query)-2]
    query += ` FROM events_forms
        INNER JOIN events ON events.id = events_forms.event_id
        INNER JOIN forms ON forms.id = events_forms.form_id`
    where, params, _ := this.Where(filters, 1)
    if where != "" {
        where = " WHERE " + where
    }
    query += where
    query += ` ORDER BY events_forms.` + this.GetOrder()
    query += ` `+ this.GetSorting()
    params = append(params, this.GetLimit())
    query += ` LIMIT $` + strconv.Itoa(len(params))
    params = append(params, this.GetOffset())
    query += ` OFFSET $` + strconv.Itoa(len(params)) + `;`

    return db.Query(query, params)
}

func (*EventsFormsModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
    query := `SELECT array_to_string(
        array(SELECT events.id || ':' || events.name FROM events GROUP BY events.id ORDER BY events.id), ';') as name;`
    events := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    query = `SELECT array_to_string(
        array(SELECT forms.id || ':' || forms.name FROM forms GROUP BY forms.id ORDER BY forms), ';') as name;`
    forms := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    return []map[string]interface{} {
        0: map[string]interface{} {
            "index": "id",
            "name": "id",
            "editable": false,
        },
        1: map[string]interface{} {
            "index": "event_id",
            "name": "event_id",
            "editable": true,
            "formatter": "select",
            "edittype": "select",
            "stype": "select",
            "search": true,
            "editrules": map[string]interface{}{"required": true},
            "editoptions": map[string]string{"value": events},
            "searchoptions": map[string]string{"value": ":Все;"+events},
        },
        2: map[string]interface{} {
            "index": "form_id",
            "name": "form_id",
            "editable": true,
            "formatter": "select",
            "edittype": "select",
            "stype": "select",
            "search": true,
            "editrules": map[string]interface{}{"required": true},
            "editoptions": map[string]string{"value": forms},
            "searchoptions": map[string]string{"value": ":Все;"+forms},
        },
    }
}
