package models

import (
    "github.com/orc/db"
    "strconv"
)

type EventsTypesModel struct {
    Entity
}

type EventsTypes struct {
    Id      int `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    EventId int `name:"event_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"events" refField:"id" refFieldShow:"name"`
    TypeId  int `name:"type_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"event_types" refField:"id" refFieldShow:"name"`
}

func (c *ModelManager) EventsTypes() *EventsTypesModel {
    model := new(EventsTypesModel)

    model.TableName = "events_types"
    model.Caption = "Мероприятия - Типы"

    model.Columns = []string{"id", "event_id", "type_id"}
    model.ColNames = []string{"ID", "Мероприятие", "Тип"}

    model.Fields = new(EventsTypes)
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

func (this *EventsTypesModel) Select(fields []string, filters map[string]interface{}, limit, offset int, sord, sidx string) (result []interface{}) {
    if len(fields) == 0 {
        return nil
    }

    query := `SELECT `

    for _, field := range fields {
        switch field {
        case "id":
            query += "events_types.id, "
            break
        case "event_id":
            query += "events.name as event_name, "
            break
        case "type_id":
            query += "event_types.name as type_name, "
            break
        }
    }

    query = query[:len(query)-2]

    query += ` FROM events_types
        INNER JOIN events ON events.id = events_types.event_id
        INNER JOIN event_types ON event_types.id = events_types.type_id`

    where, params := this.Where(filters)
    query += where

    if sidx != "" {
        query += ` ORDER BY events_types.`+sidx
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
