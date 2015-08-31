package db

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "github.com/orc/utils"
    "log"
    "reflect"
    "strconv"
    "strings"
    "time"
    "errors"
)

var DB *sql.DB = nil

func Exec(query string, params []interface{}) sql.Result {
    log.Println(query)

    stmt, err := DB.Prepare(query)
    utils.HandleErr("[queries.Exec] Prepare: ", err, nil)
    defer stmt.Close()
    result, err := stmt.Exec(params...)
    utils.HandleErr("[queries.Exec] Exec: ", err, nil)

    return result
}

func Query(query string, params []interface{}) []interface{} {
    log.Println(query)

    stmt, err := DB.Prepare(query)
    utils.HandleErr("[queries.Query] Prepare: ", err, nil)
    defer stmt.Close()
    rows, err := stmt.Query(params...)
    utils.HandleErr("[queries.Query] Query: ", err, nil)
    defer rows.Close()

    rowsInf := Exec(query, params)
    columns, err := rows.Columns()
    utils.HandleErr("[queries.Query] Columns: ", err, nil)
    size, err := rowsInf.RowsAffected()
    utils.HandleErr("[queries.Query] RowsAffected: ", err, nil)

    return ConvertData(columns, size, rows)
}

func QueryRow(query string, params []interface{}) *sql.Row {
    log.Println(query)

    stmt, err := DB.Prepare(query)
    utils.HandleErr("[queries.QueryRow] Prepare: ", err, nil)
    defer stmt.Close()
    result := stmt.QueryRow(params...)
    utils.HandleErr("[queries.QueryRow] Query: ", err, nil)

    return result
}

func QueryCreateSecuence(tableName string) {
    Exec("CREATE SEQUENCE "+tableName+"_id_seq;", nil)
}

func QueryCreateTable(m interface{}) {
    model := reflect.ValueOf(m)
    tableName := model.Elem().FieldByName("tableName").String()

    QueryCreateSecuence(tableName)
    query := "CREATE TABLE IF NOT EXISTS %s ("
    mF := model.Elem().FieldByName("fields").Elem().Type()
    for i := 0; i < mF.Elem().NumField(); i++ {
        query += mF.Elem().Field(i).Tag.Get("name") + " "
        query += mF.Elem().Field(i).Tag.Get("type") + " "
        query += mF.Elem().Field(i).Tag.Get("null") + " "
        switch mF.Elem().Field(i).Tag.Get("extra") {
        case "PRIMARY":
            query += "PRIMARY KEY DEFAULT NEXTVAL('"
            query += tableName + "_id_seq'), "
            break
        case "REFERENCES":
            query += "REFERENCES " + mF.Elem().Field(i).Tag.Get("refTable") + "(" + mF.Elem().Field(i).Tag.Get("refField") + ") ON DELETE CASCADE, "
            break
        case "UNIQUE":
            query += "UNIQUE, "
            break
        default:
            query += ", "
        }
    }
    query = query[0 : len(query)-2]; query += ");"
    Exec(fmt.Sprintf(query, tableName), nil)
}

func QueryDeleteByIds(tableName, ids string) {
    query := "DELETE FROM %s WHERE id IN (%s)"
    Exec(fmt.Sprintf(query, tableName, ids), nil)
}

func IsExists(tableName string, fields []string, params []interface{}) bool {
    query := "SELECT %s FROM %s WHERE %s;"
    f := strings.Join(fields, ", ")
    p := strings.Join(MakePairs(fields), " AND ")

    var result string
    row := QueryRow(fmt.Sprintf(query, f, tableName, p), params)
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
                    break
                case int:
                    record[columns[i]] = col.(int)
                    break
                case int64:
                    record[columns[i]] = int(col.(int64))
                    break
                case float64:
                    record[columns[i]] = col.(float64)
                    break
                case string:
                    record[columns[i]] = col.(string)
                    break
                // case []byte:
                //     record[columns[i]] = string(col.([]byte))
                //     break
                case []int8:
                    record[columns[i]] = col.([]string)
                    break
                case time.Time:
                    record[columns[i]] = col
                    break
                case []uint8:
                    data := strings.Split(strings.Trim(string(col.([]uint8)), "{}"), ",")
                    if len(data) == 1 {
                        record[columns[i]] = data[0]
                    } else {
                        record[columns[i]] = data
                    }
                    break
                default:
                    utils.HandleErr("ConvertData: ", errors.New("Unexpected type."), nil)
                }
            }
            answer[j] = record
        }
        j++
    }
    rows.Close()
    return answer
}
