package models

import (
    "github.com/klenin/orc/db"
    "strconv"
)

type Registration struct {
    id      int  `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    faceId  int  `name:"face_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"faces" refField:"id" refFieldShow:"id"`
    eventId int  `name:"event_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"events" refField:"id" refFieldShow:"name"`
    status  bool `name:"status" type:"boolean" null:"NOT NULL" extra:""`
}

func (this *Registration) GetId() int {
    return this.id
}

func (this *Registration) GetFaceId() int {
    return this.faceId
}

func (this *Registration) SetFaceId(faceId int) {
    this.faceId = faceId
}

func (this *Registration) GetEventId() int {
    return this.eventId
}

func (this *Registration) SetEventId(eventId int) {
    this.eventId = eventId
}

func (this *Registration) SetStatus(status bool) {
    this.status = status
}

func (this *Registration) GetStatus() bool {
    return this.status
}

type RegistrationsModel struct {
    Entity
}

func (*ModelManager) Registrations() *RegistrationsModel {
    model := new(RegistrationsModel)
    model.SetTableName("registrations").
        SetCaption("Регистрации").
        SetColumns([]string{"id", "face_id", "event_id", "status"}).
        SetColNames([]string{"ID", "Лицо", "Мероприятие", "Статус"}).
        SetFields(new(Registration)).
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

func (*RegistrationsModel) Delete(id int) {
    query := `DELETE FROM param_values WHERE param_values.reg_id = $1;`
    db.Query(query, []interface{}{id})

    // query = `DELETE
    //     FROM faces
    //     WHERE faces.id in
    //     (SELECT registrations.face_id
    //         FROM registrations WHERE registrations.id = $1);`
    // db.Query(query, []interface{}{id})

    query = `DELETE FROM registrations WHERE id = $1;`
    db.Query(query, []interface{}{id})
}

func (this *RegistrationsModel) Select(fields []string, filters map[string]interface{}) (result []interface{}) {
    if len(fields) == 0 {
        return nil
    }

    query := `SELECT `

    for _, field := range fields {
        switch field {
        case "id":
            query += "registrations.id, "
            break
        case "event_id":
            query += "events.name as event_name, "
            break
        case "face_id":
            query += "array_to_string(array_agg(param_values.value), ' ') as face_name, "
            break
        case "status":
            query += "registrations.status, "
            break
        }
    }

    query = query[:len(query)-2]
    query += ` FROM param_values
        INNER JOIN registrations ON registrations.id = param_values.reg_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN params ON params.id = param_values.param_id`
    where, params, _ := this.Where(filters, 1)
    if where != "" {
        query += ` WHERE ` + where + ` AND params.id in (5, 6, 7)  AND events.id = 1 GROUP BY registrations.id, events.id`
    } else {
        query += ` WHERE params.id in (5, 6, 7) AND events.id = 1 GROUP BY registrations.id, events.id`
    }
    query += ` ORDER BY registrations.` + this.orderBy
    query += ` `+ this.GetSorting()
    params = append(params, this.GetLimit())
    query += ` LIMIT $` + strconv.Itoa(len(params))
    params = append(params, this.GetOffset())
    query += ` OFFSET $` + strconv.Itoa(len(params)) + ";"

    return db.Query(query, params)
}

func (*RegistrationsModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
    query := `SELECT array_to_string(
        array(SELECT events.id || ':' || events.name FROM events GROUP BY events.id ORDER BY events.id), ';') as name;`
    events := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    query = `SELECT array_to_string(
        array(
            SELECT f.id || ':' || f.id || '-' || array_to_string(
            array(
                SELECT param_values.value
                FROM param_values
                INNER JOIN registrations ON registrations.id = param_values.reg_id
                INNER JOIN faces ON faces.id = registrations.face_id
                INNER JOIN events ON events.id = registrations.event_id
                INNER JOIN params ON params.id = param_values.param_id
                WHERE param_values.param_id IN (5, 6, 7) AND events.id = 1 AND faces.id = f.id ORDER BY param_values.param_id
            ), ' ')
            FROM param_values
            INNER JOIN registrations as reg ON reg.id = param_values.reg_id
            INNER JOIN faces as f ON f.id = reg.face_id
            INNER JOIN events ON events.id = reg.event_id
            INNER JOIN params as p ON p.id = param_values.param_id
            GROUP BY f.id ORDER BY f.id
        ), ';') as name;`

    faces := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    return []map[string]interface{} {
        0: map[string]interface{} {
            "index": "id",
            "name": "id",
            "editable": false,
        },
        1: map[string]interface{} {
            "index": "face_id",
            "name": "face_id",
            "editable": true,
            "formatter": "select",
            "edittype": "select",
            "stype": "select",
            "search": true,
            "editrules": map[string]interface{}{"required": true},
            "editoptions": map[string]string{"value": faces},
            "searchoptions": map[string]string{"value": ":Все;"+faces},
        },
        2: map[string]interface{} {
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
        3: map[string]interface{} {
            "index": "status",
            "name": "status",
            "editable": true,
            "editrules": map[string]interface{}{"required": true},
            "formatter": "checkbox",
            "formatoptions": map[string]interface{}{"disabled": true},
            "edittype": "checkbox",
            "editoptions": map[string]interface{}{"value": "true:false"},
        },
    }
}
