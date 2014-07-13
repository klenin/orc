package utils

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
)

func HandleErr(message string, err error, w http.ResponseWriter) {
	if err != nil {
		println(err.Error())
		fmt.Printf(message+"%v\n", err.Error())
		if w != nil {
			http.Error(w, fmt.Sprintf(message+"%v\n", err.Error()), http.StatusMethodNotAllowed)
		} else {
			os.Exit(1)
		}
	}
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

func ArrayStringToInterface(array []string) []interface{} {
	result := make([]interface{}, len(array))
	for i, v := range array {
		result[i] = interface{}(v)
	}
	return result
}

func IsExist(array []string, value string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}
