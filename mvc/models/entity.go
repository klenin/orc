package models

import (
    "github.com/orc/utils"
    "reflect"
)

type ModelManager struct{}

type Entity struct {
    TableName string
    Caption   string

    Fields interface{}

    Columns  []string
    ColNames []string

    Sub      bool
    SubTable []string
    SubField string
}

func (this Entity) GetTableName() string {
    return this.TableName
}

func (this Entity) GetCaption() string {
    return this.Caption
}

func (this Entity) GetSub() bool {
    return this.Sub
}

func (this Entity) GetSubTable(index int) string {
    return this.SubTable[index]
}

func (this Entity) GetSubField() string {
    return this.SubField
}

func (this Entity) GetColumns() []string {
    return this.Columns
}

func (this Entity) GetColumnByIdx(index int) string {
    return this.Columns[index]
}

func (this Entity) GetColumnSlice(index int) []string {
    return this.Columns[index:]
}

func (this Entity) GetColNames() []string {
    return this.ColNames
}

func (this Entity) GetFields() interface{} {
    return this.Fields
}

func (this Entity) LoadModelData(data map[string]interface{}) {
    rv := reflect.ValueOf(this.Fields)
    rt := rv.Type()
    n := rt.Elem().NumField()

    for key, val := range data {
        for i := 0; i < n; i++ {
            tag := rt.Elem().Field(i).Tag.Get("name")
            if tag == key {
                rv.Elem().Field(i).Set(
                    reflect.ValueOf(
                        utils.ConvertTypeForModel(rt.Elem().Field(i).Tag.Get("type"), val)))
            }
        }
    }
}

type VirtEntity interface {
    LoadModelData(data map[string]interface{})

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
