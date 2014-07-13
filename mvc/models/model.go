package models

import (
	"github.com/orc/db"
)

type ModelManager struct{}

type Field struct {
	Name    string
	Caption string
	Type    string
	Ref     bool
}

type Entity struct {
	TableName string
	Caption   string
	Fields    map[string]*Field
	Columns   []string
	ColNames  []string
	Ref       bool
	RefData   map[string]interface{}
	RefFields []string
	Sub       bool
	SubTable  []string
	SubField  string
}

func (this Entity) Select(where []string, condition string, fields []string) ([]interface{}, map[string]interface{}) {
	result1 := db.Select(this.TableName, where, condition, fields)
	if this.Ref {
		result2 := this.RefData
		return result1, result2
	}
	return result1, nil
}

func (this Entity) Insert(fields []string, params []interface{}) {
	query := db.QueryInsert(this.TableName, fields)
	db.Query(query, params)
}

func (this Entity) Update(fields []string, params []interface{}, where string) {
	query := db.QueryUpdate(this.TableName, where, fields)
	db.Query(query, params)
}

func (this Entity) Delete(field string, params []interface{}) {
	query := db.QueryDelete(this.TableName, field, len(params))
	db.Query(query, params)
}
