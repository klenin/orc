package models

import (
    "github.com/orc/utils"
    "github.com/orc/db"
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
    TableName string
    Caption   string

    Fields interface{}

    Columns  []string
    ColNames []string

    Sub      bool
    SubTable []string
    SubField string

    WherePart map[string]interface{}
    Condition ConditionEnumElem
    OrderBy   string
    Limit     interface{}
    Offset    int
}

func (this *Entity) GetTableName() string {
    return this.TableName
}

func (this *Entity) GetCaption() string {
    return this.Caption
}

func (this *Entity) GetSub() bool {
    return this.Sub
}

func (this *Entity) GetSubTable(index int) string {
    return this.SubTable[index]
}

func (this *Entity) GetSubField() string {
    return this.SubField
}

func (this *Entity) GetColumns() []string {
    return this.Columns
}

func (this *Entity) GetColumnByIdx(index int) string {
    return this.Columns[index]
}

func (this *Entity) GetColumnSlice(index int) []string {
    return this.Columns[index:]
}

func (this *Entity) GetColNames() []string {
    return this.ColNames
}

func (this *Entity) GetFields() interface{} {
    return this.Fields
}

func (this *Entity) GetConditionName() string {
    switch this.Condition {
    case OR:
        return "OR"
    case AND:
        return "AND"
    }
    panic("Entity.GetConditionName: Invalid condition")
}

func (this *Entity) LoadModelData(data map[string]interface{}) {
    rv := reflect.ValueOf(this.Fields)
    rt := rv.Type()

    for key, val := range data {
        for i := 0; i < rt.Elem().NumField(); i++ {
            tag := rt.Elem().Field(i).Tag.Get("name")
            if tag == key {
                println(tag)
                value := utils.ConvertTypeForModel(rt.Elem().Field(i).Tag.Get("type"), val)
                if value == nil {
                    continue
                }
                rv.Elem().Field(i).Set(reflect.ValueOf(value))
            }
        }
    }
}

func (this Entity) LoadWherePart(data map[string]interface{}) {
    rv := reflect.ValueOf(this.WherePart)
    rt := reflect.ValueOf(this.Fields).Type()

    for key, val := range data {
        for i := 0; i < rt.Elem().NumField(); i++ {

            if rt.Elem().Field(i).Tag.Get("name") != key {
                continue
            }

            if val != nil && reflect.TypeOf(val).Name() == "" { // hope that v is array of interfaces
                rv.Interface().(map[string]interface{})[key] = make([]interface{}, 0)
                arr := make([]interface{}, 0)
                for _, vv := range val.([]interface{}) {
                    v_ := utils.ConvertTypeForModel(rt.Elem().Field(i).Tag.Get("type"), vv)
                    if v_ == nil {
                        continue
                    }
                    arr = append(arr, v_)
                }
                rv.Interface().(map[string]interface{})[key] = arr
                continue
            }

            rv.Interface().(map[string]interface{})[key] = utils.ConvertTypeForModel(rt.Elem().Field(i).Tag.Get("type"), val)
        }
    }
}

func (this Entity) GenerateWherePart(counter int) (string, []interface{}) {
    var key []string
    var val []interface{}

    for k, v := range this.WherePart {
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

func (this *Entity) SetOrder(orderBy string) {
    rt := reflect.ValueOf(this.Fields).Type()
    for i := 0; i < rt.Elem().NumField(); i++ {
        if rt.Elem().Field(i).Tag.Get("name") == orderBy {
            reflect.ValueOf(this).Elem().FieldByName("OrderBy").Set(reflect.ValueOf(orderBy))
            break
        }
    }
}

func (this *Entity) SetCondition(c ConditionEnumElem) {
    reflect.ValueOf(this).Elem().FieldByName("Condition").Set(reflect.ValueOf(c))
}

func (this *Entity) SetLimit(limit interface{}) {
    switch limit.(type) {
    case string:
        if limit.(string) != "ALL" {
            panic("[Entity::SetLimit] Invalid value")
        }
        reflect.ValueOf(this).Elem().FieldByName("Limit").Set(reflect.ValueOf(limit))
        break
    case int:
        if limit.(int) < 0 {
            panic("[Entity::SetLimit] Invalid value")
        }
        reflect.ValueOf(this).Elem().FieldByName("Limit").Set(reflect.ValueOf(limit))
        break
    default:
        panic("[Entity::SetLimit] Invalid type")
    }
}

func (this *Entity) SetOffset(offset int) {
    if offset < 0 {
        panic("[Entity::SerOffset] Invalid value")
    }
    reflect.ValueOf(this).Elem().FieldByName("Offset").SetInt(int64(offset))
}

func (this *Entity) GetModelRefDate() (fields []string, result map[string]interface{}) {
    result = make(map[string]interface{})
    rt := reflect.ValueOf(this.Fields).Type()

    for i := 0; i < rt.Elem().NumField(); i++ {
        refFieldShow := rt.Elem().Field(i).Tag.Get("refFieldShow")
        if refFieldShow != "" {
            fields = append(fields, refFieldShow)
            refField := rt.Elem().Field(i).Tag.Get("refField")
            m := new(ModelManager).GetModel(rt.Elem().Field(i).Tag.Get("refTable"))
            data := db.Select(m, []string{refField, refFieldShow})
            result[rt.Elem().Field(i).Tag.Get("name")] = make([]interface{}, len(data))
            result[rt.Elem().Field(i).Tag.Get("name")] = data
        }
    }

    return fields, result
}

func (this *Entity) Select(fields []string, filters map[string]interface{}, limit, offset int, sord, sidx string) (result []interface{}) {
    if len(fields) == 0 {
        return nil
    }

    where, params := this.Where(filters)

    query := `SELECT `+strings.Join(fields, ", ")+` FROM `+this.GetTableName()+where

    if sidx != "" {
        query += ` ORDER BY `+this.GetTableName()+"."+sidx
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

func (this *Entity) Where(filters map[string]interface{}) (where string, params []interface{}) {
    where = ""
    if filters == nil {
        return where, nil
    }
    i := 1

    groupOp := filters["groupOp"].(string)
    rules := filters["rules"].([]interface{})

    if len(rules) > 10 {
        log.Println("More 10 rules for serching!")
    }

    firstElem := true
    where = " WHERE "

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
    return where, params
}

func (this *Entity) GetColModel() []map[string]interface{} {
    return nil
}

func (this *Entity) GetColModelForUser(user_id int) []map[string]interface{} {
    return nil
}

type VirtEntity interface {
    LoadModelData(data map[string]interface{})
    LoadWherePart(data map[string]interface{})
    GenerateWherePart(counter int) (string, []interface{})

    GetConditionName() string
    SetCondition(c ConditionEnumElem)

    SetOrder(orderBy string)
    SetLimit(limit interface{})
    SetOffset(offset int)

    GetTableName() string
    GetCaption() string

    GetFields() interface{}

    GetSub() bool
    GetSubTable(index int) string
    GetSubField() string

    GetColumns() []string
    GetColNames() []string
    GetColumnByIdx(index int) string
    GetColumnSlice(index int) []string

    GetModelRefDate() (fields []string, result map[string]interface{})
    Where(filters map[string]interface{}) (where string, params []interface{})
    Select(fields []string, filters map[string]interface{}, limit, offset int, sord, sidx string) (result []interface{})

    GetColModel() ([]map[string]interface{})
    GetColModelForUser(user_id int) ([]map[string]interface{})
}

func (this *ModelManager) GetModel(tableName string) VirtEntity {
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
    case "reg_param_vals":
        return this.RegParamVals()
    case "groups":
        return this.Groups()
    case "group_registrations":
        return this.GroupRegistrations()
    }
    return nil
}
