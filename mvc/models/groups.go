package models

import (
    "github.com/klenin/orc/db"
    "log"
    "strconv"
)

type Group struct {
    id    int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    name  string `name:"name" type:"text" null:"NOT NULL" extra:"UNIQUE"`
    owner int    `name:"face_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"faces" refField:"id" refFieldShow:"id"`
}

func (this *Group) GetId() int {
    return this.id
}

func (this *Group) SetName(name string) {
    this.name = name
}

func (this *Group) GetName() string {
    return this.name
}

func (this *Group) SetOwner(owner int) {
    this.owner = owner
}

func (this *Group) GetOwner() int {
    return this.owner
}

type GroupsModel struct {
    Entity
}

func (*ModelManager) Groups() *GroupsModel {
    model := new(GroupsModel)
    model.SetTableName("groups").
        SetCaption("Группы").
        SetColumns([]string{"id", "name", "face_id"}).
        SetColNames([]string{"ID", "Название", "Владелец"}).
        SetFields(new(Group)).
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

func (this *GroupsModel) Update(isAdmin bool, userId int, params, where map[string]interface{}) {
    faceId := -1
    query := `SELECT groups.face_id FROM groups
        INNER JOIN faces ON faces.id = groups.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 AND groups.id = $2;`
    db.QueryRow(query, []interface{}{userId, where["id"]}).Scan(&faceId)

    if !isAdmin && faceId == -1 {
        log.Println("Нет прав редактировать эту группу")
        return
    }

    this.LoadModelData(params).LoadWherePart(where).QueryUpdate().Scan()
}

func (this *GroupsModel) Add(userId int, params map[string]interface{}) error {
    var faceId int
    query := `SELECT faces.id
        FROM registrations
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN users ON faces.user_id = users.id
        WHERE users.id = $1 AND events.id = $2;`
    db.QueryRow(query, []interface{}{userId, 1}).Scan(&faceId)
    params["face_id"] = faceId
    this.LoadModelData(params).QueryInsert("").Scan()
    return nil
}

func (*GroupsModel) Delete(id int) {
    query := `DELETE
        FROM persons
        WHERE persons.group_id = $1;`
    db.Query(query, []interface{}{id})

    query = `DELETE FROM groups WHERE id = $1;`
    db.Query(query, []interface{}{id})
}

func (this *GroupsModel) Select(fields []string, filters map[string]interface{}) (result []interface{}) {
    if len(fields) == 0 {
        return nil
    }

    query := `SELECT `

    for _, field := range fields {
        switch field {
        case "id":
            query += "groups.id, "
            break
        case "name":
            query += "groups.name as group_name, "
            break
        case "face_id":
            query += "array_to_string(array_agg(param_values.value), ' ') as face_name, "
            break
        }
    }

    query = query[:len(query)-2]
    query += ` FROM param_values
        INNER JOIN registrations ON registrations.id = param_values.reg_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN params ON params.id = param_values.param_id
        INNER JOIN groups ON groups.face_id = faces.id`
    where, params, _ := this.Where(filters, 1)
    if where != "" {
        query += ` WHERE ` + where + ` AND params.id in (5, 6, 7) AND events.id = 1 GROUP BY groups.id`
    } else {
        query += ` WHERE params.id in (5, 6, 7) AND events.id = 1 GROUP BY groups.id`
    }
    query += ` ORDER BY groups.` + this.orderBy
    query += ` `+ this.GetSorting()
    params = append(params, this.GetLimit())
    query += ` LIMIT $` + strconv.Itoa(len(params))
    params = append(params, this.GetOffset())
    query += ` OFFSET $` + strconv.Itoa(len(params)) + ";"

    return db.Query(query, params)
}

func (*GroupsModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
    var query, faces string
    if isAdmin {
        query = `SELECT array_to_string(
            array(
                SELECT f.id || ':' || f.id || '-' || array_to_string(
                array(
                    SELECT param_values.value
                    FROM param_values
                    INNER JOIN registrations ON registrations.id = param_values.reg_id
                    INNER JOIN faces ON faces.id = registrations.face_id
                    INNER JOIN events ON events.id = registrations.event_id
                    INNER JOIN params ON params.id = param_values.param_id
                    WHERE param_values.param_id IN (5, 6, 7) AND events.id = 1 AND faces.id = f.id ORDER BY param_values.param_id
                ), ' ')
                FROM param_values
                INNER JOIN registrations as reg ON reg.id = param_values.reg_id
                INNER JOIN faces as f ON f.id = reg.face_id
                INNER JOIN events ON events.id = reg.event_id
                INNER JOIN params as p ON p.id = param_values.param_id
                INNER JOIN users ON users.id = f.user_id GROUP BY f.id ORDER BY f.id
            ), ';') as name;`
    } else {
        query = `SELECT array_to_string(
            array(
                SELECT f.id || ':' || array_to_string(
                array(
                    SELECT param_values.value
                    FROM param_values
                    INNER JOIN registrations ON registrations.id = param_values.reg_id
                    INNER JOIN faces ON faces.id = registrations.face_id
                    INNER JOIN events ON events.id = registrations.event_id
                    INNER JOIN params ON params.id = param_values.param_id
                    WHERE param_values.param_id IN (5, 6, 7) AND events.id = 1 AND faces.id = f.id ORDER BY param_values.param_id
                ), ' ')
                FROM param_values
                INNER JOIN registrations as reg ON reg.id = param_values.reg_id
                INNER JOIN faces as f ON f.id = reg.face_id
                INNER JOIN events ON events.id = reg.event_id
                INNER JOIN params as p ON p.id = param_values.param_id
                INNER JOIN users ON users.id = f.user_id GROUP BY f.id ORDER BY f.id
            ), ';') as name;`
    }

    faces = db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

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
            "index": "face_id",
            "name": "face_id",
            "editable": true,
            "formatter": "select",
            "edittype": "select",
            "stype": "select",
            "search": true,
            "editrules": map[string]interface{}{"required": true},
            "editoptions": map[string]string{"value": faces},
            "searchoptions": map[string]string{"value": ":Все;"+faces},
        },
    }
}
