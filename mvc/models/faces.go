package models

import (
    "github.com/orc/db"
    "log"
    "strconv"
    "strings"
)

type FaceModel struct {
    Entity
}

type Face struct {
    Id     int `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    UserId int `name:"user_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"users" refField:"id" refFieldShow:"login"`
}

func (c *ModelManager) Faces() *FaceModel {
    model := new(FaceModel)

    model.TableName = "faces"
    model.Caption = "Физические лица"

    model.Columns = []string{"id", "user_id"}
    model.ColNames = []string{"ID", "Пользователь"}

    model.Fields = new(Face)
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

func (this *FaceModel) Select(fields []string, filters map[string]interface{}, limit, offset int, sord, sidx string) (result []interface{}) {
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

    query += `array_to_string(array_agg(param_values.value), ' ') as name
        FROM reg_param_vals
        INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
        INNER JOIN params ON params.id = param_values.param_id
        INNER JOIN users ON users.id = faces.user_id`

    where, params, _ := this.Where(filters, 1)

    if where != "" {
        query += ` WHERE ` + where + ` AND params.id in (5, 6, 7) GROUP BY faces.id, users.id`
    } else {
        query += ` WHERE params.id in (5, 6, 7) GROUP BY faces.id, users.id`
    }

    if sidx != "" {
        query += ` ORDER BY faces.`+sidx
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

func (this *FaceModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
    query := `SELECT array_to_string(
        array(SELECT users.id || ':' || users.login FROM users GROUP BY users.id ORDER BY users.id), ';') as name;`
    logins := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    query = `SELECT array_to_string(
        array(SELECT faces.id || ':' || faces.id || '-' || array_to_string(array_agg(param_values.value), ' ')
        FROM reg_param_vals
        INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
        INNER JOIN params ON params.id = param_values.param_id
        WHERE params.id in (5, 6, 7) GROUP BY faces.id ORDER BY faces.id), ';') as name;`

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

func (this *FaceModel) WhereByParams(filters map[string]interface{}, num int) (where string, params []interface{}, num1 int) {
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
            FROM reg_param_vals
            INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN events ON events.id = registrations.event_id
            INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
            INNER JOIN params ON params.id = param_values.param_id
            INNER JOIN users ON users.id = faces.user_id WHERE `
        where += where1 + ")"
        params = append(params, params1...)
    }

    return where, params, i
}
