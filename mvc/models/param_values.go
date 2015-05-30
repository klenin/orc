package models

import (
    "github.com/orc/db"
    "strconv"
)

type ParamValuesModel struct {
    Entity
}

type ParamValues struct {
    Id      int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    ParamId int    `name:"param_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"params" refField:"id" refFieldShow:"name"`
    Value   string `name:"value" type:"text" null:"NULL" extra:""`
    Date    string `name:"date" type:"timestamp" null:"NOT NULL" extra:""`
}

func (c *ModelManager) ParamValues() *ParamValuesModel {
    model := new(ParamValuesModel)

    model.TableName = "param_values"
    model.Caption = "Значение параметров"

    model.Columns = []string{"id", "param_id", "value", "date"}
    model.ColNames = []string{"ID", "Параметр", "Значение", "Дата"}

    model.Fields = new(ParamValues)
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

func (this *ParamValuesModel) Select(fields []string, filters map[string]interface{}, limit, offset int, sord, sidx string) (result []interface{}) {
    if len(fields) == 0 {
        return nil
    }

    query := `SELECT `

    for _, field := range fields {
        switch field {
        case "id":
            query += "param_values.id, "
            break
        case "param_id":
            query += "forms.name || ': ' || params.name as name, "
            break
        case "value":
            query += "param_values.value, "
            break
        case "date":
            query += "param_values.date, "
            break
        }
    }

    query = query[:len(query)-2]

    query += ` FROM param_values
        INNER JOIN params ON params.id = param_values.param_id
        INNER JOIN forms ON forms.id = params.form_id`

    where, params, _ := this.Where(filters, 1)
    if where != "" {
        where = " WHERE " + where
    }
    query += where

    if sidx != "" {
        query += ` ORDER BY param_values.`+sidx
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

func (this *ParamValuesModel) GetColModel() []map[string]interface{} {
    query := `SELECT array_to_string(
        array(SELECT params.id || ': ' || forms.name || ' - ' || params.name
        FROM params
        INNER JOIN forms ON forms.id = params.form_id GROUP BY params.id, forms.name ORDER BY params.id), ';') as name;`
    params := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    return []map[string]interface{} {
        0: map[string]interface{} {
            "index": "id",
            "name": "id",
            "editable": false,
            "width": 20,
        },
        1: map[string]interface{} {
            "index": "param_id",
            "name": "param_id",
            "editable": true,
            "formatter": "select",
            "edittype": "select",
            "stype": "select",
            "search": true,
            "editrules": map[string]interface{}{"required": true},
            "editoptions": map[string]string{"value": params},
            "searchoptions": map[string]string{"value": ":Все;"+params},
        },
        2: map[string]interface{} {
            "index": "value",
            "name": "value",
            "editable": true,
            "editrules": map[string]interface{}{"required": true},
        },
        3: map[string]interface{} {
            "index": "date",
            "name": "date",
            "editable": true,
            "formatter": nil,
            "editrules": map[string]interface{}{"date": true, "required": true},
            "editoptions": map[string]interface{}{"dataInit": nil},
            "formatoptions": map[string]string{"srcformat": "Y-m-d", "newformat": "Y-m-d"},
            "searchoptions": map[string]interface{}{"sopt": []string{"eq", "ne"}, "dataInit": nil},
            "type": "date",
        },
    }
}
