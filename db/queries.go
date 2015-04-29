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

func QueryCreateTable_(m interface{}) {
    model := reflect.ValueOf(m).Elem()
    tableName := model.FieldByName("TableName").String()

    QueryCreateSecuence(tableName)
    query := "CREATE TABLE IF NOT EXISTS %s ("
    mF := model.Elem().FieldByName("Fields").Elem().Type()
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
    query = query[0 : len(query)-2]
    query += ");"
    Exec(fmt.Sprintf(query, tableName), nil)
}

func QueryInsert_(m interface{}, extra string) *sql.Row {
    var i int

    query := "INSERT INTO %s ("
    tableName := reflect.ValueOf(m).Elem().FieldByName("TableName").String()

    tFields := reflect.ValueOf(m).Elem().FieldByName("Fields").Elem().Type().Elem()
    vFields := reflect.ValueOf(m).Elem().FieldByName("Fields").Elem().Elem()

    n := tFields.NumField()
    p := make([]interface{}, n-1)

    for i = 1; i < n; i++ {
        query += tFields.Field(i).Tag.Get("name") + ", "
        v, ok := utils.ConvertTypeModel(tFields.Field(i).Tag.Get("type"), vFields.Field(i))
        if !ok && tFields.Field(i).Tag.Get("null") == "NULL" {
            continue
        }
        p[i-1] = v
    }
    query = query[0 : len(query)-2]
    query += ") VALUES (%s) %s;"

    // if i < 2 {
    //     return
    // }

    return QueryRow(fmt.Sprintf(query, tableName, strings.Join(MakeParams(n-1), ", "), extra), p)
}

func QueryUpdate_(m interface{}) *sql.Row {
    model := reflect.ValueOf(m).Elem()
    tableName := model.FieldByName("TableName").String()
    i, j := 1, 1

    query := "UPDATE %s SET "

    tFields := model.FieldByName("Fields").Elem().Type().Elem()
    vFields := model.FieldByName("Fields").Elem().Elem()

    p := make([]interface{}, 0)

    for ; j < tFields.NumField(); j++ {
        v, ok := utils.ConvertTypeModel(tFields.Field(j).Tag.Get("type"), vFields.Field(j))
        if ok == false {
            continue
        }
        query += tFields.Field(j).Tag.Get("name") + "=$" + strconv.Itoa(i) + ", "
        p = append(p, v)
        i++
    }
    query = query[0 : len(query)-2]

    if i < 2 {
        return nil
    }

    if model.FieldByName("WherePart").Len() != 0 {
        query += " WHERE %s;"
        v := model.MethodByName("GenerateWherePart").Call([]reflect.Value{reflect.ValueOf(i)})
        return QueryRow(fmt.Sprintf(query, tableName, v[0]), append(p, v[1].Interface().([]interface{})...))
    } else {
        query += ";"
        return QueryRow(fmt.Sprintf(query, tableName), p)
    }
}

func QueryDeleteByIds(tableName, ids string) {
    query := "DELETE FROM %s WHERE id IN (%s)"
    Exec(fmt.Sprintf(query, tableName, ids), nil)
}

func IsExists_(tableName string, fields []string, params []interface{}) bool {
    query := "SELECT %s FROM %s WHERE %s;"
    f := strings.Join(fields, ", ")
    p := strings.Join(MakePairs(fields), " AND ")

    var result string
    row := QueryRow(fmt.Sprintf(query, f, tableName, p), params)
    err := row.Scan(&result)

    return err != sql.ErrNoRows && result != ""
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

func Select(m interface{}, fields []string) []interface{} {
    model := reflect.ValueOf(m).Elem()
    tableName := model.FieldByName("TableName").String()

    orderBy := " ORDER BY " + model.FieldByName("OrderBy").Interface().(string)

    var limit string
    switch model.FieldByName("Limit").Interface().(type) {
    case string:
        limit = " LIMIT " + model.FieldByName("Limit").Interface().(string)
        break
    case int:
        limit = " LIMIT " + strconv.Itoa(model.FieldByName("Limit").Interface().(int))
        break
    }

    offset := " OFFSET " + strconv.Itoa(model.FieldByName("Offset").Interface().(int))
    extra := orderBy + limit + offset

    query := "SELECT %s FROM %s"

    if model.FieldByName("WherePart").Len() != 0 {
        query += " WHERE %s" + extra + ";"
        v := model.MethodByName("GenerateWherePart").Call([]reflect.Value{reflect.ValueOf(1)})
        return Query(fmt.Sprintf(query, strings.Join(fields, ", "), tableName, v[0]), v[1].Interface().([]interface{}))
    } else {
        query += extra + ";"
        return Query(fmt.Sprintf(query, strings.Join(fields, ", "), tableName), nil)
    }
}

func SelectRow(m interface{}, fields []string) *sql.Row {
    model := reflect.ValueOf(m).Elem()
    tableName := model.FieldByName("TableName").String()

    query := "SELECT %s FROM %s"

    if model.FieldByName("WherePart").Len() != 0 {
        query += " WHERE %s;"
        v := model.MethodByName("GenerateWherePart").Call([]reflect.Value{reflect.ValueOf(1)})
        return QueryRow(fmt.Sprintf(query, strings.Join(fields, ", "), tableName, v[0]), v[1].Interface().([]interface{}))
    } else {
        query += ";"
        return QueryRow(fmt.Sprintf(query, strings.Join(fields, ", "), tableName), nil)
    }
}

func SelectCount(tableName string) int {
    return int(Query("SELECT COUNT(*) FROM "+tableName+";", nil)[0].(map[string]interface{})["count"].(int))
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
