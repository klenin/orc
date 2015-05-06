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
}

func (c *ModelManager) ParamValues() *ParamValuesModel {
    model := new(ParamValuesModel)

    model.TableName = "param_values"
    model.Caption = "Значение параметров"

    model.Columns = []string{"id", "param_id", "value"}
    model.ColNames = []string{"ID", "Параметр", "Значение"}

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

func (this *ParamValuesModel) GetModelRefDate() (fields []string, result map[string]interface{}) {
    fields = []string{"name"}

    query := `SELECT params.id, forms.name || ': ' || params.name as name
        FROM params
        INNER JOIN forms ON forms.id = params.form_id ORDER BY params.id`

    return fields, map[string]interface{}{"param_id": db.Query(query, nil)}
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
        }
    }

    query = query[:len(query)-2]

    query += ` FROM param_values
        INNER JOIN params ON params.id = param_values.param_id
        INNER JOIN forms ON forms.id = params.form_id`

    where, params := this.Where(filters)
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
