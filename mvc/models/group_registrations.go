package models

import (
    "github.com/orc/db"
    "strconv"
)

type GroupRegistrationModel struct {
    Entity
}

type GroupRegistration struct {
    Id      int  `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    EventId int  `name:"event_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"events" refField:"id" refFieldShow:"name"`
    GroupId int  `name:"group_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"groups" refField:"id" refFieldShow:"name"`
    Status  bool `name:"status" type:"boolean" null:"NOT NULL" extra:""`
}





func (*ModelManager) GroupRegistrations() *GroupRegistrationsModel {
    model := new(GroupRegistrationsModel)
    model.SetTableName("group_registrations").
        SetCaption("Регистрации групп").
        SetColumns([]string{"id", "event_id", "group_id", "status"}).
        SetColNames([]string{"ID", "Мероприятие", "Группа", "Статус"}).
        SetFields(new(GroupRegistration)).
        SetCondition(AND).
        SetOrder("id").
        SetLimit("ALL").
        SetOffset(0).
        SetSorting("ASC").
        SetWherePart(make(map[string]interface{}, 0)).
        SetSub(true).
        SetSubTables([]string{"persons"}).
        SetSubField("group_id")

    return model
}

func (this *GroupRegistrationModel) Delete(id int) {
    query := `DELETE
        FROM param_values
        WHERE param_values.reg_id in
            (SELECT regs_groupregs.reg_id
                FROM regs_groupregs WHERE regs_groupregs.groupreg_id = $1);`
    db.Query(query, []interface{}{id})

    query = `DELETE
        FROM registrations
        WHERE registrations.id in
            (SELECT regs_groupregs.reg_id
                FROM regs_groupregs WHERE regs_groupregs.groupreg_id = $1);`
    db.Query(query, []interface{}{id})

    query = `DELETE FROM group_registrations WHERE id = $1;`
    db.Query(query, []interface{}{id})
}

func (this *GroupRegistrationModel) Select(fields []string, filters map[string]interface{}, limit, offset int, sord, sidx string) (result []interface{}) {
    if len(fields) == 0 {
        return nil
    }

    query := `SELECT `

    for _, field := range fields {
        switch field {
        case "id":
            query += "group_registrations.id, "
            break
        case "event_id":
            query += "events.name as event_name, "
            break
        case "group_id":
            query += "groups.name as group_name, "
            break
        case "status":
            query += "group_registrations.status, "
            break
        }
    }

    query = query[:len(query)-2]
    query += ` FROM group_registrations
        INNER JOIN events ON events.id = group_registrations.event_id
        INNER JOIN groups ON groups.id = group_registrations.group_id`
    where, params, _ := this.Where(filters, 1)
    if where != "" {
        where = " WHERE " + where
    }
    query += where

    if sidx != "" {
        query += ` ORDER BY group_registrations.`+sidx
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

func (this *GroupRegistrationModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
    query := `SELECT array_to_string(
        array(SELECT events.id || ':' || events.name FROM events GROUP BY events.id ORDER BY events.id), ';') as name;`
    events := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    query = `SELECT array_to_string(
        array(SELECT groups.id || ':' || groups.name FROM groups GROUP BY groups.id ORDER BY groups), ';') as name;`
    groups := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

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
            "index": "group_id",
            "name": "group_id",
            "editable": true,
            "formatter": "select",
            "edittype": "select",
            "stype": "select",
            "search": true,
            "editrules": map[string]interface{}{"required": true},
            "editoptions": map[string]string{"value": groups},
            "searchoptions": map[string]string{"value": ":Все;"+groups},
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
