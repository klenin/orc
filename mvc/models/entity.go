package models

import (
    "github.com/orc/utils"
    "reflect"
    "strconv"
    "strings"
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
}
