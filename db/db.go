package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/orc/utils"
	"os"
	//"reflect"
	"strconv"
	"strings"
	"time"
)

func HandleErr(message string, err error) {
	if err != nil {
		fmt.Printf(message+"%v\n", err)
		os.Exit(1)
	}
}

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
	"Типы_мероприятий",
	"Мероприятия-Типы",
	"Команды",
	"Персоны",
	"Персоны-Мероприятия",
	"Пользователи",
	"Команды-Персоны",
	"Формы",
	"Параметры",
	"Формы-Типы_мероприятий",
	"Значения_параметров",
}

func Exec(query string, params []interface{}) sql.Result {
	fmt.Println(query)

	stmt, err := DB.Prepare(query)
	utils.HandleErr("[db.Exec] Prepare: ", err)

	result, err := stmt.Exec(params...)
	utils.HandleErr("[db.Exec] Exec: ", err)
	return result
}

func Query(query string, params []interface{}) *sql.Rows {
	fmt.Println(query)

	stmt, err := DB.Prepare(query)
	utils.HandleErr("[db.Query] Prepare: ", err)

	result, err := stmt.Query(params...)
	utils.HandleErr("[db.Query] Query: ", err)
	return result
}

func QueryRow(query string, params []interface{}) *sql.Row {
	fmt.Println(query)

	stmt, err := DB.Prepare(query)
	utils.HandleErr("[db.QueryRow] Prepare: ", err)

	result := stmt.QueryRow(params...)
	utils.HandleErr("[db.QueryRow] Query: ", err)
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

func Select(tableName string, where []string, fields []string) []interface{} {
	//j := 0
	var key string
	var val []interface{}
	if len(where) != 0 {
		key = where[0] + "=$1"
		val = []interface{}{where[1]}
	}
	//fmt.Println(tableName)
	//fmt.Println(fields)
	query := QuerySelect(tableName, key, fields)
	rows := Query(query, val)
	rowsInf := Exec(query, val)

	columns, _ := rows.Columns()
	row := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i, _ := range row {
		row[i] = &values[i]
	}

	l, err := rowsInf.RowsAffected()
	utils.HandleErr("[Entity.Select] RowsAffected: ", err)
	return ConvertData(columns, l, rows)
}

func ConvertData(columns []string, l int64, rows *sql.Rows) []interface{} {
	j := 0
	row := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i, _ := range row {
		row[i] = &values[i]
	}

	answer := make([]interface{}, l)

	for rows.Next() {
		rows.Scan(row...)
		answer[j] = make(map[string]interface{}, len(values))
		record := make(map[string]interface{}, len(values))
		for i, col := range values {
			if col != nil {
				//fmt.Printf("\n%s: type= %s\n", columns[i], reflect.TypeOf(col))
				switch col.(type) {
				default:
					utils.HandleErr("Entity.Select: Unexpected type.", nil)
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
