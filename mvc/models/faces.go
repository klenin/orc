package models

import (
    "github.com/klenin/orc/db"
    "log"
    "strconv"
    "strings"
)

type Face struct {
    id     int `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    userId int `name:"user_id" type:"int" null:"NULL" extra:"REFERENCES" refTable:"users" refField:"id" refFieldShow:"login"`
}

func (this *Face) GetId() int {
    return this.id
}

func (this *Face) GetUserId() int {
    return this.userId
}

func (this *Face) SetUserId(userId int) {
    this.userId = userId
}

type FacesModel struct {
    Entity
}

func (*ModelManager) Faces() *FacesModel {
    model := new(FacesModel)
    model.SetTableName("faces").
        SetCaption("Физические лица").
        SetColumns([]string{"id", "user_id"}).
        SetColNames([]string{"ID", "Пользователь"}).
        SetFields(new(Face)).
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

func (this *FacesModel) Select(fields []string, filters map[string]interface{}) (result []interface{}) {
    if len(fields) == 0 {
        return nil
    }

    query := `SELECT `

    for _, field := range fields {
        switch field {
        case "id":
            query += "faces.id, "
            break
        case "user_id":
            query += "users.login, "
            break
        }
    }

    query = query[:len(query)-2]
    query += ` FROM faces INNER JOIN users ON users.id = faces.user_id`
    where, params, _ := this.Where(filters, 1)
    if where != "" {
        query += ` WHERE ` + where
    }
    query += ` ORDER BY faces.` + this.orderBy
    query += ` `+ this.GetSorting()
    params = append(params, this.GetLimit())
    query += ` LIMIT $` + strconv.Itoa(len(params))
    params = append(params, this.GetOffset())
    query += ` OFFSET $` + strconv.Itoa(len(params)) + `;`

    return db.Query(query, params)
}

func (*FacesModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
    query := `SELECT array_to_string(
        array(SELECT users.id || ':' || users.login FROM users GROUP BY users.id ORDER BY users.id), ';') as name;`
    logins := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

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

    faces := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    return []map[string]interface{} {
        0: map[string]interface{} {
            "index": "id",
            "name": "id",
            "editable": false,
            "formatter": "select",
            "edittype": "select",
            "stype": "select",
            "search": true,
            "editoptions": map[string]string{"value": faces},
            "searchoptions": map[string]string{"value": ":Все;"+faces},
        },
        1: map[string]interface{} {
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
        },
    }
}

func (this *FacesModel) WhereByParams(filters map[string]interface{}, num int) (where string, params []interface{}, num1 int) {
    where = ""
    if filters == nil {
        return where, nil, -1
    }
    i := num

    groupOp := filters["groupOp"].(string)
    rules := filters["rules"].([]interface{})
    var groups []interface{}
    if filters["groups"] != nil {
        groups = filters["groups"].([]interface{})
    }

    if len(rules) > 10 {
        log.Println("More 10 rules for serching!")
    }

    firstElem := true

    model := new(ModelManager).GetModel("param_values")

    for _, v := range rules {
        if !firstElem {
            if groupOp != "AND" && groupOp != "OR" {
                log.Println("`groupOp` parameter is not allowed!")
                continue
            }
            where += " " + groupOp + " "
        } else {
            firstElem = false
        }

        rule := v.(map[string]interface{})

        switch rule["op"].(string) {
        case "eq":// equal
            where += model.GetTableName()+"."+rule["field"].(string) + "::text = $"+strconv.Itoa(i)
            params = append(params, rule["data"])
            i += 1
            break
        case "ne":// not equal
            where += model.GetTableName()+"."+rule["field"].(string) + "::text <> $"+strconv.Itoa(i)
            params = append(params, rule["data"])
            i += 1
            break
        case "bw":// begins with
            where += model.GetTableName()+"."+rule["field"].(string) + "::text LIKE $"+strconv.Itoa(i)+"||'%'"
            params = append(params, rule["data"])
            i += 1
            break
        case "bn":// does not begin with
            where += model.GetTableName()+"."+rule["field"].(string) + "::text NOT LIKE $"+strconv.Itoa(i)+"||'%'"
            params = append(params, rule["data"])
            i += 1
            break
        case "ew":// ends with
            where += model.GetTableName()+"."+rule["field"].(string) + "::text LIKE '%'||$"+strconv.Itoa(i)
            params = append(params, rule["data"])
            i += 1
            break
        case "en":// does not end with
            where += model.GetTableName()+"."+rule["field"].(string) + "::text NOT LIKE '%'||$"+strconv.Itoa(i)
            params = append(params, rule["data"])
            i += 1
            break
        case "cn":// contains
            where += model.GetTableName()+"."+rule["field"].(string) + "::text LIKE '%'||$"+strconv.Itoa(i)+"||'%'"
            params = append(params, rule["data"])
            i += 1
            break
        case "nc":// does not contain
            where += model.GetTableName()+"."+rule["field"].(string) + "::text NOT LIKE '%'||$"+strconv.Itoa(i)+"||'%'"
            params = append(params, rule["data"])
            i += 1
            break
        case "nu":// is null
            where += model.GetTableName()+"."+rule["field"].(string) + "::text IS NULL"
            break
        case "nn":// is not null
            where += model.GetTableName()+"."+rule["field"].(string) + "::text IS NOT NULL"
            break
        case "in":// is in
            where += model.GetTableName()+"."+rule["field"].(string) + "::text IN ("
            result := strings.Split(rule["data"].(string), ",")
            for k := range result {
                where += "$"+strconv.Itoa(i)+", "
                params = append(params, result[k])
                i += 1
            }
            where = where[:len(where)-2]
            where += ")"
            break
        case "ni":// is not in
            where += model.GetTableName()+"."+rule["field"].(string) + "::text NOT IN ("
            result := strings.Split(rule["data"].(string), ",")
            for k := range result {
                where += "$"+strconv.Itoa(i)+", "
                params = append(params, result[k])
                i += 1
            }
            where = where[:len(where)-2]
            where += ")"
            break
        default:
            panic("`op` parameter is not allowed!")
        }
    }

    for _, v := range groups {
        filters1 := v.(map[string]interface{})
        where1, params1, num1 :=  this.WhereByParams(filters1, i)
        i = num1
        if !firstElem {
            if groupOp != "AND" && groupOp != "OR" {
                log.Println("`groupOp` parameter is not allowed!")
                continue
            }
            where += " " + groupOp + " "
        } else {
            firstElem = false
        }
        where += "faces.id in ("
        where += `SELECT faces.id
            FROM param_values
            INNER JOIN registrations ON registrations.id = param_values.reg_id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN events ON events.id = registrations.event_id
            INNER JOIN params ON params.id = param_values.param_id
            INNER JOIN users ON users.id = faces.user_id WHERE `
        where += where1 + ")"
        params = append(params, params1...)
    }

    return where, params, i
}
