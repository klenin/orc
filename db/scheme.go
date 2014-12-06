package db

import (
	"fmt"
)

func InitSchema() {
	for _, v := range Tables {
		query := fmt.Sprintf("CREATE SEQUENCE %s_id_seq;", v)
		Query(query, nil)
	}
	Query(Events, nil)
	Query(Event_types, nil)
	Query(Events_types, nil)
	Query(Teams, nil)
	Query(Persons, nil)
	Query(Users, nil)
	Query(Teams_persons, nil)
	Query(Persons_events, nil)
	Query(Forms, nil)
	Query(Params, nil)
	Query(Forms_types, nil)
	Query(Param_values, nil)
}

func DropSchema() {
	for _, v := range Tables {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", v)
		Query(query, nil)
		query = fmt.Sprintf("DROP SEQUENCE IF EXISTS %s_id_seq;", v)
		Query(query, nil)
	}
}

func GetCurrId(tableName string) string {
	var id string
	query := fmt.Sprintf("SELECT currval('%s');", tableName+"_id_seq")
	QueryRow(query, nil).Scan(&id)
	return id
}

func GetNextId(tableName string) string {
	var id string
	query := fmt.Sprintf("SELECT nextval('%s');", tableName+"_id_seq")
	QueryRow(query, nil).Scan(&id)
	return id
}
