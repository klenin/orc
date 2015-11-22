package models

import (
    "github.com/klenin/orc/db"
    "strconv"
)

type ParamValue struct {
    id      int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    paramId int    `name:"param_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"params" refField:"id" refFieldShow:"name"`
    value   string `name:"value" type:"text" null:"NULL" extra:""`
    regId   int    `name:"reg_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"registrations" refField:"id" refFieldShow:"id"`
    date    string `name:"date" type:"timestamp" null:"NOT NULL" extra:""`
    userId  int    `name:"user_id" type:"int" null:"NULL" extra:"REFERENCES" refTable:"users" refField:"id" refFieldShow:"login"`
}

func (this *ParamValue) GetId() int {
    return this.id
}

func (this *ParamValue) SetParamId(paramId int) {
    this.paramId = paramId
}

func (this *ParamValue) GetParamId() int {
    return this.paramId
}

func (this *ParamValue) SetValue(value string) {
    this.value = value
}

func (this *ParamValue) GetValue() string {
    return this.value
}

func (this *ParamValue) SetRegId(regId int) {
    this.regId = regId
}

func (this *ParamValue) GetRegId() int {
    return this.regId
}

func (this *ParamValue) SetDate(date string) {
    this.date = date
}

func (this *ParamValue) GetDate() string {
    return this.date
}

func (this *ParamValue) SetUserId(userId int) {
    this.userId = userId
}

func (this *ParamValue) GetUserId() int {
    return this.userId
}

type ParamValuesModel struct {
    Entity
}

func (*ModelManager) ParamValues() *ParamValuesModel {
    model := new(ParamValuesModel)
    model.SetTableName("param_values").
        SetCaption("Значение параметров").
        SetColumns([]string{"id", "param_id", "value", "reg_id", "date", "user_id"}).
        SetColNames([]string{"ID", "Параметр", "Значение", "Регистрация", "Дата", "Кто редактировал"}).
        SetFields(new(ParamValue)).
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

func (this *ParamValuesModel) Select(fields []string, filters map[string]interface{}) (result []interface{}) {
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
        case "reg_id":
            query += "registrations.id || ' - ' || events.name, "
            break
        case "value":
            query += "param_values.value, "
            break
        case "date":
            query += "param_values.date, "
            break
        case "user_id":
            query += "users.login, "
            break
        }
    }

    query = query[:len(query)-2]
    query += ` FROM registrations
        INNER JOIN param_values ON param_values.reg_id = registrations.id
        INNER JOIN params ON params.id = param_values.param_id
        INNER JOIN forms ON forms.id = params.form_id
        INNER JOIN users ON users.id = param_values.user_id
        INNER JOIN events ON events.id = registrations.event_id`
    where, params, _ := this.Where(filters, 1)
    if where != "" {
        where = " WHERE " + where
    }
    query += ` ORDER BY param_values.` + this.orderBy
    query += ` `+ this.GetSorting()
    params = append(params, this.GetLimit())
    query += ` LIMIT $` + strconv.Itoa(len(params))
    params = append(params, this.GetOffset())
    query += ` OFFSET $` + strconv.Itoa(len(params)) + ";"

    return db.Query(query, params)
}

func (*ParamValuesModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
    var query, params string

    if isAdmin {
        query = `SELECT array_to_string(
            array(SELECT params.id || ': ' || forms.name || ' - ' || params.name
            FROM params
            INNER JOIN forms ON forms.id = params.form_id GROUP BY params.id, forms.name ORDER BY params.id), ';') as name;`
        params = db.Query(query, nil)[0].(map[string]interface{})["name"].(string)
    } else {
        query = `SELECT array_to_string(
            array(SELECT params.id || ': ' || forms.name || ' - ' || params.name
            FROM params
            INNER JOIN forms ON forms.id = params.form_id
            WHERE params.id IN (4, 5, 6, 7)
            GROUP BY params.id, forms.name
            ORDER BY params.id), ';') as name;`
        params = db.Query(query, nil)[0].(map[string]interface{})["name"].(string)
    }

    query = `SELECT array_to_string(
        array(SELECT users.id || ':' || users.login FROM users GROUP BY users.id ORDER BY users.id), ';') as name;`
    logins := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    query = `SELECT array_to_string(
        array(SELECT registrations.id || ':' || registrations.id || ' - ' || events.name FROM registrations
        INNER JOIN events ON events.id = registrations.event_id
        GROUP BY registrations.id, events.name ORDER BY registrations.id), ';') as name;`
    regs := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    var hideUser bool

    if isAdmin {
        hideUser = false
    } else {
        hideUser = true
    }

    return []map[string]interface{} {
        0: map[string]interface{} {
            "index": "id",
            "name": "id",
            "editable": false,
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
            "index": "reg_id",
            "name": "reg_id",
            "editable": true,
            "formatter": "select",
            "edittype": "select",
            "stype": "select",
            "search": true,
            "editrules": map[string]interface{}{"required": true},
            "editoptions": map[string]string{"value": regs},
            "searchoptions": map[string]string{"value": ":Все;"+regs},
        },
        4: map[string]interface{} {
            "index": "date",
            "name": "date",
            "editable": true,
            "formatter": nil,
            "editrules": map[string]interface{}{"required": true},
            "editoptions": map[string]interface{}{"dataInit": nil},
            "formatoptions": map[string]string{"srcformat": "Y-m-d", "newformat": "Y-m-d"},
            "searchoptions": map[string]interface{}{"sopt": []string{"eq", "ne"}, "dataInit": nil},
            "type": "timestamp",
        },
        5: map[string]interface{} {
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
            "hidden": hideUser,
        },
    }
}
