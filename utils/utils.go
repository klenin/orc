package utils

import (
    "crypto/md5"
    "crypto/rand"
    "encoding/base64"
    "encoding/hex"
    "fmt"
    "github.com/lib/pq"
    "log"
    "net/http"
    "os"
    "reflect"
    "regexp"
    "strconv"
)

func HandleErr(message string, err error, w http.ResponseWriter) bool {
    if err != nil {
        log.Println(message + err.Error())

        if err, ok := err.(*pq.Error); ok {
            log.Println("pq error:", err.Code.Name())
        }

        if w != nil {
            http.Error(w, fmt.Sprintf(message+"%v\n", err.Error()), http.StatusMethodNotAllowed)
        } else {
            os.Exit(1)
        }

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

func ConvertTypeModel(type_ string, value reflect.Value) (interface{}, bool) {
    switch type_ {
    case "int":
        println("val: ", strconv.Itoa(int(value.Int())))
        return value.Int(), value.Int() != 0
    case "text", "date", "time":
        println("val: ", value.String())
        return value.String(), value.String() != ""
    case "boolean":
        println("val: ", value.Bool())
        return value.Bool(), true
    }
    panic("convertTypeModel: unknown type")
}

func ConvertTypeForModel(type_ string, value interface{}) interface{} {
    switch value.(type) {
    case string:
        // value from grid
        if value.(string) == "_empty" {
            return -1
        }

        switch type_ {
        case "int":
            if value.(string) == "_empty" {
                return -1
            }
            v, err := strconv.Atoi(value.(string))
            if HandleErr("[utils.ConvertTypeForModel] strconv.Atoi: ", err, nil) {
                return nil
            }
            return v
        case "text", "date", "time":
            return value
        case "boolean":
            v, err := strconv.ParseBool(value.(string))
            if HandleErr("[utils.ConvertTypeForModel] strconv.Atoi: ", err, nil) {
                return nil
            }
            return v
        }

    case interface{}:
        switch type_ {
        case "int":
            return value.(int)
        case "text", "date", "time":
            return value.(string)
        case "boolean":
            return value.(bool)
        }
    }
    panic("utils.ConvertTypeForModel: unknown type: " + type_)
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
