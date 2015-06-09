package models

import (
    "github.com/orc/db"
    "strconv"
)

type EventsFormsModel struct {
    Entity
}

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
    model.Condition = AND
    model.OrderBy = "id"
    model.Limit = "ALL"
    model.Offset = 0

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

func (this *EventsFormsModel) Select(fields []string, filters map[string]interface{}, limit, offset int, sord, sidx string) (result []interface{}) {
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

    if sidx != "" {
        query += ` ORDER BY events_forms.`+sidx
    }

    query += ` `+ sord

    if limit != -1 {
        params = append(params, limit)
        query += ` LIMIT $`+strconv.Itoa(len(params))
    }

    if offset != -1 {
        params = append(params, offset)
        query += ` OFFSET $`+strconv.Itoa(len(params))
    }

    query += `;`

    return db.Query(query, params)
}

func (this *EventsFormsModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
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
