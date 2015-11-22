package db

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "github.com/klenin/orc/utils"
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

var QueryGetCaptFaceIdAndRegId = `SELECT faces.id, registrations.id
    FROM regs_groupregs
    INNER JOIN group_registrations
        ON group_registrations.id = regs_groupregs.groupreg_id
    INNER JOIN registrations ON registrations.id = regs_groupregs.reg_id
    INNER JOIN groups ON groups.id = group_registrations.group_id
    INNER JOIN faces ON faces.id = groups.face_id
        AND faces.id = registrations.face_id
    INNER JOIN events ON events.id = registrations.event_id
    INNER JOIN events_forms ON events_forms.event_id = events.id
    INNER JOIN forms ON forms.id = events_forms.form_id
    INNER JOIN params ON forms.id = params.form_id
    INNER JOIN param_types ON param_types.id = params.param_type_id
    INNER JOIN param_values ON params.id = param_values.param_id
        AND param_values.reg_id = registrations.id
    WHERE group_registrations.id = $1 AND forms.personal = FALSE;`

var QueryGetCaptRegIdByGroupRegIdAndFaceId = `SELECT registrations.id
    FROM regs_groupregs
    INNER JOIN group_registrations
        ON group_registrations.id = regs_groupregs.groupreg_id
    INNER JOIN registrations ON registrations.id = regs_groupregs.reg_id
    INNER JOIN faces ON faces.id = registrations.face_id
    INNER JOIN groups ON groups.face_id = faces.id
        AND groups.id = group_registrations.group_id
    INNER JOIN events ON events.id = registrations.event_id
    INNER JOIN events_forms ON events_forms.event_id = events.id
    INNER JOIN forms ON forms.id = events_forms.form_id
    INNER JOIN params ON forms.id = params.form_id
    INNER JOIN param_types ON param_types.id = params.param_type_id
    INNER JOIN param_values ON params.id = param_values.param_id
        AND param_values.reg_id = registrations.id
    WHERE group_registrations.id = $1 AND faces.id = $2
        AND forms.personal = FALSE;`

var QueryGetRegIdByGroupRegIdAndFaceId = `SELECT registrations.id
    FROM regs_groupregs
    INNER JOIN group_registrations
        ON group_registrations.id = regs_groupregs.groupreg_id
    INNER JOIN registrations ON registrations.id = regs_groupregs.reg_id
    INNER JOIN faces ON faces.id = registrations.face_id
    INNER JOIN persons ON persons.face_id = faces.id
        AND persons.group_id = group_registrations.group_id
    INNER JOIN events ON events.id = registrations.event_id
    INNER JOIN events_forms ON events_forms.event_id = events.id
    INNER JOIN forms ON forms.id = events_forms.form_id
    INNER JOIN params ON forms.id = params.form_id
    INNER JOIN param_types ON param_types.id = params.param_type_id
    INNER JOIN param_values ON params.id = param_values.param_id
        AND param_values.reg_id = registrations.id
    WHERE group_registrations.id = $1 AND faces.id = $2
        AND forms.personal = TRUE;`
