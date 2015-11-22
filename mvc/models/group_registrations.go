package models

import (
    "github.com/klenin/orc/db"
    "strconv"
)

type GroupRegistration struct {
    id      int  `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    eventId int  `name:"event_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"events" refField:"id" refFieldShow:"name"`
    groupId int  `name:"group_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"groups" refField:"id" refFieldShow:"name"`
    status  bool `name:"status" type:"boolean" null:"NOT NULL" extra:""`
}

func (this *GroupRegistration) GetId() int {
    return this.id
}

func (this *GroupRegistration) GetEventId() int {
    return this.eventId
}

func (this *GroupRegistration) SetEventId(eventId int) {
    this.eventId = eventId
}

func (this *GroupRegistration) GetGroupId() int {
    return this.groupId
}

func (this *GroupRegistration) SetGroupId(groupId int) {
    this.groupId = groupId
}

func (this *GroupRegistration) SetStatus(status bool) {
    this.status = status
}

func (this *GroupRegistration) GetStatus() bool {
    return this.status
}

type GroupRegistrationsModel struct {
    Entity
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

func (*GroupRegistrationsModel) Delete(id int) {
    // TODO: TT
    query := `with TT AS (
        DELETE
        FROM registrations
        WHERE registrations.id in (
            SELECT rs_gs.reg_id
            FROM regs_groupregs rs_gs
            WHERE rs_gs.groupreg_id = $1
            AND array_length(array(
                SELECT regs_groupregs.groupreg_id
                FROM regs_groupregs
                WHERE regs_groupregs.reg_id = rs_gs.reg_id
            ), 1) = 1
        ) returning registrations.id
    )
    DELETE FROM param_values WHERE param_values.reg_id in (
        SELECT rs_gs.reg_id
        FROM regs_groupregs rs_gs
        WHERE rs_gs.groupreg_id = $1
        AND array_length(array(
            SELECT regs_groupregs.groupreg_id
            FROM regs_groupregs
            WHERE regs_groupregs.reg_id = rs_gs.reg_id
        ), 1) = 1
    );`
    db.Query(query, []interface{}{id})

    query = `DELETE FROM group_registrations WHERE id = $1;`
    db.Query(query, []interface{}{id})
}

func (this *GroupRegistrationsModel) Select(fields []string, filters map[string]interface{}) (result []interface{}) {
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
    query += ` ORDER BY group_registrations.` + this.orderBy
    query += ` `+ this.GetSorting()
    params = append(params, this.GetLimit())
    query += ` LIMIT $` + strconv.Itoa(len(params))
    params = append(params, this.GetOffset())
    query += ` OFFSET $` + strconv.Itoa(len(params)) + ";"

    return db.Query(query, params)
}

func (*GroupRegistrationsModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
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
