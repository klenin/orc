package models

import (
    "errors"
    "github.com/klenin/orc/db"
    "github.com/klenin/orc/mailer"
    "github.com/klenin/orc/utils"
    "log"
    "strconv"
)

const HASH_SIZE = 32

type Person struct {
    id      int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    faceId  int    `name:"face_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"faces" refField:"id" refFieldShow:"id"`
    groupId int    `name:"group_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"groups" refField:"id" refFieldShow:"name"`
    token   string `name:"token" type:"text" null:"NOT NULL" extra:""`
    status  bool   `name:"status" type:"boolean" null:"NOT NULL" extra:""`
}

func (this *Person) GetId() int {
    return this.id
}

func (this *Person) GetFaceId() int {
    return this.faceId
}

func (this *Person) SetFaceId(faceId int) {
    this.faceId = faceId
}

func (this *Person) GetGroupId() int {
    return this.groupId
}

func (this *Person) SetGroupId(groupId int) {
    this.groupId = groupId
}

func (this *Person) SetToken(token string) {
    this.token = token
}

func (this *Person) GetToken() string {
    return this.token
}

func (this *Person) SetStatus(status bool) {
    this.status = status
}

func (this *Person) GetStatus() bool {
    return this.status
}

func (*ModelManager) Persons() *PersonsModel {
    model := new(PersonsModel)
    model.SetTableName("persons").
        SetCaption("Участники").
        SetColumns([]string{"id", "face_id", "group_id", "status"}).
        SetColNames([]string{"ID", "Физическое лицо", "Группа", "Статус"}).
        SetFields(new(Person)).
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

type PersonsModel struct {
    Entity
}

func (this *PersonsModel) Add(userId int, params map[string]interface{}) error {
    var to, address string

    token := utils.GetRandSeq(HASH_SIZE)
    params["token"] = token

    log.Println("face_id: ", params["face_id"])
    log.Println("group_id: ", params["group_id"])
    log.Println("status: ", params["status"])

    if db.IsExists("persons", []string{"face_id", "group_id"}, []interface{}{params["face_id"], params["group_id"]}) {
        return errors.New("Участник уже состоит в группе")
    }

    query := `SELECT param_values.value
        FROM param_values
        INNER JOIN registrations ON registrations.id = param_values.reg_id
        INNER JOIN params ON params.id = param_values.param_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE params.id in (5, 6, 7) AND users.id = $1 AND events.id = 1 ORDER BY params.id;`
    data := db.Query(query, []interface{}{userId})
    headName := ""
    if len(data) < 3 {
        return errors.New("Данные о руководителе группы отсутсвуют")

    } else {
        headName = data[0].(map[string]interface{})["value"].(string)
        headName += " " + data[1].(map[string]interface{})["value"].(string)
        headName += " " + data[2].(map[string]interface{})["value"].(string)
    }

    groupId, err := strconv.Atoi(params["group_id"].(string))
    if err != nil {
        return err
    }

    var groupName string
    db.QueryRow("SELECT name FROM groups WHERE id = $1;", []interface{}{groupId}).Scan(&groupName)

    query = `SELECT param_values.value
        FROM param_values
        INNER JOIN registrations ON registrations.id = param_values.reg_id
        INNER JOIN params ON params.id = param_values.param_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE params.id in (4, 5, 6, 7) AND faces.id = $1 AND events.id = 1 ORDER BY params.id;`
    data = db.Query(query, []interface{}{params["face_id"]})
    if len(data) < 4 {
        return errors.New("Данные о приглашаемом участнике отсутсвуют.")

    } else {
        address = data[0].(map[string]interface{})["value"].(string)
        to = data[1].(map[string]interface{})["value"].(string)
        to += " " + data[2].(map[string]interface{})["value"].(string)
        to += " " + data[3].(map[string]interface{})["value"].(string)
    }

    if !mailer.InviteToGroup(to, address, token, headName, groupName) {
        return errors.New("Участник скорее всего указал неправильный email, отправить письмо-приглашенине невозможно")
    }

    this.LoadModelData(params).QueryInsert("").Scan()
    return nil
}

func (this *PersonsModel) Select(fields []string, filters map[string]interface{}) (result []interface{}) {
    if len(fields) == 0 {
        return nil
    }

    query := `SELECT `

    for _, field := range fields {
        switch field {
        case "id":
            query += "persons.id, "
            break
        case "group_id":
            query += "groups.name as group_name, "
            break
        case "status":
            query += "persons.status, "
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
        INNER JOIN persons ON persons.face_id = faces.id
        INNER JOIN groups ON groups.id = persons.group_id`
    where, params, _ := this.Where(filters, 1)
    if where != "" {
        query += ` WHERE ` + where + ` AND params.id in (5, 6, 7) AND events.id = 1 GROUP BY persons.id, groups.id`
    } else {
        query += ` WHERE params.id in (5, 6, 7) AND events.id = 1 GROUP BY persons.id, groups.id`
    }
    query += ` ORDER BY persons.` + this.orderBy
    query += ` `+ this.GetSorting()
    params = append(params, this.GetLimit())
    query += ` LIMIT $` + strconv.Itoa(len(params))
    params = append(params, this.GetOffset())
    query += ` OFFSET $` + strconv.Itoa(len(params)) + ";"

    return db.Query(query, params)
}

func (*PersonsModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
    var query, groups, faces string

    if isAdmin {
        query = `SELECT array_to_string(
            array(SELECT groups.id || ':' || groups.name
            FROM groups
            GROUP BY groups.id ORDER BY groups), ';') as name;`
        groups = db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

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
                GROUP BY f.id ORDER BY f.id
            ), ';') as name;`

        faces = db.Query(query, nil)[0].(map[string]interface{})["name"].(string)
    } else {
        query = `SELECT array_to_string(
            array(SELECT groups.id || ':' || groups.name FROM groups
            INNER JOIN faces ON faces.id = groups.face_id
            INNER JOIN users ON users.id = faces.user_id
            WHERE users.id = $1 AND groups.id NOT IN (SELECT group_registrations.group_id FROM group_registrations)
            GROUP BY groups.id ORDER BY groups), ';') as name;`
        groups = db.Query(query, []interface{}{userId})[0].(map[string]interface{})["name"].(string)

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
                GROUP BY f.id ORDER BY f.id
            ), ';') as name;`

        faces = db.Query(query, nil)[0].(map[string]interface{})["name"].(string)
    }

    return []map[string]interface{} {
        0: map[string]interface{} {
            "index": "id",
            "name": "id",
            "editable": false,
        },
        1: map[string]interface{} {
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
