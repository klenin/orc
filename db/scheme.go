package db

import (
	"fmt"
)

func InitScheme() {
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
	Query(Teams_users, nil)
	Query(Persons_events, nil)
	Query(Forms, nil)
	Query(Params, nil)
	Query(Forms_types, nil)
	Query(Param_values, nil)
}

func Boom() {
	Query(Select1, nil)
	Query(Select2, nil)

	Query(Insert0, nil)
	Query(Insert1, nil)
	Query(Insert2, nil)
	Query(Insert3, nil)
	Query(Insert4, nil)
	Query(Insert5, nil)
	Query(Insert6, nil)
	Query(Insert7, nil)
	Query(Insert8, nil)
	Query(Insert9, nil)
	Query(Insert10, nil)
	Query(Insert11, nil)
	Query(Insert13, nil)
	Query(Insert14, nil)
	Query(Insert15, nil)
	Query(Insert16, nil)
	Query(Insert17, nil)
	Query(Insert18, nil)
	Query(Insert19, nil)
	Query(Insert20, nil)
	Query(Insert21, nil)
	Query(Insert22, nil)
	Query(Insert23, nil)
	Query(Insert24, nil)
	Query(Insert25, nil)

	Query(Insert26, nil)
	Query(Insert27, nil)
	Query(Insert28, nil)
	Query(Insert29, nil)
	Query(Insert30, nil)
	Query(Insert31, nil)
	Query(Insert32, nil)
	Query(Insert33, nil)
	Query(Insert34, nil)
	Query(Insert35, nil)
	Query(Insert36, nil)

	Query(Insert37, nil)
	Query(Insert38, nil)
	Query(Insert39, nil)
	Query(Insert40, nil)

	Query(Insert41, nil)
	Query(Insert42, nil)
	Query(Insert43, nil)
	Query(Insert44, nil)
}

func DropScheme() {
	for _, v := range Tables {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", v)
		Query(query, nil)
	}

	for _, v := range Tables {
		query := fmt.Sprintf("DROP SEQUENCE IF EXISTS %s_id_seq;", v)
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
