package models

import (
    "github.com/orc/db"
    "strconv"
)

type ParamsModel struct {
    Entity
}

type Param struct {
    Id          int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    Name        string `name:"name" type:"text" null:"NOT NULL" extra:""`
    FormId      int    `name:"form_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"forms" refField:"id" refFieldShow:"name"`
    ParamTypeId int    `name:"param_type_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"param_types" refField:"id" refFieldShow:"name"`
    Identifier  int    `name:"identifier" type:"int" null:"NOT NULL" extra:"UNIQUE"`
}

func (c *ModelManager) Params() *ParamsModel {
    model := new(ParamsModel)

    model.TableName = "params"
    model.Caption = "Параметры"

    model.Columns = []string{"id", "name", "param_type_id", "form_id", "identifier"}
    model.ColNames = []string{"ID", "Название", "Тип", "Форма", "Идентификатор"}

    model.Fields = new(Param)
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

func (this *ParamsModel) Select(fields []string, filters map[string]interface{}, limit, offset int, sord, sidx string) (result []interface{}) {
    if len(fields) == 0 {
        return nil
    }

    query := `SELECT `

    for _, field := range fields {
        switch field {
        case "id":
            query += "params.id, "
            break
        case "name":
            query += "params.name as param_name, "
            break
        case "param_type_id":
            query += "param_types.name as type_name, "
            break
        case "form_id":
            query += "forms.name as form_name, "
            break
        case "identifier":
            query += "params.identifier, "
            break
        }
    }

    query = query[:len(query)-2]

    query += ` FROM params
        INNER JOIN param_types ON param_types.id = params.param_type_id
        INNER JOIN forms ON forms.id = params.form_id`

    where, params := this.Where(filters)
    query += where

    if sidx != "" {
        query += ` ORDER BY params.`+sidx
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
