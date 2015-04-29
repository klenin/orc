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

    where, params := this.Where(filters)
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
