package models

import (
    "github.com/orc/db"
    "strconv"
)

type EventsDocsModel struct {
    Entity
}

type EventsDocs struct {
    Id      int `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    EventId int `name:"event_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"events" refField:"id" refFieldShow:"name"`
    DocId   int `name:"doc_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"docs" refField:"id" refFieldShow:"name"`
}

func (c *ModelManager) EventsDocs() *EventsDocsModel {
    model := new(EventsDocsModel)

    model.TableName = "events_docs"
    model.Caption = "Мероприятия - Документы"

    model.Columns = []string{"id", "event_id", "doc_id"}
    model.ColNames = []string{"ID", "Мероприятие", "Документ"}

    model.Fields = new(EventsDocs)
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

func (this *EventsDocsModel) Select(fields []string, filters map[string]interface{}, limit, offset int, sord, sidx string) (result []interface{}) {
    if len(fields) == 0 {
        return nil
    }

    query := `SELECT `

    for _, field := range fields {
        switch field {
        case "id":
            query += "events_docs.id, "
            break
        case "event_id":
            query += "events.name as event_name, "
            break
        case "doc_id":
            query += "docs.name as doc_name, "
            break
        }
    }

    query = query[:len(query)-2]

    query += ` FROM events_docs
        INNER JOIN events ON events.id = events_docs.event_id
        INNER JOIN docs ON docs.id = events_docs.doc_id`

    where, params, _ := this.Where(filters, 1)
    if where != "" {
        where = " WHERE " + where
    }
    query += where

    if sidx != "" {
        query += ` ORDER BY events_docs.`+sidx
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

func (this *EventsDocsModel) GetColModel() []map[string]interface{} {
    query := `SELECT array_to_string(
        array(SELECT events.id || ':' || events.name FROM events GROUP BY events.id ORDER BY events.id), ';') as name;`
    events := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    query = `SELECT array_to_string(
        array(SELECT docs.id || ':' || docs.name FROM docs GROUP BY docs.id ORDER BY docs.id), ';') as name;`
    docs := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

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
            "index": "doc_id",
            "name": "doc_id",
            "editable": true,
            "formatter": "select",
            "edittype": "select",
            "stype": "select",
            "search": true,
            "editrules": map[string]interface{}{"required": true},
            "editoptions": map[string]string{"value": docs},
            "searchoptions": map[string]string{"value": ":Все;"+docs},
        },
    }
}
