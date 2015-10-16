package models

import (
    "github.com/klenin/orc/db"
    "strconv"
)

type Param struct {
    id          int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    name        string `name:"name" type:"text" null:"NOT NULL" extra:""`
    formId      int    `name:"form_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"forms" refField:"id" refFieldShow:"name"`
    paramTypeId int    `name:"param_type_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"param_types" refField:"id" refFieldShow:"name"`
    identifier  int    `name:"identifier" type:"int" null:"NOT NULL" extra:"UNIQUE"`
    required    bool   `name:"required" type:"boolean" null:"NOT NULL" extra:""`
    editable    bool   `name:"editable" type:"boolean" null:"NOT NULL" extra:""`
}

func (this *Param) GetId() int {
    return this.id
}

func (this *Param) SetName(name string) {
    this.name = name
}

func (this *Param) GetName() string {
    return this.name
}

func (this *Param) GetFormId() int {
    return this.formId
}

func (this *Param) SetFormId(formId int) {
    this.formId = formId
}

func (this *Param) SetParamTypeId(paramTypeId int) {
    this.paramTypeId = paramTypeId
}

func (this *Param) GetParamTypeId() int {
    return this.paramTypeId
}

func (this *Param) SetIdentifier(identifier int) {
    this.identifier = identifier
}

func (this *Param) GetIdentifier() int {
    return this.identifier
}

func (this *Param) SetRequired(required bool) {
    this.required = required
}

func (this *Param) GetRequired() bool {
    return this.required
}

func (this *Param) SetEditable(editable bool) {
    this.editable = editable
}

func (this *Param) GetEditable() bool {
    return this.editable
}

type ParamsModel struct {
    Entity
}

func (*ModelManager) Params() *ParamsModel {
    model := new(ParamsModel)
    model.SetTableName("params").
        SetCaption("Параметры").
        SetColumns([]string{"id", "name", "param_type_id", "form_id", "identifier", "required", "editable"}).
        SetColNames([]string{"ID", "Название", "Тип", "Форма", "Идентификатор", "Требование", "Редактирование"}).
        SetFields(new(Param)).
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

func (this *ParamsModel) Select(fields []string, filters map[string]interface{}) (result []interface{}) {
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
    query += ` ORDER BY params.` + this.orderBy
    query += ` `+ this.GetSorting()
    params = append(params, this.GetLimit())
    query += ` LIMIT $` + strconv.Itoa(len(params))
    params = append(params, this.GetOffset())
    query += ` OFFSET $` + strconv.Itoa(len(params)) + ";"

    return db.Query(query, params)
}

func (*ParamsModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
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
