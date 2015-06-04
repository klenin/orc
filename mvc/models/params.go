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
    Required    bool   `name:"required" type:"boolean" null:"NOT NULL" extra:""`
    Editable    bool   `name:"editable" type:"boolean" null:"NOT NULL" extra:""`

}

func (c *ModelManager) Params() *ParamsModel {
    model := new(ParamsModel)

    model.TableName = "params"
    model.Caption = "Параметры"

    model.Columns = []string{"id", "name", "param_type_id", "form_id", "identifier", "required", "editable"}
    model.ColNames = []string{"ID", "Название", "Тип", "Форма", "Идентификатор", "Требование", "Редактирование"}

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
        case "required":
            query += "params.required, "
            break
        }
    }

    query = query[:len(query)-2]

    query += ` FROM params
        INNER JOIN param_types ON param_types.id = params.param_type_id
        INNER JOIN forms ON forms.id = params.form_id`

    where, params, _ := this.Where(filters, 1)
    if where != "" {
        where = " WHERE " + where
    }
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

func (this *ParamsModel) GetColModel() []map[string]interface{} {
    query := `SELECT array_to_string(
        array(SELECT param_types.id || ':' || param_types.name FROM param_types GROUP BY param_types.id ORDER BY param_types.id), ';') as name;`
    types := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    query = `SELECT array_to_string(
        array(SELECT forms.id || ':' || forms.name FROM forms GROUP BY forms.id ORDER BY forms), ';') as name;`
    forms := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    return []map[string]interface{} {
        0: map[string]interface{} {
            "index": "id",
            "name": "id",
            "editable": false,
            "width": 20,
        },
        1: map[string]interface{} {
            "index": "name",
            "name": "name",
            "editable": true,
            "editrules": map[string]interface{}{"required": true},
        },
        2: map[string]interface{} {
            "index": "param_type_id",
            "name": "param_type_id",
            "editable": true,
            "formatter": "select",
            "edittype": "select",
            "stype": "select",
            "search": true,
            "editrules": map[string]interface{}{"required": true},
            "editoptions": map[string]string{"value": types},
            "searchoptions": map[string]string{"value": ":Все;"+types},
        },
        3: map[string]interface{} {
            "index": "form_id",
            "name": "form_id",
            "editable": true,
            "formatter": "select",
            "edittype": "select",
            "stype": "select",
            "search": true,
            "editrules": map[string]interface{}{"required": true},
            "editoptions": map[string]string{"value": forms},
            "searchoptions": map[string]string{"value": ":Все;"+forms},
        },
        4: map[string]interface{} {
            "index": "identifier",
            "name": "identifier",
            "editable": true,
            "editrules": map[string]interface{}{"required": true},
        },
        5: map[string]interface{} {
            "index": "required",
            "name": "required",
            "editable": true,
            "editrules": map[string]interface{}{"required": true},
            "formatter": "checkbox",
            "formatoptions": map[string]interface{}{"disabled": true},
            "edittype": "checkbox",
            "editoptions": map[string]interface{}{"value": "true:false"},
        },
        6: map[string]interface{} {
            "index": "editable",
            "name": "editable",
            "editable": true,
            "editrules": map[string]interface{}{"required": true},
            "formatter": "checkbox",
            "formatoptions": map[string]interface{}{"disabled": true},
            "edittype": "checkbox",
            "editoptions": map[string]interface{}{"value": "true:false"},
        },
    }
}
