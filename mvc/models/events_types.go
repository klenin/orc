package models

import (
    "github.com/klenin/orc/db"
    "strconv"
)

type EventType struct {
    id      int `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    eventId int `name:"event_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"events" refField:"id" refFieldShow:"name"`
    typeId  int `name:"type_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"event_types" refField:"id" refFieldShow:"name"`
}

func (this *EventType) GetId() int {
    return this.id
}

func (this *EventType) GetEventId() int {
    return this.eventId
}

func (this *EventType) SetEventId(eventId int) {
    this.eventId = eventId
}

func (this *EventType) GetTypeId() int {
    return this.typeId
}

func (this *EventType) SetTypeId(typeId int) {
    this.typeId = typeId
}

type EventsTypesModel struct {
    Entity
}

func (*ModelManager) EventsTypes() *EventsTypesModel {
    model := new(EventsTypesModel)
    model.SetTableName("events_types").
        SetCaption("Мероприятия - Типы").
        SetColumns([]string{"id", "event_id", "type_id"}).
        SetColNames([]string{"ID", "Мероприятие", "Тип"}).
        SetFields(new(EventType)).
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

func (this *EventsTypesModel) Select(fields []string, filters map[string]interface{}) (result []interface{}) {
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
    where, params, _ := this.Where(filters, 1)
    if where != "" {
        where = " WHERE " + where
    }
    query += where
    query += ` ORDER BY events_types.` + this.orderBy
    query += ` `+ this.GetSorting()
    params = append(params, this.GetLimit())
    query += ` LIMIT $` + strconv.Itoa(len(params))
    params = append(params, this.GetOffset())
    query += ` OFFSET $` + strconv.Itoa(len(params)) + `;`

    return db.Query(query, params)
}

func (*EventsTypesModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
    query := `SELECT array_to_string(
        array(SELECT events.id || ':' || events.name FROM events GROUP BY events.id ORDER BY events.id), ';') as name;`
    events := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    query = `SELECT array_to_string(
        array(SELECT event_types.id || ':' || event_types.name FROM event_types GROUP BY event_types.id ORDER BY event_types.id), ';') as name;`
    types := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

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
            "index": "type_id",
            "name": "type_id",
            "editable": true,
            "formatter": "select",
            "edittype": "select",
            "stype": "select",
            "search": true,
            "editrules": map[string]interface{}{"required": true},
            "editoptions": map[string]string{"value": types},
            "searchoptions": map[string]string{"value": ":Все;"+types},
        },
    }
}
