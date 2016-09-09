package utils

import (
    "crypto/md5"
    "crypto/rand"
    "encoding/base64"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "github.com/lib/pq"
    "log"
    "net/http"
    // "os"
    "reflect"
    "regexp"
    "strconv"
)

func HandleErr(message string, err error, w http.ResponseWriter) bool {
    if err != nil {
        log.Println(message + err.Error())

        if err, ok := err.(*pq.Error); ok {
            log.Println("pq error:", err.Code.Name())

            if w == nil {
                return true
            }

            switch err.Code.Name() {
            case "unique_violation":
                http.Error(w, "Нарушение ограничения уникальности", http.StatusNotModified)

                return true
            case "datetime_field_overflow":
                http.Error(w, "Выход за границы временных значений", http.StatusNotModified)

                return true
            }
        }

        if w == nil {
            return true
        }

        fmt.Fprintf(w, "%s", message+err.Error())
        // http.Error(w, fmt.Sprintf(message+"%v\n", err.Error()), -1)
        // os.Exit(1)

        return true
    }

    return false
}

func ArrayInterfaceToString(array []interface{}) []string {
    s := reflect.ValueOf(array)
    result := make([]string, s.Len())
    for i, v := range array {
        switch v.(type) {
        case int:
            result[i] = strconv.Itoa(v.(int))
        case int64:
            result[i] = strconv.Itoa(int(v.(int64)))
        case float64:
            result[i] = strconv.Itoa(int(v.(float64)))
        case string:
            result[i] = v.(string)
        }
    }

    return result
}

func UpdateOrNot(type_ string, value reflect.Value) (interface{}, bool) {
    switch type_ {
    case "int":
        log.Println("val: ", strconv.Itoa(int(value.Int())))

        return value.Int(), value.Int() != 0
    case "text", "date", "time", "timestamp":
        log.Println("val: ", value.String())

        return value.String(), value.String() != ""
    case "boolean":
        log.Println("val: ", value.Bool())

        return value.Bool(), true
    default:
        panic("UpdateOrNot: unknown type")
    }
}

func CheckTypeValue(type_ string, value interface{}) interface{} {
    if value == nil {
        return nil
    }

    switch value.(type) {
    case string:
        // value from grid
        if value.(string) == "_empty" {
            return nil
        }

        switch type_ {
        case "int":
            v, err := strconv.Atoi(value.(string))
            if err != nil {
                return nil
            }
            log.Println("CheckTypeValue-int: ", strconv.Itoa(v))

            return v
        case "text", "date", "time", "timestamp":
            log.Println("CheckTypeValue-text-date-time-timestamp: ", value.(string))

            return value
        case "boolean":
            v, err := strconv.ParseBool(value.(string))
            if err != nil {
                return nil
            }
            log.Println("CheckTypeValue-boolean: ", v)

            return v
        }
        break

    case interface{}:
        switch type_ {
        case "int":
            log.Println("__CheckTypeValue-int: ", strconv.Itoa(value.(int)))

            return value.(int)
        case "text", "date", "time", "timestamp":
            log.Println("__CheckTypeValue-text-date-time: ", value.(string))

            return value.(string)
        case "boolean":
            log.Println("__CheckTypeValue-boolean: ", value.(bool))

            return value.(bool)
        }
        panic("utils.CheckTypeValue: interface - unknown type: " + type_)
    }
    panic("utils.CheckTypeValue: unknown type: " + type_)
}

func MatchRegexp(pattern, str string) bool {
    result, err := regexp.MatchString(pattern, str)
    HandleErr("utils.MatchRegexp: ", err, nil)

    return result
}

func GetMD5Hash(text string) string {
    hasher := md5.New()
    hasher.Write([]byte(text))

    return hex.EncodeToString(hasher.Sum(nil))
}

func GetRandSeq(size int) string {
    rb := make([]byte, size)
    _, err := rand.Read(rb)
    HandleErr("utils.GetRandSeq: ", err, nil)

    return base64.URLEncoding.EncodeToString(rb)
}

func SendJSReply(result interface{}, responseWriter http.ResponseWriter) {
    response, err := json.Marshal(result)
    if HandleErr("[utils.SendJSReply] Marshal: ", err, responseWriter) {
        fmt.Fprintf(responseWriter, "%s", err.Error())
    } else {
        responseWriter.Header().Set("Content-Type", "application/json")
        fmt.Fprintf(responseWriter, "%s", string(response))
    }
}

func ParseJS(r *http.Request, rw http.ResponseWriter) (request map[string]interface{}, err error) {
    decoder := json.NewDecoder(r.Body)
    err = decoder.Decode(&request)
    if err != nil {
        return nil, err
    }

    return request, nil
}
