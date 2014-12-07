package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/orc/utils"
	"strconv"
	"strings"
	"time"
)

const user string = "admin"
const dbname string = "orc"
const password string = "admin"

var DB, _ = sql.Open(
	"postgres",
	"host=localhost"+
		" user="+user+
		" dbname="+dbname+
		" password="+password+
		" sslmode=disable")

var Tables = []string{
	"events",
	"event_types",
	"events_types",
	"teams",
	"persons",
	"persons_events",
	"users",
	"teams_persons",
	"forms",
	"params",
	"forms_types",
	"param_values",
}

var TableNames = []string{
	"Мероприятия",
	"Типы мероприятий",
	"Мероприятия-Типы",
	"Команды",
	"Персоны",
	"Персоны-Мероприятия",
	"Пользователи",
	"Команды-Персоны",
	"Формы",
	"Параметры",
	"Формы-Типы мероприятий",
	"Значения параметров",
}

func Exec(query string, params []interface{}) sql.Result {
	stmt, err := DB.Prepare(query)
	utils.HandleErr("[db.Exec] Prepare: ", err, nil)

	result, err := stmt.Exec(params...)
	utils.HandleErr("[db.Exec] Exec: ", err, nil)
	return result
}

func Query(query string, params []interface{}) *sql.Rows {
	println(query)
	stmt, err := DB.Prepare(query)
	utils.HandleErr("[db.Query] Prepare: ", err, nil)

	result, err := stmt.Query(params...)
	utils.HandleErr("[db.Query] Query: ", err, nil)
	return result
}

func QueryRow(query string, params []interface{}) *sql.Row {
	println(query)
	stmt, err := DB.Prepare(query)
	utils.HandleErr("[db.QueryRow] Prepare: ", err, nil)

	result := stmt.QueryRow(params...)
	utils.HandleErr("[db.QueryRow] Query: ", err, nil)
	return result
}

func QuerySelect(tableName, where string, fields []string) string {
	query := "SELECT %s FROM %s"
	f := strings.Join(fields, ", ")
	if where != "" {
		query += " WHERE %s;"
		return fmt.Sprintf(query, f, tableName, where)
	} else {
		return fmt.Sprintf(query, f, tableName)
	}
}

func QueryInsert(tableName string, fields []string) string {
	query := "INSERT INTO %s (%s) VALUES (%s);"
	f := strings.Join(fields, ", ")
	p := strings.Join(MakeParams(len(fields)), ", ")
	return fmt.Sprintf(query, tableName, f, p)
}

func QueryUpdate(tableName, where string, fields []string) string {
	query := "UPDATE %s SET %s WHERE %s;"
	p := strings.Join(MakePairs(fields), ", ")
	return fmt.Sprintf(query, tableName, p, where)
}

func QueryDelete(tableName, fieldName string, countParams int) string {
	query := "DELETE FROM %s WHERE %s IN (%s)"
	params := strings.Join(MakeParams(countParams), ", ")
	return fmt.Sprintf(query, tableName, fieldName, params)
}

func IsExists(tableName, fieldName string, value string) bool {
	var result string
	query := QuerySelect(tableName, fieldName+"=$1", []string{fieldName})
	row := QueryRow(query, []interface{}{value})
	err := row.Scan(&result)
	return err != sql.ErrNoRows
}

func MakeParams(n int) []string {
	var result = make([]string, n)
	for i := 0; i < n; i++ {
		result[i] = "$" + strconv.Itoa(i+1)
	}
	return result
}

func MakePairs(fields []string) []string {
	var result = make([]string, len(fields))
	for i := 0; i < len(fields); i++ {
		result[i] = fields[i] + "=$" + strconv.Itoa(i+1)
	}
	return result
}

/**
 * condition: the AND condition and the OR condition
 * where: [fieldName1, paramVal1, fieldName2, paramVal2, ...]
 */
func Select(tableName string, where []string, condition string, fields []string) []interface{} {
	var key []string
	var val []interface{}
	var paramName = 1
	if len(where) != 0 {
		for i := 0; i < len(where)-1; i += 2 {
			key = append(key, where[i]+"=$"+strconv.Itoa(paramName))
			val = append(val, where[i+1])
			paramName++
		}
	}
	query := QuerySelect(tableName, strings.Join(key, " "+condition+" "), fields)
	rows := Query(query, val)
	rowsInf := Exec(query, val)
	columns, _ := rows.Columns()
	size, err := rowsInf.RowsAffected()
	utils.HandleErr("[Entity.Select] RowsAffected: ", err, nil)
	return ConvertData(columns, size, rows)
}

func ConvertData(columns []string, size int64, rows *sql.Rows) []interface{} {
	row := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	answer := make([]interface{}, size)

	for i, _ := range row {
		row[i] = &values[i]
	}

	j := 0
	for rows.Next() {
		rows.Scan(row...)
		record := make(map[string]interface{}, len(values))
		for i, col := range values {
			if col != nil {
				//fmt.Printf("\n%s: type= %s\n", columns[i], reflect.TypeOf(col))
				switch col.(type) {
				case bool:
					record[columns[i]] = col.(bool)
				case int:
					record[columns[i]] = col.(int)
				case int64:
					record[columns[i]] = col.(int64)
				case float64:
					record[columns[i]] = col.(float64)
				case string:
					record[columns[i]] = col.(string)
				case []byte:
					record[columns[i]] = string(col.([]byte))
				case []int8:
					record[columns[i]] = col.([]string)
				case time.Time:
					record[columns[i]] = col
				default:
					utils.HandleErr("Entity.Select: Unexpected type.", nil, nil)
				}
			}
			answer[j] = record
		}
		j++
	}
	return answer
}

func InnerJoin(
	selectFields []string,
	selectRef string,

	fromTable string,
	fromTableRef string,
	fromField []string,

	joinTables []string,
	joinRef []string,
	joinField []string,

	where string) string {

	query := "SELECT "
	for i := 0; i < len(selectFields); i++ {
		query += selectRef + "." + selectFields[i] + ", "
	}
	query = query[0 : len(query)-2]
	query += " FROM " + fromTable + " " + fromTableRef
	for i := 0; i < len(joinTables); i++ {
		query += " INNER JOIN " + joinTables[i] + " " + joinRef[i]
		query += " ON " + joinRef[i] + "." + joinField[i] + " = " + fromTableRef + "." + fromField[i]
	}
	query += " " + where
	return query
}
