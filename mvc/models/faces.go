package models

import (
    "github.com/orc/db"
    "strconv"
)

type FaceModel struct {
    Entity
}

type Face struct {
    Id     int `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    UserId int `name:"user_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"users" refField:"id" refFieldShow:"login"`
}

func (c *ModelManager) Faces() *FaceModel {
    model := new(FaceModel)

    model.TableName = "faces"
    model.Caption = "Лица"

    model.Columns = []string{"id", "user_id"}
    model.ColNames = []string{"ID", "Пользователь"}

    model.Fields = new(Face)
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

func (this *FaceModel) Select(fields []string, filters map[string]interface{}, limit, offset int, sord, sidx string) (result []interface{}) {
    if len(fields) == 0 {
        return nil
    }

    query := `SELECT `

    for _, field := range fields {
        switch field {
        case "id":
            query += "faces.id, "
            break
        case "user_id":
            query += "users.login, "
            break
        }
    }

    query += `array_to_string(array_agg(param_values.value), ' ') as name
        FROM reg_param_vals
        INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
        INNER JOIN params ON params.id = param_values.param_id
        INNER JOIN users ON users.id = faces.user_id`

    where, params := this.Where(filters)
    if where != "" {
        query += where + ` AND params.id in (5, 6, 7) GROUP BY faces.id, users.id`
    } else {
        query += ` WHERE params.id in (5, 6, 7) GROUP BY faces.id, users.id`
    }

    if sidx != "" {
        query += ` ORDER BY faces.`+sidx
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

func (this *FaceModel) GetColModel() []map[string]interface{} {
    query := `SELECT array_to_string(
        array(SELECT users.id || ':' || users.login FROM users GROUP BY users.id ORDER BY users.id), ';') as name;`
    logins := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    query = `SELECT array_to_string(
        array(SELECT faces.id || ':' || array_to_string(array_agg(param_values.value), ' ')
        FROM reg_param_vals
        INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
        INNER JOIN params ON params.id = param_values.param_id
        WHERE params.id in (5, 6, 7) GROUP BY faces.id ORDER BY faces.id), ';') as name;`

    faces := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    return []map[string]interface{} {
        0: map[string]interface{} {
            "index": "id",
            "name": "id",
            "editable": false,
            "formatter": "select",
            "edittype": "select",
            "stype": "select",
            "search": true,
            "editoptions": map[string]string{"value": faces},
            "searchoptions": map[string]string{"value": ":Все;"+faces},
        },
        1: map[string]interface{} {
            "index": "user_id",
            "name": "user_id",
            "editable": true,
            "formatter": "select",
            "edittype": "select",
            "stype": "select",
            "search": true,
            "editrules": map[string]interface{}{"required": true},
            "editoptions": map[string]string{"value": logins},
            "searchoptions": map[string]string{"value": ":Все;"+logins},
        },
    }
}
