package models

import (
    "github.com/orc/db"
    "strconv"
)

type RegistrationModel struct {
    Entity
}

type Registration struct {
    Id      int  `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    FaceId  int  `name:"face_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"faces" refField:"id" refFieldShow:"id"`
    EventId int  `name:"event_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"events" refField:"id" refFieldShow:"name"`
    Status  bool `name:"status" type:"boolean" null:"NOT NULL" extra:""`
}

func (c *ModelManager) Registrations() *RegistrationModel {
    model := new(RegistrationModel)

    model.TableName = "registrations"
    model.Caption = "Регистрации"

    model.Columns = []string{"id", "face_id", "event_id", "status"}
    model.ColNames = []string{"ID", "Лицо", "Мероприятие", "Статус"}

    model.Fields = new(Registration)
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

func (this *RegistrationModel) GetModelRefDate() (fields []string, result map[string]interface{}) {
    fields = []string{"name", "name"}

    result = make(map[string]interface{})

    result["event_id"] = db.Select(new(ModelManager).GetModel("events"), []string{"id", "name"})

    query := `SELECT faces.id as id, array_to_string(array_agg(param_values.value), ' ') as name
        FROM reg_param_vals
        INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
        INNER JOIN params ON params.id = param_values.param_id
        WHERE params.id in (5, 6, 7) GROUP BY faces.id ORDER BY faces.id;`

    result["face_id"] = db.Query(query, nil)

    return fields, result
}

func (this *RegistrationModel) Select(fields []string, filters map[string]interface{}, limit, offset int, sord, sidx string) (result []interface{}) {
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

    query += ` FROM reg_param_vals
        INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
        INNER JOIN params ON params.id = param_values.param_id`

    where, params := this.Where(filters)

    if where != "" {
        query += where + ` AND params.id in (5, 6, 7) GROUP BY registrations.id, events.id`
    } else {
        query += ` WHERE params.id in (5, 6, 7) GROUP BY registrations.id, events.id`
    }

    if sidx != "" {
        query += ` ORDER BY registrations.`+sidx
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
