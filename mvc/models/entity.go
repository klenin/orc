package models

import (
    "database/sql"
    "fmt"
    "github.com/klenin/orc/utils"
    "github.com/klenin/orc/db"
    "reflect"
    "strconv"
    "strings"
    "log"
)

type ModelManager struct{}

type ConditionEnumElem int

const (
    OR ConditionEnumElem = iota
    AND ConditionEnumElem = iota
)

type Entity struct {
    tableName  string
    caption    string
    fields     interface{}
    columns    []string
    colNames   []string
    sub        bool
    subTables  []string
    subField   string
    wherePart  map[string]interface{}
    condition  ConditionEnumElem
    orderBy    string
    limit      interface{}
    offset     int
    sorting    string
}

func (this *Entity) SetTableName(name string) *Entity {
    this.tableName = name

    return this
}

func (this *Entity) GetTableName() string {
    return this.tableName
}

func (this *Entity) SetCaption(caption string) *Entity {
    this.caption = caption

    return this
}

func (this *Entity) GetCaption() string {
    return this.caption
}

func (this *Entity) SetSub(sub bool) *Entity {
    this.sub = sub

    return this
}

func (this *Entity) GetSub() bool {
    return this.sub
}

func (this *Entity) SetSubTables(subTables []string) *Entity {
    this.subTables = subTables

    return this
}

func (this *Entity) GetSubTable(index int) string {
    return this.subTables[index]
}

func (this *Entity) SetSubField(fieldName string) *Entity {
    this.subField = fieldName

    return this
}

func (this *Entity) GetSubField() string {
    return this.subField
}

func (this *Entity) SetColumns(columns []string) *Entity {
    this.columns = columns

    return this
}

func (this *Entity) GetColumns() []string {
    return this.columns
}

func (this *Entity) GetColumnByIdx(index int) string {
    return this.columns[index]
}

func (this *Entity) SetColNames(colNames []string) *Entity {
    this.colNames = colNames

    return this
}

func (this *Entity) GetColNames() []string {
    return this.colNames
}

func (this *Entity) SetFields(fields interface{}) *Entity {
    this.fields = fields

    return this
}

func (this *Entity) GetFields() interface{} {
    return this.fields
}

func (this *Entity) SetCondition(condition ConditionEnumElem) *Entity {
    this.condition = condition

    return this
}

func (this *Entity) GetConditionName() string {
    switch this.condition {
    case OR:
        return "OR"
    case AND:
        return "AND"
    }
    panic("Entity.GetConditionName: Invalid condition")
}

func (this *Entity) SetOrder(orderBy string) *Entity {
    this.orderBy = orderBy

    return this
}

func (this *Entity) GetOrder() string {
    return this.orderBy
}

func (this *Entity) SetLimit(limit interface{}) *Entity {
    switch limit.(type) {
    case string:
        if limit.(string) != "ALL" {
            panic("[Entity::SetLimit] Invalid value")
        }
        this.limit = limit

        return this
    case int:
        if limit.(int) < 0 {
            panic("[Entity::SetLimit] Invalid value")
        }
        this.limit = limit

        return this

    default:
        panic("[Entity::SetLimit] Invalid type")
    }
}

func (this *Entity) GetLimit() interface{} {
    return this.limit
}

func (this *Entity) SetOffset(offset int) *Entity {
    if offset < 0 {
        panic("[Entity::SetOffset] Invalid value")
    }
    this.offset = offset

    return this
}

func (this *Entity) GetOffset() int {
    return this.offset
}

func (this *Entity) SetSorting(sorting string) *Entity {
    if sorting == "ASC" || sorting == "DESC" || sorting == "asc" || sorting == "desc" {
        this.sorting = sorting
    } else {
        panic("[Entity::SetSorting] Invalid value")
    }

    return this
}

func (this *Entity) GetSorting() string {
    return this.sorting
}

func (this *Entity) SetWherePart(where map[string]interface{}) *Entity {
    this.wherePart = where

    return this
}

func (this *Entity) LoadModelData(data map[string]interface{}) *Entity {
    refOfValue := reflect.ValueOf(this.fields); refOfType := refOfValue.Type()
    for key, val := range data {
        for i := 0; i < refOfType.Elem().NumField(); i++ {
            tag := refOfType.Elem().Field(i).Tag.Get("name")
            if tag == key {
                method := "Set"+strings.ToUpper(string(refOfType.Elem().Field(i).Name[0]))+string(refOfType.Elem().Field(i).Name[1:])
                refMethod := refOfValue.MethodByName(method)
                log.Println("method: ", method)
                if !refMethod.IsValid() {
                    log.Println("Method is not exists!")
                    continue
                }
                value := utils.CheckTypeValue(refOfType.Elem().Field(i).Tag.Get("type"), val)
                if value == nil {
                    continue
                }
                refOfValue.MethodByName(method).Call([]reflect.Value{reflect.ValueOf(value)})
            }
        }
    }

    return this
}

func (this *Entity) LoadWherePart(data map[string]interface{}) *Entity {
    if data == nil || len(data) == 0 {
        return this
    }

    rv := reflect.ValueOf(this.wherePart)
    rt := reflect.ValueOf(this.fields).Type()

    for key, val := range data {
        for i := 0; i < rt.Elem().NumField(); i++ {

            if rt.Elem().Field(i).Tag.Get("name") != key {
                continue
            }

            if val != nil && reflect.TypeOf(val).Name() == "" { // hope that v is array of interfaces
                rv.Interface().(map[string]interface{})[key] = make([]interface{}, 0)
                arr := make([]interface{}, 0)
                for _, vv := range val.([]interface{}) {
                    v_ := utils.CheckTypeValue(rt.Elem().Field(i).Tag.Get("type"), vv)
                    if v_ == nil {
                        continue
                    }
                    arr = append(arr, v_)
                }
                rv.Interface().(map[string]interface{})[key] = arr
                continue
            }
            rv.Interface().(map[string]interface{})[key] = utils.CheckTypeValue(rt.Elem().Field(i).Tag.Get("type"), val)
        }
    }

    return this
}

func (this Entity) GenerateWherePart(counter int) (string, []interface{}) {
    var key []string
    var val []interface{}

    for k, v := range this.wherePart {
        if reflect.TypeOf(v).Name() == "" {
            for _, vv := range v.([]interface{}) {
                key = append(key, k+"=$"+strconv.Itoa(counter))
                val = append(val, vv)
                counter++
            }
            continue
        }
        key = append(key, k+"=$"+strconv.Itoa(counter))
        val = append(val, v)
        counter++
    }

    return strings.Join(key, " "+this.GetConditionName()+" "), val
}

func (this *Entity) Select_(fields []string) []interface{} {
    query := "SELECT %s FROM %s"
    if len(this.wherePart) != 0 {
        where, params := this.GenerateWherePart(1)
        params = append(params, this.orderBy)
        query += " WHERE %s ORDER BY $" + strconv.Itoa(len(params))

        switch this.limit.(type) {
        case string:
            query += " LIMIT ALL"
            break
        case int:
            query += " LIMIT $" + strconv.Itoa(len(params))
            params = append(params, this.limit)
            break
        default:
            panic("Invalid type of limit")
        }

        params = append(params, this.offset)
        query += " OFFSET $" + strconv.Itoa(len(params)) + ";"

        return db.Query(fmt.Sprintf(query, strings.Join(fields, ", "), this.tableName, where), params)
    } else {
        query += " ORDER BY $1 LIMIT $2 OFFSET $3;"

        return db.Query(
            fmt.Sprintf(query, strings.Join(fields, ", "), this.tableName),
            []interface{}{this.GetOrder(), this.GetLimit(), this.GetOffset()})
    }
}

func (this *Entity) Select(fields []string, filters map[string]interface{}) (result []interface{}) {
    if len(fields) == 0 {
        return nil
    }

    where, params, _ := this.Where(filters, 1)
    if where != "" {
        where = " WHERE " + where
    }
    query := `SELECT ` + strings.Join(fields, ", ") + ` FROM ` + this.GetTableName() + where
    query += ` ORDER BY ` + this.GetTableName() + "." + this.orderBy
    query += ` `+ this.sorting

    switch this.limit.(type) {
    case string:
        query += " LIMIT ALL"
        break
    case int:
        query += " LIMIT $" + strconv.Itoa(len(params))
        params = append(params, this.GetLimit())
        break
    default:
        panic("Invalid type of limit")
    }

    params = append(params, this.GetOffset())
    query += ` OFFSET $` + strconv.Itoa(len(params)) + `;`

    return db.Query(query, params)
}

func (this *Entity) Where(filters map[string]interface{}, num int) (where string, params []interface{}, num1 int) {
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
            where += this.GetTableName()+"."+rule["field"].(string) + "::text = $"+strconv.Itoa(i)
            params = append(params, rule["data"])
            i += 1
            break
        case "ne":// not equal
            where += this.GetTableName()+"."+rule["field"].(string) + "::text <> $"+strconv.Itoa(i)
            params = append(params, rule["data"])
            i += 1
            break
        case "bw":// begins with
            where += this.GetTableName()+"."+rule["field"].(string) + "::text LIKE $"+strconv.Itoa(i)+"||'%'"
            params = append(params, rule["data"])
            i += 1
            break
        case "bn":// does not begin with
            where += this.GetTableName()+"."+rule["field"].(string) + "::text NOT LIKE $"+strconv.Itoa(i)+"||'%'"
            params = append(params, rule["data"])
            i += 1
            break
        case "ew":// ends with
            where += this.GetTableName()+"."+rule["field"].(string) + "::text LIKE '%'||$"+strconv.Itoa(i)
            params = append(params, rule["data"])
            i += 1
            break
        case "en":// does not end with
            where += this.GetTableName()+"."+rule["field"].(string) + "::text NOT LIKE '%'||$"+strconv.Itoa(i)
            params = append(params, rule["data"])
            i += 1
            break
        case "cn":// contains
            where += this.GetTableName()+"."+rule["field"].(string) + "::text LIKE '%'||$"+strconv.Itoa(i)+"||'%'"
            params = append(params, rule["data"])
            i += 1
            break
        case "nc":// does not contain
            where += this.GetTableName()+"."+rule["field"].(string) + "::text NOT LIKE '%'||$"+strconv.Itoa(i)+"||'%'"
            params = append(params, rule["data"])
            i += 1
            break
        case "nu":// is null
            where += this.GetTableName()+"."+rule["field"].(string) + "::text IS NULL"
            break
        case "nn":// is not null
            where += this.GetTableName()+"."+rule["field"].(string) + "::text IS NOT NULL"
            break
        case "in":// is in
            where += this.GetTableName()+"."+rule["field"].(string) + "::text IN ("
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
            where += this.GetTableName()+"."+rule["field"].(string) + "::text NOT IN ("
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
        where1, params1, num1 :=  this.Where(filters1, i)
        i = num1
        if where != "" {
            if !firstElem {
                if groupOp != "AND" && groupOp != "OR" {
                    log.Println("`groupOp` parameter is not allowed!")
                    continue
                }
                where += " " + groupOp + " "
            } else {
                firstElem = false
            }
            where += "(" + where1 + ")"
            params = append(params, params1...)
        }
    }

    return where, params, i
}

func (this *Entity) Delete(id int) {
    query := `DELETE FROM ` + this.GetTableName() + ` WHERE id = $1;`
    db.Exec(query, []interface{}{id})
}

func (this *Entity) QueryUpdate() *sql.Row {
    j := 1
    query := "UPDATE %s SET "
    tFields := reflect.ValueOf(this.fields).Type().Elem()
    vFields := reflect.ValueOf(this.fields).Elem()
    params := make([]interface{}, 0)

    for i := 1; i < tFields.NumField(); i++ {
        value, ok := utils.UpdateOrNot(tFields.Field(i).Tag.Get("type"), vFields.Field(i))
        if !ok {
            continue
        }
        query += tFields.Field(i).Tag.Get("name") + "=$" + strconv.Itoa(j) + ", "
        params = append(params, value)
        j++
    }
    query = query[0 : len(query)-2]

    if len(this.wherePart) != 0 {
        query += " WHERE %s;"
        v1, v2 := this.GenerateWherePart(j)

        return db.QueryRow(fmt.Sprintf(query, this.tableName, v1), append(params, v2...))
    } else {
        query += ";"

        return db.QueryRow(fmt.Sprintf(query, this.tableName), params)
    }
}

func (this *Entity) Update(isAdmin bool, userId int, params, where map[string]interface{}) {
    this.LoadModelData(params)
    this.LoadWherePart(where)
    this.QueryUpdate().Scan()
}

func (this *Entity) QueryInsert(extra string) *sql.Row {
    i := 1
    query := "INSERT INTO %s ("
    tFields := reflect.ValueOf(this.fields).Type().Elem()
    vFields := reflect.ValueOf(this.fields).Elem()
    params := make([]interface{}, 0)

    for i = 1; i < tFields.NumField(); i++ {
        value, ok := utils.UpdateOrNot(tFields.Field(i).Tag.Get("type"), vFields.Field(i))
        if !ok && tFields.Field(i).Tag.Get("null") == "NULL" {
            value = nil
        }
        query += tFields.Field(i).Tag.Get("name") + ", "
        params = append(params, value)
    }
    query = query[0 : len(query)-2]; query += ") VALUES (%s) %s;"

    return db.QueryRow(fmt.Sprintf(query, this.tableName, strings.Join(db.MakeParams(i-1), ", "), extra), params)
}

func (this *Entity) Add(userId int, params map[string]interface{}) error {
    this.LoadModelData(params)
    this.QueryInsert("").Scan()

    return nil
}

func (this *Entity) SelectRow(fields []string) *sql.Row {
    query := "SELECT %s FROM %s"

    if len(this.wherePart) != 0 {
        query += " WHERE %s;"
        where, params := this.GenerateWherePart(1)

        return db.QueryRow(fmt.Sprintf(query, strings.Join(fields, ", "), this.tableName, where), params)
    } else {
        query += ";"

        return db.QueryRow(fmt.Sprintf(query, strings.Join(fields, ", "), this.tableName), nil)
    }
}

func (this *Entity) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
    return nil
}

func (this *Entity) WhereByParams(filters map[string]interface{}, num int) (where string, params []interface{}, num1 int) {
    return "", nil, -1
}

type EntityInterface interface {
    GetTableName() string
    SetTableName(string) *Entity

    GetCaption() string
    SetCaption(string) *Entity

    GetSub() bool
    SetSub(bool) *Entity

    GetSubField() string
    SetSubField(string) *Entity

    GetColumns() []string
    SetColumns([]string) *Entity

    GetColNames() []string
    SetColNames([]string) *Entity

    GetFields() interface{}
    SetFields(interface{}) *Entity

    GetSubTable(int) string
    SetSubTables([]string) *Entity

    GetColumnByIdx(int) string

    GetConditionName() string

    SetCondition(ConditionEnumElem) *Entity

    SetOrder(string) *Entity
    GetOrder() string

    SetLimit(interface{}) *Entity
    GetLimit() interface{}

    SetOffset(int) *Entity
    GetOffset() int

    SetSorting(string) *Entity
    GetSorting() string

    SetWherePart(map[string]interface{}) *Entity

    LoadModelData(map[string]interface{}) *Entity
    LoadWherePart(map[string]interface{}) *Entity
    GenerateWherePart(int) (string, []interface{})

    Where(filters map[string]interface{}, num int) (where string, params []interface{}, num1 int)
    WhereByParams(filters map[string]interface{}, num int) (where string, params []interface{}, num1 int)

    Select(fields []string, filters map[string]interface{}) ([]interface{})
    Select_(fields []string) []interface{}
    SelectRow(fields []string) *sql.Row

    QueryUpdate() *sql.Row
    QueryInsert(string) *sql.Row

    Delete(id int)
    Add(userId int, params map[string]interface{}) error
    Update(isAdmin bool, userId int, params, where map[string]interface{})

    GetColModel(isAdmin bool, userId int) ([]map[string]interface{})
}

func (this *ModelManager) GetModel(tableName string) EntityInterface {
    switch tableName {
    case "events":
        return this.Events()
    case "event_types":
        return this.EventTypes()
    case "events_types":
        return this.EventsTypes()
    case "persons":
        return this.Persons()
    case "users":
        return this.Users()
    case "forms":
        return this.Forms()
    case "params":
        return this.Params()
    case "events_forms":
        return this.EventsForms()
    case "param_values":
        return this.ParamValues()
    case "param_types":
        return this.ParamTypes()
    case "registrations":
        return this.Registrations()
    case "faces":
        return this.Faces()
    case "groups":
        return this.Groups()
    case "group_registrations":
        return this.GroupRegistrations()
    case "regs_groupregs":
        return this.RegsGroupRegs()
    default:
        panic("Table is dont exists!")
    }

    return nil
}
